package builder

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

type templateData struct {
	Module      string
	PackageName string
	Types       []Writable
	TypeNames   []string
	Params      []Writable
	Service     string
	Methods     []*Method
}

func (b *Builder) generateResourceTypes(tag *base.Tag, schemas []*base.SchemaProxy) error {
	types := b.schemasToTypes(schemas)

	typesBuf := bytes.NewBuffer(nil)
	if err := b.templates.ExecuteTemplate(typesBuf, "types.py.tmpl", templateData{
		PackageName: strcase.ToSnake(tag.Name),
		Module:      b.cfg.Module,
		Types:       types,
	}); err != nil {
		return err
	}

	dir := path.Join(b.cfg.Out, strcase.ToSnake(tag.Name))
	typeFileName := path.Join(dir, "types.py")
	typesFile, err := openGeneratedFile(typeFileName)
	if err != nil {
		return err
	}
	defer typesFile.Close()

	if _, err := typesFile.Write([]byte(typesBuf.String())); err != nil {
		return err
	}

	return nil
}

func (b *Builder) generateResource(tagName string, paths *v3.Paths) error {
	if tagName == "" {
		return fmt.Errorf("empty tag name")
	}

	tag := b.tagByTagName(tagName)

	dir := path.Join(b.cfg.Out, strcase.ToSnake(tag.Name))
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	resolvedSchemas := b.schemasByTag[tagName]
	if err := b.generateResourceTypes(tag, b.schemasByTag[tagName]); err != nil {
		return err
	}

	typeNames := make([]string, 0, len(resolvedSchemas))
	for _, s := range resolvedSchemas {
		if name := b.getReferenceSchema(s); name != "" {
			typeNames = append(typeNames, name)
		}
	}
	slices.Sort(typeNames)

	bodyTypes := b.pathsToBodyTypes(paths)
	innerTypes := bodyTypes

	paramTypes := b.pathsToParamTypes(paths)
	innerTypes = append(innerTypes, paramTypes...)

	responseTypes := b.pathsToResponseTypes(paths)
	innerTypes = append(innerTypes, responseTypes...)

	methods, err := b.pathsToMethods(paths)
	if err != nil {
		return fmt.Errorf("convert paths to methods: %w", err)
	}

	slog.Info("generating file",
		slog.String("tag", tag.Name),
		slog.Int("schema_structs", len(typeNames)),
		slog.Int("body_structs", len(bodyTypes)),
		slog.Int("path_params_structs", len(paramTypes)),
		slog.Int("methods", len(methods)),
	)

	serviceBuf := bytes.NewBuffer(nil)
	if err := b.templates.ExecuteTemplate(serviceBuf, "resource.py.tmpl", templateData{
		PackageName: strcase.ToSnake(tag.Name),
		Module:      b.cfg.Module,
		TypeNames:   typeNames,
		Params:      innerTypes,
		Service:     strcase.ToCamel(tag.Name),
		Methods:     methods,
	}); err != nil {
		return err
	}

	serviceFileName := path.Join(dir, "resource.py")
	serviceFile, err := openGeneratedFile(serviceFileName)
	if err != nil {
		return err
	}
	defer serviceFile.Close()

	if _, err := serviceFile.Write([]byte(serviceBuf.String())); err != nil {
		return err
	}

	return nil
}

func (b *Builder) writeClientFile(fname string, tags []string) error {
	f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(0o755))
	if err != nil {
		return fmt.Errorf("create %q: %w", fname, err)
	}
	defer f.Close()

	type resource struct {
		Name    string
		Package string
	}

	resources := make([]resource, 0, len(tags))
	for i := range tags {
		if p := b.pathsByTag[tags[i]]; p.PathItems.Len() == 0 {
			continue
		}
		resources = append(resources, resource{
			Name:    strcase.ToCamel(tags[i]),
			Package: strcase.ToSnake(tags[i]),
		})
	}

	slices.SortFunc(resources, func(a, b resource) int {
		return strings.Compare(a.Name, b.Name)
	})

	if err := b.templates.ExecuteTemplate(f, "client.py.tmpl", map[string]any{
		"PackageName": b.cfg.PkgName,
		"Module":      b.cfg.Module,
		"Version":     b.spec.Info.Version,
		"Resources":   resources,
	}); err != nil {
		return fmt.Errorf("generate client: %w", err)
	}

	return nil
}

func (b *Builder) writeClientPackage(fname string) error {
	if err := os.MkdirAll(path.Dir(fname), os.ModePerm); err != nil {
		return err
	}

	f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(0o755))
	if err != nil {
		return fmt.Errorf("create %q: %w", fname, err)
	}
	defer f.Close()

	if err := b.templates.ExecuteTemplate(f, "client.go.tmpl", map[string]any{
		"Name":        b.cfg.Name,
		"PackageName": b.cfg.PkgName,
		"Module":      b.cfg.Module,
		"Version":     b.spec.Info.Version,
	}); err != nil {
		return fmt.Errorf("generate client: %w", err)
	}

	return nil
}

func openGeneratedFile(filename string) (*os.File, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("get current working directory: %w", err)
	}

	p := filepath.Join(cwd, filename)
	f, err := os.OpenFile(p, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(0o755))
	if err != nil {
		return nil, fmt.Errorf("create %q: %w", p, err)
	}

	return f, nil
}

func (b *Builder) tagByTagName(name string) *base.Tag {
	idx := slices.IndexFunc(b.spec.Tags, func(tag *base.Tag) bool {
		return strings.EqualFold(tag.Name, name)
	})
	tag := &base.Tag{
		Name: name,
	}
	if idx != -1 {
		tag = b.spec.Tags[idx]
	}
	return tag
}

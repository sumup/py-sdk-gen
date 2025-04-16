# Contributing

py-sdk-gen relies roughly on the same architecture as high-level languages. [Parsing](https://en.wikipedia.org/wiki/Parsing) and [Lexical analysis](https://en.wikipedia.org/wiki/Lexical_analysis) is offloaded to [kin-openapi/openapi3]("github.com/getkin/kin-openapi/openapi3") that is responsible for loading and parsing the OpenAPI specs.

Loaded specs are transformed into intermediate representation that's specific to py-sdk-gen and in the final step the intermediate representation is transformed into the Golang code of the SDK.

## Code generation

Code generation is handled by a combination of templates and raw string formatting in the code. Unfortunately, there's no easy way to craft AST and then write it out.

# moldable 🔧

`moldable` builds precise interfaces from any package so you can plug in `mockery`, `gomock`, `moq`, or any other mock tool you like.

> [!WARNING]
> This project is in active development and may contain bugs or breaking changes. We recommend testing thoroughly in your environment before using in production. Issues and contributions are welcome!

## Why?

Most mock generators need an interface to do anything. When a library only exposes concrete structs you must hand-write that interface and keep it in sync with every upstream change. `moldable` creates it in one command and gets out of your way. After that you can use whatever mocking framework you like.

## Features

- Processes entire packages in one command. No need to list every struct; `moldable` automatically finds every exported struct that has methods and builds the matching interface.
- Keeps every setting in a single YAML file so you can generate many packages at once, choose where files land, decide how interfaces are named, and commit that file to version control for identical results on any machine.
- Uses Go's official `go/ast` and `go/types` packages, so the generated file is always syntactically correct.
- Renders every method signature exactly as found in the source (parameter names, types, results, variadic dots).
- Supports generics: type parameters on structs and methods are reproduced with their constraints.
- Preserves type aliases, pointers, slices, maps, channels, embedded structs and any nested combination of them.
- Builds the correct import block automatically, choosing non-conflicting local aliases when the same base name appears from different packages.
- Relies on Go's native package loader, so it honours `go.mod` boundaries, works with vendored code, Go workspaces, and private modules without extra flags.
- Lets you tailor the printed code: pick the package name that appears in generated files, use template names like `{package}.generated.go`, add a suffix (`Client` → `ClientContract`), and place everything in a clean output directory tree.

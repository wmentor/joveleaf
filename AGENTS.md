## Code Style

**IMPORTANT**: All code comments must be written in English!

**IMPORTANT**: When returning errors returned by other functions/methods, use `errors.Wrap` from `github.com/pkg/errors` or `fmt.Errorf` to add context.

**IMPORTANT**: For custom errors in new packages, define them as global variables with `Err` prefix (e.g., `ErrRefRequired`) in `internal/<package>/errors.go`.

**IMPORTANT**: README.md MUST be written in English!

**IMPORTANT**: Use `any` instead of `interface{}` for modern Go style (Go 1.18+). They are equivalent, but `any` is more concise and readable.

# Unit tests

**IMPORTANT**: All tests must start with `t.Parallel()` to enable concurrent test execution.

**IMPORTANT**: All helper functions for tests must be marked with `t.Helper()` to provide better error context.

# Tools

Use `go doc` to explore packages and symbols:

- `go doc <pkg>` -- list public structs, functions, and methods of a package.
- `go doc <pkg>.<Sym>` -- show documentation for a specific symbol.
- `go doc -src <pkg>.<Sym>` -- show full source code of a symbol.
- `go doc -u <pkg>` -- include unexported symbols.
- `go doc -all <pkg>` -- show full documentation for all symbols in a package.

Prefer `gopls` over `grep` for finding usages and references:

- `gopls references <file>:<line>:<col>` -- find all callers/references of a symbol.
- `gopls definition <file>:<line>:<col>` -- jump to a symbol's definition.
- `gopls implementation <file>:<line>:<col>` -- find implementations of an interface.
- `gopls callers <file>:<line>:<col>` -- find callers of a function.
- `gopls call_hierarchy <file>:<line>:<col>` -- display full call hierarchy.
- `gopls workspace_symbol <query>` -- search symbols across the workspace.
- `gopls signature <file>:<line>:<col>` -- show a function's signature.
- `gopls symbols <file>` -- list symbols in a file.

# Important Rules

**IMPORTANT**: Never edit `CHANGELOG.md` manually. Changelog entries are automatically generated and managed by the release process.

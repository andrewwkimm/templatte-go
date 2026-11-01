# Repository instructions

## Workflow

- Use the Makefile targets. Run `make build` after changes; it tidies, applies `go fix`, lints with fixes, compiles, and tests. Report every skipped or failed check.
- When the gopls MCP tools are available, check `go_diagnostics` after each edit and use `go_symbol_references` before changing an exported identifier.
- Every exported identifier is public API. Preserve public behavior unless the request explicitly changes it.
- Do not commit, push, merge, tag, publish, or delete remote data without explicit approval.

## Conventions

- Tests use testify (`require` when the test cannot continue, `assert` otherwise) and `cmp.Diff` for structured values, with hand-written fakes.
- Add a dependency only when the standard library cannot reasonably do the job.
- Add a `depguard` rule in `.golangci.yml` only when a concrete package boundary exists to enforce.
- Before proposing a release tag, run `make build`. After `v1.0.0`, check the change with `gorelease` and use the version it accepts.

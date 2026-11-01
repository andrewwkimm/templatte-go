# templatte-go

My personal, opinionated Go library template; the Go version of [templatte](https://github.com/andrewwkimm/templatte).

## Targets

| Target | Does |
|---|---|
| `make build` | Tidy, apply `go fix`, lint with fixes (including formatting), compile, and test |
| `make ci` | The same checks without changing files, as the CI workflow runs them |
| `make format` | Format only |
| `make fix` | Apply `go fix` modernizations |
| `make test` | Run tests with the race detector in random order |
| `make vuln` | Check for known vulnerabilities |

## Releasing

A Go module release is a Git tag. Record the change in `CHANGELOG.md`, then tag `vX.Y.Z` and push the tag. A breaking change after `v1` needs a new major version, which changes the module path to end in `/v2`.

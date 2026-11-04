# Validation — RFC 11, 4 November 2026

Environment: Linux amd64, official Go 1.23.12 toolchain, dependency-free module.

Executed successfully:

- `go test -race ./...` — all shared lab and CLI package tests passed.
- `go vet ./...` — passed.
- `go build -buildvcs=false ./...` — passed. Plain `go build ./...` first failed
  while reading the enclosing workspace's VCS metadata. Disabling VCS stamping
  was the environment-specific workaround; it does not disable compilation checks.
- `go test -cover ./lab ./cmd/sessionlab` — shared lab coverage 87.2%; session CLI
  package 27.5% (decoder unit tests, with main exercised separately below).
- Built the CLI binary with `-buildvcs=false`, then invoked it as a subprocess:
  default JSON trace, HTML generation, rejection of three malformed input files
  with exit status 1 and no stdout, and HTML injection escaping all passed.
- Sample trace assertions: 11 steps; rejection steps exactly 2, 7, 11; ordered
  drain; zero active connections after primary close.

Generated `examples/sessions-result.json` and `docs/rfc11/index.html` from the
compiled Go executable. No Go server, banking integration, historical wire
interoperability, or load benchmark was run. The static HTML was checked by
generation and content assertions; browser visual/interaction testing was not run.

The existing GitHub Actions workflow runs tests, vet and build on push. Remote
CI result must be inspected separately; local checks do not prove CI success.
New repository creation and Pages configuration require GitHub browser sign-in
which was not available during this run. The source was prepared for public
check-in to the existing shared-suite repository; verify the published commit
and source before treating publication as complete.

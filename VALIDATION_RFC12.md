# RFC 12 validation — 5 November 2026

Official Go 1.23.12, Linux amd64, dependency-free shared module.

Checks executed successfully:

- `go test -race ./...` — all shared suite and CLI tests passed.
- `go vet ./...` — passed.
- `go build -buildvcs=false ./...` — passed. VCS stamping was disabled locally
  because the shared scratch directory is not a clean repository checkout.
- `go test -run '^$' -fuzz FuzzHandshakeAccounting -fuzztime=5s -parallel=2 ./lab`
  — passed; 36,729 executions, no accounting or partial-frame bound failures.
- Compiled CLI subprocess smoke checks — default trace accounting and exact
  sequence of actions, HTML generation, and negative scenario exit code 2.
- CLI unit tests exercise malformed/trailing/unknown-field JSON, null, input and
  configuration bounds, HTML injection escaping, and status 0/1/2 behavior.
- Browser preview: all 10 ticks expanded and collapsed with the bulk control;
  displayed metrics match the Go report (4, 2, 1, 0).

Generated `examples/handshake-result.json` and `docs/rfc12/index.html` directly
from the compiled executable. This validates a teaching state machine, not
historical hardware interoperability, production network latency, durable
banking transactions or real backpressure on a physical interface.

The GitHub Actions workflow runs race tests, vet and build on push. Remote CI
status is checked separately after source publication; local results alone do
not prove CI success. The public static deployment returned success at
https://rfc12-handshake-lab.sadhu-sreenivas.chatgpt.site .
A dedicated repository remains blocked by the GitHub browser sign-in screen;
the actual source is checked into the existing public shared-suite repository.

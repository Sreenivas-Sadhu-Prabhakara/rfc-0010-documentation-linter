# RFC 8 validation — 2026-11-01

Go 1.23.12 linux/amd64. gofmt, go test -race ./..., go vet ./..., and go build -buildvcs=false all passed.
Default late-commit scenario verified caller timed_out/business committed. Success scenario verified caller success/business committed. Contradictory terminal evidence verified exit 1 with no JSON success result.
Existing RFC 5–7 tests passed in the same suite. No public hosting or historical interoperability verified.

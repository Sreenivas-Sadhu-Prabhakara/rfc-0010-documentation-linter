# RFC 9 executed validation — 2026-11-02

Go 1.23.12 linux/amd64. gofmt applied. go test -race ./..., go vet ./..., and go build -buildvcs=false all passed. Full existing lab suite also passed. CLI JSON asserted two established leases before close, zero after.

Tests exercise concurrent capacity saturation, unique per-host numbers, ownership isolation, confirmed-parent requirement, cascading cleanup, reused numbers with new lease IDs, stale confirmations, invalid kinds and limits.

No remote repository commit, hosted deployment, peer protocol or historical interoperability verified.

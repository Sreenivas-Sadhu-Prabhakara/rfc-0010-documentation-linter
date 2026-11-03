# Executed validation

Go 1.23.12 linux/amd64.

- gofmt applied.
- go test -race ./... passed.
- go vet ./... passed.
- go build -buildvcs=false passed (workspace VCS metadata unavailable).
- Conversion fuzzing passed: 24,309 executions over approximately 3 seconds.
- Default interaction CLI verified: HELP in both modes, 9 vs 1 requests, 1890ms vs 210ms modeled completion.
- Replacement conversion verified: José -> Jos? with byte-offset warning.
- Strict conversion verified exit 1.
- Pipeline JSON output reassembled byte-for-byte to README.md.
- Injected handler failure verified exit 1.

No historical network compatibility or production deployment claimed.

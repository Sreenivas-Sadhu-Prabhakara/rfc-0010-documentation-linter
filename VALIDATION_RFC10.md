# RFC 10 validation — 2026-11-03

Go 1.23.12 linux/amd64. gofmt, full go test -race ./..., go vet ./..., build with -buildvcs=false all passed.
Library tests cover metadata, valid short proposals/repeated titles, duplicate IDs, missing replacement targets, self/multi-note cycles and invalid dates.
CLI checks verified valid fixture, null, unknown fields, trailing JSON, malformed JSON and >1 MiB file rejection. All negative cases returned exit 1.
No public check-in or hosted deployment completed.

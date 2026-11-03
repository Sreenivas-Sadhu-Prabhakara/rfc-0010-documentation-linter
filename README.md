# RFC 10 — Documentation Linter

Run `go run ./cmd/notelint -file examples/notes.json`. See [RFC_10.md](RFC_10.md) for semantics and limitations.

# RFC Go Labs 5–10

Three runnable, dependency-free Go CLI exercises. Requires Go 1.22 or newer.
These are modern educational models, not historical protocol implementations.

## Run

```sh
go test -race ./...
go vet ./...
go build -o labs ./cmd/labs
./labs interaction -rtt 200ms
./labs convert -text 'José' -policy replacement
./labs convert -text 'José' -policy strict # expected exit 1
./labs pipeline -text 'hello' -capacity 2
./labs pipeline -file README.md -capacity 1
./labs pipeline -file README.md -fail-at 0 # expected exit 1
```

## RFC 5 — interaction

Inspired by local feedback and message assembly in https://www.rfc-editor.org/rfc/rfc5.html.
Default scenario inserts each letter of HELLO, backspaces twice, inserts P and submits.
Both modes emit HELP. Remote mode sends 9 requests; local sends 1.
The default 200ms RTT plus 10ms processing yields 1890ms versus 210ms simulated completion.
Events are queued at time zero, remote mode is stop-and-wait. This is a deterministic
cost model, not a real-time typing benchmark. display_delay_ms is the last editing
feedback time relative to scenario start, not a percentile or average. Local editing
has zero modeled CPU cost. No sleeps or goroutines are needed for this model.
An optional -events file accepts [{"kind":"insert","text":"hello"},{"kind":"submit"}].
Backspace removes one Unicode rune, not one grapheme. Input limit is 4096 runes.
No remote code executes and no real payment is submitted.

## RFC 6 — conversion

Inspired by conversion flags and overflow discussion in https://www.rfc-editor.org/rfc/rfc6.html.
Toy destination profile is printable ASCII 0x20–0x7E, not EBCDIC or any bank profile.
Strict mode rejects unsupported input. Replacement mode emits ? plus warnings with
UTF-8 byte offsets. Invalid UTF-8 always fails. Empty text succeeds at limit zero.
Overflow and rejected conversions return no usable output. Nothing trims, normalizes,
or changes case. Identity fields should use strict policy unless explicitly designed otherwise.

## RFC 7 — bounded network/handler handoff

Inspired by https://www.rfc-editor.org/rfc/rfc7.html: separate network preparation and
handler responsibilities, 1006-byte text chunks, bounded buffers, link/destination metadata.
RFC Editor status is Unknown; no update/obsoleting relationship appears in its record.
The source is a reconstruction with unreadable portions and unresolved questions.
Our frames use JSON (payload is base64), modern sequence numbers and a bounded Go channel.
They do not reproduce the original 16-bit heading, marking, EBCDIC conversion, interrupt
handling, RFNM, multiplexing fairness, full-duplex behavior or receiver reassembly.
Empty input emits zero frames. Capacity must be 1–1024. Handler failures stop further
handler calls and cancel the producer. Queued frames are discarded on failure.
The callback must return promptly: arbitrary blocked callback code cannot be interrupted.
Transport success is not business completion. Inputs are loaded into memory by the CLI;
queue bounding does not bound whole-process input memory. A streaming reader is a future extension.

## Validation

Tests cover final-command equivalence and exact timing, Unicode editing, limits,
conversion outcomes, exact/over-boundary fragmentation, ordered reconstruction,
handler failure and cancellation. Fuzz target checks conversion bounds and preservation.
Run `go test -fuzz=FuzzConvert -fuzztime=10s ./lab` for fuzzing.
See VALIDATION.md for checks actually executed in the build environment.

## RFC 8

See [RFC_8.md](RFC_8.md). Run `go run ./cmd/acklab` for receipt, timeout and late-commit evidence.

## RFC 9

See [RFC_9.md](RFC_9.md). Run `go run ./cmd/linklab` for owned link allocation, confirmation and cleanup.

## RFC 10

See [RFC_10.md](RFC_10.md). Run `go run ./cmd/notelint -file examples/notes.json`.

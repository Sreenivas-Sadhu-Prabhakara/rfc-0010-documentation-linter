# RFC 10 — documentation linter

Public, MIT-licensed Go educational lab. Requires Go 1.22+.

```bash
git clone https://github.com/Sreenivas-Sadhu-Prabhakara/rfc-0010-documentation-linter.git
cd rfc-0010-documentation-linter
go run ./cmd/notelint -file examples/notes.json
go test -race ./...
go vet ./...
```

This repository includes the shared RFC 5–10 lab library and test harness so every command and test runs independently after cloning. Its focus is RFC 10. These are modern educational models, not historical protocol implementations or production banking software.

Source, tests, examples and original validation reports are committed. GitHub Actions is configured; a workflow file alone does not establish a successful remote run. Earlier validation reports record local checks.

## Publication

Published to this public repository on 3 November 2026. Browser hosting is separate; no live Go backend is claimed. Historical pre-publication remarks in included reports refer to the earlier ZIP delivery.

## Shared lab guide

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

## RFC 11 — session lifecycle and inbound routing

The shared suite now includes RFC 11. Read [RFC_11.md](RFC_11.md) for the historical context, current RFC Editor status, applications, model boundaries and hosting limitations.

```sh
go run ./cmd/sessionlab
go run ./cmd/sessionlab -html > docs/rfc11/index.html
go test -race ./...
```

Open `docs/rfc11/index.html` locally to explore the saved Go-generated trace. It is a static report, not a live backend. Source and meaningful tests were added on 4 November 2026 under this repository's MIT license. The static trace is publicly hosted at https://rfc11-session-lifecycle.sadhu-sreenivas.chatgpt.site . GitHub Actions passed for source commit 41c8c35e2d3e8f5f9aec0dd55ea9f402ac299cd4. A dedicated RFC 11 repository and GitHub Pages configuration remain pending GitHub sign-in; public static hosting through Sites is complete. Covered: RFC 11; next: RFC 12. See [VALIDATION_RFC11.md](VALIDATION_RFC11.md).

## RFC 12 — ready/valid handshake trace validator

Covered: RFC 12; next: RFC 13. See [RFC_12.md](RFC_12.md) for official status, context, applications, tradeoffs and run instructions. The new Go model verifies stable offers across stalls, accepted frame boundaries and reset accounting. It is an educational adaptation, not historical hardware compliance or production banking software.

```sh
go run ./cmd/handshakelab
go run ./cmd/handshakelab -html > docs/rfc12/index.html
go test -race ./...
```

Public static viewer: https://rfc12-handshake-lab.sadhu-sreenivas.chatgpt.site . It displays saved Go output, not a running Go backend. A dedicated RFC 12 repo remains blocked by GitHub browser sign-in; source is published here under the MIT license. See [VALIDATION_RFC12.md](VALIDATION_RFC12.md).

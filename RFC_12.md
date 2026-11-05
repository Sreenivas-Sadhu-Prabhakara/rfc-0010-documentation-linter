# RFC 12 — Ready/valid handshake trace lab

Covered: **RFC 12**. Next: **RFC 13**. No documents skipped.

On 5 November 2026, the latest public shared-suite repository was inspected.
`RFC_11.md` confirms RFC 11 was covered, implemented and hosted. Main pointed to
443fdb07aa87d58c91af6808c6a3feefd8ee3d4d before this work. The latest module,
license, README, scenario runner and Go HTML CLI source were read before adapting
the existing project's JSON-to-static-report approach. RFC 12 introduces a new
state machine rather than changing the RFC 11 implementation.

## Official source and current catalog relationships

- M. Wingfield, **26 August 1969**, *IMP-HOST INTERFACE FLOW DIAGRAMS*.
- Official text: https://www.rfc-editor.org/rfc/rfc12.txt
- Official record: https://www.rfc-editor.org/info/rfc12/
- Machine-readable catalog: https://www.rfc-editor.org/rfc/rfc12.json
- Publication and current status: **UNKNOWN**, source **Legacy**.
- `updates`, `updated_by`, `obsoletes`, `obsoleted_by`: all empty in the
  official metadata retrieved on 5 November 2026.

The text labels itself REVISED; that does not establish an update relationship
to another RFC. The absence of a formal obsoleting relationship does not make
this historical hardware interface a present-day production recommendation.
It predates the modern standards-track category scheme.

## Historical purpose and ideas

Early ARPANET hosts connected through Interface Message Processors (IMPs).
RFC 12 turns the interface logic from Appendix B of BBN Report 1822 into flow
diagrams. Its scope is the ordering of hardware operations, not an application
API or a new host-to-host protocol.

The receive direction initializes counters and flags, requests bits, checks
readiness and the last-bit condition, shifts data into a register, and waits
until a completed word has been taken. Padding and message-end conditions
affect how input finishes. The send direction requests a host word, loads the
shift register and waits for the IMP's bit-level readiness transitions while
shifting the word and tracking its end. These are sequences with explicit
wait states, not an assumption that offering data means it has been consumed.

Today the transferable idea is to specify both the data and the conditions
under which ownership transfers. A slow receiver must not cause a producer to
silently replace unaccepted data. A frame boundary must be accepted before it
can count as a completed frame. Transport acceptance still does not establish
durable banking business completion.

## Three useful applications

1. **Banking streaming adapter trace validator — selected.** Detect a producer
   changing its offered statement chunk while a receiver is stalled. Separate
   accepted chunks from completed batches and show what reset discards.
2. **Cloud queue admission lab.** Compare an offered job, an admitted job and
   completed processing. Useful for explaining why queue acceptance and business
   commit need separate evidence; would require an additional worker model.
3. **Hardware handshake teaching tool.** Explore word-level readiness and framing
   before studying actual asynchronous bit signals. A hardware-accurate emulator
   would need electrical timing and the full BBN interface definitions.

## Built exercise and tradeoffs

`lab/handshake.go` is a deterministic tick-based trace validator. Each tick has
`ready`, `valid`, `data`, `last` and `reset`. A word transfers only when
`ready && valid`. A stalled valid offer must retain both its data and Last flag
and keep Valid asserted until transfer or reset. A transferred Last completes
the frame. Reset discards accepted words in the partial frame, clears any held
unaccepted offer, and preserves counts of already completed frames.

The accounting invariant checked by fuzzing at every tick is:

```
accepted words = completed frame words + reset-discarded words + partial words
```

This is a **modern educational adaptation**, not RFC 12 compliance. The original
asynchronous bit request/high-low handshake is not the synchronous ready/valid
convention used here. We do not reproduce signal names, bit timing, shift
registers, padding, interrupts, checksum, IMP errors or physical hardware.
Last is our modern word-level end marker. Data uses short strings as symbolic
words, not any historical machine word size. No network or bank is contacted.

The model is deliberately single-threaded: replay is deterministic and tests
need no sleeps. Ticks are logical observations, not elapsed milliseconds, a
latency benchmark or a fairness/liveness guarantee. Two consecutive ticks with
ready and valid both high are two transfers, even if their data is identical.
Stable Valid is a new offered word on each ready tick, not edge-triggered delivery.

Limits: 1000 ticks, 128 data bytes per tick, and a configured frame limit of
1..64 words (default 4). Last must be asserted by the word limit. A violation
does not mutate the receiver state: the report keeps the last valid prefix so
a later recovery/reset can be inspected. The violation remains recorded.
This diagnostic rejection behavior does not claim real hardware can undo an
invalid cycle. Reset takes priority and ignores other signals on that tick.

Data/Last are ignored when Valid is false, except the defensive data-size check.
Unaccepted held offers are not counted as accepted or reset-discarded words.
`finished` means no partial frame or held offer remains; inspect `violations`
as well before treating a trace as valid. Empty arrays are valid idle traces.
Changing data during a stall, withdrawing Valid or changing Last is a violation.

## Run and acceptance criteria

Requires Go 1.22+. From the shared repository root:

```sh
go run ./cmd/handshakelab
go run ./cmd/handshakelab -file examples/handshake.json -max-words 4
go run ./cmd/handshakelab -html > docs/rfc12/index.html
go build -o handshakelab ./cmd/handshakelab
./handshakelab -file examples/handshake-invalid.json
go test -race ./...
go vet ./...
go build ./...
go test -run '^$' -fuzz FuzzHandshakeAccounting -fuzztime=5s -parallel=2 ./lab
```

CLI exit codes: 0 for a valid finished trace; 1 for malformed input,
unknown JSON fields, configuration, file or output errors; 2 for a trace with
violations or unfinished work (the diagnostic report is still emitted).
The invalid sample intentionally exits 2. Running that sample via `go run`
prints `exit status 2`; Go's launcher may itself return a different code, so
use the compiled binary to depend on exact process exit codes.

Only one JSON array is accepted; trailing JSON, top-level null and inputs
larger than 1 MiB are rejected. The default scenario contains 10 ticks with
4 accepted words, 2 completed frames, 1 reset-discarded word and 0 violations.
It ends with no held offer or partial frame. The negative scenario changes a
stalled offer, records one violation, then demonstrates a valid recovery.

Tests cover stable stalls, exactly ordered frame completion, immutable trace
snapshots, changed/withdrawn offers, Last mutations, reset accounting, exact and
over-limit inputs, unfinished work, consecutive ready/valid ticks, decoder
validation, CLI exit codes and HTML escaping. The fuzz target checks accepted
word accounting and bounded partial frames over generated sequences.

## Source publication and hosting

Source is MIT licensed in the existing public shared-suite repository:
https://github.com/Sreenivas-Sadhu-Prabhakara/rfc-0010-documentation-linter

The browser's current GitHub new-repository page requires sign-in. The GitHub
connector supports pushes to existing repositories but exposes no repository
creation operation. A dedicated `rfc-0012-handshake-lab` repository remains an
unfinished step; the source is published in the shared suite instead.

`docs/rfc12/index.html` is generated from the Go executable. It is a portable
static trace viewer, not a live Go backend. Sites is used for public static
hosting. The public deployment succeeded on 5 November 2026:
https://rfc12-handshake-lab.sadhu-sreenivas.chatgpt.site

The HTML can also be opened locally. See `VALIDATION_RFC12.md` for executed checks.

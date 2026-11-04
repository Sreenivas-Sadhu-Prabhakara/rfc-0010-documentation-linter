# RFC 11 — Session lifecycle and inbound routing lab

Covered: **RFC 11**. Next document: **RFC 12**. No numbers skipped.
Progress verified on 4 November 2026: the public RFC 10 repository contains the
shared RFC 5–10 suite. Its latest main commit before this work was
8d167d5107582ef3c806f1fd6e01493a6e9f6f2c (a daily-log addition).
The latest remote `lab/links.go`, module file, license and workflow were read;
the reused local allocator and module matched the remote Git blob hashes.

## Official source and status

G. Deloche, UCLA, August 1969: *Implementation of the Host - Host Software
Procedures in GORDO*.

- Official record: https://www.rfc-editor.org/info/rfc11/
- Text: https://www.rfc-editor.org/rfc/rfc11.txt
- Machine-readable metadata: https://www.rfc-editor.org/rfc/rfc11.json
- RFC Editor publication/current status: **UNKNOWN**, source **Legacy**.
- Obsoleted by **RFC 33**. RFC 11 obsoletes no document; its updates and
  updated-by lists are empty in the official metadata checked 4 November 2026.
- RFC 33's record: https://www.rfc-editor.org/info/rfc33/ (it is itself listed
  Unknown and updated by RFCs 36 and 47). This is relationship context only;
  we will study RFC 33 when the sequence reaches it.

Unknown is the catalog status, not a claim that the old design is current.
RFC 11 predates today's standards-track categories. It is an obsolete early
host-networking design, not a production deployment recommendation.

## What it was trying to solve

Several user programs on UCLA's SDS Sigma 7 needed to share network hardware.
The GORDO operating system integrated a Network process that managed logical
connections and a Handler that drove the IMP/HOST hardware interface.
This made network use a software service rather than each application managing
the physical interface itself.

A connection appears full duplex to its user while using two directional
links internally. The paper separates a host-to-host control connection,
a primary connection for interactive use and login, and auxiliary connections
for bulk data. Opening and closing involve control messages and acknowledgments.
Its allocation tables connect host IDs, connection IDs, and incoming links;
buffers and shared pages coordinate applications, Network, and Handler work.
An emergency ring temporarily holds incoming messages when a user is not ready.

The useful architectural lesson today is explicit lifecycle state, unambiguous
correlation, bounded resource use and routing isolation. A transport-level
acknowledgment does not demonstrate business completion. The historical checksum,
ASCII commands, shared pages and interrupt flow are not modern API security or
reliability mechanisms.

## Three possible exercises

1. **Banking adapter session lab — selected:** reserve resources, wait for a
   matching acknowledgment, separate interactive and batch traffic, route safely,
   and clean up on closure. Small enough for deterministic tests.
2. **Cloud worker handoff:** compare a bounded queue with an overflow ring when
   workers are slow. Better for throughput experiments; requires a scheduler and
   a clearly specified fairness policy.
3. **Multiplexed gateway routing audit:** generate colliding incoming link numbers
   on different hosts and verify tenant routing. More targeted, but misses the
   handshake and cleanup lessons.

## Implementation and tradeoffs

`lab/sessions.go` wraps the existing, freshly verified RFC 9 `Links` allocator.
The wrapper adds unique primary sessions per user/host pair, acknowledgment
correlation, host-plus-incoming-link routing, bounded per-session inboxes and
cascading cleanup. A mutex makes each operation atomic. This favors understandable
invariants over maximum concurrency. Capacity includes pending reservations.

The Go CLI executes a finite JSON scenario and emits JSON or a self-contained
HTML trace viewer. The viewer contains saved Go output; it is not a live Go
backend and does not reimplement the state machine in JavaScript.

Educational changes from RFC 11 are deliberate:

- Opaque, monotonically increasing generation IDs reject stale ACKs and stale
  data after link-number reuse. These are modern extensions, not historic fields.
- One supplied user label represents the local user; host identifies the peer.
  Labels are not authentication and there is no remote-user negotiation.
- Negative acknowledgments abandon a pending request instead of automatic retry.
  Closing releases local state immediately, without simulating remote close ACKs.
- A full inbox returns an explicit error. There is no emergency ring, automatic
  retry, RFNM, fairness scheduler, wire transport or business transaction logic.
- 1024 bytes per message and 1..64 queued messages are lab limits. They do not
  reproduce the historical 8095-bit message layout, checksum or padding.
- There is no historical login, privileged instruction, interrupt or GORDO memory
  emulation. No payment, credential or real bank data is used.

## Run

Requires Go 1.22+. From the repository root:

```sh
go run ./cmd/sessionlab
go run ./cmd/sessionlab -file examples/sessions.json -capacity 8 -depth 1
go run ./cmd/sessionlab -html > docs/rfc11/index.html
go test -race ./...
go vet ./...
go build ./...
```

Open `docs/rfc11/index.html` in a browser; expand steps or use Expand all.
No web server or package installation is needed for this static report.

The default trace has 11 steps and three expected rejected operations: early
data, a second message into the full inbox, and a late acknowledgment after
closure. Step 8 drains the first batch; step 9 explicitly retries the second.
Step 10 closes the primary and leaves zero active sessions. The inbox is volatile
and queued data is discarded at close. No durable processing is implied.

Scenario operations: `open`, `ack`, `reject`, `deliver`, `drain`, `close`.
Use `id` to correlate both ACKs and delivered data with the connection generation;
`parent` refers to the established primary ID when opening an auxiliary.
Each successful open allocates the next ID starting at 1. Errors do not allocate.
Unknown JSON fields, malformed input, trailing JSON, invalid configuration,
more than 1000 events or input over 1 MiB fail with exit status 1.
Per-operation semantic failures are recorded in the output and execution
continues with exit status 0; inspect the `error` fields when scripting checks.

## Acceptance criteria and validation

Tests verify that data is inaccessible before establishment; matching ACK replay
is safe; ownership, peer mismatches and conflicting ACKs are rejected; equal link
numbers on distinct hosts remain isolated; full/oversized messages are rejected;
drain returns only owned messages; primary close removes established and pending
children; resource reuse never accepts stale generations; concurrent duplicate
primary opens produce one reservation; and old snapshots remain immutable.
See VALIDATION_RFC11.md for checks actually executed.

## Publication and hosting

Dedicated repository creation and Pages enablement were blocked by the browser's
GitHub sign-in screen. The connector supports source writes to existing
repositories but exposes no repository-creation or Pages-setting operation.
To preserve and publish actual code, this lab was added to the existing public
MIT-licensed shared-suite repository, `rfc-0010-documentation-linter`.

The remaining hosting step is to enable GitHub Pages for `main` / `docs` in
repository settings, then verify the resulting live `/rfc11/` page. A dedicated
`rfc-0011-session-lab` repository also remains to be created and populated when
repository creation is available. Neither operation is claimed complete.

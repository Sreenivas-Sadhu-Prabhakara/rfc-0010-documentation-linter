# RFC 8 — Acknowledgment evidence lab

Source: https://www.rfc-editor.org/info/rfc8/ and https://www.rfc-editor.org/rfc/rfc8.pdf
Verified 2026-11-01. G. Deloche, ARPA Network Functional Specifications, May 1969.
The scan title page says 5 May 1969. Status Unknown, no update/obsoleting
relationships listed. The PDF is handwritten and parts are hard to read; no
claim of complete transcription or historical interoperability is made.

The host-interface and RFNM discussion motivates separating network evidence
from business evidence. The business states below are OUR modern additions.
RFNM is not proof of a payment commit. This lab does not emulate RFNM or ARPANET.

## Run

    go run ./cmd/acklab
    go run ./cmd/acklab -events examples/ack-success.json
    go test -race ./...

Default: receipt at 10ms, caller timeout at 100ms, commit evidence at 200ms.
The result retains caller=timed_out but records business=committed. A timeout
without later business evidence leaves business=unknown. Never blindly resubmit
an actual transaction based on this simulator; there is no real payment system.

Events are a JSON array of {"at_ms":10,"kind":"received"}. Supported kinds:
received, timeout, committed, rejected. Time is a nonnegative, nondecreasing
integer in milliseconds. Ties follow array order. Repeated evidence is safe;
contradictory business outcomes and unknown events fail. Terminal business
evidence can arrive without a transport receipt, because receipt telemetry
may be absent. Caller outcome remains fixed once answered or timed out.

This deterministic offline model has no real timers, retry engine, persistence,
networking, authentication, or bank integration. There is no goroutine to cancel.
Input files are bounded and unknown JSON fields rejected. Invalid input exits 1.

## Applications considered

- Acknowledgment evidence lab: small, banking-relevant; synthetic evidence only.
- Host interface status explorer: useful error classification; less business focus.
- Flow-control simulator: networking relevance; overlaps our RFC 7 pipeline.

## Publication state

Code built and tested locally. Public repository creation and hosting were NOT
completed. Connected GitHub tooling has no creation/Pages controls; the prior
browser authentication attempt was cancelled. No hosted URL is claimed.

## Progress

RFC 7 was read and its pipeline built in the previous session.
This session covers RFC 8; next is RFC 9. Publishing RFC 5–8 remains outstanding.

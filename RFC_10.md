# RFC 10 — Documentation Conventions

Official text: https://www.rfc-editor.org/rfc/rfc10.html
Record: https://www.rfc-editor.org/info/rfc10/ (checked 2026-11-03).
Steve Crocker, 29 July 1969. Category Unknown; obsoletes RFC 3,
obsoleted by RFC 16; updated by RFCs 24, 27 and 30.

This administrative note encourages timely, tentative contributions, retains
basic identifying metadata, assigns serial numbers through Steve Crocker,
and revises distribution information. It is not a network protocol or today's
RFC submission procedure. We will cover the related documents in numerical order.

Applications considered: metadata linter (selected), distribution-list planner,
and replacement-history viewer. A read-only linter provides useful checks without
sending messages or imposing heavyweight approval on short ideas.

## Run

    go run ./cmd/notelint -file examples/notes.json
    go test -race ./...

The tool accepts an array of local proposal records. Each requires a positive
unique id, nonblank title/author/affiliation/body, and YYYY-MM-DD date. Duplicate
titles are allowed and a short question is acceptable. Optional replaces points
to an existing record. Missing targets, self-references and cycles fail. Output
is JSON; invalid records exit 1. No input file is changed and nothing is sent.

Our JSON schema, ISO dates, graph checking, 1 MiB/10000-note bounds and replacement
field are modern additions, NOT RFC 10 requirements. Body validation only checks
nonblank text; it does not evaluate grammar, sentence counts or technical merit.
Metadata is not proof of approval or authenticated authorship. Empty arrays pass.
This is an offline documentation utility, not historical RFC compliance.

Public repository check-in and hosting remain incomplete; no new deployment
or public URL is claimed. Previous browser login was cancelled; available GitHub
tools lack repository creation and Pages configuration. This deliverable is source.

Progress: previous RFC 9. Covered RFC 10; next RFC 11.

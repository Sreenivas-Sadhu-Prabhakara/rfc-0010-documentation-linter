# RFC 9 — Host Software: link ownership lab

Verified source: https://www.rfc-editor.org/info/rfc9/ (2026-11-02).
G. Deloche, May 1969. RFC Editor status Unknown; no updates or obsoleting
relationships listed. The document is a scan. Its official page supplies
page images/OCR; the standalone PDF fetch returned HTTP 429 during this run.
We do not claim a complete transcription or historical protocol conformance.

The document discusses primary/auxiliary links, control link zero, network
service calls, and tables indexed by host, link and user. Allocation is distinct
from establishment after remote confirmation. This inspires the modern lab.

Applications considered: owned link allocator (selected); host status index;
connection-handshake trace viewer. The allocator is small and tests resource
ownership/capacity; unlike the trace viewer it does not model peer dialogue.

Run:

    go run ./cmd/linklab
    go test -race ./...

The demo allocates a primary link, confirms it, allocates and confirms an
auxiliary, then closes the primary. JSON shows two leases before and zero after.

Our additions: mutex protection, monotonically increasing opaque lease IDs,
owner checks, configurable global capacity 1–64, and cascading child cleanup.
Only an established primary owned by the same user and host permits auxiliary
allocation. Network numbers 1–255 are unique per host; zero is never allocated.
Pending allocations count against capacity. Closed network numbers can be reused,
but stale lease IDs cannot confirm or close a new lease. Multiple primary links
per user are allowed. Confirmation is a local method, not real peer negotiation.

The API rejects invalid inputs; the CLI runs a fixed reproducible demo. This is
an in-memory library without persistence, retries, authentication or sockets.
User strings model ownership; they are NOT authenticated identities. No actual
banking data or services are used. Business completion is not represented.

Checks include concurrent capacity pressure, unique number allocation, owner
isolation, parent confirmation, cascade cleanup, reuse and stale confirmations.
Public GitHub check-in/hosting remains outstanding: no new repository or hosted
URL was created. The previous browser sign-in was cancelled and connected tools
lack repository creation/Pages controls. Do not describe a ZIP as published.

Progress: RFC 8 was previously covered; this session covers RFC 9, next RFC 10.

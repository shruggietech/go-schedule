# S130 Storage and Release Contract

## Dependency consumption

Root and desktop select the same SQLite/libc pair from published metadata. Both retain Go1.26.0 and their /v2 identities. Desktop resolves the root through its local parent replacement. Normal -mod=readonly resolution and affected tests succeed without changing any of the four manifest/checksum files.

## Current behavior

Direct initialization is transactional; committed task/audit records survive reopening; refused unsupported state is not converted or deleted. Foreign keys, one production writer and existing WAL remain. Prepared join/ordered operations with 0/1 bindings return the intended results in memory and WAL. Cancellation waiting on a real SQLite lock returns the context error while the blocking lock is still held and does not commit partial state. All existing access/security/cancellation/audit suites remain mandatory.

## Release continuation

The candidate stays v2.0.0. This PR implements the dependency correction, with Closes #404 and Refs #412; it does not merge, tag, stage or promote. After owner integration, fresh exact-main CI and metadata checks precede exact immutable tag/draft approval. Actual downloaded assets, identities and hashes precede separate promotion approval. Existing workflows publish those accepted bytes unchanged. #412 closes only for actual public delivery/discovery. #400 is superseded only after replacement integration. Public v1.7.0 and historical comparison records remain unchanged now.

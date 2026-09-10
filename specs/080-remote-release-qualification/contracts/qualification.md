# Contract: v1.4 Remote Source Qualification

## Required Platform Job

The pull request must run a job named `v1.4 remote qualification (<platform>)` on Windows, macOS, and Linux with fail-fast disabled. Each job builds package-shaped binaries and runs the S080 integration journey plus the detailed remote, enrollment, authorization, client, CLI, profile, and connection suites under the race detector.

## Lifecycle Contract

The package-shaped journey must prove one chronological sequence: default-off start, configured TLS start, local pairing creation, remote phrase exchange, pinned identity verification, capability-appropriate use, underprivileged denial, actor-attributed audit, credential revocation, listener disablement, local restart, and replacement-binary restart with retained identity and state.

Connection refusal alone is not default-off evidence. The same process must answer through its protected local IPC endpoint while the remote endpoint is absent.

## Documentation Contract

The published guide must include exact procedures for service-bound configuration, certificate and key handling, private-network HTTPS, SSH-tunneled HTTPS, reverse proxy, direct public HTTPS, desktop pairing, CLI pairing, direct JSON use, capability choice, credential rotation and revocation, disabling remote access, binary upgrade, backup, and incident recovery.

Examples must be labeled as examples when addresses, paths, certificate tooling, firewall commands, proxy products, or infrastructure are operator choices. No example may disable TLS verification or place a phrase, bearer, or private key in a committed file, shell history, URL, log, or process argument.

## Completion Boundary

S080 qualifies the reviewed source commit. It does not create or qualify a v1.4.0 tag, release archive, installer, checksum manifest, or public release. Those immutable artifacts require separate operator authorization after merge.

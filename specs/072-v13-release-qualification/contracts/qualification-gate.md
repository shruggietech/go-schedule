# Contract: v1.3 Qualification Gate

## Pull-request check identity

The workflow MUST expose one matrix job whose visible name begins with `v1.3 release qualification` and includes the hosted operating-system value. Its matrix MUST contain:

- `ubuntu-latest`
- `macos-latest`
- `windows-latest`

The matrix MUST use `fail-fast: false` and MUST NOT conditionally skip a supported operating system.

## Package-shaped default-state contract

On each matrix entry, the gate MUST run the S072 package-shaped integration test under the race detector. The test MUST:

1. Build the current `goschedd` source into an isolated candidate directory.
2. Launch the candidate with isolated configuration, data, and local IPC paths.
3. Wait for local health through the supported API client.
4. Assert zero notification channels and zero notification deliveries.
5. Assert MCP HTTP is disabled and exposes no endpoint, credential fingerprint, client identity, or access evidence.
6. Stop the daemon and wait for process termination.
7. Launch the same candidate against retained state and repeat steps 3 through 6.

The child daemon MUST be non-interactive. Windows launches MUST use the repository hidden-console helper.

## Detailed notification contract

The named gate MUST run the existing notification and migration suites under the race detector. Together they remain authoritative for webhook success and bounded failure behavior, retry and restart behavior, disabled-channel behavior, secret and task-data redaction, source-run outcome preservation, and forward migration from the pre-notification schema.

## Detailed MCP contract

The named gate MUST run the existing MCP Observe, stdio, and HTTP suites under the race detector. Together they remain authoritative for both supported protocol revisions, five resources, four continuation templates, zero mutation tools, official SDK interoperability, credential, Host, and Origin rejection, rotation, revocation, disablement, disconnect, restart lifecycle, redaction, hostile-content isolation, and bounded errors.

## Pass contract

A platform job passes only when its package-shaped journey and all focused suites pass. The release boundary passes only when all three platform jobs and the canonical eight local verification gates pass. S072 produces no tag, GitHub release, or public artifact.

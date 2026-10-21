# S110 CI correction: Windows native peer fixture

Published issue: [#339](https://github.com/shruggietech/go-schedule/issues/339), parent #277, milestone 15.

## Outcome

Make the native credential helper's forged-parent fixture reliable on Windows while retaining actual OS peer authentication, direct-entry refusal and zero credential-backend access for an untrusted peer.

## Evidence and authority

Parent: #277. Follow-up to #309, discovered during owner-authorized S110 CI monitoring in PR #338. Pending head cb614525cfe650330a0a2847568069a11122ae2c failed [the Windows race job](https://github.com/shruggietech/go-schedule/actions/runs/39669525921/job/119463385548). TestReadHelperRejectsDirectEntryAndUntrustedLivePeer fails at context_test.go:123 with use of closed network connection after 2.01s. The other 21 checks succeeded. The preceding head 8badade passed all 22 checks, so a passing repetition does not remove this failure.

The fixture launches serveAuthenticatedHelper and assumes listener.Accept will return before correct rejection closes the client. go-winio can discard an immediately closed client and wait for another connection; the fixture's two-second context then closes the listener. This is a fixture setup race, not evidence that untrusted credentials were disclosed.

A deterministic Windows race reproduction against unchanged source waits for the helper's rejection before calling Accept. Insert the following immediately after starting the existing helper goroutine, before listener.Accept, preserving the code for the later assertion:

```go
rejected := <-done
done <- rejected
```

Run go test -race ./internal/clientsecret -run ^TestReadHelperRejectsDirectEntryAndUntrustedLivePeer$ -count=1. It reproduces the same closed-listener failure. No actual user credential is accessed; the fixture uses readBackend.

## Bounded scope

Establish both real native endpoints before invoking peer rejection. Extract the existing established-connection authentication handler only as needed, preserving its kernel peer PID, actual-parent and executable-image checks, cancellation, owned close and request framing. Keep direct flag/standard-stream refusal, legitimate helper behavior, and the assertion that the forged peer never reaches the backend. Do not skip the Windows case, accept arbitrary setup failures as success, serialize the repository suite or loosen authentication.

## Functional acceptance criteria

- The live forged-parent fixture establishes its native channel before rejection can close that channel, so correct early refusal cannot be misreported as a setup failure.
- The actual untrusted peer is rejected before credential storage access; direct standard-stream invocation remains refused and legitimate authenticated exchange retains its behavior.
- The fixture closes its native endpoints and joins its owned work on completion and failure within the existing bounded context.

## Dependencies

Parent #277; follows the delivered #309 private-channel boundary. Independent of #262, #271 and #335. The bounded correction belongs to PR #338 under the existing owner authorization to resolve CI failures. Both authorized Codex rounds are already complete; no third round is requested. Preserve the failing hosted record and distinguish the later correction head and checks.

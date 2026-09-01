# Contract: Windows Installed Core Verification

## Restricted pipe descriptor

Input:

- configured local group name;
- resolved configured-group SID and group account type;
- direct member records containing SID and SID usage.

Output:

- protected DACL with SYSTEM and Built-in Administrators full-access ACEs;
- one read/write ACE for the configured group;
- one read/write ACE for each unique valid direct user member SID, ordered deterministically;
- restricted `AccessInfo` naming only the configured group.

Failure contract:

- missing, non-group, or invalid configured-group SID fails before opening the pipe;
- membership enumeration failure fails before opening the pipe;
- invalid direct-user SID fails before opening the pipe;
- no failure path broadens to compatibility mode.

## Executor diagnostic

| Outcome | Exit code | Output contract |
| --- | --- | --- |
| Success | `0` | Captured child output, bounded by configured cap |
| Child exits nonzero | concrete nonzero value | Captured child output; unchanged existing behavior |
| Process cannot start | absent | `process start failed for "<executable>": <wrapped OS error>` when the child emitted no output |
| Unsupported Windows `run_as` | absent | Existing `run_as:` diagnostic |

The process-start diagnostic must not include task arguments, stdin, or environment values.

## Native probe completion

The installed-core probe is `proven` only when all of the following are true:

1. Candidate and installed identity match.
2. `goschedd` is running as LocalSystem.
3. The normal client reaches the real installed named pipe.
4. A manually triggered production run succeeds with expected record and marker.
5. A scheduled production run succeeds with expected record and a distinct marker observation.
6. A nonzero child exit retains its concrete exit code and expected output.
7. A missing executable produces a no-exit-code process-start diagnostic.
8. Evidence contains required host, security, task, run, filesystem, and log fields.

Any unmet assertion yields `failed`; missing host/provenance prerequisites yield `unavailable`. Neither state may be represented as proven.

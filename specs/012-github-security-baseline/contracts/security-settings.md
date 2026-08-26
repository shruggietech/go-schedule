# Contract: Hosted Security Activation

Remote activation occurs only after the autopilot pre-push halt and explicit
operator authorization.

## Requested settings

| Control | Requested value | Verification |
| --- | --- | --- |
| Private vulnerability reporting | enabled | Repository private-reporting API and route inspection |
| Secret scanning | enabled | `security_and_analysis.secret_scanning.status` |
| Push protection | enabled | `security_and_analysis.secret_scanning_push_protection.status` |
| Non-provider patterns | enabled | `security_and_analysis.secret_scanning_non_provider_patterns.status` |
| Validity checks | enabled | `security_and_analysis.secret_scanning_validity_checks.status` |
| CodeQL advanced setup | workflow committed and successful | Workflow/check run plus code-scanning analysis API when token scope permits |

## Evidence classification

- `enabled`: read-back or hosted check proves the control active.
- `unavailable`: GitHub explicitly rejects or omits the capability for the
  repository plan; record the exact response and compensating control.
- `unverified`: the request may have succeeded, but the available token or API
  cannot prove it.
- `failed`: a supported control or workflow did not activate or complete.

## Safety rules

- Never expose alert, secret, or vulnerability contents in evidence.
- Report only aggregate alert counts when accessible.
- Do not create a sample advisory, commit a test secret, dismiss an alert, or
  change organization settings.
- A missing token scope is `unverified`, not `unavailable` and not `enabled`.
- A plan limitation must name the affected control individually.

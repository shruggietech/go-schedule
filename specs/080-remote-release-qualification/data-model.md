# Data Model: Remote Access Release Qualification

## Service Configuration Binding

| Field | Rule |
| --- | --- |
| Configuration path | Existing readable and valid file at installation time, stored as an absolute path |
| Daemon arguments | Empty for default installation or exactly `--config`, `<absolute-path>` for configured installation |
| Credential material | Never accepted or stored in service arguments |
| Replacement behavior | A later install replaces the service definition through the existing service adapter |

## Qualification Stage

| Stage | Required state | Evidence |
| --- | --- | --- |
| Default | Local IPC healthy, remote listener absent | Package-shaped process and connection probe |
| Enabled | TLS listener healthy, expected daemon identity present | Manifest and TLS-pinned client |
| Paired | One active actor and credential with selected capability | Enrollment result and local administration read |
| Used | Authorized request succeeds and denied request has no side effect | Remote response and audit record |
| Revoked | Former bearer is rejected, local IPC remains healthy | Stable remote error and local health |
| Disabled | Restarted daemon has no listener, retained local state remains | Connection refusal plus local manifest and task read |
| Upgraded | Replacement binary reads the same configuration and state | Preserved daemon identity and records |

## Deployment Mode

| Mode | Recommendation | Product boundary | Operator boundary |
| --- | --- | --- | --- |
| Private-network HTTPS | Recommended | TLS, bearer authentication, capabilities, limits, audit | Private routing, DNS, firewall, certificate lifecycle |
| SSH-tunneled HTTPS | Recommended | Same application boundary | SSH access and tunnel lifecycle |
| Reverse-proxied HTTPS | Advanced supported | Application TLS and immediate-peer trust | Proxy, public and backend certificates, forwarding, monitoring |
| Direct public HTTPS | Advanced supported | Explicit bind and exposure acknowledgement | Public DNS, trusted certificate, renewal, firewall, monitoring, incident response |

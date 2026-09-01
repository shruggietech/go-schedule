# Windows Lifecycle Evidence Contract

Evidence records date, Windows version, scenario, artifact origin/path/hash,
and prior v0.9.0 provenance when applicable. Every MSI operation records label,
mode, exit code, absolute verbose-log path, and log SHA-256.

Relevant phases record local group name/SID/members, intended account
membership, service existence/state/start type/account, product registration,
install directory, PATH state, forbidden access-denied/26421/rollback signals,
local-group action evidence, and ordering before service startup.

After the intended account refreshes its login session, a non-elevated access
probe records that the current token contains the group SID and that the client
reaches the restricted daemon. Elevated client success is not accepted as this
evidence.

Fresh and upgrade scenarios run on separately reset disposable hosts. Every
required observation is proven, failed, or unavailable with a reason. Missing
signals and unavailable prerequisites cannot produce a proven status.

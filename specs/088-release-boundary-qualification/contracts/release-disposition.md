# Release Disposition Contract

Retain the version heading and tagged changelog URL, with four to six single-line highlights. Describe cumulative improvements without claiming the draft is already public.

Use Sandbox without restart, host virtualization changes, active-host product installation, or security-setting changes. Untested native items may ship under the explicit maintainer waiver, but remain untested, not passed. Known failures are not silently approved by a missing-test waiver.

The preparation PR uses Refs #231, Refs #228, and Refs #226, not closing keywords. Its body discloses the waiver, merge-before-staging dependency, actual validation, and remaining criteria. Submit UTF-8 Markdown through a body file and verify stored content.

Existing validators remain authoritative for fully qualified promotion. This preparation does not make partial evidence pass their full-qualification contract. Any future waiver-aware promotion path must preserve exact-asset and reviewed-source verification and expose waived observations explicitly.

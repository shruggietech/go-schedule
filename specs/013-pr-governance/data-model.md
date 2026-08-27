# Process Model: Lightweight Pull-Request Workflow

This feature adds no application or hosted-settings data.

| Stage | Outcome |
| --- | --- |
| review branch | Local work stays separate from `main` |
| local verification | `sh scripts/verify.sh all` passes before publication |
| publication halt | Maintainer authorizes branch push and PR creation |
| pull request | AI reviewers receive a durable review venue |
| review response | Warranted feedback is implemented; other feedback is answered |
| merge | Maintainer performs final review and chooses whether to merge |
| housekeeping | Merged branch is removed and local `main` is synchronized |

No branch-protection state, approval count, required-context set, or
conversation-resolution state is introduced.

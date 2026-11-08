# Runtime-event evidence contract

Existing incident capture remains opt-in. Only an actually launched, failed Windows PowerShell child may query Application events. Both the query and local parser require the exact child PID, observed launch/exit interval, .NET Runtime provider, event 1023 or 1025, and matching application basename. Malformed data is failed evidence, no match is absent, unavailable query is unavailable, and output/count overflow is truncated. No such condition changes the child's native result.

The collector writes runtime-events.json through the existing private exclusive artifact API. It retains no arbitrary message/stack, machine, account, command-line or environment data. Event evidence does not supply a runtime-module hash, a dump, a root cause or an upstream-fix claim. Non-Windows and disabled/launch-failure paths do not query Windows.

# Data model

A runtime-event receipt names the owned target PID/application, observed UTC launch/exit interval, source channel, collection status and duration, and at most sixteen sanitized events. Each event retains provider, event ID, record ID, PID, UTC timestamp, application basename, CoreCLR version and .NET version. The receipt is an additional bounded artifact in the existing incident envelope; dump Collection.Status and Original native result remain independent.

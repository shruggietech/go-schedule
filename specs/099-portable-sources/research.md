# Research

The v1 bundle package already owns canonicalization, validation, comparison, and plan binding. Existing store APIs can create standalone triggers and trigger sets with fresh keys, update watchers while preserving a target path, and replace notification assignments for one scope. Extending that contract is smaller and safer than introducing an independent transfer service.

Watchers cannot be instantiated from exported data alone because the store requires an absolute local path. Preview therefore accepts a target-local path map held with its single-use plan. Channel names are the only available nonsecret external reference for notification assignments; ambiguous or absent names are conflicts.

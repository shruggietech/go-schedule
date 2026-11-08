# Evaluation data model

An evaluation owns one absolute target root, a mode, a per-path read/parse cache and an ordered finding list. Workflow nodes retain YAML line identities. Lifecycle entries derive directory identity, state, delivery, disposition and checked/open task counts; inventory rows are matched to those identities. A new evaluation discards all cached content.

Findings carry file/rule context; read/decode/child failures remain findings and return nonzero status. Timing observations retain source/command/start/end/native result and are diagnostic records, not persisted runtime state or user data.

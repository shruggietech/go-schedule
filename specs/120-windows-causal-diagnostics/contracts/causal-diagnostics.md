# Causal diagnostics contract

The ordinary script/API interface and successful output remain compatible. Opt-in outer process diagnostics contain only owned fixture executable/harness identity, PID, runtime/source hashes, completed phases, bounded original output, actual native exit and separate collection/cleanup status. They do not dump environment variables, real credentials or arbitrary execution inputs. A root exit and a blocked inherited stream are separate facts.

Affected request-aware persistence operations use the caller's cancellation context. Legacy persistence callers retain their existing signatures through explicit Background delegates. Canceled work cannot silently continue through a blocked connection acquisition, complete-state assertions cannot accept partial data, and audit intent must precede mutation. An actual mutation's terminal audit update uses an independently bounded completion path; failed persistence leaves truthful uncertainty and an attributable error.

Controlled diagnostic regressions acquire only the test-owned pool/process. Native IPC cancellation must be observed by the server before the held connection is released. Always release acquired resources and stop owned descendants. No public profiling listener, installed-service operation or global crash setting is introduced. Every result names its exact source and changed conditions.

Integration observation allows 120 seconds for execution, two seconds for inherited-stream drain and one MiB per stream. Overflow, drain timeout and cleanup failure remain separate from the root native exit. Terminal audit completion uses a fresh 30-second context after handler execution, preserving uncertain intent if completion cannot persist.

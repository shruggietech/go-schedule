# Windows Release Gate Fixtures

The `passing` fixture is deliberately inert. Its `.msi` file is plain text and
every observation says it is automated fixture data. It validates schema,
identity, scenario, attachment, and mutation behavior only. It is never native
Windows 11 evidence and cannot authorize release promotion because its source
commit and workflow run are synthetic.

Go mutation tests clone the complete record in memory and alter identity,
outcomes, timing, display bounds, environment roles, paths, and digests. This
keeps each negative case focused while the checked-in passing record remains a
human-reviewable example of schema version 1.

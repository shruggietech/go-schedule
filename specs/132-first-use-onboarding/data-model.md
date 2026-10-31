# S132 Data Model

Existing RunRecord.id is the selected run identity; sourceRunId refers to a separate upstream run. Both remain text with existing access boundaries and no persistence change.

Example workspace holds source.txt, snapshot.txt and snapshot.cmd or snapshot.sh. A fresh owned workspace is required. One snapshot task retains its ID across success, missing-input failure and restored-input success. Runs retain independent IDs and terminal evidence; recovery never rewrites previous records.

CREATE TABLE actors (
	id           TEXT PRIMARY KEY,
	kind         TEXT NOT NULL CHECK(kind IN ('local_os','desktop','cli','json','mcp')),
	display_name TEXT NOT NULL,
	capability   TEXT NOT NULL CHECK(capability IN ('observe','operate','manage','enroll')),
	state        TEXT NOT NULL CHECK(state IN ('active','expired','revoked')),
	builtin      INTEGER NOT NULL DEFAULT 0 CHECK(builtin IN (0,1)),
	created_at   TEXT NOT NULL,
	updated_at   TEXT NOT NULL,
	expires_at   TEXT
);

CREATE TABLE alerts (
	id           TEXT PRIMARY KEY,
	task_id      TEXT,
	severity     TEXT NOT NULL,
	kind         TEXT NOT NULL,
	message      TEXT NOT NULL,
	created_at   TEXT NOT NULL,
	acknowledged INTEGER NOT NULL DEFAULT 0
, run_id TEXT);

CREATE TABLE audit_events (
	id             TEXT PRIMARY KEY,
	actor_id       TEXT,
	daemon_id      TEXT NOT NULL,
	operation      TEXT NOT NULL,
	target_kind    TEXT NOT NULL,
	target_id      TEXT NOT NULL DEFAULT '',
	result         TEXT NOT NULL CHECK(result IN ('uncertain','succeeded','failed','denied')),
	correlation_id TEXT NOT NULL,
	occurred_at    TEXT NOT NULL,
	completed_at   TEXT
);

CREATE TABLE client_credentials (
	id          TEXT PRIMARY KEY,
	actor_id    TEXT NOT NULL UNIQUE REFERENCES actors(id) ON DELETE CASCADE,
	digest      BLOB NOT NULL UNIQUE,
	fingerprint TEXT NOT NULL,
	state       TEXT NOT NULL CHECK(state IN ('active','revoked')),
	created_at  TEXT NOT NULL,
	updated_at  TEXT NOT NULL,
	last_used_at TEXT,
	expires_at  TEXT,
	revoked_at  TEXT
);

CREATE TABLE completion_chains (
	id             TEXT PRIMARY KEY,
	source_task_id TEXT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
	target_task_id TEXT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
	on_outcome     TEXT NOT NULL,
	created_at     TEXT NOT NULL,
	updated_at     TEXT NOT NULL,
	UNIQUE(source_task_id, target_task_id, on_outcome)
);

CREATE TABLE completion_deliveries (
	id             TEXT PRIMARY KEY,
	chain_id       TEXT NOT NULL,
	source_task_id TEXT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
	target_task_id TEXT NOT NULL,
	source_run_id  TEXT NOT NULL REFERENCES runs(id) ON DELETE CASCADE,
	state          TEXT NOT NULL,
	attempts       INTEGER NOT NULL DEFAULT 0,
	created_at     TEXT NOT NULL,
	claimed_at     TEXT,
	completed_at   TEXT,
	target_run_id  TEXT REFERENCES runs(id) ON DELETE SET NULL,
	resolution     TEXT NOT NULL DEFAULT '',
	UNIQUE(chain_id, source_run_id)
);

CREATE TABLE daemon_identity (
	singleton       INTEGER PRIMARY KEY CHECK(singleton = 1),
	installation_id TEXT NOT NULL UNIQUE,
	display_name    TEXT NOT NULL,
	created_at      TEXT NOT NULL,
	updated_at      TEXT NOT NULL
);

CREATE TABLE external_trigger_sets (
	id             TEXT PRIMARY KEY,
	name           TEXT NOT NULL,
	target_task_id TEXT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
	created_at     TEXT NOT NULL,
	updated_at     TEXT NOT NULL
);

CREATE TABLE external_triggers (
	id             TEXT PRIMARY KEY,
	name           TEXT NOT NULL,
	key            TEXT NOT NULL UNIQUE,
	target_task_id TEXT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
	enabled        INTEGER NOT NULL DEFAULT 1,
	created_at     TEXT NOT NULL,
	updated_at     TEXT NOT NULL
, set_id TEXT REFERENCES external_trigger_sets(id) ON DELETE CASCADE, set_position INTEGER);

CREATE TABLE filesystem_watchers (
	id             TEXT PRIMARY KEY,
	name           TEXT NOT NULL,
	kind           TEXT NOT NULL CHECK(kind IN ('file','directory')),
	path           TEXT NOT NULL,
	pattern        TEXT NOT NULL DEFAULT '',
	recursive      INTEGER NOT NULL DEFAULT 0,
	debounce_ns    INTEGER NOT NULL,
	stability_ns   INTEGER NOT NULL,
	target_task_id TEXT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
	enabled        INTEGER NOT NULL DEFAULT 1,
	created_at     TEXT NOT NULL,
	updated_at     TEXT NOT NULL
);

CREATE TABLE groups (
	id          TEXT PRIMARY KEY,
	name        TEXT NOT NULL,
	parent_id   TEXT REFERENCES groups(id) ON DELETE CASCADE,
	enabled     INTEGER NOT NULL DEFAULT 1,
	created_at  TEXT NOT NULL,
	updated_at  TEXT NOT NULL
);

CREATE TABLE notification_assignments (
	id                         TEXT PRIMARY KEY,
	channel_id                 TEXT NOT NULL REFERENCES notification_channels(id) ON DELETE CASCADE,
	task_id                    TEXT REFERENCES tasks(id) ON DELETE CASCADE,
	group_id                   TEXT REFERENCES groups(id) ON DELETE CASCADE,
	on_success                 INTEGER NOT NULL DEFAULT 0 CHECK(on_success IN (0,1)),
	on_failure                 INTEGER NOT NULL DEFAULT 0 CHECK(on_failure IN (0,1)),
	created_at                 TEXT NOT NULL,
	updated_at                 TEXT NOT NULL,
	failure_threshold          INTEGER NOT NULL DEFAULT 1 CHECK(failure_threshold BETWEEN 1 AND 100),
	on_failure_to_start        INTEGER NOT NULL DEFAULT 0 CHECK(on_failure_to_start IN (0,1)),
	duration_threshold_seconds INTEGER NOT NULL DEFAULT 0 CHECK(duration_threshold_seconds BETWEEN 0 AND 2592000),
	on_recovery                INTEGER NOT NULL DEFAULT 0 CHECK(on_recovery IN (0,1)),
	reminder_interval_seconds  INTEGER NOT NULL DEFAULT 0 CHECK(reminder_interval_seconds BETWEEN 0 AND 2592000),
	quiet_period_seconds       INTEGER NOT NULL DEFAULT 0 CHECK(quiet_period_seconds BETWEEN 0 AND 2592000),
	CHECK ((task_id IS NOT NULL) != (group_id IS NOT NULL)),
	CHECK (on_success = 1 OR on_failure = 1 OR on_failure_to_start = 1 OR duration_threshold_seconds > 0)
);

CREATE TABLE "notification_channels" (
 id TEXT PRIMARY KEY,
 name TEXT NOT NULL,
 kind TEXT NOT NULL CHECK(kind IN ('webhook','smtp')),
 endpoint TEXT NOT NULL,
 endpoint_summary TEXT NOT NULL,
 authorization TEXT NOT NULL DEFAULT '',
 enabled INTEGER NOT NULL DEFAULT 1,
 created_at TEXT NOT NULL,
 updated_at TEXT NOT NULL,
 health_interval_seconds INTEGER NOT NULL DEFAULT 0 CHECK(health_interval_seconds = 0 OR health_interval_seconds BETWEEN 60 AND 86400),
 smtp_config TEXT NOT NULL DEFAULT ''
);

CREATE TABLE notification_condition_states (
	task_id              TEXT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
	channel_id           TEXT NOT NULL REFERENCES notification_channels(id) ON DELETE CASCADE,
	policy_fingerprint   TEXT NOT NULL,
	consecutive_failures INTEGER NOT NULL DEFAULT 0,
	active_condition     TEXT NOT NULL DEFAULT '' CHECK(active_condition IN ('','failure','failure_to_start','consecutive_failure','duration_exceeded')),
	active_since         TEXT,
	last_notified_at     TEXT,
	last_evaluated_run_id TEXT NOT NULL DEFAULT '',
	updated_at           TEXT NOT NULL,
	PRIMARY KEY(task_id,channel_id)
);

CREATE TABLE notification_daemon_health_states (
	channel_id       TEXT PRIMARY KEY REFERENCES notification_channels(id) ON DELETE CASCADE,
	next_due_at      TEXT NOT NULL,
	last_delivery_id TEXT NOT NULL DEFAULT '',
	updated_at       TEXT NOT NULL
);

CREATE TABLE notification_deliveries (
	id                   TEXT PRIMARY KEY,
	channel_id           TEXT REFERENCES notification_channels(id) ON DELETE SET NULL,
	channel_name         TEXT NOT NULL,
	destination_summary  TEXT NOT NULL,
	endpoint             TEXT NOT NULL DEFAULT '',
	authorization        TEXT NOT NULL DEFAULT '',
	event_kind           TEXT NOT NULL CHECK(event_kind IN ('run.completed','test','daemon.health')),
	task_id              TEXT,
	run_id               TEXT,
	task_name            TEXT NOT NULL DEFAULT '',
	group_id             TEXT NOT NULL DEFAULT '',
	group_name           TEXT NOT NULL DEFAULT '',
	payload              BLOB NOT NULL,
	state                TEXT NOT NULL CHECK(state IN ('pending','claimed','succeeded','failed')),
	attempts             INTEGER NOT NULL DEFAULT 0,
	next_attempt_at      TEXT NOT NULL,
	created_at           TEXT NOT NULL,
	claimed_at           TEXT,
	completed_at         TEXT,
	last_status          INTEGER NOT NULL DEFAULT 0,
	last_error           TEXT NOT NULL DEFAULT '',
	condition_kind       TEXT NOT NULL DEFAULT '' CHECK(condition_kind IN ('','success','failure','failure_to_start','duration_exceeded','recovery','daemon_health','consecutive_failure')),
	condition_summary    TEXT NOT NULL DEFAULT ''
, kind TEXT NOT NULL DEFAULT 'webhook' CHECK(kind IN ('webhook','smtp')), smtp_config TEXT NOT NULL DEFAULT '', recipients TEXT NOT NULL DEFAULT '[]', diagnostic TEXT NOT NULL DEFAULT '');

CREATE TABLE pairing_sessions (
	id                 TEXT PRIMARY KEY,
	display_name       TEXT NOT NULL,
	kind               TEXT NOT NULL CHECK(kind IN ('desktop','cli','json','mcp')),
	capability         TEXT NOT NULL CHECK(capability IN ('observe','operate','manage','enroll')),
	salt               BLOB NOT NULL,
	verifier           BLOB NOT NULL,
	attempts_remaining INTEGER NOT NULL CHECK(attempts_remaining >= 0),
	state              TEXT NOT NULL CHECK(state IN ('active','consumed','cancelled','expired','exhausted')),
	created_at         TEXT NOT NULL,
	expires_at         TEXT NOT NULL,
	completed_at       TEXT
, grant_expires_at TEXT);

CREATE TABLE portable_identities (
	object_kind TEXT NOT NULL,
	object_id   TEXT NOT NULL,
	portable_id TEXT NOT NULL UNIQUE,
	created_at  TEXT NOT NULL,
	PRIMARY KEY (object_kind, object_id)
);

CREATE TABLE runs (
	id            TEXT PRIMARY KEY,
	task_id       TEXT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
	scheduled_for TEXT NOT NULL,
	started_at    TEXT,
	ended_at      TEXT,
	outcome       TEXT NOT NULL,
	exit_code     INTEGER,
	output        TEXT NOT NULL DEFAULT '',
	trigger       TEXT NOT NULL
, source_task_id TEXT, source_run_id TEXT, output_truncated INTEGER NOT NULL DEFAULT 0, source_trigger_id TEXT, source_watcher_id TEXT, start_failed INTEGER NOT NULL DEFAULT 0 CHECK(start_failed IN (0,1)));

CREATE TABLE schedules (
	id            TEXT PRIMARY KEY,
	kind          TEXT NOT NULL,
	rrule         TEXT,
	anchor        TEXT,
	run_at        TEXT,
	trigger_id    TEXT,
	human_summary TEXT NOT NULL DEFAULT ''
, expression TEXT NOT NULL DEFAULT '', calendar_adjustment TEXT NOT NULL DEFAULT '', elapsed_epoch TEXT);

CREATE TABLE schema_version (version INTEGER NOT NULL);

CREATE TABLE "tasks" (
	id                  TEXT PRIMARY KEY,
	name                TEXT NOT NULL,
	group_id            TEXT REFERENCES groups(id) ON DELETE SET NULL,
	command             TEXT NOT NULL,
	args_json           TEXT NOT NULL DEFAULT '[]',
	working_dir         TEXT NOT NULL DEFAULT '',
	env_json            TEXT NOT NULL DEFAULT '{}',
	run_as              TEXT NOT NULL DEFAULT '',
	enabled             INTEGER NOT NULL DEFAULT 1,
	timezone            TEXT NOT NULL DEFAULT 'Local',
	schedule_id         TEXT REFERENCES schedules(id) ON DELETE CASCADE,
	overlap_policy      TEXT NOT NULL DEFAULT 'queue_one',
	catchup_policy      TEXT NOT NULL DEFAULT 'one',
	state               TEXT NOT NULL DEFAULT 'active',
	created_at          TEXT NOT NULL,
	updated_at          TEXT NOT NULL,
	missing_date_policy TEXT NOT NULL DEFAULT 'skip',
	time_basis          TEXT NOT NULL DEFAULT 'wall_clock',
	dst_gap_policy      TEXT NOT NULL DEFAULT 'next_valid',
	dst_overlap_policy  TEXT NOT NULL DEFAULT 'first',
	stdin               TEXT NOT NULL DEFAULT ''
);

CREATE UNIQUE INDEX idx_actors_builtin ON actors(builtin) WHERE builtin=1;

CREATE INDEX idx_actors_state ON actors(state,created_at,id);

CREATE INDEX idx_alerts_ack ON alerts(acknowledged);

CREATE INDEX idx_audit_events_actor ON audit_events(actor_id,occurred_at,id);

CREATE INDEX idx_audit_events_operation ON audit_events(operation,occurred_at,id);

CREATE INDEX idx_audit_events_order ON audit_events(occurred_at,id);

CREATE INDEX idx_audit_events_result ON audit_events(result,occurred_at,id);

CREATE INDEX idx_client_credentials_state ON client_credentials(state,actor_id);

CREATE INDEX idx_completion_chains_source ON completion_chains(source_task_id, on_outcome);

CREATE INDEX idx_completion_chains_target ON completion_chains(target_task_id);

CREATE INDEX idx_completion_deliveries_state ON completion_deliveries(state, created_at);

CREATE INDEX idx_external_trigger_sets_target ON external_trigger_sets(target_task_id);

CREATE UNIQUE INDEX idx_external_triggers_set_position ON external_triggers(set_id,set_position) WHERE set_id IS NOT NULL;

CREATE INDEX idx_external_triggers_target ON external_triggers(target_task_id);

CREATE INDEX idx_filesystem_watchers_enabled ON filesystem_watchers(enabled);

CREATE INDEX idx_filesystem_watchers_target ON filesystem_watchers(target_task_id);

CREATE INDEX idx_groups_parent ON groups(parent_id);

CREATE INDEX idx_notification_assignments_group ON notification_assignments(group_id);

CREATE UNIQUE INDEX idx_notification_assignments_group_channel ON notification_assignments(group_id,channel_id) WHERE group_id IS NOT NULL;

CREATE INDEX idx_notification_assignments_task ON notification_assignments(task_id);

CREATE UNIQUE INDEX idx_notification_assignments_task_channel ON notification_assignments(task_id,channel_id) WHERE task_id IS NOT NULL;

CREATE INDEX idx_notification_channels_enabled ON notification_channels(enabled);

CREATE INDEX idx_notification_deliveries_channel ON notification_deliveries(channel_id,created_at);

CREATE INDEX idx_notification_deliveries_pending ON notification_deliveries(state,next_attempt_at,created_at);

CREATE INDEX idx_notification_deliveries_run ON notification_deliveries(run_id);

CREATE INDEX idx_notification_deliveries_task ON notification_deliveries(task_id,created_at);

CREATE INDEX idx_pairing_sessions_state ON pairing_sessions(state,expires_at,id);

CREATE INDEX idx_portable_identities_portable ON portable_identities(portable_id);

CREATE INDEX idx_runs_outcome_effective_time ON runs(outcome, COALESCE(ended_at, scheduled_for));

CREATE INDEX idx_runs_scheduled ON runs(scheduled_for);

CREATE INDEX idx_runs_task ON runs(task_id);

CREATE INDEX idx_tasks_group ON tasks(group_id);

CREATE INDEX idx_tasks_state ON tasks(state);

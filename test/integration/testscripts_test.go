package integration

// Behavioral tests for the maintainer test scripts under test/scripts/.
//
// These drive the scripts as real subprocesses against a throwaway database and
// assert on the rows that come back, because the thing worth testing is what a
// scheduled task actually records -- not whether a function returns.
//
// Every test skips with a stated reason when its interpreter or sqlite3 is
// missing. A skip is not a pass: the reason is printed so that a run on a
// machine lacking pwsh cannot be mistaken for a run that exercised it. Silent
// coverage gaps are how a cross-platform defect reaches a release.

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

// twin is one of the two matched implementations of a script.
type twin struct {
	name string // "powershell" or "posix"
	// command builds the argv for running the given script with args.
	command func(script string, args ...string) []string
	// flag translates a canonical option name to this twin's spelling.
	flag func(canonical string) string
}

func powershellTwin() twin {
	return twin{
		name: "powershell",
		command: func(script string, args ...string) []string {
			return append([]string{"pwsh", "-NoProfile", "-File", script}, args...)
		},
		// -FooBar
		flag: func(c string) string {
			parts := strings.Split(c, "-")
			var b strings.Builder
			b.WriteByte('-')
			for _, p := range parts {
				b.WriteString(strings.ToUpper(p[:1]) + p[1:])
			}
			return b.String()
		},
	}
}

func posixTwin() twin {
	return twin{
		name: "posix",
		command: func(script string, args ...string) []string {
			return append([]string{"bash", script}, args...)
		},
		flag: func(c string) string { return "--" + c },
	}
}

// scriptPath resolves a script under test/scripts/ relative to this package.
func scriptPath(t *testing.T, base string) string {
	t.Helper()
	p, err := filepath.Abs(filepath.Join("..", "scripts", base))
	if err != nil {
		t.Fatalf("resolving script path: %v", err)
	}
	if _, err := os.Stat(p); err != nil {
		t.Fatalf("script not found at %s: %v", p, err)
	}
	return p
}

// requireSqlite skips unless a usable sqlite3 is on PATH, and returns its path.
func requireSqlite(t *testing.T) string {
	t.Helper()
	p, err := exec.LookPath("sqlite3")
	if err != nil {
		t.Skip("SKIP: sqlite3 is not on PATH; these tests drive scripts that require it")
	}
	return p
}

// requireTwin skips unless the twin's interpreter is available.
func requireTwin(t *testing.T, tw twin) {
	t.Helper()
	var interp string
	switch tw.name {
	case "powershell":
		interp = "pwsh"
	default:
		interp = "bash"
	}
	if _, err := exec.LookPath(interp); err != nil {
		t.Skipf("SKIP: %s is not on PATH; the %s twin cannot be exercised here", interp, tw.name)
	}
}

// runScript runs a script twin in an isolated data directory and returns its
// combined output and exit code. A non-zero exit is a result, not an error:
// the exit-code contract is one of the things under test.
func runScript(t *testing.T, tw twin, dataDir, script string, args ...string) (string, int) {
	t.Helper()
	argv := tw.command(script, args...)
	cmd := exec.Command(argv[0], argv[1:]...) //nolint:gosec // fixed argv built above
	cmd.Env = append(os.Environ(), "GOSCHEDULE_TEST_DIR="+dataDir)
	out, err := cmd.CombinedOutput()
	code := 0
	if err != nil {
		var ee *exec.ExitError
		if ok := asExitError(err, &ee); ok {
			code = ee.ExitCode()
		} else {
			t.Fatalf("running %v: %v\n%s", argv, err, out)
		}
	}
	return string(out), code
}

func asExitError(err error, target **exec.ExitError) bool {
	ee, ok := err.(*exec.ExitError)
	if ok {
		*target = ee
	}
	return ok
}

// queryScalar runs a single-value query against a database via sqlite3.
func queryScalar(t *testing.T, sqlite, db, query string) string {
	t.Helper()
	out, err := exec.Command(sqlite, db, query).Output() //nolint:gosec // test-local
	if err != nil {
		t.Fatalf("querying %q: %v", query, err)
	}
	return strings.TrimSpace(string(out))
}

func queryInt(t *testing.T, sqlite, db, query string) int {
	t.Helper()
	s := queryScalar(t, sqlite, db, query)
	n, err := strconv.Atoi(s)
	if err != nil {
		t.Fatalf("expected an integer from %q, got %q", query, s)
	}
	return n
}

func allTwins() []twin { return []twin{powershellTwin(), posixTwin()} }

// TestScriptsHeartbeatSingleShot: the default mode records exactly one beat.
// This is the behavior the whole design rests on -- the scheduler supplies the
// cadence, so one invocation must mean one beat.
func TestScriptsHeartbeatSingleShot(t *testing.T) {
	sqlite := requireSqlite(t)
	script := map[string]string{
		"powershell": scriptPath(t, "Test-Heartbeat.ps1"),
		"posix":      scriptPath(t, "Test-Heartbeat.sh"),
	}
	for _, tw := range allTwins() {
		t.Run(tw.name, func(t *testing.T) {
			requireTwin(t, tw)
			dir := t.TempDir()
			out, code := runScript(t, tw, dir, script[tw.name], tw.flag("label"), "unit")
			if code != 0 {
				t.Fatalf("expected exit 0, got %d\n%s", code, out)
			}
			db := filepath.Join(dir, "heartbeat.db")
			if n := queryInt(t, sqlite, db, "SELECT COUNT(*) FROM beat;"); n != 1 {
				t.Errorf("expected exactly 1 beat, got %d", n)
			}
			if got := queryScalar(t, sqlite, db, "SELECT label FROM beat;"); got != "unit" {
				t.Errorf("label = %q, want %q", got, "unit")
			}
			// With no interval declared there is no expected moment, and drift
			// must be absent rather than a misleading zero.
			if got := queryScalar(t, sqlite, db, "SELECT expected_source FROM beat;"); got != "none" {
				t.Errorf("expected_source = %q, want %q", got, "none")
			}
			if n := queryInt(t, sqlite, db, "SELECT COUNT(*) FROM beat WHERE drift_ms IS NULL;"); n != 1 {
				t.Error("drift should be NULL when no expected moment is available")
			}
		})
	}
}

// TestScriptsHeartbeatAnchorDrift: supplying an anchor plus an interval yields
// a true dispatch-latency figure, labelled with its source.
func TestScriptsHeartbeatAnchorDrift(t *testing.T) {
	sqlite := requireSqlite(t)
	script := map[string]string{
		"powershell": scriptPath(t, "Test-Heartbeat.ps1"),
		"posix":      scriptPath(t, "Test-Heartbeat.sh"),
	}
	anchor := time.Now().UTC().Add(-5 * time.Minute).Format("2006-01-02T15:04:05Z")
	for _, tw := range allTwins() {
		t.Run(tw.name, func(t *testing.T) {
			requireTwin(t, tw)
			dir := t.TempDir()
			_, code := runScript(t, tw, dir, script[tw.name],
				tw.flag("interval-seconds"), "60", tw.flag("anchor-iso"), anchor)
			if code != 0 {
				t.Fatalf("expected exit 0, got %d", code)
			}
			db := filepath.Join(dir, "heartbeat.db")
			if got := queryScalar(t, sqlite, db, "SELECT expected_source FROM beat;"); got != "anchor" {
				t.Errorf("expected_source = %q, want %q", got, "anchor")
			}
			if n := queryInt(t, sqlite, db, "SELECT COUNT(*) FROM beat WHERE drift_ms IS NOT NULL;"); n != 1 {
				t.Error("drift should be recorded when an anchor and interval are supplied")
			}
			// The anchor is exactly 5 intervals back, so the reconstructed grid
			// lands on "now" and drift is just startup cost. Anything near a
			// whole interval would mean the grid arithmetic picked wrong k.
			if ms := queryInt(t, sqlite, db, "SELECT ABS(drift_ms) FROM beat;"); ms > 30000 {
				t.Errorf("drift %dms is more than half the interval; grid arithmetic is wrong", ms)
			}
		})
	}
}

// TestScriptsNoAnchorMeansNoDrift is the regression test for the defect that
// v0.5.1 fixes. Before it, an interval alone caused the run's start to be
// snapped to the nearest interval boundary counted from the Unix epoch. That is
// correct only when the schedule happens to sit on that grid -- and this
// scheduler anchors interval schedules to task creation time, so a task created
// at :06 fires at :06 forever and epoch snapping reported a constant ~6s
// "drift" for a scheduler that was in fact on time to within a quarter second.
//
// Reporting nothing is the correct behavior here. A confident wrong number is
// worse than an absent one, because nothing about its presentation tells you
// which one you got.
func TestScriptsNoAnchorMeansNoDrift(t *testing.T) {
	sqlite := requireSqlite(t)
	script := map[string]string{
		"powershell": scriptPath(t, "Test-Heartbeat.ps1"),
		"posix":      scriptPath(t, "Test-Heartbeat.sh"),
	}
	for _, tw := range allTwins() {
		t.Run(tw.name, func(t *testing.T) {
			requireTwin(t, tw)
			dir := t.TempDir()
			// An interval but no anchor: the exact shape that used to lie.
			if _, code := runScript(t, tw, dir, script[tw.name],
				tw.flag("interval-seconds"), "60"); code != 0 {
				t.Fatalf("expected exit 0, got %d", code)
			}
			db := filepath.Join(dir, "heartbeat.db")
			if got := queryScalar(t, sqlite, db, "SELECT expected_source FROM beat;"); got != "none" {
				t.Errorf("expected_source = %q, want %q -- an interval alone must not "+
					"produce a drift figure, because epoch snapping cannot know the "+
					"schedule's phase", got, "none")
			}
			if n := queryInt(t, sqlite, db, "SELECT COUNT(*) FROM beat WHERE drift_ms IS NOT NULL;"); n != 0 {
				t.Error("drift must be absent without an anchor, not fabricated from the epoch grid")
			}
			// And the interval is still recorded, because gap detection needs it.
			if n := queryInt(t, sqlite, db, "SELECT interval_seconds FROM beat;"); n != 60 {
				t.Errorf("interval_seconds = %d, want 60", n)
			}
		})
	}
}

// TestScriptsReadTimeAnchor: drift computed at read time from raw start
// timestamps. This is the primary path, because the anchor cannot be known
// before the task exists -- the scheduler derives an interval schedule's phase
// from the task's creation moment, so supplying it to the recorder is a
// chicken-and-egg problem. Deriving at read time also means a wrong anchor is
// recoverable: re-run the query, do not re-run the experiment.
func TestScriptsReadTimeAnchor(t *testing.T) {
	requireSqlite(t)
	beat := map[string]string{
		"powershell": scriptPath(t, "Test-Heartbeat.ps1"),
		"posix":      scriptPath(t, "Test-Heartbeat.sh"),
	}
	reader := map[string]string{
		"powershell": scriptPath(t, "Test-ReadTestDB.ps1"),
		"posix":      scriptPath(t, "Test-ReadTestDB.sh"),
	}
	anchor := time.Now().UTC().Add(-10 * time.Minute).Format("2006-01-02T15:04:05Z")
	for _, tw := range allTwins() {
		t.Run(tw.name, func(t *testing.T) {
			requireTwin(t, tw)
			dir := t.TempDir()
			// Recorded with NO anchor -- the ordinary case.
			if _, code := runScript(t, tw, dir, beat[tw.name],
				tw.flag("interval-seconds"), "60"); code != 0 {
				t.Fatalf("recording failed with exit %d", code)
			}
			out, code := runScript(t, tw, dir, reader[tw.name],
				tw.flag("query"), "drift", tw.flag("interval-seconds"), "60",
				tw.flag("anchor-iso"), anchor, tw.flag("quiet"))
			if code != 0 {
				t.Fatalf("read-time drift exited %d: %s", code, out)
			}
			if !strings.Contains(out, "read-time") {
				t.Errorf("drift output should identify itself as read-time derived; got: %s", out)
			}
		})
	}
}

// TestScriptsAnchorTimestampForms exercises the RFC 3339 spellings a maintainer
// will actually paste in. This exists because v0.5.1 shipped a POSIX twin that
// parsed timestamps with `date -d`, which is GNU-only -- macOS ships BSD date,
// where that is not a parse flag at all. It could not reproduce on a GNU-date
// host, so only a macOS runner caught it. The offset form is included because it
// takes a different branch of the BSD fallback than the Z form.
func TestScriptsAnchorTimestampForms(t *testing.T) {
	requireSqlite(t)
	beat := map[string]string{
		"powershell": scriptPath(t, "Test-Heartbeat.ps1"),
		"posix":      scriptPath(t, "Test-Heartbeat.sh"),
	}
	reader := map[string]string{
		"powershell": scriptPath(t, "Test-ReadTestDB.ps1"),
		"posix":      scriptPath(t, "Test-ReadTestDB.sh"),
	}
	now := time.Now().UTC().Add(-10 * time.Minute)
	forms := map[string]string{
		"utc-z":  now.Format("2006-01-02T15:04:05Z"),
		"offset": now.In(time.FixedZone("test", -4*3600)).Format("2006-01-02T15:04:05-07:00"),
	}
	for _, tw := range allTwins() {
		for name, ts := range forms {
			t.Run(tw.name+"/"+name, func(t *testing.T) {
				requireTwin(t, tw)
				dir := t.TempDir()
				if _, code := runScript(t, tw, dir, beat[tw.name],
					tw.flag("interval-seconds"), "60"); code != 0 {
					t.Fatalf("recording failed with exit %d", code)
				}
				out, code := runScript(t, tw, dir, reader[tw.name],
					tw.flag("query"), "drift", tw.flag("interval-seconds"), "60",
					tw.flag("anchor-iso"), ts, tw.flag("quiet"))
				if code != 0 {
					t.Fatalf("anchor %q (%s) rejected with exit %d: %s", ts, name, code, out)
				}
			})
		}
	}
}

// TestScriptsAnchorNeedsInterval: an anchor without an interval cannot
// reconstruct a grid, and saying so is better than guessing one.
func TestScriptsAnchorNeedsInterval(t *testing.T) {
	requireSqlite(t)
	beat := map[string]string{
		"powershell": scriptPath(t, "Test-Heartbeat.ps1"),
		"posix":      scriptPath(t, "Test-Heartbeat.sh"),
	}
	reader := map[string]string{
		"powershell": scriptPath(t, "Test-ReadTestDB.ps1"),
		"posix":      scriptPath(t, "Test-ReadTestDB.sh"),
	}
	for _, tw := range allTwins() {
		t.Run(tw.name, func(t *testing.T) {
			requireTwin(t, tw)
			dir := t.TempDir()
			// Seed one beat with no interval, so none can be inferred either.
			if _, code := runScript(t, tw, dir, beat[tw.name]); code != 0 {
				t.Fatalf("recording failed with exit %d", code)
			}
			_, code := runScript(t, tw, dir, reader[tw.name],
				tw.flag("query"), "drift", tw.flag("anchor-iso"), "2026-07-23T12:00:00Z")
			if code != 2 {
				t.Errorf("anchor without a usable interval: exit = %d, want 2", code)
			}
		})
	}
}

// TestScriptsExitCodeContract: 0 success, non-zero as requested, and 2 reserved
// for usage and prerequisite failures. Getting this wrong makes a run row in
// `gosched runs` mean nothing.
func TestScriptsExitCodeContract(t *testing.T) {
	sqlite := requireSqlite(t)
	script := map[string]string{
		"powershell": scriptPath(t, "Test-Heartbeat.ps1"),
		"posix":      scriptPath(t, "Test-Heartbeat.sh"),
	}
	for _, tw := range allTwins() {
		t.Run(tw.name, func(t *testing.T) {
			requireTwin(t, tw)

			t.Run("induced failure still records", func(t *testing.T) {
				dir := t.TempDir()
				_, code := runScript(t, tw, dir, script[tw.name], tw.flag("fail-with"), "3")
				if code != 3 {
					t.Errorf("exit = %d, want 3", code)
				}
				db := filepath.Join(dir, "heartbeat.db")
				if n := queryInt(t, sqlite, db, "SELECT COUNT(*) FROM beat;"); n != 1 {
					t.Error("the beat must still be recorded when the run reports failure")
				}
				if got := queryScalar(t, sqlite, db, "SELECT outcome FROM beat;"); got != "failed" {
					t.Errorf("outcome = %q, want %q", got, "failed")
				}
			})

			t.Run("reserved codes rejected", func(t *testing.T) {
				for _, reserved := range []string{"0", "2"} {
					dir := t.TempDir()
					_, code := runScript(t, tw, dir, script[tw.name], tw.flag("fail-with"), reserved)
					if code != 2 {
						t.Errorf("--fail-with %s: exit = %d, want 2 (usage error)", reserved, code)
					}
				}
			})

			t.Run("missing prerequisite is 2 not 1", func(t *testing.T) {
				dir := t.TempDir()
				bogus := filepath.Join(dir, "definitely-not-sqlite3")
				_, code := runScript(t, tw, dir, script[tw.name], tw.flag("sqlite-exe"), bogus)
				if code != 2 {
					t.Errorf("exit = %d, want 2; an unmet prerequisite is a usage-class failure, "+
						"and conflating it with a runtime failure sends a maintainer debugging "+
						"the wrong thing", code)
				}
			})
		})
	}
}

// TestScriptsLoopIsBounded: continuous mode must never be unbounded. A runaway
// loop launched under a scheduler is a resource incident.
func TestScriptsLoopIsBounded(t *testing.T) {
	sqlite := requireSqlite(t)
	script := map[string]string{
		"powershell": scriptPath(t, "Test-Heartbeat.ps1"),
		"posix":      scriptPath(t, "Test-Heartbeat.sh"),
	}
	for _, tw := range allTwins() {
		t.Run(tw.name, func(t *testing.T) {
			requireTwin(t, tw)
			dir := t.TempDir()
			_, code := runScript(t, tw, dir, script[tw.name],
				tw.flag("loop"), tw.flag("max-beats"), "3", tw.flag("interval-seconds"), "1")
			if code != 0 {
				t.Fatalf("expected exit 0, got %d", code)
			}
			db := filepath.Join(dir, "heartbeat.db")
			if n := queryInt(t, sqlite, db, "SELECT COUNT(*) FROM beat;"); n != 3 {
				t.Errorf("expected exactly 3 beats, got %d", n)
			}
			if n := queryInt(t, sqlite, db, "SELECT COUNT(DISTINCT session_id) FROM beat;"); n != 1 {
				t.Errorf("a single loop invocation is one session, got %d", n)
			}
		})
	}
}

// TestScriptsConcurrentWriters: overlap-policy testing produces simultaneous
// writers by construction, so contention must be invisible. A lost record here
// would look exactly like a scheduler that failed to fire.
func TestScriptsConcurrentWriters(t *testing.T) {
	sqlite := requireSqlite(t)
	tw := powershellTwin()
	if runtime.GOOS != "windows" {
		tw = posixTwin()
	}
	requireTwin(t, tw)
	name := "Test-Heartbeat.ps1"
	if tw.name == "posix" {
		name = "Test-Heartbeat.sh"
	}
	script := scriptPath(t, name)

	const writers = 5
	dir := t.TempDir()
	// Seed the schema once so the race is on inserts, not on table creation.
	if _, code := runScript(t, tw, dir, script); code != 0 {
		t.Fatalf("seed run failed with exit %d", code)
	}

	var wg sync.WaitGroup
	codes := make([]int, writers)
	for i := 0; i < writers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, codes[i] = runScript(t, tw, dir, script)
		}(i)
	}
	wg.Wait()

	for i, c := range codes {
		if c != 0 {
			t.Errorf("concurrent writer %d exited %d; contention must be waited out, not surfaced", i, c)
		}
	}
	db := filepath.Join(dir, "heartbeat.db")
	if n := queryInt(t, sqlite, db, "SELECT COUNT(*) FROM beat;"); n != writers+1 {
		t.Errorf("expected %d beats, got %d -- records were lost under contention", writers+1, n)
	}
}

// TestScriptsTwinParity is the test the analyze gate added. Twin divergence is
// the likeliest defect in this feature: the two implementations are written
// separately, and nothing but this asserts they agree. It compares the rows
// each twin records for the same invocation, field for field, excluding only
// the fields that are legitimately per-run.
func TestScriptsTwinParity(t *testing.T) {
	sqlite := requireSqlite(t)
	ps, posix := powershellTwin(), posixTwin()
	requireTwin(t, ps)
	requireTwin(t, posix)

	dir := t.TempDir()
	if _, code := runScript(t, ps, dir, scriptPath(t, "Test-Heartbeat.ps1"),
		"-IntervalSeconds", "60", "-Label", "parity", "-AnchorIso", "2026-07-23T12:00:00Z"); code != 0 {
		t.Fatalf("powershell twin exited %d", code)
	}
	if _, code := runScript(t, posix, dir, scriptPath(t, "Test-Heartbeat.sh"),
		"--interval-seconds", "60", "--label", "parity", "--anchor-iso", "2026-07-23T12:00:00Z"); code != 0 {
		t.Fatalf("posix twin exited %d", code)
	}

	db := filepath.Join(dir, "heartbeat.db")
	if n := queryInt(t, sqlite, db, "SELECT COUNT(*) FROM beat;"); n != 2 {
		t.Fatalf("expected one beat per twin, got %d", n)
	}
	// Both twins must agree on every field that describes the same machine and
	// the same declared schedule. Divergence in any of these means the twins
	// have drifted and results from one cannot be compared with the other.
	for _, col := range []string{"label", "hostname", "expected_source", "interval_seconds", "outcome"} {
		n := queryInt(t, sqlite, db,
			"SELECT COUNT(DISTINCT "+col+") FROM beat;")
		if n != 1 {
			got := queryScalar(t, sqlite, db, "SELECT GROUP_CONCAT(DISTINCT "+col+") FROM beat;")
			t.Errorf("twins disagree on %s: %s", col, got)
		}
	}
}

// TestScriptsSystemSnapshot: the host-inspection script records one snapshot,
// and probe degradation never costs the snapshot itself.
func TestScriptsSystemSnapshot(t *testing.T) {
	sqlite := requireSqlite(t)
	script := map[string]string{
		"powershell": scriptPath(t, "Test-GetSystemInfo.ps1"),
		"posix":      scriptPath(t, "Test-GetSystemInfo.sh"),
	}
	for _, tw := range allTwins() {
		t.Run(tw.name, func(t *testing.T) {
			requireTwin(t, tw)
			dir := t.TempDir()
			out, code := runScript(t, tw, dir, script[tw.name],
				tw.flag("invocation-source"), "unit", tw.flag("skip-ports"))
			if code != 0 {
				t.Fatalf("expected exit 0 even with degraded probes, got %d\n%s", code, out)
			}
			db := filepath.Join(dir, "system.db")
			if n := queryInt(t, sqlite, db, "SELECT COUNT(*) FROM snapshot;"); n != 1 {
				t.Fatalf("expected exactly 1 snapshot, got %d", n)
			}
			// Required columns must never be NULL; optional ones may be, and
			// that is the documented meaning of "the probe could not run".
			for _, col := range []string{"hostname", "iso_local", "iso_utc", "os_platform", "script_flavor"} {
				if n := queryInt(t, sqlite, db,
					"SELECT COUNT(*) FROM snapshot WHERE "+col+" IS NULL;"); n != 0 {
					t.Errorf("%s must never be NULL", col)
				}
			}
			want := "powershell"
			if tw.name == "posix" {
				want = "posix"
			}
			if got := queryScalar(t, sqlite, db, "SELECT script_flavor FROM snapshot;"); got != want {
				t.Errorf("script_flavor = %q, want %q", got, want)
			}
			// --skip-ports must actually skip.
			if n := queryInt(t, sqlite, db, "SELECT COUNT(*) FROM snapshot_port;"); n != 0 {
				t.Errorf("skip-ports was requested but %d port rows were written", n)
			}
		})
	}
}

// TestScriptsReaderQueries: every canned query runs and returns something
// parseable against a populated database.
func TestScriptsReaderQueries(t *testing.T) {
	requireSqlite(t)
	script := map[string]string{
		"powershell": scriptPath(t, "Test-ReadTestDB.ps1"),
		"posix":      scriptPath(t, "Test-ReadTestDB.sh"),
	}
	beat := map[string]string{
		"powershell": scriptPath(t, "Test-Heartbeat.ps1"),
		"posix":      scriptPath(t, "Test-Heartbeat.sh"),
	}
	heartbeatQueries := []string{
		"summary", "recent", "cadence", "drift", "jitter", "overlaps",
		"failures", "restarts", "hosts", "schema",
	}
	for _, tw := range allTwins() {
		t.Run(tw.name, func(t *testing.T) {
			requireTwin(t, tw)
			dir := t.TempDir()
			if _, code := runScript(t, tw, dir, beat[tw.name],
				tw.flag("loop"), tw.flag("max-beats"), "2", tw.flag("interval-seconds"), "1"); code != 0 {
				t.Fatalf("seeding beats failed with exit %d", code)
			}
			for _, q := range heartbeatQueries {
				out, code := runScript(t, tw, dir, script[tw.name],
					tw.flag("database"), "heartbeat", tw.flag("query"), q, tw.flag("quiet"))
				if code != 0 {
					t.Errorf("query %q exited %d\n%s", q, code, out)
				}
			}
			// An unknown query is a usage error, not a runtime one.
			if _, code := runScript(t, tw, dir, script[tw.name],
				tw.flag("query"), "no-such-query"); code != 2 {
				t.Errorf("unknown query: exit = %d, want 2", code)
			}
			// --list works without touching a database at all.
			if _, code := runScript(t, tw, dir, script[tw.name], tw.flag("list")); code != 0 {
				t.Errorf("--list exited %d, want 0", code)
			}
		})
	}
}

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

type fixture struct {
	t                    *testing.T
	mu                   sync.Mutex
	server               *httptest.Server
	client               *client
	dir                  string
	state                string
	release              release
	pages                [][]release
	assets               []asset
	contents             map[int64][]byte
	requests             []string
	patches              []map[string]any
	uploads              int
	listStatus           int
	badList              bool
	releaseOverride      *release
	failUploadAfterStore bool
	latest               release
	ignoreLatest         bool
	redirect             string
	redirectID           int64
}

func newFixture(t *testing.T, draft bool) *fixture {
	t.Helper()
	f := &fixture{t: t, dir: t.TempDir(), contents: map[int64][]byte{}}
	f.state = filepath.Join(f.dir, "selection.json")
	f.release = release{ID: 42, Tag: "v1.2.3", Name: "Reviewed release", Body: "Native checks remain UNTESTED.", Target: "0123456789012345678901234567890123456789", Draft: draft}
	for i := 0; i < 8; i++ {
		name := fmt.Sprintf("artifact-%d.bin", i)
		data := []byte("reviewed-bytes-" + name)
		a := asset{ID: int64(100 + i), Name: name, Size: int64(len(data)), State: "uploaded"}
		f.assets = append(f.assets, a)
		f.contents[a.ID] = data
	}
	f.server = httptest.NewServer(http.HandlerFunc(f.serve))
	t.Cleanup(f.server.Close)
	f.release.UploadURL = f.server.URL + "/repos/owner/repo/releases/42/assets{?name,label}"
	f.pages = [][]release{{{ID: 1, Tag: "v0.0.1"}}, {f.release}}
	var err error
	f.client, err = newClient(f.server.URL, "fixture-secret-token", &http.Client{Timeout: 2 * time.Second})
	if err != nil {
		t.Fatal(err)
	}
	return f
}
func (f *fixture) change(fn func()) { f.mu.Lock(); defer f.mu.Unlock(); fn() }
func (f *fixture) serve(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.requests = append(f.requests, r.Method+" "+r.URL.RequestURI())
	if r.Header.Get("Authorization") != "Bearer fixture-secret-token" {
		f.t.Error("API request lacks authentication")
		w.WriteHeader(401)
		return
	}
	path := r.URL.Path
	switch {
	case path == "/repos/owner/repo/releases" && r.Method == "GET":
		if f.listStatus != 0 {
			w.WriteHeader(f.listStatus)
			fmt.Fprint(w, "fixture-secret-token must not appear in diagnostics")
			return
		}
		if f.badList {
			fmt.Fprint(w, "not-json")
			return
		}
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page == 0 {
			page = 1
		}
		if page < len(f.pages) {
			w.Header().Set("Link", fmt.Sprintf("<%s/repos/owner/repo/releases?per_page=100&page=%d>; rel=\"next\"", f.server.URL, page+1))
		}
		if page > len(f.pages) {
			json.NewEncoder(w).Encode([]release{})
			return
		}
		json.NewEncoder(w).Encode(f.pages[page-1])
	case path == "/repos/owner/repo/releases/42" && r.Method == "GET":
		if f.releaseOverride != nil {
			json.NewEncoder(w).Encode(f.releaseOverride)
		} else {
			json.NewEncoder(w).Encode(f.release)
		}
	case path == "/repos/owner/repo/releases/42/assets" && r.Method == "GET":
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page == 0 {
			page = 1
		}
		start := (page - 1) * 4
		end := start + 4
		if end < len(f.assets) {
			w.Header().Set("Link", fmt.Sprintf("<%s/repos/owner/repo/releases/42/assets?per_page=100&page=%d>; rel=\"next\"", f.server.URL, page+1))
		} else {
			end = len(f.assets)
		}
		if start > end {
			start = end
		}
		json.NewEncoder(w).Encode(append([]asset{}, f.assets[start:end]...))
	case strings.HasPrefix(path, "/repos/owner/repo/releases/assets/") && r.Method == "GET":
		if r.Header.Get("Accept") != "application/octet-stream" {
			f.t.Error("asset bytes require octet-stream accept")
		}
		id, _ := strconv.ParseInt(strings.TrimPrefix(path, "/repos/owner/repo/releases/assets/"), 10, 64)
		if f.redirect != "" && id == f.redirectID {
			http.Redirect(w, r, f.redirect, http.StatusFound)
			return
		}
		data, ok := f.contents[id]
		if !ok {
			w.WriteHeader(404)
			return
		}
		w.Write(data)
	case path == "/repos/owner/repo/releases/42/assets" && r.Method == "POST":
		f.uploads++
		if r.URL.Query().Get("name") != checksumName {
			f.t.Error("upload not bound to checksum filename")
		}
		for _, a := range f.assets {
			if a.Name == checksumName {
				w.WriteHeader(422)
				return
			}
		}
		data, _ := io.ReadAll(r.Body)
		a := asset{ID: 999, Name: checksumName, Size: int64(len(data)), State: "uploaded"}
		f.assets = append(f.assets, a)
		f.contents[a.ID] = data
		if f.failUploadAfterStore {
			f.failUploadAfterStore = false
			w.WriteHeader(502)
			return
		}
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(a)
	case path == "/repos/owner/repo/releases/42" && r.Method == "PATCH":
		var payload map[string]any
		json.NewDecoder(r.Body).Decode(&payload)
		f.patches = append(f.patches, payload)
		if value, ok := payload["draft"]; ok {
			f.release.Draft = value.(bool)
		}
		if !f.ignoreLatest {
			f.latest = f.release
		}
		json.NewEncoder(w).Encode(f.release)
	case path == "/repos/owner/repo/releases/latest" && r.Method == "GET":
		if f.latest.ID == 0 {
			w.WriteHeader(404)
			return
		}
		json.NewEncoder(w).Encode(f.latest)
	default:
		w.WriteHeader(404)
	}
}
func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
func (f *fixture) selectRelease() selection {
	f.t.Helper()
	must(f.t, f.client.resolve(context.Background(), "owner/repo", "v1.2.3", f.state))
	s, err := loadSelection(f.state)
	must(f.t, err)
	return s
}
func (f *fixture) initial() selection {
	f.t.Helper()
	s := f.selectRelease()
	must(f.t, f.client.download(context.Background(), &s, filepath.Join(f.dir, "initial"), f.state))
	s, err := loadSelection(f.state)
	must(f.t, err)
	return s
}
func (f *fixture) checksums() []byte {
	f.t.Helper()
	var assets = append([]asset{}, f.assets...)
	sort.Slice(assets, func(i, j int) bool { return assets[i].Name < assets[j].Name })
	var out strings.Builder
	for _, a := range assets {
		if a.Name != checksumName {
			fmt.Fprintf(&out, "%s  %s\n", digest(f.contents[a.ID]), a.Name)
		}
	}
	return []byte(out.String())
}
func (f *fixture) addChecksum(data []byte) {
	f.assets = append(f.assets, asset{ID: 999, Name: checksumName, Size: int64(len(data)), State: "uploaded"})
	f.contents[999] = data
}
func (f *fixture) checksumFile() string {
	path := filepath.Join(f.dir, checksumName)
	must(f.t, os.WriteFile(path, f.checksums(), 0600))
	return path
}
func (f *fixture) accepted() selection {
	f.t.Helper()
	s := f.initial()
	must(f.t, f.client.ensureChecksum(context.Background(), &s, f.checksumFile(), f.state))
	return s
}
func (f *fixture) verified() selection {
	f.t.Helper()
	s := f.accepted()
	must(f.t, f.client.download(context.Background(), &s, filepath.Join(f.dir, "final"), f.state))
	s, err := loadSelection(f.state)
	must(f.t, err)
	return s
}
func (f *fixture) assertReadOnly() {
	f.t.Helper()
	if f.uploads != 0 || len(f.patches) != 0 {
		f.t.Fatalf("unexpected mutations uploads=%d patches=%v", f.uploads, f.patches)
	}
}

func TestResolveAuthenticatedDraftAcrossAllPages(t *testing.T) {
	f := newFixture(t, true)
	s := f.selectRelease()
	if s.ReleaseID != 42 || !s.WasDraft || len(s.Assets) != 8 {
		t.Fatalf("wrong selection: %+v", s)
	}
	for _, request := range f.requests {
		if strings.Contains(request, "/tags/") {
			t.Fatalf("tag endpoint queried: %s", request)
		}
	}
	raw, err := os.ReadFile(f.state)
	must(t, err)
	if bytes.HasPrefix(raw, []byte{239, 187, 191}) || bytes.Contains(raw, []byte("fixture-secret-token")) || bytes.Contains(raw, []byte("upload_url")) {
		t.Fatal("selection leaked secret/URL or has BOM")
	}
	if err = f.client.resolve(context.Background(), "owner/repo", "v1.2.3", f.state); err == nil {
		t.Fatal("resolve overwrote a selected identity")
	}
	f.assertReadOnly()
}
func TestResolveRefusesMissingAmbiguousOrChangedIdentity(t *testing.T) {
	for _, name := range []string{"missing", "duplicate-later-page", "invalid-id", "changed-id", "changed-tag", "changed-body", "prerelease", "no-disclosure", "auth", "bad-json"} {
		t.Run(name, func(t *testing.T) {
			f := newFixture(t, true)
			switch name {
			case "missing":
				f.pages = [][]release{{}}
			case "duplicate-later-page":
				f.pages = append(f.pages, []release{f.release})
			case "invalid-id":
				f.pages[1][0].ID = 0
			case "changed-id":
				r := f.release
				r.ID = 43
				f.releaseOverride = &r
			case "changed-tag":
				r := f.release
				r.Tag = "v9.9.9"
				f.releaseOverride = &r
			case "changed-body":
				r := f.release
				r.Body += " changed"
				f.releaseOverride = &r
			case "prerelease":
				f.pages[1][0].Prerelease = true
			case "no-disclosure":
				f.pages[1][0].Body = "Complete"
			case "auth":
				f.listStatus = 403
			case "bad-json":
				f.badList = true
			}
			err := f.client.resolve(context.Background(), "owner/repo", "v1.2.3", f.state)
			if err == nil {
				t.Fatal("identity refusal missing")
			}
			if strings.Contains(err.Error(), "fixture-secret-token") {
				t.Fatal("error leaked token")
			}
			if _, err = os.Stat(f.state); !os.IsNotExist(err) {
				t.Fatal("refused selection persisted state")
			}
			f.assertReadOnly()
		})
	}
}
func TestDraftPromotionRetainsOneIdentityAndIsIdempotent(t *testing.T) {
	f := newFixture(t, true)
	original := metadata(f.release)
	s := f.verified()
	must(t, f.client.publish(context.Background(), &s))
	must(t, f.client.publish(context.Background(), &s))
	if f.uploads != 1 || len(f.patches) != 1 || !reflect.DeepEqual(f.patches[0], map[string]any{"draft": false, "make_latest": "true"}) {
		t.Fatalf("unexpected mutation contract uploads=%d patches=%v", f.uploads, f.patches)
	}
	if metadata(f.release) != original || f.latest.ID != 42 || f.release.Draft {
		t.Fatal("publication changed content identity or latest")
	}
}
func TestChecksumResumeAndPublishedRecovery(t *testing.T) {
	for _, name := range []string{"matching-draft", "differing-draft", "interrupted-upload", "published", "published-missing"} {
		t.Run(name, func(t *testing.T) {
			published := strings.HasPrefix(name, "published")
			f := newFixture(t, !published)
			if name == "matching-draft" || name == "published" {
				f.addChecksum(f.checksums())
			}
			if name == "differing-draft" {
				f.addChecksum([]byte("different checksum\n"))
			}
			s := f.initial()
			path := f.checksumFile()
			if name == "interrupted-upload" {
				f.failUploadAfterStore = true
				if err := f.client.ensureChecksum(context.Background(), &s, path, f.state); err == nil {
					t.Fatal("unconfirmed upload reported success")
				}
				var err error
				s, err = loadSelection(f.state)
				must(t, err)
			}
			err := f.client.ensureChecksum(context.Background(), &s, path, f.state)
			if name == "differing-draft" || name == "published-missing" {
				if err == nil {
					t.Fatal("invalid checksum accepted")
				}
				f.assertReadOnly()
				return
			}
			must(t, err)
			if name == "interrupted-upload" && f.uploads != 1 {
				t.Fatalf("upload repeated %d times", f.uploads)
			}
			if name != "interrupted-upload" && f.uploads != 0 {
				t.Fatal("existing checksum overwritten")
			}
			must(t, f.client.download(context.Background(), &s, filepath.Join(f.dir, "final"), f.state))
			must(t, f.client.publish(context.Background(), &s))
			if published && !reflect.DeepEqual(f.patches, []map[string]any{{"make_latest": "true"}}) {
				t.Fatalf("published recovery changed contents: %v", f.patches)
			}
		})
	}
}
func TestDownloadRefusesChangedAssetIdentityAndBytes(t *testing.T) {
	for _, name := range []string{"missing", "added", "renamed", "replaced-id", "size", "incomplete", "duplicate-name", "duplicate-id", "unsafe-name", "same-size-changed-bytes"} {
		t.Run(name, func(t *testing.T) {
			f := newFixture(t, true)
			s := f.accepted()
			f.change(func() {
				switch name {
				case "missing":
					f.assets = f.assets[1:]
				case "added":
					f.assets = append(f.assets, asset{ID: 666, Name: "extra", Size: 1, State: "uploaded"})
				case "renamed":
					f.assets[0].Name = "other.bin"
				case "replaced-id":
					f.assets[0].ID = 777
				case "size":
					f.assets[0].Size++
				case "incomplete":
					f.assets[0].State = "starter"
				case "duplicate-name":
					f.assets[1].Name = f.assets[0].Name
				case "duplicate-id":
					f.assets[1].ID = f.assets[0].ID
				case "unsafe-name":
					f.assets[0].Name = "../escape"
				case "same-size-changed-bytes":
					f.contents[100] = bytes.Repeat([]byte{'X'}, len(f.contents[100]))
				}
			})
			if err := f.client.download(context.Background(), &s, filepath.Join(f.dir, "final"), f.state); err == nil {
				t.Fatal("changed asset accepted")
			}
			if len(f.patches) != 0 {
				t.Fatal("changed assets published")
			}
		})
	}
}
func TestPublishRefusesMetadataChangesOrWrongLatestIdentity(t *testing.T) {
	for _, name := range []string{"tag", "notes", "name", "target", "prerelease", "asset", "same-tag-wrong-id", "same-size-changed-bytes"} {
		t.Run(name, func(t *testing.T) {
			f := newFixture(t, true)
			s := f.verified()
			f.change(func() {
				switch name {
				case "tag":
					f.release.Tag = "v9.9.9"
				case "notes":
					f.release.Body += " changed"
				case "name":
					f.release.Name += " changed"
				case "target":
					f.release.Target = "main"
				case "prerelease":
					f.release.Prerelease = true
				case "asset":
					f.assets[0].ID++
				case "same-size-changed-bytes":
					f.contents[100] = bytes.Repeat([]byte{'X'}, len(f.contents[100]))
				case "same-tag-wrong-id":
					f.ignoreLatest = true
					f.latest = f.release
					f.latest.Draft = false
					f.latest.ID = 9999
				}
			})
			if err := f.client.publish(context.Background(), &s); err == nil {
				t.Fatal("changed publication accepted")
			}
			if name != "same-tag-wrong-id" && len(f.patches) != 0 {
				t.Fatal("identity changed before mutation but still published")
			}
		})
	}
}
func TestAssetRedirectDoesNotForwardCredentials(t *testing.T) {
	f := newFixture(t, true)
	var authorization string
	var mu sync.Mutex
	download := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		authorization = r.Header.Get("Authorization")
		mu.Unlock()
		w.Write(f.contents[100])
	}))
	defer download.Close()
	f.redirect = download.URL + "/signed?secret=must-not-be-recorded"
	f.redirectID = 100
	s := f.initial()
	if !s.InitialDownloaded {
		t.Fatal("redirect download failed")
	}
	mu.Lock()
	defer mu.Unlock()
	if authorization != "" {
		t.Fatal("authentication forwarded to another download origin")
	}
	raw, err := os.ReadFile(f.state)
	must(t, err)
	if bytes.Contains(raw, []byte("must-not-be-recorded")) {
		t.Fatal("signed URL persisted")
	}
}
func TestSelectionAndCommandRefusals(t *testing.T) {
	f := newFixture(t, true)
	s := f.selectRelease()
	if err := f.client.publish(context.Background(), &s); err == nil {
		t.Fatal("unverified selection published")
	}
	if err := f.client.ensureChecksum(context.Background(), &s, "missing", f.state); err == nil {
		t.Fatal("checksum accepted before download")
	}
	for _, raw := range []string{`{}`, `{"version":2}`, `{"version":1,"unknown":true}`} {
		must(t, os.WriteFile(f.state, []byte(raw), 0600))
		if _, err := loadSelection(f.state); err == nil {
			t.Fatal("invalid selection accepted")
		}
	}
	for _, args := range [][]string{nil, {"unknown"}, {"resolve"}, {"download", "--state", "x"}, {"publish", "--state", "x", "unexpected"}} {
		var output bytes.Buffer
		if run(args, func(string) string { return "" }, &output, &output) != 2 {
			t.Fatalf("arguments not rejected: %v", args)
		}
	}
	var output bytes.Buffer
	if run([]string{"resolve", "-h"}, func(string) string { return "" }, &output, &output) != 0 {
		t.Fatal("help failed")
	}
	if run([]string{"publish", "--state", "x"}, func(string) string { return "" }, &output, &output) != 1 {
		t.Fatal("missing authentication accepted")
	}
	f.assertReadOnly()
}

func TestAPIErrorRetainsBoundedPermissionMessageWithoutSecrets(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Sensitive-Header", "must-not-appear")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"message": "Workflow write permission denied for fixture-secret-token at https://example.invalid/private?signature=secret\n" + strings.Repeat("x", 900)})
	}))
	defer server.Close()
	client, err := newClient(server.URL, "fixture-secret-token", &http.Client{})
	must(t, err)
	if client.http.Timeout <= 0 {
		t.Fatal("client default timeout is not finite")
	}
	var out release
	err = client.json(context.Background(), "GET", "repos/owner/repo/releases/42", nil, &out)
	if err == nil {
		t.Fatal("permission failure was accepted")
	}
	message := err.Error()
	for _, want := range []string{"GET", "HTTP 403", "Workflow write permission denied", "[redacted]", "[URL omitted]"} {
		if !strings.Contains(message, want) {
			t.Fatalf("error lost %q: %s", want, message)
		}
	}
	for _, secret := range []string{"fixture-secret-token", "https://", "signature=secret", "Sensitive-Header", "must-not-appear"} {
		if strings.Contains(message, secret) {
			t.Fatalf("error leaked %q", secret)
		}
	}
	if len([]rune(message)) > 600 {
		t.Fatal("API error is not bounded")
	}
}

func TestChecksumUploadRefusesOtherOriginAndReleaseID(t *testing.T) {
	for _, name := range []string{"origin", "release-id"} {
		t.Run(name, func(t *testing.T) {
			f := newFixture(t, true)
			s := f.initial()
			f.change(func() {
				if name == "origin" {
					f.release.UploadURL = "https://example.invalid/repos/owner/repo/releases/42/assets{?name,label}"
				} else {
					f.release.UploadURL = f.server.URL + "/repos/owner/repo/releases/43/assets{?name,label}"
				}
			})
			if err := f.client.ensureChecksum(context.Background(), &s, f.checksumFile(), f.state); err == nil {
				t.Fatal("unsafe upload URL accepted")
			}
			f.assertReadOnly()
		})
	}
}
func TestFailedRepeatDownloadInvalidatesFinalVerification(t *testing.T) {
	f := newFixture(t, true)
	s := f.verified()
	f.change(func() { f.contents[100] = bytes.Repeat([]byte{'X'}, len(f.contents[100])) })
	if err := f.client.download(context.Background(), &s, filepath.Join(f.dir, "retry-final"), f.state); err == nil {
		t.Fatal("changed bytes accepted")
	}
	persisted, err := loadSelection(f.state)
	must(t, err)
	if persisted.FinalVerified {
		t.Fatal("failed repeat download retained stale final verification")
	}
	if err = f.client.publish(context.Background(), &persisted); err == nil {
		t.Fatal("stale verification published")
	}
	if len(f.patches) != 0 {
		t.Fatal("failed revalidation changed publication")
	}
}

// The staged body comes directly from these source-owned notes. Validate that
// exact content through read-only production selection, not a synthetic phrase.
func TestPreparedV160NotesPermitStandingWaiverSelection(t *testing.T) {
	notes, err := os.ReadFile("../../.github/release-notes/v1.6.0.md")
	must(t, err)
	f := newFixture(t, true)
	f.release.Tag = "v1.6.0"
	f.release.Body = string(notes)
	f.pages[1][0] = f.release
	must(t, f.client.resolve(context.Background(), "owner/repo", "v1.6.0", f.state))
	selected, err := loadSelection(f.state)
	must(t, err)
	if selected.Tag != "v1.6.0" || selected.ReleaseID != 42 || !selected.WasDraft || selected.Metadata != metadata(f.release) {
		t.Fatalf("prepared notes lost exact selection identity: %+v", selected)
	}
	f.assertReadOnly()
}

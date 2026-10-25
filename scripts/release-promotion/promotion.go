package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const checksumName = "SHA256SUMS.txt"

type release struct {
	ID         int64  `json:"id"`
	Tag        string `json:"tag_name"`
	Name       string `json:"name"`
	Body       string `json:"body"`
	Target     string `json:"target_commitish"`
	Draft      bool   `json:"draft"`
	Prerelease bool   `json:"prerelease"`
	UploadURL  string `json:"upload_url"`
}

type asset struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Size   int64  `json:"size"`
	State  string `json:"state"`
	SHA256 string `json:"verified_sha256,omitempty"`
}

type selection struct {
	Version           int     `json:"version"`
	Repository        string  `json:"repository"`
	Tag               string  `json:"tag"`
	ReleaseID         int64   `json:"release_id"`
	Metadata          string  `json:"content_sha256"`
	WasDraft          bool    `json:"was_draft"`
	Assets            []asset `json:"assets"`
	InitialDownloaded bool    `json:"initial_downloaded"`
	ChecksumAccepted  bool    `json:"checksum_accepted"`
	FinalVerified     bool    `json:"final_verified"`
}

type client struct {
	base  *url.URL
	token string
	http  *http.Client
}
type apiError struct {
	status    int
	operation string
	message   string
}

func (e *apiError) Error() string {
	result := fmt.Sprintf("GitHub %s returned HTTP %d", e.operation, e.status)
	if e.message != "" {
		return result + ": " + e.message
	}
	return result + "; check authentication, repository access and the selected release"
}

var messageURL = regexp.MustCompile("(?i)https?://[^\\s<>\\\"']+")

func sanitizedMessage(body io.Reader, token string) string {
	var details struct{ Message string }
	if err := json.NewDecoder(io.LimitReader(body, 8192)).Decode(&details); err != nil {
		return ""
	}
	text := details.Message
	if token != "" {
		text = strings.ReplaceAll(text, token, "[redacted]")
	}
	text = messageURL.ReplaceAllString(text, "[URL omitted]")
	text = strings.Join(strings.Fields(text), " ")
	runes := []rune(text)
	if len(runes) > 512 {
		text = string(runes[:512]) + "..."
	}
	return text
}

func newClient(origin, token string, h *http.Client) (*client, error) {
	base, err := url.Parse(strings.TrimRight(origin, "/"))
	if err != nil || base.Host == "" || base.User != nil || base.RawQuery != "" || base.Fragment != "" || (base.Scheme != "https" && base.Scheme != "http") {
		return nil, fmt.Errorf("GITHUB_API_URL must be an HTTP API origin without credentials, query or fragment")
	}
	if h == nil {
		h = &http.Client{}
	}
	copied := *h
	if copied.Timeout <= 0 {
		copied.Timeout = 15 * time.Minute
	}
	copied.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) >= 10 {
			return fmt.Errorf("too many asset download redirects")
		}
		if len(via) > 0 && (req.URL.Scheme != via[0].URL.Scheme || req.URL.Host != via[0].URL.Host) {
			req.Header.Del("Authorization")
		}
		return nil
	}
	return &client{base: base, token: token, http: &copied}, nil
}
func (c *client) endpoint(path string) string {
	return strings.TrimRight(c.base.String(), "/") + "/" + path
}
func (s selection) releasePath() string {
	return fmt.Sprintf("repos/%s/releases/%d", s.Repository, s.ReleaseID)
}
func metadata(r release) string {
	// Length-prefixed fields cannot collide when content contains delimiters.
	content := ""
	for _, field := range []string{r.Tag, r.Name, r.Body, r.Target, strconv.FormatBool(r.Prerelease)} {
		content += strconv.Itoa(len(field)) + ":" + field
	}
	return digest([]byte(content))
}
func digest(data []byte) string { sum := sha256.Sum256(data); return hex.EncodeToString(sum[:]) }
func (c *client) request(ctx context.Context, method, rawURL string, body io.Reader, accept, contentType string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, rawURL, body)
	if err != nil {
		return nil, fmt.Errorf("cannot construct GitHub %s request", method)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", accept)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GitHub %s request failed; inspect network access or retry the unchanged selection", method)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		message := sanitizedMessage(resp.Body, c.token)
		resp.Body.Close()
		return nil, &apiError{status: resp.StatusCode, operation: method, message: message}
	}
	return resp, nil
}
func (c *client) json(ctx context.Context, method, path string, body []byte, out any) error {
	resp, err := c.request(ctx, method, c.endpoint(path), bytes.NewReader(body), "application/vnd.github+json", "application/json")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if err := decodeJSON(resp.Body, out); err != nil {
		return fmt.Errorf("GitHub %s response is invalid JSON: %w", method, err)
	}
	return nil
}
func decodeJSON(r io.Reader, out any) error {
	d := json.NewDecoder(io.LimitReader(r, 8<<20))
	if err := d.Decode(out); err != nil {
		return err
	}
	if d.Decode(new(any)) != io.EOF {
		return fmt.Errorf("unexpected trailing JSON content")
	}
	return nil
}

// Each page is visited before selecting a match so a later ambiguous tag cannot be ignored.
func (c *client) pages(ctx context.Context, path string, each func(json.RawMessage) error) error {
	next := c.endpoint(path)
	seen := map[string]bool{}
	for next != "" {
		if seen[next] {
			return fmt.Errorf("GitHub pagination repeated a page")
		}
		seen[next] = true
		u, err := url.Parse(next)
		if err != nil || u.Scheme != c.base.Scheme || u.Host != c.base.Host || u.User != nil || !strings.HasPrefix(u.Path, strings.TrimRight(c.base.Path, "/")+"/repos/") {
			return fmt.Errorf("GitHub pagination left the selected API origin")
		}
		resp, err := c.request(ctx, "GET", next, nil, "application/vnd.github+json", "")
		if err != nil {
			return err
		}
		var rows []json.RawMessage
		err = decodeJSON(resp.Body, &rows)
		resp.Body.Close()
		if err != nil {
			return fmt.Errorf("GitHub list response is invalid: %w", err)
		}
		if rows == nil {
			return fmt.Errorf("GitHub list response must be an array")
		}
		for _, row := range rows {
			if err = each(row); err != nil {
				return err
			}
		}
		next = ""
		for _, link := range strings.Split(resp.Header.Get("Link"), ",") {
			parts := strings.Split(link, ";")
			if len(parts) < 2 {
				continue
			}
			isNext := false
			for _, p := range parts[1:] {
				if strings.TrimSpace(p) == `rel="next"` {
					isNext = true
				}
			}
			if isNext {
				target := strings.Trim(strings.TrimSpace(parts[0]), "<>")
				v, e := url.Parse(target)
				if e != nil {
					return fmt.Errorf("GitHub pagination link is invalid")
				}
				next = u.ResolveReference(v).String()
			}
		}
	}
	return nil
}
func (c *client) getRelease(ctx context.Context, s selection) (release, error) {
	var r release
	err := c.json(ctx, "GET", s.releasePath(), nil, &r)
	if err != nil {
		return r, err
	}
	if r.ID != s.ReleaseID || r.Tag != s.Tag || metadata(r) != s.Metadata || r.Prerelease || !strings.Contains(strings.ToLower(r.Body), "untested") {
		return r, fmt.Errorf("release %d (%s) content identity changed; refusing further operations", s.ReleaseID, s.Tag)
	}
	if !s.WasDraft && r.Draft {
		return r, fmt.Errorf("published release %d became a draft; select and review it again", s.ReleaseID)
	}
	return r, nil
}
func (c *client) listAssets(ctx context.Context, s selection) ([]asset, error) {
	var assets []asset
	err := c.pages(ctx, s.releasePath()+"/assets?per_page=100", func(row json.RawMessage) error {
		var a asset
		if err := json.Unmarshal(row, &a); err != nil {
			return err
		}
		assets = append(assets, a)
		return nil
	})
	if err != nil {
		return nil, err
	}
	if err := validateAssets(assets); err != nil {
		return nil, err
	}
	sort.Slice(assets, func(i, j int) bool { return assets[i].Name < assets[j].Name })
	return assets, nil
}
func validateAssets(assets []asset) error {
	ids := map[int64]bool{}
	names := map[string]bool{}
	for _, a := range assets {
		if a.ID <= 0 || a.Size < 0 || a.State != "uploaded" || a.Name == "" || a.Name == "." || a.Name == ".." || strings.ContainsAny(a.Name, "/\\:*?\"<>|\r\n\x00") || strings.TrimRight(a.Name, " .") != a.Name {
			return fmt.Errorf("release asset %q has an invalid ID, unsafe name or incomplete upload", a.Name)
		}
		upper := strings.ToUpper(strings.Split(a.Name, ".")[0])
		reserved := upper == "CON" || upper == "PRN" || upper == "AUX" || upper == "NUL" || len(upper) == 4 && (strings.HasPrefix(upper, "COM") || strings.HasPrefix(upper, "LPT")) && upper[3] >= '1' && upper[3] <= '9'
		if reserved || ids[a.ID] || names[strings.ToLower(a.Name)] {
			return fmt.Errorf("release asset %q has a duplicate identity or unavailable filename", a.Name)
		}
		ids[a.ID] = true
		names[strings.ToLower(a.Name)] = true
	}
	return nil
}
func sameAssets(expected, observed []asset, allowNewChecksum bool) error {
	byName := map[string]asset{}
	for _, a := range expected {
		byName[a.Name] = a
	}
	for _, a := range observed {
		old, ok := byName[a.Name]
		if !ok {
			if allowNewChecksum && a.Name == checksumName {
				continue
			}
			return fmt.Errorf("release asset %q was added after selection", a.Name)
		}
		if old.ID != a.ID || old.Size != a.Size || old.State != a.State {
			return fmt.Errorf("release asset %q changed identity after selection", a.Name)
		}
		delete(byName, a.Name)
	}
	if len(byName) > 0 {
		return fmt.Errorf("release assets disappeared after selection")
	}
	return nil
}
func (c *client) resolve(ctx context.Context, repository, tag, statePath string) error {
	if !repositoryPattern.MatchString(repository) || !tagPattern.MatchString(tag) {
		return fmt.Errorf("repository or exact release tag is invalid")
	}
	if _, err := os.Lstat(statePath); err == nil {
		return fmt.Errorf("selection file already exists; reuse its commands instead of selecting another release")
	} else if !os.IsNotExist(err) {
		return err
	}
	var matches []release
	err := c.pages(ctx, "repos/"+repository+"/releases?per_page=100", func(row json.RawMessage) error {
		var r release
		if err := json.Unmarshal(row, &r); err != nil {
			return err
		}
		if r.Tag == tag {
			matches = append(matches, r)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if len(matches) != 1 {
		return fmt.Errorf("exact tag %s has %d authenticated release matches; expected exactly one", tag, len(matches))
	}
	r := matches[0]
	if r.ID <= 0 || r.Prerelease || !strings.Contains(strings.ToLower(r.Body), "untested") {
		return fmt.Errorf("release for %s requires a positive ID, full-release state and untested native-check disclosure", tag)
	}
	s := selection{Version: 1, Repository: repository, Tag: tag, ReleaseID: r.ID, Metadata: metadata(r), WasDraft: r.Draft}
	if _, err = c.getRelease(ctx, s); err != nil {
		return err
	}
	s.Assets, err = c.listAssets(ctx, s)
	if err != nil {
		return err
	}
	if len(s.Assets) == 0 {
		return fmt.Errorf("release %d has no staged assets", s.ReleaseID)
	}
	return saveSelection(statePath, s)
}
func loadSelection(path string) (selection, error) {
	var s selection
	data, err := os.ReadFile(path)
	if err != nil {
		return s, fmt.Errorf("read release selection: %w", err)
	}
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err = d.Decode(&s); err != nil {
		return s, fmt.Errorf("read release selection JSON: %w", err)
	}
	if d.Decode(new(any)) != io.EOF {
		return s, fmt.Errorf("release selection has trailing content")
	}
	if s.Version != 1 || s.ReleaseID <= 0 || !repositoryPattern.MatchString(s.Repository) || !tagPattern.MatchString(s.Tag) || len(s.Metadata) != 64 {
		return s, fmt.Errorf("release selection has an unsupported version or invalid identity")
	}
	if err = validateAssets(s.Assets); err != nil {
		return s, err
	}
	if s.InitialDownloaded {
		for _, a := range s.Assets {
			if _, err = hex.DecodeString(a.SHA256); err != nil || len(a.SHA256) != 64 {
				return s, fmt.Errorf("release selection lacks verified digest for %q", a.Name)
			}
		}
	}
	if s.ChecksumAccepted && !s.InitialDownloaded || s.FinalVerified && !s.ChecksumAccepted {
		return s, fmt.Errorf("release selection has an invalid validation sequence")
	}
	return s, nil
}
func saveSelection(path string, s selection) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	file, err := os.CreateTemp(filepath.Dir(path), ".release-selection-*")
	if err != nil {
		return err
	}
	name := file.Name()
	defer os.Remove(name)
	if err = file.Chmod(0600); err == nil {
		_, err = file.Write(append(data, '\n'))
	}
	if err == nil {
		err = file.Sync()
	}
	closeErr := file.Close()
	if err == nil {
		err = closeErr
	}
	if err != nil {
		return fmt.Errorf("write release selection: %w", err)
	}
	if err = os.Rename(name, path); err != nil {
		return fmt.Errorf("replace release selection: %w", err)
	}
	return nil
}
func (c *client) assetResponse(ctx context.Context, s selection, a asset) (*http.Response, error) {
	return c.request(ctx, "GET", c.endpoint(fmt.Sprintf("repos/%s/releases/assets/%d", s.Repository, a.ID)), nil, "application/octet-stream", "")
}
func (c *client) download(ctx context.Context, s *selection, directory, statePath string) error {
	// A failed repeat validation must not leave an earlier final verification usable.
	if s.FinalVerified {
		s.FinalVerified = false
		if err := saveSelection(statePath, *s); err != nil {
			return err
		}
	}
	if _, err := c.getRelease(ctx, *s); err != nil {
		return err
	}
	assets, err := c.listAssets(ctx, *s)
	if err != nil {
		return err
	}
	if err = sameAssets(s.Assets, assets, false); err != nil {
		return err
	}
	if s.InitialDownloaded && !s.ChecksumAccepted {
		return fmt.Errorf("accept the complete checksum before the final download")
	}
	if err = os.Mkdir(directory, 0700); err != nil {
		return fmt.Errorf("create a new download directory %q: %w", directory, err)
	}
	for i, a := range assets {
		resp, err := c.assetResponse(ctx, *s, a)
		if err != nil {
			return fmt.Errorf("download asset %q (ID %d): %w", a.Name, a.ID, err)
		}
		temp, err := os.CreateTemp(directory, ".asset-*")
		if err != nil {
			resp.Body.Close()
			return err
		}
		hash := sha256.New()
		n, copyErr := io.Copy(io.MultiWriter(temp, hash), io.LimitReader(resp.Body, a.Size+1))
		resp.Body.Close()
		closeErr := temp.Close()
		if copyErr != nil || closeErr != nil || n != a.Size {
			os.Remove(temp.Name())
			return fmt.Errorf("download asset %q did not return its complete selected bytes", a.Name)
		}
		got := hex.EncodeToString(hash.Sum(nil))
		if s.InitialDownloaded && got != s.Assets[i].SHA256 {
			os.Remove(temp.Name())
			return fmt.Errorf("release asset %q bytes changed after the initial download", a.Name)
		}
		if err = os.Rename(temp.Name(), filepath.Join(directory, a.Name)); err != nil {
			os.Remove(temp.Name())
			return err
		}
		assets[i].SHA256 = got
	}
	if s.InitialDownloaded {
		s.FinalVerified = true
	} else {
		s.InitialDownloaded = true
	}
	s.Assets = assets
	return saveSelection(statePath, *s)
}
func checksumAsset(assets []asset) (asset, bool) {
	for _, a := range assets {
		if a.Name == checksumName {
			return a, true
		}
	}
	return asset{}, false
}
func (c *client) checksumBytes(ctx context.Context, s selection, a asset) ([]byte, error) {
	if a.Size > 1<<20 {
		return nil, fmt.Errorf("checksum asset is unexpectedly large")
	}
	resp, err := c.assetResponse(ctx, s, a)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, a.Size+1))
	if err != nil || int64(len(data)) != a.Size {
		return nil, fmt.Errorf("checksum download is incomplete")
	}
	return data, nil
}
func validateChecksum(data []byte, assets []asset) error {
	expected := map[string]string{}
	for _, a := range assets {
		if a.Name != checksumName {
			expected[a.Name] = a.SHA256
		}
	}
	if len(expected) != 8 {
		return fmt.Errorf("checksum requires exactly eight validated staged assets")
	}
	for _, line := range strings.Split(strings.TrimSuffix(string(data), "\n"), "\n") {
		sum, name, ok := strings.Cut(line, "  ")
		want, exists := expected[name]
		if !ok || !exists || sum != want {
			return fmt.Errorf("checksum does not describe every exact staged asset once")
		}
		delete(expected, name)
	}
	if len(expected) != 0 {
		return fmt.Errorf("checksum omits staged assets")
	}
	return nil
}
func (c *client) uploadChecksum(ctx context.Context, s selection, r release, data []byte) (asset, error) {
	raw := strings.Split(r.UploadURL, "{")[0]
	u, err := url.Parse(raw)
	// GitHub.com uploads use uploads.github.com; enterprise uses the authenticated API host.
	allowedHost := c.base.Host
	if c.base.Host == "api.github.com" {
		allowedHost = "uploads.github.com"
	}
	if err != nil || u.Scheme != c.base.Scheme || u.Host != allowedHost || u.User != nil || u.RawQuery != "" || u.Fragment != "" || !strings.HasSuffix(u.Path, "/"+s.releasePath()+"/assets") {
		return asset{}, fmt.Errorf("selected release upload URL does not identify its authenticated release ID")
	}
	query := u.Query()
	query.Set("name", checksumName)
	u.RawQuery = query.Encode()
	resp, err := c.request(ctx, "POST", u.String(), bytes.NewReader(data), "application/vnd.github+json", "text/plain")
	if err != nil {
		return asset{}, fmt.Errorf("checksum upload did not confirm success; rerun against the unchanged release to observe any completed upload: %w", err)
	}
	defer resp.Body.Close()
	var a asset
	if err = decodeJSON(resp.Body, &a); err != nil {
		return a, fmt.Errorf("checksum upload returned an invalid response; rerun to observe its status")
	}
	if err = validateAssets([]asset{a}); err != nil || a.Name != checksumName || a.Size != int64(len(data)) {
		return a, fmt.Errorf("checksum upload did not confirm the expected asset identity")
	}
	return a, nil
}
func (c *client) ensureChecksum(ctx context.Context, s *selection, path, statePath string) error {
	if !s.InitialDownloaded {
		return fmt.Errorf("download and validate initial assets before accepting checksums")
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	data, err := io.ReadAll(io.LimitReader(file, (1<<20)+1))
	file.Close()
	if err != nil {
		return err
	}
	if len(data) > 1<<20 {
		return fmt.Errorf("generated checksum file is unexpectedly large")
	}
	if err = validateChecksum(data, s.Assets); err != nil {
		return err
	}
	r, err := c.getRelease(ctx, *s)
	if err != nil {
		return err
	}
	assets, err := c.listAssets(ctx, *s)
	if err != nil {
		return err
	}
	if err = sameAssets(s.Assets, assets, true); err != nil {
		return err
	}
	existing, found := checksumAsset(assets)
	if !found {
		if !r.Draft {
			return fmt.Errorf("published release lacks its complete checksum; latest pointer was not changed")
		}
		existing, err = c.uploadChecksum(ctx, *s, r, data)
		if err != nil {
			return err
		}
		assets, err = c.listAssets(ctx, *s)
		if err != nil {
			return err
		}
		if err = sameAssets(s.Assets, assets, true); err != nil {
			return err
		}
		observed, ok := checksumAsset(assets)
		if !ok || observed.ID != existing.ID {
			return fmt.Errorf("uploaded checksum asset changed identity")
		}
		existing = observed
	}
	observedBytes, err := c.checksumBytes(ctx, *s, existing)
	if err != nil {
		return err
	}
	if !bytes.Equal(observedBytes, data) {
		return fmt.Errorf("existing checksum asset differs from the exact staged bytes; it was not overwritten")
	}
	byName := map[string]asset{}
	for _, a := range s.Assets {
		byName[a.Name] = a
	}
	for i, a := range assets {
		if a.Name == checksumName {
			assets[i].SHA256 = digest(data)
		} else {
			assets[i].SHA256 = byName[a.Name].SHA256
		}
	}
	s.Assets = assets
	s.ChecksumAccepted = true
	s.FinalVerified = false
	return saveSelection(statePath, *s)
}

// Verify bytes immediately before mutation as well as during the final download.
func (c *client) verifyRemoteBytes(ctx context.Context, s selection) error {
	for _, a := range s.Assets {
		resp, err := c.assetResponse(ctx, s, a)
		if err != nil {
			return fmt.Errorf("verify asset %q before publication: %w", a.Name, err)
		}
		hash := sha256.New()
		n, copyErr := io.Copy(hash, io.LimitReader(resp.Body, a.Size+1))
		resp.Body.Close()
		if copyErr != nil || n != a.Size || hex.EncodeToString(hash.Sum(nil)) != a.SHA256 {
			return fmt.Errorf("release asset %q bytes changed before publication; nothing was published", a.Name)
		}
	}
	return nil
}

func (c *client) publish(ctx context.Context, s *selection) error {
	if !s.InitialDownloaded || !s.ChecksumAccepted || !s.FinalVerified {
		return fmt.Errorf("complete initial download, checksum acceptance and final verification before publication")
	}
	if _, ok := checksumAsset(s.Assets); !ok {
		return fmt.Errorf("selected release lacks an accepted checksum asset")
	}
	r, err := c.getRelease(ctx, *s)
	if err != nil {
		return err
	}
	assets, err := c.listAssets(ctx, *s)
	if err != nil {
		return err
	}
	if err = sameAssets(s.Assets, assets, false); err != nil {
		return err
	}
	if err := c.verifyRemoteBytes(ctx, *s); err != nil {
		return err
	}
	var latest release
	err = c.json(ctx, "GET", "repos/"+s.Repository+"/releases/latest", nil, &latest)
	if err != nil {
		var status *apiError
		if !errors.As(err, &status) || status.status != http.StatusNotFound {
			return err
		}
	}
	alreadyLatest := latest.ID == s.ReleaseID && latest.Tag == s.Tag && !latest.Draft && !latest.Prerelease
	if r.Draft || !alreadyLatest {
		body := []byte(`{"make_latest":"true"}`)
		if r.Draft {
			body = []byte(`{"draft":false,"make_latest":"true"}`)
		}
		var updated release
		if err = c.json(ctx, "PATCH", s.releasePath(), body, &updated); err != nil {
			return fmt.Errorf("publication did not confirm success; rerun the unchanged release to recover its latest pointer: %w", err)
		}
		if updated.ID != s.ReleaseID || updated.Tag != s.Tag || updated.Draft || metadata(updated) != s.Metadata {
			return fmt.Errorf("publication response changed release identity or contents")
		}
	}
	confirmed, err := c.getRelease(ctx, *s)
	if err != nil {
		return err
	}
	if confirmed.Draft {
		return fmt.Errorf("selected release remains draft after publication")
	}
	assets, err = c.listAssets(ctx, *s)
	if err != nil {
		return err
	}
	if err = sameAssets(s.Assets, assets, false); err != nil {
		return err
	}
	if err = c.json(ctx, "GET", "repos/"+s.Repository+"/releases/latest", nil, &latest); err != nil {
		return err
	}
	if latest.ID != s.ReleaseID || latest.Tag != s.Tag || latest.Draft || latest.Prerelease {
		return fmt.Errorf("latest release does not identify selected release %d (%s); rerun the unchanged release to repair it", s.ReleaseID, s.Tag)
	}
	return nil
}

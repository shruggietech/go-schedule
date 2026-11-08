package incident

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"path"
	"regexp"
	"strings"
	"time"
	"unicode/utf16"
	"unicode/utf8"
)

// RuntimeEventLimit bounds retained event records.
const RuntimeEventLimit = 16

// RuntimeEventByteLimit bounds each raw event-query stream before decoding.
const RuntimeEventByteLimit = 256 << 10
const runtimeQueryTimeout = 3 * time.Second

// RuntimeEvent contains attributed runtime facts, never arbitrary event text.
type RuntimeEvent struct {
	Provider    string    `json:"provider"`
	EventID     int       `json:"event_id"`
	RecordID    uint64    `json:"record_id"`
	PID         int       `json:"pid"`
	Time        time.Time `json:"time"`
	Application string    `json:"application"`
	CoreCLR     string    `json:"coreclr_version"`
	Framework   string    `json:"dotnet_version"`
}

// RuntimeEventReceipt is independent of the incident's native result and dump status.
type RuntimeEventReceipt struct {
	Kind        string         `json:"kind"`
	Version     int            `json:"version"`
	Source      string         `json:"source"`
	PID         int            `json:"pid"`
	Application string         `json:"application"`
	Started     time.Time      `json:"started"`
	Exited      time.Time      `json:"exited"`
	Status      string         `json:"status"`
	Detail      string         `json:"detail,omitempty"`
	DurationMS  int64          `json:"duration_ms"`
	Events      []RuntimeEvent `json:"events"`
}

// RuntimeEventQueryResult keeps platform query failure separate from event evidence.
type RuntimeEventQueryResult struct {
	Data                   []byte
	Error                  error
	Unavailable, Truncated bool
}

// CollectRuntimeEvents retains a bounded private receipt, not a native dump.
func (c *Capture) CollectRuntimeEvents(target Target, started, exited time.Time, query func(context.Context, string) RuntimeEventQueryResult) (Artifact, error) {
	if _, err := c.path("runtime-events.json"); err != nil {
		return Artifact{Name: "runtime-events.json", Status: "unavailable"}, err
	}
	return c.saveRuntimeEvents(collectRuntimeEvents(target, started, exited, query))
}
func (c *Capture) saveRuntimeEvents(receipt RuntimeEventReceipt) (Artifact, error) {
	a := Artifact{Name: "runtime-events.json", Status: "unavailable"}
	data, err := json.MarshalIndent(receipt, "", "  ")
	if err != nil {
		return a, err
	}
	if err = c.Write(a.Name, append(data, '\n')); err != nil {
		return a, err
	}
	return c.Inspect(a.Name, MetadataLimit), nil
}
func executableBase(s string) string { return path.Base(strings.ReplaceAll(s, string(rune(92)), "/")) }
func collectRuntimeEvents(target Target, started, exited time.Time, query func(context.Context, string) RuntimeEventQueryResult) (r RuntimeEventReceipt) {
	now := time.Now()
	defer func() { r.DurationMS = time.Since(now).Milliseconds() }()
	r = RuntimeEventReceipt{Kind: "go-schedule.runtime-events", Version: 1, Source: "Windows Application event log", PID: target.PID, Application: executableBase(target.Executable), Started: started.UTC(), Exited: exited.UTC(), Status: "failed", Events: []RuntimeEvent{}}
	if query == nil || target.PID <= 0 || target.Executable == "" || started.IsZero() || exited.Before(started) {
		r.Detail = "Target PID, executable and ordered execution interval are required."
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), runtimeQueryTimeout)
	defer cancel()
	xpath := fmt.Sprintf(`*[System[Provider[@Name='.NET Runtime'] and (EventID=1023 or EventID=1025) and Execution[@ProcessID=%d] and TimeCreated[@SystemTime >= '%s' and @SystemTime <= '%s']]]`, target.PID, r.Started.Format(time.RFC3339Nano), r.Exited.Format(time.RFC3339Nano))
	result := query(ctx, xpath)
	if result.Truncated || len(result.Data) > RuntimeEventByteLimit {
		r.Status = "truncated"
		r.Detail = "Event query output exceeded its bounded byte allowance."
		return
	}
	if result.Error != nil {
		if result.Unavailable {
			r.Status = "unavailable"
		}
		r.Detail = "Event query did not complete: " + result.Error.Error()
		return
	}
	events, truncated, err := decodeRuntimeEvents(result.Data, target, r.Started, r.Exited)
	if err != nil {
		r.Detail = err.Error()
		return
	}
	r.Events = events
	switch {
	case truncated:
		r.Status = "truncated"
		r.Detail = "Event count exceeded sixteen; retained events are partial evidence."
	case len(events) == 0:
		r.Status = "absent"
		r.Detail = "No matching retained event was available at query time; this does not establish that no crash occurred."
	default:
		r.Status = "collected"
		r.Detail = "Event-log correlation only; no runtime-module hash, native dump or causal diagnosis is supplied."
	}
	return
}

type runtimeXML struct {
	XMLName xml.Name
	Events  []runtimeXMLEvent            `xml:"Event"`
	Other   []struct{ XMLName xml.Name } `xml:",any"`
}
type runtimeXMLEvent struct {
	XMLName xml.Name
	System  struct {
		Provider struct {
			Name string `xml:"Name,attr"`
		}
		EventID       int
		EventRecordID uint64
		Execution     struct {
			PID int `xml:"ProcessID,attr"`
		}
		TimeCreated struct {
			Time string `xml:"SystemTime,attr"`
		}
	}
	Data []string `xml:"EventData>Data"`
}

var runtimeVersion = regexp.MustCompile(`^[0-9]+([.][0-9]+){1,4}([+-][[:alnum:].-]+)?$`)

func decodeRuntimeEvents(data []byte, target Target, started, exited time.Time) ([]RuntimeEvent, bool, error) {
	text, decoded, err := runtimeXMLText(data)
	if err != nil {
		return nil, false, err
	}
	decoder := xml.NewDecoder(bytes.NewReader(text))
	decoder.CharsetReader = func(label string, input io.Reader) (io.Reader, error) {
		if decoded && (strings.EqualFold(label, "utf-16") || strings.EqualFold(label, "utf-16le")) {
			return input, nil
		}
		return nil, errors.New("unsupported event XML encoding")
	}
	var doc runtimeXML
	if err := decoder.Decode(&doc); err != nil {
		return nil, false, errors.New("event XML is malformed")
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return nil, false, errors.New("event XML must contain exactly one document")
	}
	if doc.XMLName.Local != "Events" || len(doc.Other) > 0 {
		return nil, false, errors.New("event XML has an unexpected envelope")
	}
	events := []RuntimeEvent{}
	truncated := len(doc.Events) > RuntimeEventLimit
	for i, e := range doc.Events {
		if i == RuntimeEventLimit {
			break
		}
		if e.System.Execution.PID != target.PID || e.System.Provider.Name != ".NET Runtime" || (e.System.EventID != 1023 && e.System.EventID != 1025) {
			continue
		}
		at, err := time.Parse(time.RFC3339Nano, e.System.TimeCreated.Time)
		if err != nil {
			return nil, false, errors.New("runtime event timestamp is invalid")
		}
		if at.Before(started) || at.After(exited) {
			continue
		}
		if e.XMLName.Space != "http://schemas.microsoft.com/win/2004/08/events/event" || e.System.EventRecordID == 0 {
			return nil, false, errors.New("runtime event identity is incomplete")
		}
		fields := map[string]string{}
		for _, datum := range e.Data {
			for _, line := range strings.Split(datum, "\n") {
				for _, key := range []string{"Application", "CoreCLR Version", ".NET Version"} {
					if v, ok := strings.CutPrefix(line, key+":"); ok {
						if _, exists := fields[key]; exists {
							return nil, false, errors.New("runtime event identity fields are ambiguous")
						}
						fields[key] = strings.TrimSpace(v)
					}
				}
			}
		}
		application, present := fields["Application"]
		if !present || application == "" {
			return nil, false, errors.New("runtime event application identity is incomplete")
		}
		app := executableBase(application)
		if !strings.EqualFold(app, executableBase(target.Executable)) {
			continue
		}
		clr, framework := fields["CoreCLR Version"], fields[".NET Version"]
		if len(app) > 260 || len(clr) > 128 || len(framework) > 128 || !runtimeVersion.MatchString(clr) || !runtimeVersion.MatchString(framework) {
			return nil, false, errors.New("runtime event version fields are invalid")
		}
		events = append(events, RuntimeEvent{Provider: e.System.Provider.Name, EventID: e.System.EventID, RecordID: e.System.EventRecordID, PID: target.PID, Time: at.UTC(), Application: app, CoreCLR: clr, Framework: framework})
	}
	return events, truncated, nil
}
func runtimeXMLText(data []byte) ([]byte, bool, error) {
	if bytes.HasPrefix(data, []byte{255, 254}) {
		data = data[2:]
		if len(data)%2 != 0 {
			return nil, false, errors.New("event UTF-16 output is incomplete")
		}
		words := make([]uint16, len(data)/2)
		for i := range words {
			words[i] = binary.LittleEndian.Uint16(data[i*2:])
		}
		for i := 0; i < len(words); i++ {
			w := words[i]
			if w >= 0xd800 && w <= 0xdbff {
				if i+1 == len(words) || words[i+1] < 0xdc00 || words[i+1] > 0xdfff {
					return nil, false, errors.New("event UTF-16 output has invalid surrogate data")
				}
				i++
			} else if w >= 0xdc00 && w <= 0xdfff {
				return nil, false, errors.New("event UTF-16 output has invalid surrogate data")
			}
		}
		return []byte(string(utf16.Decode(words))), true, nil
	}
	data = bytes.TrimPrefix(data, []byte{239, 187, 191})
	if !utf8.Valid(data) {
		return nil, false, errors.New("event output has unsupported encoding")
	}
	return data, false, nil
}

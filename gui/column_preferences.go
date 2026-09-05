package gui

import (
	"encoding/json"
	"math"
	"sync"

	"fyne.io/fyne/v2"
)

const columnLayoutPreferenceVersion = 1

const (
	scheduleColumnLayoutPreferenceKey = "tables.schedule.columns.v1"
	activityColumnLayoutPreferenceKey = "tables.activity.columns.v1"
)

type storedColumnLayout struct {
	Version     int       `json:"version"`
	Columns     []string  `json:"columns"`
	Proportions []float32 `json:"proportions"`
}

// columnProfile owns one view's normalized column proportions. The mutex
// protects callbacks from tests and UI events that may observe a refresh while
// a preference-backed adjustment is completing.
type columnProfile struct {
	mu       sync.RWMutex
	columns  []structuredColumn
	defaults []float32
	values   []float32
	prefs    fyne.Preferences
	key      string
}

func newColumnProfile(columns []structuredColumn, prefs fyne.Preferences, key string) *columnProfile {
	profile := &columnProfile{
		columns: append([]structuredColumn(nil), columns...),
		prefs:   prefs,
		key:     key,
	}
	profile.defaults = defaultColumnProportions(columns)
	profile.values = append([]float32(nil), profile.defaults...)
	if prefs == nil || key == "" {
		return profile
	}
	var stored storedColumnLayout
	if err := json.Unmarshal([]byte(prefs.String(key)), &stored); err != nil || !validStoredColumnLayout(stored, columns) {
		return profile
	}
	if normalized := normalizeColumnProportions(stored.Proportions); len(normalized) == len(columns) {
		profile.values = normalized
	}
	return profile
}

func defaultColumnProportions(columns []structuredColumn) []float32 {
	values := make([]float32, len(columns))
	for i, column := range columns {
		values[i] = column.Preferred
		if values[i] <= 0 || !finiteFloat32(values[i]) {
			values[i] = 1
		}
	}
	return normalizeColumnProportions(values)
}

func validStoredColumnLayout(stored storedColumnLayout, columns []structuredColumn) bool {
	if stored.Version != columnLayoutPreferenceVersion || len(stored.Columns) != len(columns) || len(stored.Proportions) != len(columns) {
		return false
	}
	for i, column := range columns {
		if stored.Columns[i] != column.Header || stored.Proportions[i] <= 0 || !finiteFloat32(stored.Proportions[i]) {
			return false
		}
	}
	return true
}

func normalizeColumnProportions(values []float32) []float32 {
	normalized := append([]float32(nil), values...)
	total := float32(0)
	for _, value := range normalized {
		if value <= 0 || !finiteFloat32(value) {
			return nil
		}
		total += value
	}
	if total <= 0 || !finiteFloat32(total) {
		return nil
	}
	for i := range normalized {
		normalized[i] /= total
		if normalized[i] <= 0 || !finiteFloat32(normalized[i]) {
			return nil
		}
	}
	return normalized
}

func finiteFloat32(value float32) bool {
	return !math.IsInf(float64(value), 0) && !math.IsNaN(float64(value))
}

func (p *columnProfile) proportions() []float32 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return append([]float32(nil), p.values...)
}

func (p *columnProfile) setProportions(values []float32, persist bool) bool {
	normalized := normalizeColumnProportions(values)
	if len(normalized) != len(p.columns) {
		return false
	}
	p.mu.Lock()
	p.values = normalized
	p.mu.Unlock()
	if persist {
		p.persist()
	}
	return true
}

func (p *columnProfile) widths(available, gap float32) []float32 {
	return profileColumnWidths(p.columns, p.proportions(), available, gap)
}

func (p *columnProfile) adjust(boundary int, delta, available, gap float32) bool {
	if boundary < 0 || boundary+1 >= len(p.columns) || !finiteFloat32(delta) || delta == 0 {
		return false
	}
	widths := p.widths(available, gap)
	if len(widths) != len(p.columns) {
		return false
	}
	usable := columnUsableWidth(len(p.columns), available, gap)
	minimumTotal := float32(0)
	for _, column := range p.columns {
		minimumTotal += maxFloat32(column.Minimum, 0)
	}
	if usable <= 0 || minimumTotal > usable {
		return false
	}
	leftMin := maxFloat32(p.columns[boundary].Minimum, 0)
	rightMin := maxFloat32(p.columns[boundary+1].Minimum, 0)
	delta = maxFloat32(delta, leftMin-widths[boundary])
	delta = minFloat32(delta, widths[boundary+1]-rightMin)
	if delta == 0 {
		return false
	}
	widths[boundary] += delta
	widths[boundary+1] -= delta
	return p.setProportions(widths, false)
}

func (p *columnProfile) persist() {
	if p.prefs == nil || p.key == "" {
		return
	}
	stored := storedColumnLayout{Version: columnLayoutPreferenceVersion, Proportions: p.proportions()}
	stored.Columns = make([]string, len(p.columns))
	for i, column := range p.columns {
		stored.Columns[i] = column.Header
	}
	encoded, err := json.Marshal(stored)
	if err == nil {
		p.prefs.SetString(p.key, string(encoded))
	}
}

func (p *columnProfile) reset() {
	p.setProportions(p.defaults, true)
}

func columnUsableWidth(columnCount int, available, gap float32) float32 {
	if available < 0 {
		available = 0
	}
	usable := available - maxFloat32(float32(columnCount-1), 0)*gap
	return maxFloat32(usable, 0)
}

func profileColumnWidths(columns []structuredColumn, proportions []float32, available, gap float32) []float32 {
	if len(columns) == 0 || len(proportions) != len(columns) {
		return responsiveColumnWidths(columns, available, gap)
	}
	usable := columnUsableWidth(len(columns), available, gap)
	minimumTotal := float32(0)
	for _, column := range columns {
		minimumTotal += maxFloat32(column.Minimum, 0)
	}
	if minimumTotal > usable {
		return responsiveColumnWidths(columns, available, gap)
	}

	widths := make([]float32, len(columns))
	remaining := usable
	remainingWeight := float32(1)
	fixed := make([]bool, len(columns))
	for {
		changed := false
		for i, column := range columns {
			if fixed[i] {
				continue
			}
			candidate := remaining * proportions[i] / remainingWeight
			minimum := maxFloat32(column.Minimum, 0)
			if candidate < minimum {
				widths[i] = minimum
				remaining -= minimum
				remainingWeight -= proportions[i]
				fixed[i] = true
				changed = true
			}
		}
		if !changed {
			break
		}
	}
	for i := range columns {
		if !fixed[i] {
			widths[i] = remaining * proportions[i] / remainingWeight
		}
	}
	allocated := float32(0)
	for _, width := range widths {
		allocated += width
	}
	widths[len(widths)-1] += usable - allocated
	return widths
}

func minFloat32(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

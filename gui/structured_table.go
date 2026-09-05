package gui

import (
	"image/color"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	structuredColumnGap float32 = 4
	columnBoundaryWidth float32 = 12
	columnKeyboardStep  float32 = 10
)

// colorNameAlternateRow is a quiet, translucent structural stripe. It remains
// translucent so Fyne's full-row hover and selection surfaces stay visible
// underneath it.
const colorNameAlternateRow fyne.ThemeColorName = "goScheduleAlternateRow"

// structuredColumn defines one aligned header/body field. Minimum is protected
// when space permits; Weight receives surplus width. At narrower supported
// sizes every minimum scales proportionally, so the row can never introduce a
// horizontal scroll requirement.
type structuredColumn struct {
	Header    string
	Minimum   float32
	Weight    float32
	Preferred float32
	Alignment fyne.TextAlign
}

// structuredCell is the normalized presentation of one source value.
type structuredCell struct {
	Text       string
	FullText   string
	Importance widget.Importance
	TextStyle  fyne.TextStyle
}

// structuredRowModel is an immutable view of one source record. Identity is
// stable across refreshes; Summary exposes complete values when cells truncate.
type structuredRowModel struct {
	Identity string
	Cells    []structuredCell
	Summary  string
}

func responsiveColumnWidths(columns []structuredColumn, available, gap float32) []float32 {
	widths := make([]float32, len(columns))
	if len(columns) == 0 {
		return widths
	}
	if available < 0 {
		available = 0
	}
	gaps := gap * float32(len(columns)-1)
	usable := available - gaps
	if usable < 0 {
		usable = 0
	}
	minimumTotal := float32(0)
	weightTotal := float32(0)
	for _, column := range columns {
		minimumTotal += maxFloat32(column.Minimum, 0)
		weightTotal += maxFloat32(column.Weight, 0)
	}
	if minimumTotal > usable && minimumTotal > 0 {
		scale := usable / minimumTotal
		for i, column := range columns {
			widths[i] = maxFloat32(column.Minimum, 0) * scale
		}
	} else {
		for i, column := range columns {
			widths[i] = maxFloat32(column.Minimum, 0)
		}
		surplus := usable - minimumTotal
		if weightTotal > 0 {
			for i, column := range columns {
				widths[i] += surplus * maxFloat32(column.Weight, 0) / weightTotal
			}
		} else {
			widths[len(widths)-1] += surplus
		}
	}
	// Put any floating-point remainder in the last column. This makes the
	// allocator an exact conservation contract instead of an approximation.
	allocated := float32(0)
	for _, width := range widths {
		allocated += width
	}
	widths[len(widths)-1] += usable - allocated
	return widths
}

type structuredColumnsLayout struct {
	columns []structuredColumn
	gap     float32
	profile *columnProfile
}

func (l *structuredColumnsLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if len(objects) == 0 {
		return
	}
	widths := responsiveColumnWidths(l.columns, size.Width, l.gap)
	if l.profile != nil {
		widths = l.profile.widths(size.Width, l.gap)
	}
	x := float32(0)
	for i := range l.columns {
		if i >= len(objects) {
			break
		}
		object := objects[i]
		width := widths[i]
		object.Move(fyne.NewPos(x, 0))
		object.Resize(fyne.NewSize(width, size.Height))
		if i < len(l.columns)-1 && len(objects) > len(l.columns)+i {
			boundary := objects[len(l.columns)+i]
			boundary.Move(fyne.NewPos(x+width+(l.gap-columnBoundaryWidth)/2, 0))
			boundary.Resize(fyne.NewSize(columnBoundaryWidth, size.Height))
		}
		x += width + l.gap
	}
}

func (l *structuredColumnsLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	height := float32(0)
	for _, object := range objects {
		height = maxFloat32(height, object.MinSize().Height)
	}
	// Width deliberately stays at zero: the owning view controls width and the
	// allocator scales columns without creating horizontal overflow.
	return fyne.NewSize(0, height)
}

// structuredRow renders a full selectable record. Fyne's parent List owns the
// hover/focus/selection layer; this widget supplies a translucent alternating
// surface plus consistently aligned labels.
type structuredRow struct {
	widget.BaseWidget
	columns    []structuredColumn
	labels     []*widget.Label
	boundaries []*columnBoundary
	content    *fyne.Container
	profile    *columnProfile
	header     bool
	alternate  bool
	identity   string
	summary    string
	onSelect   func(string)
	onActivate func(string)
}

func newStructuredRow(columns []structuredColumn, header bool, onSelect, onActivate func(string)) *structuredRow {
	return newStructuredRowWithProfile(columns, header, nil, nil, onSelect, onActivate)
}

func newStructuredRowWithProfile(columns []structuredColumn, header bool, profile *columnProfile, onResize func(int, float32, bool), onSelect, onActivate func(string)) *structuredRow {
	row := &structuredRow{
		columns:    append([]structuredColumn(nil), columns...),
		header:     header,
		profile:    profile,
		onSelect:   onSelect,
		onActivate: onActivate,
	}
	objects := make([]fyne.CanvasObject, 0, len(columns)*2-1)
	row.labels = make([]*widget.Label, len(columns))
	for i, column := range columns {
		text := ""
		if header {
			text = column.Header
		}
		label := widget.NewLabelWithStyle(text, column.Alignment, fyne.TextStyle{Bold: header})
		label.Truncation = fyne.TextTruncateEllipsis
		row.labels[i] = label
		objects = append(objects, label)
	}
	if header && profile != nil {
		row.boundaries = make([]*columnBoundary, max(len(columns)-1, 0))
		for i := range row.boundaries {
			boundary := newColumnBoundary(columns[i].Header, columns[i+1].Header, i, onResize)
			row.boundaries[i] = boundary
			objects = append(objects, boundary)
		}
	}
	row.content = container.New(&structuredColumnsLayout{columns: row.columns, gap: structuredColumnGap, profile: profile}, objects...)
	if header {
		row.summary = structuredHeaderSummary(columns)
	}
	row.ExtendBaseWidget(row)
	return row
}

// columnBoundary is the shared pointer and keyboard interaction at one header
// boundary. It transfers width through the owning table instead of storing a
// second copy of layout state.
type columnBoundary struct {
	widget.BaseWidget
	left, right string
	index       int
	onResize    func(int, float32, bool)
	focused     bool
}

func newColumnBoundary(left, right string, index int, onResize func(int, float32, bool)) *columnBoundary {
	boundary := &columnBoundary{left: left, right: right, index: index, onResize: onResize}
	boundary.ExtendBaseWidget(boundary)
	return boundary
}

func (b *columnBoundary) AccessibilityLabel() string {
	return "Resize " + b.left + " and " + b.right + " columns"
}

func (b *columnBoundary) AccessibilityRole() fyne.AccessibleRole { return fyne.AccessibleRoleButton }

func (b *columnBoundary) Cursor() desktop.Cursor { return desktop.HResizeCursor }

func (b *columnBoundary) Dragged(event *fyne.DragEvent) {
	if b.onResize != nil {
		b.onResize(b.index, event.Dragged.DX, false)
	}
}

func (b *columnBoundary) DragEnd() {
	if b.onResize != nil {
		b.onResize(b.index, 0, true)
	}
}

func (b *columnBoundary) FocusGained() { b.focused = true; b.Refresh() }

func (b *columnBoundary) FocusLost() { b.focused = false; b.Refresh() }

func (b *columnBoundary) TypedRune(rune) {}

func (b *columnBoundary) TypedKey(event *fyne.KeyEvent) {
	if b.onResize == nil {
		return
	}
	switch event.Name {
	case fyne.KeyLeft:
		b.onResize(b.index, -columnKeyboardStep, true)
	case fyne.KeyRight:
		b.onResize(b.index, columnKeyboardStep, true)
	}
}

func (b *columnBoundary) CreateRenderer() fyne.WidgetRenderer {
	line := canvas.NewRectangle(b.Theme().Color(theme.ColorNameSeparator, fyne.CurrentApp().Settings().ThemeVariant()))
	return &columnBoundaryRenderer{boundary: b, line: line}
}

type columnBoundaryRenderer struct {
	boundary *columnBoundary
	line     *canvas.Rectangle
}

func (r *columnBoundaryRenderer) Destroy() {}

func (r *columnBoundaryRenderer) Layout(size fyne.Size) {
	width := float32(1)
	if r.boundary.focused {
		width = 3
	}
	r.line.Move(fyne.NewPos((size.Width-width)/2, 2))
	r.line.Resize(fyne.NewSize(width, maxFloat32(size.Height-4, 0)))
}

func (r *columnBoundaryRenderer) MinSize() fyne.Size { return fyne.NewSize(columnBoundaryWidth, 0) }

func (r *columnBoundaryRenderer) Objects() []fyne.CanvasObject { return []fyne.CanvasObject{r.line} }

func (r *columnBoundaryRenderer) Refresh() {
	name := theme.ColorNameSeparator
	if r.boundary.focused {
		name = theme.ColorNamePrimary
	}
	r.line.FillColor = r.boundary.Theme().Color(name, fyne.CurrentApp().Settings().ThemeVariant())
	r.line.Refresh()
}

func (r *structuredRow) bind(model structuredRowModel, index int) {
	if model.Identity == "" || len(model.Cells) != len(r.labels) {
		r.identity = ""
		r.summary = ""
		r.alternate = false
		for _, label := range r.labels {
			label.SetText("")
			label.Importance = widget.MediumImportance
			label.TextStyle = fyne.TextStyle{}
			label.Refresh()
		}
		r.Refresh()
		return
	}
	r.identity = model.Identity
	r.summary = model.Summary
	r.alternate = index%2 == 1
	for i, cell := range model.Cells {
		text := cell.Text
		if text == "" {
			text = "Unknown"
		}
		r.labels[i].Importance = cell.Importance
		r.labels[i].TextStyle = cell.TextStyle
		r.labels[i].SetText(text)
	}
	r.Refresh()
}

// Tapped forwards a full-row selection to the owning virtualized list.
func (r *structuredRow) Tapped(*fyne.PointEvent) {
	if r.identity != "" && r.onSelect != nil {
		r.onSelect(r.identity)
	}
}

// DoubleTapped activates the stable identity currently bound to this recycled
// row widget.
func (r *structuredRow) DoubleTapped(*fyne.PointEvent) {
	if r.identity != "" && r.onActivate != nil {
		r.onActivate(r.identity)
	}
}

// AccessibilityLabel exposes the full row rather than its possibly truncated
// visual cells.
func (r *structuredRow) AccessibilityLabel() string { return r.summary }

// AccessibilityRole identifies a row as a composite group of labeled values.
func (r *structuredRow) AccessibilityRole() fyne.AccessibleRole { return fyne.AccessibleRoleContainer }

func (r *structuredRow) CreateRenderer() fyne.WidgetRenderer {
	background := canvas.NewRectangle(color.Transparent)
	return &structuredRowRenderer{row: r, background: background, objects: []fyne.CanvasObject{background, r.content}}
}

type structuredRowRenderer struct {
	row        *structuredRow
	background *canvas.Rectangle
	objects    []fyne.CanvasObject
}

func (r *structuredRowRenderer) Destroy() {}

func (r *structuredRowRenderer) Layout(size fyne.Size) {
	r.background.Resize(size)
	r.row.content.Resize(size)
}

func (r *structuredRowRenderer) MinSize() fyne.Size { return r.row.content.MinSize() }

func (r *structuredRowRenderer) Objects() []fyne.CanvasObject { return r.objects }

func (r *structuredRowRenderer) Refresh() {
	variant := fyne.CurrentApp().Settings().ThemeVariant()
	switch {
	case r.row.header:
		r.background.FillColor = r.row.Theme().Color(theme.ColorNameHeaderBackground, variant)
	case r.row.alternate:
		r.background.FillColor = r.row.Theme().Color(colorNameAlternateRow, variant)
	default:
		r.background.FillColor = color.Transparent
	}
	r.background.Refresh()
	r.row.content.Refresh()
}

// structuredList combines a fixed header, virtualized vertical body, and an
// optional complete-value disclosure. It deliberately exposes the underlying
// list to existing view tests and Fyne keyboard behavior.
type structuredList struct {
	columns          []structuredColumn
	profile          *columnProfile
	rows             []structuredRowModel
	header           *structuredRow
	list             *widget.List
	disclosure       *widget.Label
	root             fyne.CanvasObject
	emptyDisclosure  string
	selectedIdentity string
	onSelected       func(string)
	onActivate       func(string)
}

func newStructuredList(columns []structuredColumn, emptyDisclosure string, onSelected, onActivate func(string)) *structuredList {
	return newStructuredListWithProfile(columns, emptyDisclosure, onSelected, onActivate, nil)
}

func newAdjustableStructuredList(columns []structuredColumn, emptyDisclosure string, onSelected, onActivate func(string), preferences fyne.Preferences, preferenceKey string) *structuredList {
	return newStructuredListWithProfile(columns, emptyDisclosure, onSelected, onActivate, newColumnProfile(columns, preferences, preferenceKey))
}

func newStructuredListWithProfile(columns []structuredColumn, emptyDisclosure string, onSelected, onActivate func(string), profile *columnProfile) *structuredList {
	table := &structuredList{
		columns:         append([]structuredColumn(nil), columns...),
		profile:         profile,
		emptyDisclosure: emptyDisclosure,
		onSelected:      onSelected,
		onActivate:      onActivate,
	}
	table.header = newStructuredRowWithProfile(table.columns, true, profile, table.resizeColumns, nil, nil)
	table.list = widget.NewList(
		func() int { return len(table.rows) },
		func() fyne.CanvasObject {
			return newStructuredRowWithProfile(table.columns, false, profile, nil, table.selectIdentity, table.activateIdentity)
		},
		func(index widget.ListItemID, object fyne.CanvasObject) {
			row := object.(*structuredRow)
			if index < 0 || index >= len(table.rows) {
				row.bind(structuredRowModel{}, 0)
				return
			}
			row.bind(table.rows[index], index)
		},
	)
	table.disclosure = widget.NewLabel(emptyDisclosure)
	table.disclosure.Wrapping = fyne.TextWrapWord
	table.disclosure.Selectable = true
	table.disclosure.Importance = widget.LowImportance
	if emptyDisclosure == "" {
		table.disclosure.Hide()
	}
	table.list.OnSelected = func(index widget.ListItemID) {
		if index < 0 || index >= len(table.rows) {
			return
		}
		row := table.rows[index]
		table.selectedIdentity = row.Identity
		if table.emptyDisclosure != "" {
			table.disclosure.SetText(row.Summary)
		}
		if table.onSelected != nil {
			table.onSelected(row.Identity)
		}
	}
	table.list.OnUnselected = func(index widget.ListItemID) {
		if index >= 0 && index < len(table.rows) && table.rows[index].Identity != table.selectedIdentity {
			return
		}
		table.selectedIdentity = ""
		if table.emptyDisclosure != "" {
			table.disclosure.SetText(table.emptyDisclosure)
		}
	}
	table.root = container.NewBorder(table.header, table.disclosure, nil, nil, table.list)
	return table
}

func (t *structuredList) resizeColumns(boundary int, delta float32, persist bool) {
	if t.profile == nil {
		return
	}
	if delta != 0 && t.profile.adjust(boundary, delta, t.header.Size().Width, structuredColumnGap) {
		t.header.Refresh()
		t.list.Refresh()
	}
	if persist {
		t.profile.persist()
	}
}

func (t *structuredList) resetColumns() {
	if t.profile == nil {
		return
	}
	t.profile.reset()
	t.header.Refresh()
	t.list.Refresh()
}

func (t *structuredList) setRows(rows []structuredRowModel) {
	selected := t.selectedIdentity
	t.list.UnselectAll()
	t.rows = append(t.rows[:0], rows...)
	t.selectedIdentity = ""
	if index := t.indexOf(selected); index >= 0 {
		t.list.Select(index)
	} else if t.emptyDisclosure != "" {
		t.disclosure.SetText(t.emptyDisclosure)
	}
	t.list.Refresh()
}

func (t *structuredList) indexOf(identity string) int {
	if identity == "" {
		return -1
	}
	for index, row := range t.rows {
		if row.Identity == identity {
			return index
		}
	}
	return -1
}

func (t *structuredList) selectedIndex() int { return t.indexOf(t.selectedIdentity) }

func (t *structuredList) selectIdentity(identity string) {
	if index := t.indexOf(identity); index >= 0 {
		t.list.Select(index)
	}
}

func (t *structuredList) activateIdentity(identity string) {
	if t.indexOf(identity) >= 0 && t.onActivate != nil {
		t.onActivate(identity)
	}
}

func structuredRowSummary(columns []structuredColumn, cells []structuredCell) string {
	parts := make([]string, 0, min(len(columns), len(cells)))
	for index := 0; index < len(columns) && index < len(cells); index++ {
		value := cells[index].FullText
		if value == "" {
			value = cells[index].Text
		}
		if value == "" {
			value = "Unknown"
		}
		parts = append(parts, columns[index].Header+": "+value)
	}
	return strings.Join(parts, " | ")
}

func structuredHeaderSummary(columns []structuredColumn) string {
	headers := make([]string, len(columns))
	for index, column := range columns {
		headers[index] = column.Header
	}
	return "Columns: " + strings.Join(headers, ", ")
}

func normalizedWords(value, fallback string, uppercase bool) string {
	words := strings.Join(strings.Fields(strings.NewReplacer("_", " ", "-", " ").Replace(value)), " ")
	if words == "" {
		return fallback
	}
	if uppercase {
		return strings.ToUpper(words)
	}
	words = strings.ToLower(words)
	runes := []rune(words)
	return strings.ToUpper(string(runes[0])) + string(runes[1:])
}

func stablePresentationIdentity(parts ...string) string {
	var identity strings.Builder
	for _, part := range parts {
		identity.WriteString(strconv.Itoa(len(part)))
		identity.WriteByte(':')
		identity.WriteString(part)
		identity.WriteByte(';')
	}
	return identity.String()
}

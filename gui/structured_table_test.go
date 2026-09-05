package gui

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func TestResponsiveColumnWidthsConserveAvailableWidth(t *testing.T) {
	columns := []structuredColumn{
		{Header: "Identity", Minimum: 120, Weight: 3},
		{Header: "State", Minimum: 80},
		{Header: "Detail", Minimum: 100, Weight: 1},
	}
	for _, available := range []float32{640, 300, 72} {
		widths := responsiveColumnWidths(columns, available, 4)
		if len(widths) != len(columns) {
			t.Fatalf("available=%v width count=%d, want %d", available, len(widths), len(columns))
		}
		got := float32(0)
		for _, width := range widths {
			if width < 0 || math.IsNaN(float64(width)) {
				t.Fatalf("available=%v invalid widths=%v", available, widths)
			}
			got += width
		}
		if len(widths) > 1 {
			got += float32(len(widths)-1) * 4
		}
		if math.Abs(float64(got-available)) > 0.01 {
			t.Fatalf("available=%v allocated=%v widths=%v", available, got, widths)
		}
	}
}

func TestResponsiveColumnWidthsProtectMinimaAndWeightSurplus(t *testing.T) {
	columns := []structuredColumn{
		{Header: "Primary", Minimum: 100, Weight: 3},
		{Header: "Fixed", Minimum: 50},
		{Header: "Flexible", Minimum: 50, Weight: 1},
	}
	wide := responsiveColumnWidths(columns, 304, 2)
	if got, want := wide, []float32{175, 50, 75}; !reflect.DeepEqual(got, want) {
		t.Fatalf("wide widths=%v, want %v", got, want)
	}
	narrow := responsiveColumnWidths(columns, 104, 2)
	if got, want := narrow, []float32{50, 25, 25}; !reflect.DeepEqual(got, want) {
		t.Fatalf("narrow widths=%v, want %v", got, want)
	}
}

func TestProfileWidthsConserveAndClampAtSupportedAndNarrowSizes(t *testing.T) {
	columns := []structuredColumn{
		{Header: "When", Minimum: 140},
		{Header: "Task", Minimum: 160},
		{Header: "State", Minimum: 100},
	}
	profile := &columnProfile{columns: columns, values: []float32{0.6, 0.2, 0.2}}
	wide := profile.widths(604, 2)
	if got, want := wide, []float32{330, 160, 110}; !closeFloat32s(got, want) {
		t.Fatalf("wide profile widths=%v, want %v", got, want)
	}
	narrow := profile.widths(204, 2)
	if got, want := narrow, []float32{70, 80, 50}; !closeFloat32s(got, want) {
		t.Fatalf("narrow profile widths=%v, want %v", got, want)
	}
}

func TestAdjustableStructuredListBoundarySupportsPointerKeyboardAndAccessibility(t *testing.T) {
	prefs := testApp.Preferences()
	const key = "test.columns.boundary"
	t.Cleanup(func() { prefs.RemoveValue(key) })
	columns := []structuredColumn{
		{Header: "When", Minimum: 100, Preferred: 1},
		{Header: "Task", Minimum: 100, Preferred: 1},
		{Header: "State", Minimum: 80, Preferred: 1},
	}
	table := newAdjustableStructuredList(columns, "", nil, nil, prefs, key)
	if got := len(table.header.boundaries); got != 2 {
		t.Fatalf("boundary count=%d, want 2", got)
	}
	table.header.Resize(fyne.NewSize(600, 36))
	boundary := table.header.boundaries[0]
	if boundary.AccessibilityLabel() != "Resize When and Task columns" || boundary.AccessibilityRole() != fyne.AccessibleRoleButton {
		t.Fatalf("boundary accessibility=%q/%q", boundary.AccessibilityLabel(), boundary.AccessibilityRole())
	}
	if boundary.Cursor() != desktop.HResizeCursor {
		t.Fatalf("boundary cursor=%v", boundary.Cursor())
	}
	before := table.profile.widths(600, structuredColumnGap)
	boundary.FocusGained()
	if !boundary.focused {
		t.Fatal("boundary did not expose focus state")
	}
	boundary.TypedKey(&fyne.KeyEvent{Name: fyne.KeyRight})
	afterKey := table.profile.widths(600, structuredColumnGap)
	if math.Abs(float64(afterKey[0]-before[0]-columnKeyboardStep)) > 0.01 {
		t.Fatalf("keyboard widths=%v, before=%v", afterKey, before)
	}
	boundary.Dragged(&fyne.DragEvent{Dragged: fyne.Delta{DX: 20}})
	boundary.DragEnd()
	afterDrag := table.profile.widths(600, structuredColumnGap)
	if math.Abs(float64(afterDrag[0]-afterKey[0]-20)) > 0.01 || afterDrag[2] != before[2] {
		t.Fatalf("drag widths=%v, keyboard=%v", afterDrag, afterKey)
	}
	if got := newColumnProfile(columns, prefs, key).widths(600, structuredColumnGap); !closeFloat32s(got, afterDrag) {
		t.Fatalf("persisted widths=%v, want %v", got, afterDrag)
	}
	boundary.FocusLost()
	if boundary.focused {
		t.Fatal("boundary retained focus state")
	}
}

func TestAdjustableStructuredListKeepsHeaderAndBodyAlignedAfterResize(t *testing.T) {
	columns := []structuredColumn{
		{Header: "When", Minimum: 100, Preferred: 1},
		{Header: "Task", Minimum: 100, Preferred: 1},
		{Header: "State", Minimum: 80, Preferred: 1},
	}
	table := newAdjustableStructuredList(columns, "", nil, nil, nil, "")
	table.header.Resize(fyne.NewSize(600, 36))
	body := table.list.CreateItem().(*structuredRow)
	body.bind(structuredRowModel{Identity: "row", Cells: []structuredCell{{Text: "a"}, {Text: "b"}, {Text: "c"}}}, 0)
	body.Resize(fyne.NewSize(600, 36))
	table.header.boundaries[1].TypedKey(&fyne.KeyEvent{Name: fyne.KeyLeft})
	body.Refresh()
	body.Resize(fyne.NewSize(600, 36))
	for i := range columns {
		if table.header.labels[i].Position() != body.labels[i].Position() || table.header.labels[i].Size() != body.labels[i].Size() {
			t.Fatalf("column %d header=(%v,%v) body=(%v,%v)", i, table.header.labels[i].Position(), table.header.labels[i].Size(), body.labels[i].Position(), body.labels[i].Size())
		}
	}
}

func TestColumnBoundariesRenderInBothThemesAtSupportedWidths(t *testing.T) {
	t.Cleanup(func() { applyBrandTheme(testApp.Settings(), defaultAppearancePreferences()) })
	columns := []structuredColumn{
		{Header: "When", Minimum: 145, Preferred: 28},
		{Header: "Task", Minimum: 160, Preferred: 31},
		{Header: "Event", Minimum: 105, Preferred: 18},
		{Header: "Outcome", Minimum: 125, Preferred: 23},
	}
	for _, mode := range []appearanceMode{appearanceDark, appearanceLight} {
		applyBrandTheme(testApp.Settings(), appearancePreferences{Mode: mode, Font: fontSystem, ScrollSensitivity: defaultScrollSensitivity})
		table := newAdjustableStructuredList(columns, "", nil, nil, nil, "")
		for _, width := range []float32{560, 1200} {
			table.header.Resize(fyne.NewSize(width, 36))
			allocated := structuredColumnGap * float32(len(columns)-1)
			for _, columnWidth := range table.profile.widths(width, structuredColumnGap) {
				allocated += columnWidth
			}
			if math.Abs(float64(allocated-width)) > 0.01 {
				t.Fatalf("mode=%q width=%v allocated=%v", mode, width, allocated)
			}
		}
		for _, boundary := range table.header.boundaries {
			renderer := boundary.CreateRenderer().(*columnBoundaryRenderer)
			renderer.Refresh()
			_, _, _, alpha := renderer.line.FillColor.RGBA()
			if alpha == 0 {
				t.Fatalf("mode=%q boundary is invisible", mode)
			}
		}
	}
}

func TestStructuredColumnsLayoutAlignsHeaderAndBody(t *testing.T) {
	columns := []structuredColumn{
		{Header: "One", Minimum: 80, Weight: 2},
		{Header: "Two", Minimum: 60},
		{Header: "Three", Minimum: 100, Weight: 1},
	}
	header := newStructuredRow(columns, true, nil, nil)
	body := newStructuredRow(columns, false, nil, nil)
	body.bind(structuredRowModel{Identity: "row", Cells: []structuredCell{{Text: "a"}, {Text: "b"}, {Text: "c"}}, Summary: "row"}, 0)
	header.Resize(fyne.NewSize(500, 36))
	body.Resize(fyne.NewSize(500, 36))
	for i := range columns {
		if header.labels[i].Position() != body.labels[i].Position() || header.labels[i].Size() != body.labels[i].Size() {
			t.Fatalf("column %d header=(%v,%v) body=(%v,%v)", i, header.labels[i].Position(), header.labels[i].Size(), body.labels[i].Position(), body.labels[i].Size())
		}
	}
}

func TestStructuredRowBindsSemanticsTruncationAndIdentity(t *testing.T) {
	columns := []structuredColumn{{Header: "Name", Minimum: 100}, {Header: "State", Minimum: 80}}
	var selected, activated []string
	row := newStructuredRow(columns, false,
		func(id string) { selected = append(selected, id) },
		func(id string) { activated = append(activated, id) },
	)
	model := structuredRowModel{
		Identity: "stable-id",
		Cells: []structuredCell{
			{Text: "A very long task name", FullText: "A very long task name", Importance: widget.MediumImportance},
			{Text: "ERROR", Importance: widget.DangerImportance, TextStyle: fyne.TextStyle{Bold: true}},
		},
		Summary: "Name: A very long task name | State: ERROR",
	}
	row.bind(model, 1)
	if !row.alternate || row.identity != "stable-id" || row.AccessibilityLabel() != model.Summary {
		t.Fatalf("bound row alternate=%v identity=%q accessible=%q", row.alternate, row.identity, row.AccessibilityLabel())
	}
	for i, label := range row.labels {
		if label.Truncation != fyne.TextTruncateEllipsis || label.Text != model.Cells[i].Text || label.Importance != model.Cells[i].Importance {
			t.Fatalf("cell %d = text %q truncation %v importance %v", i, label.Text, label.Truncation, label.Importance)
		}
	}
	test.Tap(row)
	test.DoubleTap(row)
	if !reflect.DeepEqual(selected, []string{"stable-id"}) || !reflect.DeepEqual(activated, []string{"stable-id"}) {
		t.Fatalf("selected=%v activated=%v", selected, activated)
	}
}

func TestStructuredRowRejectsInvalidCellCountAndUnboundActivation(t *testing.T) {
	row := newStructuredRow([]structuredColumn{{Header: "A", Minimum: 1}, {Header: "B", Minimum: 1}}, false, func(string) {
		t.Fatal("unbound row selected")
	}, func(string) {
		t.Fatal("unbound row activated")
	})
	row.bind(structuredRowModel{Identity: "invalid", Cells: []structuredCell{{Text: "only one"}}}, 0)
	if row.identity != "" {
		t.Fatalf("invalid model left identity %q", row.identity)
	}
	test.Tap(row)
	test.DoubleTap(row)
}

func TestStructuredListKeepsHeaderOnEmptyAndReconcilesStableIdentity(t *testing.T) {
	columns := []structuredColumn{{Header: "Name", Minimum: 100}, {Header: "State", Minimum: 80}}
	var selected []string
	table := newStructuredList(columns, "Select a row for complete values.", func(id string) {
		selected = append(selected, id)
	}, nil)
	if got := table.list.Length(); got != 0 {
		t.Fatalf("empty row count=%d", got)
	}
	if got := []string{table.header.labels[0].Text, table.header.labels[1].Text}; !reflect.DeepEqual(got, []string{"Name", "State"}) {
		t.Fatalf("headers=%v", got)
	}
	rows := []structuredRowModel{
		{Identity: "a", Cells: []structuredCell{{Text: "Alpha"}, {Text: "Enabled"}}, Summary: "Name: Alpha | State: Enabled"},
		{Identity: "b", Cells: []structuredCell{{Text: "Beta"}, {Text: "Disabled"}}, Summary: "Name: Beta | State: Disabled"},
	}
	table.setRows(rows)
	table.list.Select(0)
	if table.selectedIdentity != "a" || table.disclosure.Text != rows[0].Summary {
		t.Fatalf("selection=%q disclosure=%q", table.selectedIdentity, table.disclosure.Text)
	}
	table.setRows([]structuredRowModel{rows[1], rows[0]})
	if table.selectedIdentity != "a" || table.selectedIndex() != 1 || table.disclosure.Text != rows[0].Summary {
		t.Fatalf("reordered selection=%q index=%d disclosure=%q", table.selectedIdentity, table.selectedIndex(), table.disclosure.Text)
	}
	table.setRows([]structuredRowModel{rows[1]})
	if table.selectedIdentity != "" || table.selectedIndex() != -1 || table.disclosure.Text != "Select a row for complete values." {
		t.Fatalf("removed selection=%q index=%d disclosure=%q", table.selectedIdentity, table.selectedIndex(), table.disclosure.Text)
	}
	if !reflect.DeepEqual(selected, []string{"a", "a"}) {
		t.Fatalf("selection callbacks=%v", selected)
	}
}

func TestStructuredListVirtualizesHundredRowsAndBindsCurrentModel(t *testing.T) {
	columns := []structuredColumn{{Header: "Name", Minimum: 100}}
	table := newStructuredList(columns, "Nothing selected.", nil, nil)
	rows := make([]structuredRowModel, 100)
	for i := range rows {
		rows[i] = structuredRowModel{
			Identity: fmt.Sprintf("row-%03d", i),
			Cells:    []structuredCell{{Text: fmt.Sprintf("Value %03d", i)}},
			Summary:  fmt.Sprintf("Name: Value %03d", i),
		}
	}
	table.setRows(rows)
	if got := table.list.Length(); got != 100 {
		t.Fatalf("row count=%d", got)
	}
	item := table.list.CreateItem().(*structuredRow)
	table.list.UpdateItem(99, item)
	if item.identity != "row-099" || item.labels[0].Text != "Value 099" || !item.alternate {
		t.Fatalf("bound last row identity=%q text=%q alternate=%v", item.identity, item.labels[0].Text, item.alternate)
	}
}

func TestStructuredRowSummaryUsesHeadersAndFullValues(t *testing.T) {
	columns := []structuredColumn{{Header: "Task"}, {Header: "Group"}}
	cells := []structuredCell{{Text: "short", FullText: "完整 task"}, {Text: "none", FullText: "Group with spaces"}}
	if got, want := structuredRowSummary(columns, cells), "Task: 完整 task | Group: Group with spaces"; got != want {
		t.Fatalf("summary=%q, want %q", got, want)
	}
}

func TestNormalizedWordsTitleCasesUnicodeWithoutCorruption(t *testing.T) {
	if got := normalizedWords("éTAT_initial", "Unknown", false); got != "État initial" {
		t.Fatalf("normalizedWords Unicode result = %q", got)
	}
}

package gui

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"fyne.io/fyne/v2"
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

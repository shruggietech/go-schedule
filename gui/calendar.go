package gui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/shruggietech/go-schedule/internal/api/server"
)

const dayKeyFmt = "2006-01-02"

// occurrencesByDay buckets occurrences by their local calendar day (YYYY-MM-DD).
func occurrencesByDay(occ []server.Occurrence) map[string][]server.Occurrence {
	m := make(map[string][]server.Occurrence)
	for _, o := range occ {
		k := o.Time.Local().Format(dayKeyFmt)
		m[k] = append(m[k], o)
	}
	return m
}

// buildCalendarGrid renders a month grid for the month containing `month`. Days
// with occurrences are marked; tapping a day invokes onSelect with that day.
func buildCalendarGrid(occ []server.Occurrence, month time.Time, onSelect func(day time.Time)) *fyne.Container {
	byDay := occurrencesByDay(occ)

	grid := container.NewGridWithColumns(7)
	for _, wd := range []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"} {
		grid.Add(widget.NewLabelWithStyle(wd, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
	}

	first := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, month.Location())
	// Leading blanks so the 1st lands under its weekday (Mon=0 ... Sun=6).
	lead := (int(first.Weekday()) + 6) % 7
	for i := 0; i < lead; i++ {
		grid.Add(widget.NewLabel(""))
	}

	daysInMonth := first.AddDate(0, 1, -1).Day()
	for d := 1; d <= daysInMonth; d++ {
		day := time.Date(month.Year(), month.Month(), d, 0, 0, 0, 0, month.Location())
		label := fmt.Sprintf("%d", d)
		if items := byDay[day.Format(dayKeyFmt)]; len(items) > 0 {
			label = fmt.Sprintf("%d •%d", d, len(items))
		}
		dd := day
		grid.Add(newToolbarButtonPlain(label, func() {
			if onSelect != nil {
				onSelect(dd)
			}
		}))
	}
	return grid
}

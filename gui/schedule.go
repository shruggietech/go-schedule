package gui

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/shruggietech/go-schedule/internal/api/server"
)

// scheduleState holds the Schedule tab's shared state. It is guarded by a mutex
// because the occurrences are written by the async loader goroutine and read by
// Fyne widget callbacks; the selected window and view mode are likewise read off
// the loader goroutine and written by UI callbacks.
type scheduleState struct {
	mu   sync.Mutex
	occ  []server.Occurrence
	days int
	view string
}

func (s *scheduleState) snapshotOcc() []server.Occurrence {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.occ
}
func (s *scheduleState) setOcc(o []server.Occurrence) { s.mu.Lock(); s.occ = o; s.mu.Unlock() }
func (s *scheduleState) getDays() int                 { s.mu.Lock(); defer s.mu.Unlock(); return s.days }
func (s *scheduleState) setDays(d int)                { s.mu.Lock(); s.days = d; s.mu.Unlock() }
func (s *scheduleState) getView() string              { s.mu.Lock(); defer s.mu.Unlock(); return s.view }
func (s *scheduleState) setView(v string)             { s.mu.Lock(); s.view = v; s.mu.Unlock() }

// buildScheduleTab shows past and upcoming runs over a window, in either a List
// (agenda) view or a toggleable Calendar (month-grid) view (FR-025). The selected
// window is preserved across view toggles (FR-027) and both views update live.
func (a *App) buildScheduleTab() fyne.CanvasObject {
	st := &scheduleState{days: 7, view: "List"}

	list := widget.NewList(
		func() int { return len(st.snapshotOcc()) },
		func() fyne.CanvasObject { return widget.NewLabel("template") },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			occ := st.snapshotOcc()
			if i < 0 || i >= len(occ) {
				return
			}
			e := occ[i]
			marker := "▷" // scheduled (future)
			if e.Kind == "past" {
				marker = "✓"
				if e.Outcome != "success" && e.Outcome != "" {
					marker = "✗"
				}
			}
			o.(*widget.Label).SetText(fmt.Sprintf("%s  %s   %s", marker, fmtTime(e.Time), e.TaskName))
		},
	)

	calBox := container.NewVBox()
	renderCalendar := func() {
		occ := st.snapshotOcc()
		detail := widget.NewLabel("Select a day to see its runs.")
		grid := buildCalendarGrid(occ, time.Now(), func(day time.Time) {
			detail.SetText(occurrencesOnDayText(occ, day))
		})
		calBox.Objects = []fyne.CanvasObject{grid, widget.NewSeparator(), detail}
		calBox.Refresh()
	}

	content := container.NewStack(list)
	render := func() {
		if st.getView() == "Calendar" {
			renderCalendar()
			content.Objects = []fyne.CanvasObject{calBox}
		} else {
			content.Objects = []fyne.CanvasObject{list}
		}
		content.Refresh()
	}

	load := func() {
		days := st.getDays()
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			from := time.Now().Add(-24 * time.Hour)
			to := time.Now().Add(time.Duration(days) * 24 * time.Hour)
			resp, err := a.backend.GetCalendar(ctx, from, to)
			fyne.Do(func() {
				if err != nil {
					a.showError(err)
					return
				}
				st.setOcc(sortByTime(resp.Occurrences))
				list.Refresh()
				if st.getView() == "Calendar" {
					renderCalendar()
				}
			})
		}()
	}
	a.registerRefresher(load)

	rangeSel := widget.NewSelect([]string{"1 day", "7 days", "30 days"}, func(s string) {
		switch s {
		case "1 day":
			st.setDays(1)
		case "30 days":
			st.setDays(30)
		default:
			st.setDays(7)
		}
		load()
	})
	rangeSel.SetSelected("7 days")

	viewSel := widget.NewSelect([]string{"List", "Calendar"}, func(s string) {
		st.setView(s)
		render() // toggling never touches the window, so it is preserved
	})
	viewSel.SetSelected("List")

	// No manual Refresh: the view updates live from the event stream (FR-023).
	toolbar := container.NewHBox(
		widget.NewLabel("View:"), viewSel,
		widget.NewLabel("Window:"), rangeSel,
	)
	return container.NewBorder(toolbar, nil, nil, nil, content)
}

// occurrencesOnDayText renders the occurrences falling on day as a text block.
func occurrencesOnDayText(occ []server.Occurrence, day time.Time) string {
	var b strings.Builder
	for _, o := range occ {
		if o.Time.Local().Format(dayKeyFmt) == day.Format(dayKeyFmt) {
			fmt.Fprintf(&b, "%s   %s\n", fmtTime(o.Time), o.TaskName)
		}
	}
	if b.Len() == 0 {
		return "No runs on " + day.Format("Mon Jan 2")
	}
	return strings.TrimRight(b.String(), "\n")
}

// sortByTime orders occurrences ascending by time (simple insertion-free sort).
func sortByTime(in []server.Occurrence) []server.Occurrence {
	out := make([]server.Occurrence, len(in))
	copy(out, in)
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j].Time.Before(out[j-1].Time); j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
}

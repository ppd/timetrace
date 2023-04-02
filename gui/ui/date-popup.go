package ui

import (
	"time"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	fwidget "fyne.io/x/fyne/widget"
	"github.com/dominikbraun/timetrace/gui/state"
)

func ShowDatePopup(initialDate time.Time, cb func(time.Time)) {
	var calendarPopup *widget.PopUp

	calendar := fwidget.NewCalendar(initialDate, func(t time.Time) {
		calendarPopup.Hide()
		cb(t)
	})

	todayButton := widget.NewButton("Today", func() {
		today, _ := state.GetState().T.Formatter().ParseDate("today")
		calendarPopup.Hide()
		cb(today)
	})

	calendarPopup = widget.NewModalPopUp(
		container.NewBorder(
			todayButton,
			nil,
			nil,
			nil,
			calendar,
		),
		state.GetState().MainWindow.Canvas(),
	)
	calendarPopup.Show()
}

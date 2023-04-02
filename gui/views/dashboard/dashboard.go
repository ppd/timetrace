package dashboard

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/dominikbraun/timetrace/gui/state"
	"github.com/dominikbraun/timetrace/gui/ui"
)

func Dashboard() fyne.CanvasObject {
	state.DashboardState().RefreshStatePeriodically()

	projectButton := widget.NewButtonWithIcon("", theme.ListIcon(), func() {
		state.CoreState().GoToProjectsView()
	})

	aboutButton := widget.NewButtonWithIcon("", theme.InfoIcon(), func() {
		state.CoreState().ChangeView(state.About)
	})

	calendarButton := widget.NewButtonWithIcon("", theme.MoreVerticalIcon(), func() {
		currentDate, _ := state.DashboardState().Date.Get()
		ui.ShowDatePopup(currentDate, func(t time.Time) {
			state.DashboardState().Date.Set(t)
		})
	})

	reportButton := widget.NewButtonWithIcon("", theme.DocumentIcon(), func() {
		state.CoreState().ChangeView(state.Report)
	})

	content := container.NewBorder(
		container.NewVBox(
			container.NewBorder(
				nil,
				nil,
				container.NewHBox(reportButton, projectButton, calendarButton),
				aboutButton,
				StartProject(),
			),
			widget.NewSeparator(),
			Status(),
			Tags(),
			widget.NewSeparator(),
		),
		nil,
		nil,
		nil,
		RecordList(),
	)

	return content
}

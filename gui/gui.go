package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/gui/state"
	"github.com/dominikbraun/timetrace/gui/views/dashboard"
	"github.com/dominikbraun/timetrace/gui/views/record"
)

func RunGui(t *core.Timetrace) {
	var window fyne.Window
	if showActiveTimetraceGUI() {
		return
	}

	go runServer(func() {
		window.Show()
	})

	theState := state.InitState(t)
	theState.RefreshState()
	theState.RefreshStatusPeriodically()

	a := app.New()
	window = a.NewWindow("Timetrace")

	if desk, ok := a.(desktop.App); ok {
		m := fyne.NewMenu("Timetrace",
			fyne.NewMenuItem("Show", func() {
				window.Show()
			}))
		desk.SetSystemTrayMenu(m)
		desk.SetSystemTrayIcon(theme.HistoryIcon())
	}

	window.Resize(fyne.NewSize(600, 400))

	// views
	dashboardView := dashboard.Dashboard()
	editRecordView := record.EditRecordView()

	// routing
	theState.ActiveView.AddListener(binding.NewDataListener(func() {
		activeView, _ := theState.ActiveView.Get()
		switch state.View(activeView) {
		case state.Main:
			window.SetContent(dashboardView)
		case state.EditRecord:
			window.SetContent(editRecordView)
		}
	}))

	window.SetCloseIntercept(func() {
		window.Hide()
	})

	window.Show()

	a.Run()
}

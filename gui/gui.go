package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/gui/state"
	"github.com/dominikbraun/timetrace/gui/views/about"
	"github.com/dominikbraun/timetrace/gui/views/dashboard"
	"github.com/dominikbraun/timetrace/gui/views/project"
	"github.com/dominikbraun/timetrace/gui/views/projects"
	"github.com/dominikbraun/timetrace/gui/views/record"
	"github.com/dominikbraun/timetrace/gui/views/report"
)

func RunGui(t *core.Timetrace) {
	var window fyne.Window
	if showActiveTimetraceGUI() {
		return
	}

	go runServer(func() {
		window.Show()
	})

	theState := state.InitCoreState(t)
	theState.UpdateProjects()

	a := app.New()
	window = a.NewWindow("Timetrace")
	state.CoreState().MainWindow = window

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
	projectsView := projects.Projects()
	editProjectView := project.Project()
	aboutView := about.About()
	reportView := report.Report()

	// routing
	theState.RequestedView.AddListener(binding.NewDataListener(func() {
		requestedView, _ := theState.RequestedView.Get()
		switch state.View(requestedView) {
		case state.Main:
			window.SetContent(dashboardView)
		case state.EditRecord:
			window.SetContent(editRecordView)
		case state.Projects:
			window.SetContent(projectsView)
		case state.EditProject:
			window.SetContent(editProjectView)
		case state.About:
			window.SetContent(aboutView)
		case state.Report:
			window.SetContent(reportView)
		}
		theState.SetActiveView(state.View(requestedView))
	}))

	window.SetCloseIntercept(func() {
		window.Hide()
	})

	window.Show()

	a.Run()
}

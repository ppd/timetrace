package dashboard

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	fwidget "fyne.io/x/fyne/widget"
	"github.com/dominikbraun/timetrace/gui/shared"
	"github.com/dominikbraun/timetrace/gui/state"
)

func StartProject() fyne.CanvasObject {
	theState := state.DashboardState()

	entry := fwidget.NewCompletionEntry([]string{})

	doesProjectExist := func() bool {
		if _, err := state.Timetrace().LoadProject(entry.Text); err != nil {
			return false
		} else {
			return true
		}
	}

	startTheProject := func(projectKey string) {
		if err := theState.StartProject(projectKey); err != nil {
			dialog.ShowError(err, state.CoreState().MainWindow)
			entry.SetText("")
		}
	}

	createAndStartTheProject := func() bool {
		projectKey := entry.Text
		dialog.ShowConfirm(
			"Create new project?",
			fmt.Sprintf("Do you want to create the new project '%s'?", projectKey),
			func(createNew bool) {
				if createNew {
					if err := theState.CreateProject(projectKey); err != nil {
						dialog.ShowError(err, state.CoreState().MainWindow)
					} else {
						state.CoreState().UpdateProjects()
						startTheProject(projectKey)
					}
				}
			},
			state.CoreState().MainWindow,
		)
		return true
	}

	startOrCreateProject := func() {
		if len(entry.Text) == 0 {
			return
		}
		if doesProjectExist() {
			startTheProject(entry.Text)
		} else {
			createAndStartTheProject()
		}
	}

	entry.ActionItem = widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() { startOrCreateProject() })
	entry.SetPlaceHolder("Start or create project")

	entry.OnChanged = func(s string) {
		if len(s) < 2 {
			entry.HideCompletion()
			return
		}
		projectLabels, _ := state.CoreState().ProjectLabels.Get()
		options := shared.FilterByContains(projectLabels, entry.Text)
		if len(options) == 0 {
			entry.ActionItem.(*widget.Button).Icon = theme.ContentAddIcon()
		} else {
			entry.ActionItem.(*widget.Button).Icon = theme.MediaPlayIcon()
		}
		entry.ActionItem.Refresh()
		entry.SetOptions(options)
		entry.ShowCompletion()
	}

	entry.CustomUpdate = func(i widget.ListItemID, o fyne.CanvasObject) {
		options := entry.Options

		if i >= len(options) {
			return
		}

		o.(*widget.Label).SetText(options[i])
	}

	entry.OnSubmitted = func(s string) { startOrCreateProject() }

	syncEnabled := func() {
		isActive, _ := theState.IsRecordActive.Get()
		if isActive {
			entry.Disable()
		} else {
			entry.Enable()
		}
	}
	theState.IsRecordActive.AddListener(
		binding.NewDataListener(syncEnabled),
	)

	state.CoreState().ActiveView.AddListener(binding.NewDataListener(func() {
		activeView, _ := state.CoreState().ActiveView.Get()
		if activeView == int(state.Main) {
			state.CoreState().MainWindow.Canvas().Focus(entry)
		}
	}))

	return entry
}

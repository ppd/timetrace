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
	state := state.GetState()
	projectLabels, _ := state.ProjectLabels.Get()

	entry := fwidget.NewCompletionEntry([]string{})

	doesProjectExist := func() bool {
		if _, err := state.T.LoadProject(entry.Text); err != nil {
			return false
		} else {
			return true
		}
	}

	startTheProject := func(projectKey string) {
		if err := state.StartProject(projectKey); err != nil {
			dialog.ShowError(err, state.MainWindow)
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
					if err := state.CreateProject(projectKey); err != nil {
						dialog.ShowError(err, state.MainWindow)
					} else {
						startTheProject(projectKey)
					}
				}
			},
			state.MainWindow,
		)
		return true
	}

	startOrCreateProject := func() {
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
		options := shared.FilterByContains(projectLabels, entry.Text)
		if len(options) == 0 {
			entry.ActionItem.(*widget.Button).Icon = theme.ContentAddIcon()
		} else {
			entry.ActionItem.(*widget.Button).Icon = theme.MediaPlayIcon()
		}
		entry.ActionItem.Refresh()
		entry.SetOptions(options)
		if fyne.CurrentApp().Driver().CanvasForObject(entry) != nil {
			entry.ShowCompletion()
		}
	}

	entry.OnSubmitted = func(s string) { startOrCreateProject() }

	syncEnabled := func() {
		isActive, _ := state.IsRecordActive.Get()
		if isActive {
			entry.Disable()
		} else {
			entry.Enable()
		}
	}
	state.IsRecordActive.AddListener(
		binding.NewDataListener(syncEnabled),
	)

	return entry
}

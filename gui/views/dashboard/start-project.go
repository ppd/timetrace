package dashboard

import (
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

	startTheProject := func() {
		if err := state.StartProject(entry.Text); err != nil {
			dialog.ShowError(err, fyne.CurrentApp().Driver().AllWindows()[0])
			entry.SetText("")
		}
	}

	entry.Entry.ActionItem = widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() { startTheProject() })
	entry.SetPlaceHolder("Start project")

	entry.OnChanged = func(s string) {
		if len(s) < 2 {
			entry.HideCompletion()
			return
		}
		options := shared.FilterByContains(projectLabels, entry.Text)
		entry.SetOptions(options)
		if fyne.CurrentApp().Driver().CanvasForObject(entry) != nil {
			entry.ShowCompletion()
		}
	}

	entry.OnSubmitted = func(s string) { startTheProject() }

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

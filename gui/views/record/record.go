package record

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	fwidget "fyne.io/x/fyne/widget"
	"github.com/dominikbraun/timetrace/gui/shared"
	"github.com/dominikbraun/timetrace/gui/state"
	"github.com/dominikbraun/timetrace/gui/ui"
)

func EditRecordView() fyne.CanvasObject {
	theState := state.GetState()

	startTimeEntry := ui.NewTimeEntry(theState.RecordToEditStart)
	endTimeEntry := ui.NewTimeEntry(theState.RecordToEditEnd)
	tagsEntry := widget.NewEntryWithData(theState.RecordToEditTags)

	projectEntry := fwidget.NewCompletionEntry([]string{})
	projectEntry.Bind(theState.RecordToEditProject)

	boundOnChanged := projectEntry.OnChanged
	projectEntry.OnChanged = func(s string) {
		boundOnChanged(s)
		if len(s) < 2 {
			projectEntry.HideCompletion()
			return
		}
		projectLabels, _ := theState.ProjectLabels.Get()
		options := shared.FilterByContains(projectLabels, projectEntry.Text)
		projectEntry.SetOptions(options)
		if fyne.CurrentApp().Driver().CanvasForObject(projectEntry) != nil && !theState.RecordToEditIsExternalChange {
			projectEntry.ShowCompletion()
		}
		theState.RecordToEditIsExternalChange = false
	}

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ConfirmIcon(), func() {
			theState.SaveRecordToEdit()
		}),
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			theState.EditRecord(theState.RecordToEdit)
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.DeleteIcon(), func() {
			dialog.ShowConfirm(
				"Delete?",
				"Are you sure you want to delete this record?",
				func(yes bool) {
					if yes {
						theState.DeleteRecordToEdit()
					}
				},
				fyne.CurrentApp().Driver().AllWindows()[0],
			)
		}),
	)

	form := container.New(
		layout.NewFormLayout(),
		widget.NewLabel("Project"),
		projectEntry,
		widget.NewLabel("Start"),
		startTimeEntry,
		widget.NewLabel("End"),
		endTimeEntry,
		widget.NewLabel("Tags"),
		tagsEntry,
	)

	content := container.NewVBox(toolbar, form)

	return content
}

package project

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/dominikbraun/timetrace/gui/state"
)

func Project() fyne.CanvasObject {
	theState := state.GetState()

	keyEntry := widget.NewEntryWithData(theState.ProjectToEditKey)
	chronosAccountEntry := widget.NewEntryWithData(theState.ProjectToEditChronosAccount)
	chronosProjectEntry := widget.NewEntryWithData(theState.ProjectToEditChronosProject)

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ConfirmIcon(), func() {
			theState.SaveProjectToEdit()
		}),
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			theState.EditProject(theState.ProjectToEdit.Key)
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.DeleteIcon(), func() {
			dialog.ShowConfirm(
				"Delete?",
				"Are you sure you want to delete this project?",
				func(yes bool) {
					if yes {
						if err := theState.DeleteProjectToEdit(); err != nil {
							dialog.ShowError(err, theState.MainWindow)
						}
					}
				},
				theState.MainWindow,
			)
		}),
	)

	form := container.New(
		layout.NewFormLayout(),
		widget.NewLabel("Key"),
		keyEntry,
		widget.NewLabel("Chronos Account"),
		chronosAccountEntry,
		widget.NewLabel("Chronos Project"),
		chronosProjectEntry,
	)

	content := container.NewVBox(toolbar, form)

	return content
}

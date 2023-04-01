package projects

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/dominikbraun/timetrace/gui/state"
)

func projectList() fyne.CanvasObject {
	theState := state.GetState()

	list := widget.NewListWithData(
		theState.FilteredProjects,
		func() fyne.CanvasObject {
			return widget.NewButton("template", func() {})
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			projectLabel, _ := i.(binding.String).Get()
			o.(*widget.Button).SetText(projectLabel)
			o.(*widget.Button).OnTapped = func() {
				if err := theState.EditProject(projectLabel); err != nil {
					dialog.ShowError(err, theState.MainWindow)
				}
			}
		},
	)

	return list
}

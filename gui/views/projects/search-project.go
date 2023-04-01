package projects

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/dominikbraun/timetrace/gui/state"
)

func searchProject() fyne.CanvasObject {
	theState := state.GetState()

	entry := widget.NewEntryWithData(theState.ProjectsFilter)

	entry.SetPlaceHolder("Search project")

	return entry
}

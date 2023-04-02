package dashboard

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/dominikbraun/timetrace/gui/state"
)

func Tags() fyne.CanvasObject {
	theState := state.DashboardState()
	tags := widget.NewEntryWithData(theState.Tags)
	tags.SetPlaceHolder("Tags for active project")
	tags.Hide()

	theState.IsRecordActive.AddListener(binding.NewDataListener(func() {
		isActive, _ := theState.IsRecordActive.Get()
		if isActive {
			tags.Show()
		} else {
			tags.Hide()
		}
	}))

	return tags
}

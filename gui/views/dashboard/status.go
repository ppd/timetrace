package dashboard

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/dominikbraun/timetrace/gui/state"
)

func Status() fyne.CanvasObject {
	theState := state.DashboardState()
	onStopped := func() {
		theState.Stop()
	}
	label := widget.NewLabelWithData(theState.Status)
	button := widget.NewButtonWithIcon("", theme.MediaStopIcon(), onStopped)

	syncEnabled := func() {
		isActive, _ := theState.IsRecordActive.Get()
		if isActive {
			button.Enable()
		} else {
			button.Disable()
		}
	}
	theState.IsRecordActive.AddListener(
		binding.NewDataListener(syncEnabled),
	)

	return container.NewBorder(
		nil,
		nil,
		nil,
		button,
		label,
	)
}

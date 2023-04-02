package report

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/dominikbraun/timetrace/gui/state"
	"github.com/dominikbraun/timetrace/gui/ui"
)

func Report() fyne.CanvasObject {
	theState := state.GetState()

	backButton := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		theState.GoToMainView()
	})

	currentDate, _ := theState.Date.Get()
	theDate := state.NewBoundTimeWithData(currentDate)
	dateEntry := ui.NewDateEntry(theDate)

	theDate.AddListener(binding.NewDataListener(func() {
		chosenDaten, _ := theDate.Get()
		theState.UpdateReport(chosenDaten)
	}))

	reportDisplay := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
	reportDisplay.Bind(theState.Report)

	toolbar := container.NewBorder(
		nil,
		nil,
		backButton,
		nil,
		dateEntry,
	)

	return container.NewBorder(
		toolbar,
		nil,
		nil,
		nil,
		reportDisplay,
	)
}

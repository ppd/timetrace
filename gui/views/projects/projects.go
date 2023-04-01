package projects

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/dominikbraun/timetrace/gui/state"
)

func Projects() fyne.CanvasObject {
	theState := state.GetState()

	backButton := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		theState.GoToMainView()
	})

	content := container.NewBorder(
		container.New(
			layout.NewFormLayout(),
			backButton,
			searchProject(),
		),
		nil,
		nil,
		nil,
		projectList(),
	)

	return content
}

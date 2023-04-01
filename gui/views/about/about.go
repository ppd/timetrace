package about

import (
	_ "embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/dominikbraun/timetrace/gui/state"
)

//go:generate sh -c "printf %s $(git describe --tags) > version.txt"
//go:embed version.txt
var version string

func About() fyne.CanvasObject {
	theState := state.GetState()

	backButton := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		theState.GoToMainView()
	})

	image := canvas.NewImageFromResource(resourceTimetraceCobraJpg)
	image.FillMode = canvas.ImageFillContain

	return container.NewBorder(
		container.New(
			layout.NewFormLayout(),
			backButton,
			widget.NewLabel("Timetrace - Kobra Edition - "+version),
		),
		nil,
		nil,
		nil,
		image,
	)
}

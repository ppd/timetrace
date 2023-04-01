package dashboard

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/dominikbraun/timetrace/gui/state"
)

func Dashboard() fyne.CanvasObject {
	projectButton := widget.NewButtonWithIcon("", theme.ListIcon(), func() {
		state.GetState().GoToProjectsView()
	})

	aboutButton := widget.NewButtonWithIcon("", theme.InfoIcon(), func() {
		state.GetState().ChangeView(state.About)
	})

	content := container.NewBorder(
		container.NewVBox(
			container.NewBorder(
				nil,
				nil,
				projectButton,
				aboutButton,
				StartProject(),
			),
			widget.NewSeparator(),
			Status(),
			Tags(),
			widget.NewSeparator(),
		),
		nil,
		nil,
		nil,
		RecordList(),
	)

	return content
}

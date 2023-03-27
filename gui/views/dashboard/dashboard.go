package dashboard

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func Dashboard() fyne.CanvasObject {

	content := container.NewBorder(
		container.NewVBox(
			StartProject(),
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

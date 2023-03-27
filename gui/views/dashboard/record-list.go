package dashboard

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/gui/state"
)

func formatTime(theTime time.Time) string {
	return fmt.Sprintf("%02d:%02d", theTime.Hour(), theTime.Minute())
}

func formatRecordLabel(record core.Record) string {
	tags := ""
	if len(record.Tags) > 0 {
		tags = fmt.Sprintf("(%s)", strings.Join(record.Tags, ", "))
	}
	return fmt.Sprintf(
		"%s - %s | %s %s",
		formatTime(record.Start),
		formatTime(*record.End),
		record.Project.Key,
		tags,
	)
}

func RecordList() fyne.CanvasObject {
	theState := state.GetState()
	list := widget.NewListWithData(
		theState.Records,
		func() fyne.CanvasObject {
			return widget.NewButton("template", func() {})
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			recordUntyped, _ := i.(binding.Untyped).Get()
			record := recordUntyped.(*core.Record)
			o.(*widget.Button).SetText(formatRecordLabel(*record))
			o.(*widget.Button).OnTapped = func() {
				theState.EditRecord(record)
			}
		},
	)
	return list
}

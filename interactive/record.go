package interactive

import (
	"errors"
	"fmt"
	"time"

	"github.com/dominikbraun/timetrace/core"
	"github.com/manifoldco/promptui"
)

func SelectRecord(t *core.Timetrace) (time.Time, error) {
	today, _ := t.Formatter().ParseDate("today")

	records, _ := t.ListRecords(today)
	if len(records) == 0 {
		return time.Time{}, errors.New("no records to edit")
	}

	items := make([]string, 0)
	for _, record := range records {
		if record.End == nil {
			continue
		}
		tags := ""
		if len(record.Tags) > 0 {
			tags = fmt.Sprintf("(%s)", t.Formatter().FormatTags(record.Tags))
		}
		items = append(items, fmt.Sprintf(
			"%s - %s | %s %s",
			t.Formatter().TimeString(record.Start),
			t.Formatter().TimeString(*record.End),
			record.Project.Key,
			tags,
		))
	}

	prompt := promptui.Select{
		Label: "Select record to edit",
		Items: items,
		Size:  10,
	}

	recordIndex, _, _ := prompt.Run()
	record, _ := t.LoadRecordByID(len(records) - recordIndex)
	recordTime := record.Start

	return recordTime, nil
}

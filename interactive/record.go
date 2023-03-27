package interactive

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dominikbraun/timetrace/core"
	"github.com/manifoldco/promptui"
)

func FormatRecordsAsListItems(records []*core.Record, timeStringFormatter func(time.Time) string) []string {
	items := make([]string, 0)
	for _, record := range records {
		if record.End == nil {
			continue
		}
		tags := ""
		if len(record.Tags) > 0 {
			tags = fmt.Sprintf("(%s)", strings.Join(record.Tags, ", "))
		}
		items = append(items, fmt.Sprintf(
			"%s - %s | %s %s",
			timeStringFormatter(record.Start),
			timeStringFormatter(*record.End),
			record.Project.Key,
			tags,
		))
	}

	return items
}

func SelectRecord(t *core.Timetrace) (time.Time, error) {
	today, err := t.Formatter().ParseDate("today")
	if err != nil {
		return time.Time{}, err
	}
	records, err := t.ListRecords(today)
	if err != nil {
		return time.Time{}, err
	}
	if len(records) == 0 {
		return time.Time{}, errors.New("no records to edit")
	}

	items := FormatRecordsAsListItems(records, t.Formatter().TimeString)

	prompt := promptui.Select{
		Label: "Select record to edit",
		Items: items,
		Size:  10,
	}

	recordIndex, _, _ := prompt.Run()
	record, _ := t.LoadRecordByID(len(items) - recordIndex)
	recordTime := record.Start

	return recordTime, nil
}

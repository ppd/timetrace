package state

import (
	"bytes"
	"time"

	"fyne.io/fyne/v2/data/binding"
	"github.com/dominikbraun/timetrace/core"
	"github.com/olekukonko/tablewriter"
)

type reportState struct {
	Report binding.String
}

var theReportState = &reportState{
	Report: binding.NewString(),
}

func ReportState() *reportState {
	return theReportState
}

func (s *reportState) UpdateReport(theDate time.Time) error {
	filters := []func(*core.Record) bool{
		core.FilterNoneNilEndTime,
		core.FilterByTimeRange(theDate, theDate),
	}

	report, err := Timetrace().Report(filters...)
	if err != nil {
		return err
	}

	projects, total := report.Table(core.TableOptions{ShowBillable: false, ShowDate: false})

	if len(projects) == 0 {
		s.Report.Set("Nothing to report for that day :)")
		return nil
	}

	buffer := &bytes.Buffer{}
	headers := []string{"Project", "Module", "Tags", "Start", "End", "Total"}
	table := tablewriter.NewWriter(buffer)
	table.SetHeader(headers)
	table.SetFooter([]string{"", "", "", "", "TOTAL", total})
	table.SetRowLine(true)
	table.SetAutoMergeCellsByColumnIndex([]int{0})
	table.AppendBulk(projects)
	table.Render()

	s.Report.Set(buffer.String())

	return nil
}

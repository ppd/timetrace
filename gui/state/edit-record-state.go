package state

import (
	"strings"
	"time"

	"fyne.io/fyne/v2/data/binding"
	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/gui/shared"
)

type editRecordState struct {
	Record           *core.Record
	Start            TimeBinding
	End              TimeBinding
	Tags             binding.String
	Project          binding.String
	IsExternalChange bool
}

func (s *editRecordState) EditRecord(record *core.Record) {
	s.IsExternalChange = true
	s.Record = record
	s.Start.Set(record.Start)
	s.End.Set(*record.End)
	s.Tags.Set(strings.Join(record.Tags, ", "))
	s.Project.Set(record.Project.Key)
	CoreState().ChangeView(EditProject)
}

var theEditRecordState = &editRecordState{
	Start:            NewBoundTimeWithData(time.Now()),
	End:              NewBoundTimeWithData(time.Now()),
	Tags:             binding.NewString(),
	Project:          binding.NewString(),
	IsExternalChange: false,
}

func EditRecordState() *editRecordState {
	return theEditRecordState
}

func (s *editRecordState) SaveRecordToEdit() {
	recordBefore := s.Record
	recordAfter := *recordBefore
	recordAfter.Start, _ = s.Start.Get()
	end, _ := s.End.Get()
	recordAfter.End = &end
	recordAfter.Project.Key, _ = s.Project.Get()
	tags, _ := s.Tags.Get()
	recordAfter.Tags = shared.SplitAndTrim(tags)
	if err := GetTimetrace().SaveRecord(recordAfter, true); err != nil {
		panic("uh oh")
	}
	if err := GetTimetrace().SyncRecordFilepath(*recordBefore, recordAfter); err != nil {
		panic("uh oh")
	}
	DashboardState().GoToDashboard()
}

func (s *editRecordState) DeleteRecordToEdit() {
	if err := GetTimetrace().DeleteRecord(*s.Record); err != nil {
		panic("uh oh")
	}
	DashboardState().GoToDashboard()
}

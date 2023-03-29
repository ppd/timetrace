package state

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2/data/binding"
	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/gui/shared"
	"github.com/dominikbraun/timetrace/gui/ui"
)

type View int

const (
	Main View = iota
	EditRecord
)

func (v View) String() string {
	return []string{"Main", "EditRecord"}[v]
}

type State struct {
	// core reference
	T *core.Timetrace
	// active view
	ActiveView binding.Int
	// main view
	Records        binding.UntypedList
	ProjectLabels  binding.StringList
	Status         binding.String
	IsRecordActive binding.Bool
	Tags           binding.String
	// edit
	RecordToEdit                 *core.Record
	RecordToEditStart            ui.TimeBinding
	RecordToEditEnd              ui.TimeBinding
	RecordToEditTags             binding.String
	RecordToEditProject          binding.String
	RecordToEditIsExternalChange bool
}

var theState State

func InitState(t *core.Timetrace) *State {
	theState = State{
		T:                   t,
		ActiveView:          binding.NewInt(),
		Records:             binding.NewUntypedList(),
		ProjectLabels:       binding.BindStringList(&[]string{}),
		Status:              binding.NewString(),
		IsRecordActive:      binding.NewBool(),
		Tags:                binding.NewString(),
		RecordToEditStart:   ui.NewBoundTimeWithData(time.Now()),
		RecordToEditEnd:     ui.NewBoundTimeWithData(time.Now()),
		RecordToEditTags:    binding.NewString(),
		RecordToEditProject: binding.NewString(),
	}
	theState.ActiveView.Set(int(Main))
	return &theState
}

func (s *State) ChangeView(view View) {
	s.ActiveView.Set(int(view))
}

func (s *State) UpdateRecords() {
	today, _ := s.T.Formatter().ParseDate("today")
	records, _ := s.T.ListRecords(today)
	recordsUntyped := make([]interface{}, 0)
	for _, record := range records {
		if record.End != nil {
			recordsUntyped = append(recordsUntyped, record)
		}
	}
	s.Records.Set(recordsUntyped)
}

func (s *State) UpdateProjectLabels() {
	s.ProjectLabels.Set(s.T.ListProjectNames())
}

func (s *State) UpdateStatus() {
	status, _ := s.T.Status()
	label := "No active project"
	workedToday := "Worked today: -"
	isActive := false
	tags := ""
	if status != nil {
		workedToday = fmt.Sprintf("Worked today: %s", s.T.Formatter().FormatDuration(status.TrackedTimeToday))
	}
	if status != nil && status.Current != nil {
		label = fmt.Sprintf(
			"%s in progress | Current project: %s",
			status.Current.Project.Key,
			s.T.Formatter().FormatDuration(status.Current.Duration()),
		)
		isActive = true
		tags = strings.Join(status.Current.Tags, ", ")
	}
	s.IsRecordActive.Set(isActive)
	s.Status.Set(fmt.Sprintf("%s | %s", label, workedToday))

	currentTags, _ := s.Tags.Get()
	if len(currentTags) == 0 {
		s.Tags.Set(tags)
	}
}

func (s *State) RefreshState() {
	s.UpdateRecords()
	s.UpdateStatus()
}

func (s *State) Stop() {
	s.StoreTags()
	s.T.Stop()
	s.RefreshState()
}

func (s *State) StoreTags() {
	record, err := s.T.LoadLatestRecord()
	if err != nil {
		panic("uh oh")
	}
	tags, _ := s.Tags.Get()
	record.Tags = shared.SplitAndTrim(tags)
	if err := s.T.SaveRecord(*record, true); err != nil {
		panic("uh oh")
	}
	s.Tags.Set("")
}

func (s *State) StartProject(projectKey string) error {
	if isActive, _ := s.IsRecordActive.Get(); isActive {
		return nil
	}
	if err := s.T.Start(projectKey, true, []string{}); err != nil {
		return err
	}
	s.UpdateStatus()
	return nil
}

func (s *State) EditRecord(record *core.Record) {
	s.RecordToEditIsExternalChange = true
	s.RecordToEdit = record
	s.RecordToEditStart.Set(record.Start)
	s.RecordToEditEnd.Set(*record.End)
	s.RecordToEditTags.Set(strings.Join(record.Tags, ", "))
	s.RecordToEditProject.Set(record.Project.Key)
	s.ActiveView.Set(int(EditRecord))
}

func (s *State) SaveRecordToEdit() {
	recordBefore := s.RecordToEdit
	recordAfter := *recordBefore
	recordAfter.Start, _ = s.RecordToEditStart.Get()
	end, _ := s.RecordToEditEnd.Get()
	recordAfter.End = &end
	recordAfter.Project.Key, _ = s.RecordToEditProject.Get()
	tags, _ := s.RecordToEditTags.Get()
	recordAfter.Tags = shared.SplitAndTrim(tags)
	if err := s.T.SaveRecord(recordAfter, true); err != nil {
		panic("uh oh")
	}
	if err := s.T.SyncRecordFilepath(*recordBefore, recordAfter); err != nil {
		panic("uh oh")
	}
	s.GoToMainView()
}

func (s *State) DeleteRecordToEdit() {
	if err := s.T.DeleteRecord(*s.RecordToEdit); err != nil {
		panic("uh oh")
	}
	s.GoToMainView()
}

func (s *State) GoToMainView() {
	s.RefreshState()
	s.ChangeView(Main)
}

func (s *State) RefreshStatePeriodically() func() {
	killMe := false

	go func() {
		for {
			if killMe {
				return
			}
			time.Sleep(time.Second * time.Duration(10))
			s.RefreshState()
		}
	}()

	return func() {
		killMe = true
	}
}

func GetState() *State {
	return &theState
}

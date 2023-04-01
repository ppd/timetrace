package state

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/gui/shared"
	"github.com/dominikbraun/timetrace/gui/ui"
)

type View int

const (
	Main View = iota
	EditRecord
	Projects
	EditProject
	About
)

func (v View) String() string {
	return []string{"Main", "EditRecord", "Projects", "EditProject", "About"}[v]
}

type State struct {
	// core reference
	T *core.Timetrace
	// main window
	MainWindow fyne.Window
	// active view
	ActiveView binding.Int
	// main view
	Records        binding.UntypedList
	ProjectLabels  binding.StringList
	Status         binding.String
	IsRecordActive binding.Bool
	Tags           binding.String
	// edit record
	RecordToEdit                 *core.Record
	RecordToEditStart            ui.TimeBinding
	RecordToEditEnd              ui.TimeBinding
	RecordToEditTags             binding.String
	RecordToEditProject          binding.String
	RecordToEditIsExternalChange bool
	// projects
	ProjectsFilter                 binding.String
	ProjectsFilterIsExternalChange bool
	FilteredProjects               binding.StringList
	// edit project
	ProjectToEdit               *core.Project
	ProjectToEditKey            binding.String
	ProjectToEditChronosProject binding.String
	ProjectToEditChronosAccount binding.String
	// refresh trigger - listener is spawned on the main (ui) thread
	TriggerRefresh binding.Untyped
}

var theState State

func InitState(t *core.Timetrace) *State {
	theState = State{
		T:                              t,
		ActiveView:                     binding.NewInt(),
		Records:                        binding.NewUntypedList(),
		ProjectLabels:                  binding.BindStringList(&[]string{}),
		Status:                         binding.NewString(),
		IsRecordActive:                 binding.NewBool(),
		Tags:                           binding.NewString(),
		RecordToEditStart:              ui.NewBoundTimeWithData(time.Now()),
		RecordToEditEnd:                ui.NewBoundTimeWithData(time.Now()),
		RecordToEditTags:               binding.NewString(),
		RecordToEditProject:            binding.NewString(),
		RecordToEditIsExternalChange:   false,
		ProjectsFilter:                 binding.NewString(),
		ProjectsFilterIsExternalChange: false,
		FilteredProjects:               binding.NewStringList(),
		ProjectToEditKey:               binding.NewString(),
		ProjectToEditChronosProject:    binding.NewString(),
		ProjectToEditChronosAccount:    binding.NewString(),
		TriggerRefresh:                 binding.NewUntyped(),
	}
	theState.ActiveView.Set(int(Main))
	theState.TriggerRefresh.AddListener(binding.NewDataListener(func() {
		theState.RefreshState()
	}))
	theState.ProjectsFilter.AddListener(binding.NewDataListener(func() {
		theState.FilterProjectList()
	}))
	theState.ProjectLabels.AddListener(binding.NewDataListener(func() {
		theState.FilterProjectList()
	}))
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

func (s *State) UpdateProjects() {
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

func (s *State) CreateProject(projectKey string) error {
	project := core.Project{
		Key: projectKey,
	}
	return s.T.SaveProject(project, false)
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
	s.ChangeView(EditRecord)
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

func (s *State) EditProject(projectKey string) error {
	project, err := s.T.LoadProject(projectKey)
	if err != nil {
		return err
	}
	s.ProjectToEdit = project
	s.ProjectToEditKey.Set(project.Key)
	s.ProjectToEditChronosProject.Set(project.ChronosProject)
	s.ProjectToEditChronosAccount.Set(project.ChronosAccount)
	s.ChangeView(EditProject)
	return nil
}

func (s *State) DeleteProjectToEdit() error {
	if err := s.T.DeleteProject(*s.ProjectToEdit); err != nil {
		return err
	}
	s.GoToProjectsView()
	return nil
}

func (s *State) SaveProjectToEdit() error {
	key, _ := s.ProjectToEditKey.Get()
	chronosAccount, _ := s.ProjectToEditChronosAccount.Get()
	chronosProject, _ := s.ProjectToEditChronosProject.Get()
	project := core.Project{
		Key:            key,
		ChronosProject: chronosProject,
		ChronosAccount: chronosAccount,
	}
	if err := s.T.SaveProject(project, true); err != nil {
		return err
	}
	s.GoToProjectsView()
	return nil
}

func (s *State) DeleteRecordToEdit() {
	if err := s.T.DeleteRecord(*s.RecordToEdit); err != nil {
		panic("uh oh")
	}
	s.GoToMainView()
}

func (s *State) GoToProjectsView() {
	s.UpdateProjects()
	s.ChangeView(Projects)
}

func (s *State) GoToMainView() {
	s.RefreshState()
	s.ChangeView(Main)
}

func (s *State) FilterProjectList() {
	filter, _ := theState.ProjectsFilter.Get()
	labels, _ := s.ProjectLabels.Get()
	s.FilteredProjects.Set(shared.FilterByContains(labels, filter))
}

func (s *State) RefreshStatePeriodically() (chan<- bool, <-chan bool) {
	stop := make(chan bool)
	done := make(chan bool)
	i := 0

	go func() {
		for {
			killMe := false
			time.Sleep(time.Second * time.Duration(10))
			s.TriggerRefresh.Set(i)
			select {
			case killMe = <-stop:
			default:
				i++
			}
			if killMe {
				done <- true
				return
			}
		}
	}()

	return stop, done
}

func GetState() *State {
	return &theState
}

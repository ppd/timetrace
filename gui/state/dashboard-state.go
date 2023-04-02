package state

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2/data/binding"
	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/gui/shared"
)

type dashboardState struct {
	// main view
	Date           TimeBinding
	Records        binding.UntypedList
	Status         binding.String
	IsRecordActive binding.Bool
	Tags           binding.String
	// refresh trigger - listener is spawned on the main (ui) thread
	TriggerRefresh binding.Untyped
}

var theDashboardState *dashboardState

func InitDashboardState() *dashboardState {
	theDashboardState = &dashboardState{
		Date:           NewBoundTimeWithData(Today()),
		Records:        binding.NewUntypedList(),
		Status:         binding.NewString(),
		IsRecordActive: binding.NewBool(),
		Tags:           binding.NewString(),
		TriggerRefresh: binding.NewUntyped(),
	}

	theDashboardState.TriggerRefresh.AddListener(binding.NewDataListener(func() {
		theDashboardState.RefreshState()
	}))

	theDashboardState.Date.AddListener(binding.NewDataListener(func() {
		theDashboardState.UpdateRecords()
	}))

	return theDashboardState
}

func DashboardState() *dashboardState {
	if theDashboardState == nil {
		InitDashboardState()
	}
	return theDashboardState
}

func (s *dashboardState) UpdateRecords() {
	theDate, _ := s.Date.Get()
	records, _ := Timetrace().ListRecords(theDate)
	recordsUntyped := make([]interface{}, 0)
	for _, record := range records {
		if record.End != nil {
			recordsUntyped = append(recordsUntyped, record)
		}
	}
	s.Records.Set(recordsUntyped)
}

func (s *dashboardState) UpdateStatus() {
	status, _ := Timetrace().Status()
	label := "No active project"
	workedToday := "Worked today: -"
	isActive := false
	tags := ""
	if status != nil {
		workedToday = fmt.Sprintf("Worked today: %s", Timetrace().Formatter().FormatDuration(status.TrackedTimeToday))
	}
	if status != nil && status.Current != nil {
		label = fmt.Sprintf(
			"%s in progress | Current project: %s",
			status.Current.Project.Key,
			Timetrace().Formatter().FormatDuration(status.Current.Duration()),
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
func (s *dashboardState) Stop() {
	s.StoreTags()
	Timetrace().Stop()
	s.RefreshState()
}

func (s *dashboardState) StoreTags() {
	record, err := Timetrace().LoadLatestRecord()
	if err != nil {
		panic("uh oh")
	}
	tags, _ := s.Tags.Get()
	record.Tags = shared.SplitAndTrim(tags)
	if err := Timetrace().SaveRecord(*record, true); err != nil {
		panic("uh oh")
	}
	s.Tags.Set("")
}

func (s *dashboardState) CreateProject(projectKey string) error {
	project := core.Project{
		Key: projectKey,
	}
	return Timetrace().SaveProject(project, false)
}

func (s *dashboardState) StartProject(projectKey string) error {
	if isActive, _ := s.IsRecordActive.Get(); isActive {
		return nil
	}
	if err := Timetrace().Start(projectKey, true, []string{}); err != nil {
		return err
	}
	s.UpdateStatus()
	return nil
}

func (s *dashboardState) RefreshState() {
	s.UpdateRecords()
	s.UpdateStatus()
}

func (s *dashboardState) GoToDashboard() {
	s.RefreshState()
	CoreState().ChangeView(Main)
}

func (s *dashboardState) RefreshStatePeriodically() (chan<- bool, <-chan bool) {
	stop := make(chan bool)
	done := make(chan bool)
	i := 0

	go func() {
		for {
			killMe := false
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
			time.Sleep(time.Second * time.Duration(10))
		}
	}()

	return stop, done
}

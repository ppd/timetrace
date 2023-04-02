package state

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"github.com/dominikbraun/timetrace/core"
)

type View int

const (
	Main View = iota
	EditRecord
	Projects
	EditProject
	About
	Report
)

func (v View) String() string {
	return []string{"Main", "EditRecord", "Projects", "EditProject", "About", "Report"}[v]
}

type coreState struct {
	// core reference
	T *core.Timetrace
	// main window
	MainWindow fyne.Window
	// active view
	ActiveView binding.Int
	// entities
	ProjectLabels binding.StringList
}

var theCoreState = &coreState{
	ActiveView:    binding.NewInt(),
	ProjectLabels: binding.NewStringList(),
}

func (s *coreState) ChangeView(view View) {
	s.ActiveView.Set(int(view))
}

func GetToday() time.Time {
	today, _ := GetTimetrace().Formatter().ParseDate("today")
	return today
}

func GetTimetrace() *core.Timetrace {
	return CoreState().T
}

func InitCoreState(t *core.Timetrace) *coreState {
	theCoreState.T = t
	return theCoreState
}

func CoreState() *coreState {
	return theCoreState
}

func (s *coreState) UpdateProjects() {
	s.ProjectLabels.Set(GetTimetrace().ListProjectNames())
}

func (s *coreState) GoToProjectsView() {
	s.UpdateProjects()
	s.ChangeView(Projects)
}

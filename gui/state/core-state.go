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
	None
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
	// requested view
	RequestedView binding.Int
	// entities
	ProjectLabels binding.StringList
}

var theCoreState = &coreState{
	ActiveView:    binding.NewInt(),
	RequestedView: binding.NewInt(),
	ProjectLabels: binding.NewStringList(),
}

func (s *coreState) ChangeView(view View) {
	s.RequestedView.Set(int(view))
}

func (s *coreState) SetActiveView(view View) {
	s.ActiveView.Set(int(view))
}

func Today() time.Time {
	today, _ := Timetrace().Formatter().ParseDate("today")
	return today
}

func Timetrace() *core.Timetrace {
	return CoreState().T
}

func InitCoreState(t *core.Timetrace) *coreState {
	theCoreState.T = t
	theCoreState.SetActiveView(None)
	return theCoreState
}

func CoreState() *coreState {
	return theCoreState
}

func (s *coreState) UpdateProjects() {
	s.ProjectLabels.Set(Timetrace().ListProjectNames())
}

func (s *coreState) GoToProjectsView() {
	s.UpdateProjects()
	s.ChangeView(Projects)
}

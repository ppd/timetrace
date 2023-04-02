package state

import (
	"fyne.io/fyne/v2/data/binding"
	"github.com/dominikbraun/timetrace/gui/shared"
)

type projectsState struct {
	// projects
	ProjectsFilter                 binding.String
	ProjectsFilterIsExternalChange bool
	FilteredProjects               binding.StringList
}

var theProjectsState *projectsState

func InitProjectsState() {
	theProjectsState = &projectsState{
		ProjectsFilter:                 binding.NewString(),
		ProjectsFilterIsExternalChange: false,
		FilteredProjects:               binding.NewStringList(),
	}

	theProjectsState.ProjectsFilter.AddListener(binding.NewDataListener(func() {
		theProjectsState.FilterProjectList()
	}))

	CoreState().ProjectLabels.AddListener(binding.NewDataListener(func() {
		theProjectsState.FilterProjectList()
	}))
}

func ProjectsState() *projectsState {
	if theProjectsState == nil {
		InitProjectsState()
	}
	return theProjectsState
}

func (s *projectsState) FilterProjectList() {
	filter, _ := s.ProjectsFilter.Get()
	labels, _ := CoreState().ProjectLabels.Get()
	s.FilteredProjects.Set(shared.FilterByContains(labels, filter))
}

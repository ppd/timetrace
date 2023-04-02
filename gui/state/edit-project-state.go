package state

import (
	"fyne.io/fyne/v2/data/binding"
	"github.com/dominikbraun/timetrace/core"
)

type editProjectState struct {
	Project        *core.Project
	Key            binding.String
	ChronosProject binding.String
	ChronosAccount binding.String
}

var theEditProjectState = &editProjectState{
	Key:            binding.NewString(),
	ChronosProject: binding.NewString(),
	ChronosAccount: binding.NewString(),
}

func EditProjectState() *editProjectState {
	return theEditProjectState
}

func (s *editProjectState) DoEdit(projectKey string) error {
	project, err := GetTimetrace().LoadProject(projectKey)
	if err != nil {
		return err
	}
	s.Project = project
	s.Key.Set(project.Key)
	s.ChronosProject.Set(project.ChronosProject)
	s.ChronosAccount.Set(project.ChronosAccount)
	CoreState().ChangeView(EditProject)
	return nil
}

func (s *editProjectState) DeleteProjectToEdit() error {
	if err := GetTimetrace().DeleteProject(*s.Project); err != nil {
		return err
	}
	CoreState().GoToProjectsView()
	return nil
}

func (s *editProjectState) SaveProjectToEdit() error {
	key, _ := s.Key.Get()
	chronosAccount, _ := s.ChronosAccount.Get()
	chronosProject, _ := s.ChronosProject.Get()
	project := core.Project{
		Key:            key,
		ChronosProject: chronosProject,
		ChronosAccount: chronosAccount,
	}
	if err := GetTimetrace().SaveProject(project, true); err != nil {
		return err
	}
	CoreState().GoToProjectsView()
	return nil
}

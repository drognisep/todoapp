package view

import (
	"fmt"
	"regexp"
	"sync"
	"time"
	"todo/data"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/drognisep/fynehelpers/layouthelp"
	"github.com/pkg/errors"
)

type TaskListSelector struct {
	widget.Select

	mux        sync.RWMutex
	multiview  *layouthelp.MultiView
	viewMap    map[string]*TaskList
	activeList *TaskList
	model      *data.Model
}

func (s *TaskListSelector) defaultView() fyne.CanvasObject {
	label := widget.NewLabel("Select or create a project to get started")
	label.Alignment = fyne.TextAlignCenter
	return container.NewVBox(layout.NewSpacer(), label, layout.NewSpacer())
}

func NewTaskListSelector(ctx *UiCtx, model *data.Model) *TaskListSelector {
	selector := &TaskListSelector{
		multiview: layouthelp.NewMultiView(),
		viewMap:   map[string]*TaskList{},
		model:     model,
	}
	selector.Select = widget.Select{
		Options:     []string{},
		PlaceHolder: "Select a Project",
		OnChanged:   selector.onChanged,
	}

	selector.multiview.Push(selector.defaultView())
	for _, p := range model.Projects {
		tl := NewTaskList(ctx, &p.Tasks)
		selector.AddTaskList(p.Name, tl)
	}

	time.AfterFunc(300*time.Millisecond, func() {
		lastOpened := model.LastProjectOpened
		_, ok := selector.viewMap[lastOpened]
		if lastOpened != "" && ok {
			selector.SetSelected(lastOpened)
		}
	})
	selector.ExtendBaseWidget(selector)
	return selector
}

func (s *TaskListSelector) onChanged(viewKey string) {
	s.mux.RLock()
	selectedView := s.viewMap[viewKey]
	s.mux.RUnlock()
	if viewKey == "" || selectedView == nil {
		s.mux.Lock()
		s.multiview.Replace(s.defaultView())
		s.activeList = nil
		s.model.LastProjectOpened = viewKey
		s.mux.Unlock()
		return
	}
	s.mux.Lock()
	s.multiview.Replace(selectedView)
	s.activeList = selectedView
	s.model.LastProjectOpened = viewKey
	s.mux.Unlock()
}

func (s *TaskListSelector) AddTaskList(key string, list *TaskList) {
	if key != "" {
		defer s.Refresh()
		s.mux.Lock()
		defer s.mux.Unlock()
		s.viewMap[key] = list
		s.Options = append(s.Options, key)
	}
}

func (s *TaskListSelector) RemoveTaskList(key string) {
	if key != "" {
		defer s.Refresh()
		s.mux.Lock()
		defer s.mux.Unlock()
		for k, v := range s.viewMap {
			if v == s.activeList {
				if k == key {
					s.multiview.Replace(s.defaultView())
					s.Selected = ""
				}
				break
			}
		}
		delete(s.viewMap, key)
		for i, opt := range s.Options {
			if opt == key {
				s.Options = append(s.Options[:i], s.Options[i+1:]...)
				break
			}
		}
	}
}

func (s *TaskListSelector) ActiveList() *TaskList {
	s.mux.RLock()
	active := s.activeList
	s.mux.RUnlock()
	return active
}

func (s *TaskListSelector) ActiveProject() string {
	s.mux.RLock()
	selected := s.Selected
	s.mux.RUnlock()
	return selected
}

func (s *TaskListSelector) ViewContainer() *fyne.Container {
	return s.multiview.Container()
}

var (
	labelRegex     = regexp.MustCompile(`^[A-Za-z0-9-_]+( [A-Za-z0-9-_]+)*$`)
	labelValidator = func(model *data.Model) func(string) error {
		return func(val string) error {
			if !labelRegex.MatchString(val) {
				return errors.New("Must be an alphanumeric string with words separated by only 1 space")
			}
			if model.HasProjectName(val) {
				return errors.New("Project name already exists")
			}
			return nil
		}
	}
)

func (s *TaskListSelector) ProjectMenu(ctx *UiCtx) *fyne.Menu {
	newProject := fyne.NewMenuItem("New", func() {
		project := &data.Project{
			Tasks: []*data.Task{},
		}
		nameEntry := widget.NewEntryWithData(binding.BindString(&project.Name))
		nameEntry.Validator = labelValidator(s.model)
		nameItem := widget.NewFormItem("Project Name", nameEntry)
		d := dialog.NewForm("New Project", "Create", "Cancel", []*widget.FormItem{nameItem}, func(doCreate bool) {
			if doCreate {
				s.model.AddProject(project)
				s.AddTaskList(project.Name, NewTaskList(ctx, &project.Tasks))
				s.SetSelected(project.Name)
			}
		}, ctx.MainWindow)
		d.Resize(fyne.NewSize(BigFloat, BigFloat))
		d.Show()
	})
	deleteProject := fyne.NewMenuItem("Delete", func() {
		activeProject := s.ActiveProject()
		message := fmt.Sprintf("Are you sure you want to delete '%s'?", activeProject)
		dialog.ShowConfirm("Are you sure?", message, func(confirmed bool) {
			if confirmed {
				s.model.RemoveProject(activeProject)
				s.RemoveTaskList(activeProject)
			}
		}, ctx.MainWindow)
	})
	renameProject := fyne.NewMenuItem("Rename", func() {
		activeProject := s.ActiveProject()
		activeList := s.ActiveList()
		newName := activeProject
		newNameEntry := widget.NewEntryWithData(binding.BindString(&newName))
		newNameEntry.Validator = labelValidator(s.model)
		newNameItem := widget.NewFormItem("New Name", newNameEntry)
		d := dialog.NewForm(fmt.Sprintf("Rename '%s'", activeProject), "Rename", "Cancel", []*widget.FormItem{newNameItem}, func(doRename bool) {
			if doRename {
				s.model.RenameProject(activeProject, newName)
				s.RemoveTaskList(activeProject)
				s.AddTaskList(newName, activeList)
				s.SetSelected(newName)
			}
		}, ctx.MainWindow)
		d.Resize(fyne.NewSize(BigFloat, BigFloat))
		d.Show()
	})
	return fyne.NewMenu("Project", newProject, deleteProject, renameProject)
}

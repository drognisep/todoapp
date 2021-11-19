package view

import (
	"sync"
	"time"
	"todo/data"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/drognisep/fynehelpers/layouthelp"
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

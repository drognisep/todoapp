package view

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/drognisep/fynehelpers/layouthelp"
)

type TaskListSelector struct {
	widget.Select

	mux         sync.RWMutex
	multiview   *layouthelp.MultiView
	viewMap     map[string]*TaskList
	activeList  *TaskList
	defaultView fyne.CanvasObject
}

func NewTaskListSelector(defaultView fyne.CanvasObject) *TaskListSelector {
	selector := &TaskListSelector{
		multiview:   layouthelp.NewMultiView(),
		viewMap:     map[string]*TaskList{},
		defaultView: defaultView,
	}
	selector.Select = widget.Select{
		Options:     []string{},
		PlaceHolder: "Select a Project",
		OnChanged: func(viewKey string) {
			selector.mux.RLock()
			selectedView := selector.viewMap[viewKey]
			selector.mux.RUnlock()
			if viewKey == "" || selectedView == nil {
				selector.mux.Lock()
				selector.multiview.Replace(defaultView)
				selector.activeList = nil
				selector.mux.Unlock()
				return
			}
			selector.mux.Lock()
			selector.multiview.Replace(selectedView)
			selector.activeList = selectedView
			selector.mux.Unlock()
		},
	}
	selector.multiview.Push(defaultView)
	selector.ExtendBaseWidget(selector)
	return selector
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
					s.multiview.Replace(s.defaultView)
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

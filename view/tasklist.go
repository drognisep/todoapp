package view

import (
	"sync"
	"todo/data"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type TaskList struct {
	widget.List

	mux   sync.RWMutex
	tasks []*data.Task
}

func NewTaskList(w fyne.Window, tasks []*data.Task) *TaskList {
	taskList := &TaskList{}
	taskList.List = widget.List{
		Length: func() int {
			return len(taskList.tasks)
		},
		CreateItem: func() fyne.CanvasObject {
			return PrototypeTaskView(w, taskList)
		},
		UpdateItem: func(id widget.ListItemID, item fyne.CanvasObject) {
			taskView := item.(*TaskView)
			taskView.update(taskList.tasks[id], func() {
				taskList.Delete(id)
			})
		},
	}
	for _, t := range tasks {
		taskList.Append(t)
	}
	taskList.ExtendBaseWidget(taskList)
	return taskList
}

func (l *TaskList) Delete(id widget.ListItemID) {
	l.mux.Lock()
	curTasks := l.tasks
	l.tasks = append(curTasks[:id], curTasks[id+1:]...)
	l.mux.Unlock()
	l.Refresh()
}

func (l *TaskList) Append(task *data.Task) {
	if task == nil {
		return
	}
	l.mux.Lock()
	l.tasks = append(l.tasks, task)
	l.mux.Unlock()
	l.Refresh()
}

func (l *TaskList) Tasks() []*data.Task {
	l.mux.RLock()
	defer l.mux.RUnlock()
	list := make([]*data.Task, len(l.tasks))
	copy(list, l.tasks)
	return list
}

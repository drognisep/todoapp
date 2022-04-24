package view

import (
	"sync"
	"todo/data"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type TaskList struct {
	widget.List

	mux      sync.RWMutex
	tasksPtr *[]*data.Task
	tasks    []*data.Task
	dragFrom widget.ListItemID
	dragTo   widget.ListItemID
	dragging bool
}

func NewTaskList(ctx *UiCtx, tasks *[]*data.Task) *TaskList {
	taskList := &TaskList{
		tasksPtr: tasks,
	}
	taskList.List = widget.List{
		Length: func() int {
			return len(taskList.tasks)
		},
		CreateItem: func() fyne.CanvasObject {
			return PrototypeTaskView(ctx, taskList)
		},
		UpdateItem: func(id widget.ListItemID, item fyne.CanvasObject) {
			taskView := item.(*TaskView)
			taskView.update(id, taskList.tasks[id], func() {
				taskList.Delete(id)
			})
		},
	}
	for _, t := range *tasks {
		taskList.Append(t)
	}
	taskList.ExtendBaseWidget(taskList)
	return taskList
}

func (l *TaskList) Delete(id widget.ListItemID) {
	l.mux.Lock()
	l.delete(id)
	l.mux.Unlock()
	l.Refresh()
}

func (l *TaskList) delete(id widget.ListItemID) {
	curTasks := l.tasks
	l.tasks = append(curTasks[:id], curTasks[id+1:]...)
	*l.tasksPtr = l.tasks
}

func (l *TaskList) Append(task *data.Task) {
	if task == nil {
		return
	}
	l.mux.Lock()
	l.append(task)
	l.mux.Unlock()
	l.Refresh()
}

func (l *TaskList) append(task *data.Task) {
	l.tasks = append(l.tasks, task)
	*l.tasksPtr = l.tasks
}

func (l *TaskList) Tasks() []*data.Task {
	l.mux.RLock()
	defer l.mux.RUnlock()
	list := make([]*data.Task, len(l.tasks))
	copy(list, l.tasks)
	return list
}

func (l *TaskList) DragFrom(id widget.ListItemID) {
	l.mux.Lock()
	l.dragging = true
	l.dragFrom = id
	l.dragTo = id
	l.mux.Unlock()
}

func (l *TaskList) DragTo(id widget.ListItemID) {
	l.mux.RLock()
	dragging := l.dragging
	dragToID := l.dragTo
	l.mux.RUnlock()
	if !dragging || dragToID == id {
		return
	}
	l.mux.Lock()
	l.dragTo = id
	l.mux.Unlock()
}

func (l *TaskList) DragEnd() {
	defer l.Refresh()
	l.mux.Lock()
	l.dragging = false
	dragFrom, dragTo := l.dragFrom, l.dragTo
	l.dragFrom, l.dragTo = -1, -1
	l.reorderNode(dragFrom, dragTo)
	l.mux.Unlock()
	l.Refresh()
}

func (l *TaskList) reorderNode(fromID, toID widget.ListItemID) {
	if fromID == toID {
		return
	}
	dragTask := l.tasks[fromID]
	l.delete(fromID)
	l.tasks = append(l.tasks[:toID], append([]*data.Task{dragTask}, l.tasks[toID:]...)...)
	*l.tasksPtr = l.tasks
}

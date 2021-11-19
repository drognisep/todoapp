package view

import (
	"sync"
	"todo/data"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type TaskList struct {
	widget.List

	mux      sync.RWMutex
	tasksPtr *[]*data.Task
	tasks    []*data.Task
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
			taskView.update(taskList.tasks[id], func() {
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
	curTasks := l.tasks
	l.tasks = append(curTasks[:id], curTasks[id+1:]...)
	*l.tasksPtr = l.tasks
	l.mux.Unlock()
	l.Refresh()
}

func (l *TaskList) Append(task *data.Task) {
	if task == nil {
		return
	}
	l.mux.Lock()
	l.tasks = append(l.tasks, task)
	*l.tasksPtr = l.tasks
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

func TaskPanel(ctx *UiCtx, model *data.Model) *fyne.Container {
	projectSelect := NewTaskListSelector(ctx, model)

	addTaskBtn := widget.NewButtonWithIcon("New Task", theme.ContentAddIcon(), func() {
		task := data.Task{}
		nameEntry := widget.NewEntryWithData(binding.BindString(&task.Name))
		nameEntry.Validator = TaskNameValidator
		d := dialog.NewForm("Add Task", "Add", "Cancel", []*widget.FormItem{widget.NewFormItem("Name", nameEntry)}, func(confirm bool) {
			if confirm {
				projectSelect.ActiveList().Append(&task)
			}
		}, ctx.MainWindow)
		d.Resize(fyne.NewSize(BigFloat, BigFloat))
		d.Show()
	})
	addTaskBtn.Importance = widget.HighImportance
	saveBtn := widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), func() {
		ctx.SaveData(model)
	})
	content := container.NewBorder(container.NewHBox(widget.NewLabel("Project"), projectSelect, addTaskBtn, saveBtn), nil, nil, nil, projectSelect.ViewContainer())
	return content
}

package main

import (
	"fmt"
	"todo/data"
	"todo/view"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	osdialog "github.com/sqweek/dialog"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			osdialog.Message("A fatal error occurred: %v", r).Title("Unexpected error").Error()
		}
	}()
	tasks, err := data.LoadTaskData()
	if err != nil {
		panic(fmt.Errorf("Failed to load task data: %v\n", err))
	}

	a := app.NewWithID("com.saylorsolutions.todoapp")
	w := a.NewWindow("TODO")
	w.Resize(fyne.NewSize(640, 480))
	w.CenterOnScreen()
	ctx := &view.UiCtx{
		App:        a,
		MainWindow: w,
	}

	taskList := view.NewTaskList(ctx, tasks)
	addBtn := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		task := data.Task{}
		nameEntry := widget.NewEntryWithData(binding.BindString(&task.Name))
		nameEntry.Validator = view.TaskNameValidator
		d := dialog.NewForm("Add Task", "Add", "Cancel", []*widget.FormItem{widget.NewFormItem("Name", nameEntry)}, func(confirm bool) {
			if confirm {
				taskList.Append(&task)
			}
		}, w)
		d.Resize(fyne.NewSize(view.BigFloat, view.BigFloat))
		d.Show()
	})
	addBtn.Importance = widget.HighImportance
	saveBtn := widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), func() {
		saveData(taskList, w)
	})
	content := container.NewBorder(container.NewHBox(addBtn, saveBtn), nil, nil, nil, taskList)

	ctrlS := desktop.CustomShortcut{
		KeyName:  fyne.KeyS,
		Modifier: desktop.ControlModifier,
	}
	w.Canvas().AddShortcut(&ctrlS, func(_ fyne.Shortcut) {
		saveData(taskList, w)
	})

	w.SetContent(content)
	w.SetCloseIntercept(func() {
		dialog.ShowConfirm("Save?", "Do you want to save before quitting?", func(doSave bool) {
			if doSave {
				saveData(taskList, w)
			}
			w.Close()
		}, w)
	})

	w.ShowAndRun()
}

func saveData(taskList *view.TaskList, w fyne.Window) {
	tasks := taskList.Tasks()
	if err := data.SaveTaskData(tasks); err != nil {
		dialog.ShowError(err, w)
	}
}

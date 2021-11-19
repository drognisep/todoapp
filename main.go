package main

import (
	"fmt"
	"runtime/debug"
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
			debug.PrintStack()
			osdialog.Message("A fatal error occurred: %v", r).Title("Unexpected error").Error()
		}
	}()
	model, err := data.LoadTaskData()
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

	ctrlS := desktop.CustomShortcut{
		KeyName:  fyne.KeyS,
		Modifier: desktop.ControlModifier,
	}
	w.Canvas().AddShortcut(&ctrlS, func(_ fyne.Shortcut) {
		saveData(ctx, model)
	})

	w.SetContent(taskPanel(ctx, model))
	w.SetOnClosed(func() {
		saveData(ctx, model)
	})
	w.ShowAndRun()
}

func taskPanel(ctx *view.UiCtx, model *data.Model) *fyne.Container {
	projectSelect := view.NewTaskListSelector(ctx, model)

	addTaskBtn := widget.NewButtonWithIcon("New Task", theme.ContentAddIcon(), func() {
		task := data.Task{}
		nameEntry := widget.NewEntryWithData(binding.BindString(&task.Name))
		nameEntry.Validator = view.TaskNameValidator
		d := dialog.NewForm("Add Task", "Add", "Cancel", []*widget.FormItem{widget.NewFormItem("Name", nameEntry)}, func(confirm bool) {
			if confirm {
				projectSelect.ActiveList().Append(&task)
			}
		}, ctx.MainWindow)
		d.Resize(fyne.NewSize(view.BigFloat, view.BigFloat))
		d.Show()
	})
	addTaskBtn.Importance = widget.HighImportance
	saveBtn := widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), func() {
		saveData(ctx, model)
	})
	content := container.NewBorder(container.NewHBox(widget.NewLabel("Project"), projectSelect, addTaskBtn, saveBtn), nil, nil, nil, projectSelect.ViewContainer())
	return content
}

func saveData(ctx *view.UiCtx, model *data.Model) {
	lbl := widget.NewLabel("Saving...")
	inf := widget.NewProgressBarInfinite()
	content := container.NewVBox(lbl, inf)
	popup := widget.NewModalPopUp(content, ctx.MainWindow.Canvas())
	popup.Show()
	defer popup.Hide()
	if err := data.SaveTaskData(model); err != nil {
		dialog.ShowError(err, ctx.MainWindow)
	}
}

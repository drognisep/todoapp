package view

import (
	"todo/data"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func TaskPanel(ctx *UiCtx, model *data.Model) *fyne.Container {
	projectSelect := NewTaskListSelector(ctx, model)
	ctx.MainWindow.SetMainMenu(fyne.NewMainMenu(projectSelect.ProjectMenu(ctx)))
	top := topControls(ctx, model, projectSelect)
	content := container.NewBorder(top, nil, nil, nil, projectSelect.ViewContainer())
	return content
}

func topControls(ctx *UiCtx, model *data.Model, projectSelect *TaskListSelector) *fyne.Container {
	addTaskBtn := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
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
	return container.NewHBox(projectSelect, addTaskBtn, saveBtn)
}
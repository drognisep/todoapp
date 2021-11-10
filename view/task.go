package view

import (
	"sync"
	"todo/data"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var _ fyne.Widget = (*TaskView)(nil)

type TaskView struct {
	widget.BaseWidget
	window fyne.Window
	parent fyne.CanvasObject

	mux      sync.Mutex
	renderer *treeViewRenderer
}

func PrototypeTaskView(w fyne.Window, parent fyne.CanvasObject) *TaskView {
	view := &TaskView{
		window: w,
		parent: parent,
	}
	view.ExtendBaseWidget(view)
	return view
}

func (t *TaskView) update(task *data.Task, onDelete func()) {
	t.mux.Lock()
	nameBinding := binding.BindString(&task.Name)
	if t.renderer != nil {
		t.renderer.check.Bind(binding.BindBool(&task.Done))
		t.renderer.label.Bind(nameBinding)
		t.renderer.deleteBtn.OnTapped = onDelete
	}
	t.renderer.label.OnDoubleTap = func() {
		newName := task.Name
		entry := widget.NewEntryWithData(binding.BindString(&newName))
		entryItem := widget.NewFormItem("Name", entry)
		d := dialog.NewForm("Edit Task", "Done", "Cancel", []*widget.FormItem{entryItem}, func(confirmed bool) {
			if confirmed {
				_ = nameBinding.Set(newName)
				t.Refresh()
				t.parent.Refresh()
			}
		}, t.window)
		d.Resize(fyne.NewSize(BigFloat, BigFloat))
		d.Show()
	}
	t.mux.Unlock()
	t.Refresh()
}

func (t *TaskView) CreateRenderer() fyne.WidgetRenderer {
	renderer := &treeViewRenderer{
		check:     widget.NewCheck("", func(_ bool) {}),
		label:     NewTappableLabel("Task description with enough space to see"),
		deleteBtn: widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {}),
	}
	renderer.objects = []fyne.CanvasObject{renderer.check, renderer.label, renderer.deleteBtn}
	t.mux.Lock()
	t.renderer = renderer
	t.mux.Unlock()
	return renderer
}

type treeViewRenderer struct {
	check     *widget.Check
	label     *TappableLabel
	deleteBtn *widget.Button
	objects []fyne.CanvasObject
}

func (t *treeViewRenderer) Destroy() {
	t.check = nil
	t.label = nil
	t.deleteBtn = nil
}

func (t *treeViewRenderer) Layout(parent fyne.Size) {
	height := t.MinSize().Height

	checkW := t.check.MinSize().Width
	btnW := t.check.MinSize().Width
	lblW := parent.Width - btnW - checkW

	t.check.Resize(fyne.NewSize(checkW, height))

	t.label.Move(fyne.NewPos(btnW, 0))
	t.label.Resize(fyne.NewSize(lblW, height))

	t.deleteBtn.Move(fyne.NewPos(btnW+lblW, 0))
	t.deleteBtn.Resize(fyne.NewSize(btnW, height))
}

func (t *treeViewRenderer) MinSize() fyne.Size {
	return runningWidth(t.check.MinSize(), t.label.MinSize(), t.deleteBtn.MinSize())
}

func (t *treeViewRenderer) Objects() []fyne.CanvasObject {
	return t.objects
}

func (t *treeViewRenderer) Refresh() {
	t.check.Refresh()
	t.label.Refresh()
	t.deleteBtn.Refresh()
}

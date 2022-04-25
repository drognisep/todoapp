package view

import (
	"fyne.io/fyne/v2/container"
	"sync"
	"todo/data"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var _ fyne.Widget = (*TaskView)(nil)
var _ fyne.Draggable = (*TaskView)(nil)

type TaskView struct {
	widget.BaseWidget
	window        fyne.Window
	listContainer *TaskList
	id            widget.ListItemID

	mux      sync.RWMutex
	renderer *treeViewRenderer
	dragging bool
}

func (t *TaskView) Dragged(*fyne.DragEvent) {
	t.mux.RLock()
	dragging := t.dragging
	id := t.id
	t.mux.RUnlock()
	if !dragging {
		t.mux.Lock()
		t.dragging = true
		t.mux.Unlock()
		t.listContainer.DragFrom(id)
	}
}

func (t *TaskView) DragEnd() {
	t.mux.Lock()
	t.dragging = false
	t.mux.Unlock()
	t.listContainer.DragEnd()
}

func (t *TaskView) MouseIn(*desktop.MouseEvent) {
	t.mux.RLock()
	id := t.id
	t.mux.RUnlock()
	t.listContainer.DragTo(id)
}
func (t *TaskView) MouseMoved(*desktop.MouseEvent) {}
func (t *TaskView) MouseOut()                      {}

func PrototypeTaskView(ctx *UiCtx, parent *TaskList) *TaskView {
	view := &TaskView{
		window:        ctx.MainWindow,
		listContainer: parent,
	}
	view.ExtendBaseWidget(view)
	return view
}

func (t *TaskView) update(id widget.ListItemID, task *data.Task, onDelete func()) {
	t.mux.RLock()
	t.id = id
	renderer := t.renderer
	t.mux.RUnlock()
	if renderer == nil {
		return
	}
	t.mux.Lock()
	renderer.mux.Lock()
	check, label, deleteBtn := renderer.check, renderer.label, renderer.deleteBtn
	for _, obj := range []fyne.CanvasObject{check, label, deleteBtn} {
		if obj == nil {
			renderer.mux.Unlock()
			t.mux.Unlock()
			return
		}
	}
	nameBinding := binding.BindString(&task.Name)
	descBinding := binding.BindString(&task.Description)
	check.Bind(binding.BindBool(&task.Done))
	label.Bind(nameBinding)
	deleteBtn.OnTapped = onDelete
	label.OnDoubleTap = func() {
		newName := task.Name
		nameEntry := widget.NewEntryWithData(binding.BindString(&newName))
		nameEntry.Wrapping = fyne.TextWrapWord
		nameEntry.Hide()
		nameLabel := NewTappableLabel(newName)
		nameLabel.OnDoubleTap = func() {
			nameLabel.Hide()
			nameEntry.Show()
		}
		newDesc := task.Description

		descriptionEntry := &widget.Entry{MultiLine: true, Wrapping: fyne.TextWrapWord}
		descriptionEntry.Bind(binding.BindString(&newDesc))
		descriptionEntry.Hide()

		btn := widget.NewButton("Add a description", func() {})
		btn.Importance = widget.HighImportance
		btn.OnTapped = func() {
			descriptionEntry.Show()
			btn.Hide()
		}
		descriptionMarkdown := NewTappableMarkdown(newDesc)
		descriptionMarkdown.OnDoubleTap = func() {
			descriptionMarkdown.Hide()
			descriptionEntry.Show()
		}
		if newDesc == "" {
			descriptionMarkdown.Hide()
		} else {
			btn.Hide()
		}
		summaryEntryItem := widget.NewFormItem("Name", container.NewVBox(nameLabel, nameEntry))
		descEntryItem := widget.NewFormItem("Description", container.NewVBox(btn, descriptionMarkdown, minHeightEntry(descriptionEntry, 300)))
		d := dialog.NewForm("Edit Task", "Done", "Cancel", []*widget.FormItem{summaryEntryItem, descEntryItem}, func(confirmed bool) {
			if confirmed {
				_ = nameBinding.Set(newName)
				_ = descBinding.Set(newDesc)
				t.Refresh()
				t.listContainer.Refresh()
			}
		}, t.window)
		d.Resize(fyne.NewSize(BigFloat, BigFloat))
		d.Show()
	}
	renderer.mux.Unlock()
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
	mux       sync.RWMutex
	check     *widget.Check
	label     *TappableLabel
	deleteBtn *widget.Button
	objects   []fyne.CanvasObject
}

func (t *treeViewRenderer) Destroy() {
}

func (t *treeViewRenderer) Layout(parent fyne.Size) {
	height := t.MinSize().Height

	t.mux.Lock()
	defer t.mux.Unlock()
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
	t.mux.RLock()
	size := runningWidth(t.check.MinSize(), t.label.MinSize(), t.deleteBtn.MinSize())
	t.mux.RUnlock()
	return size
}

func (t *treeViewRenderer) Objects() []fyne.CanvasObject {
	t.mux.RLock()
	objs := t.objects
	t.mux.RUnlock()
	return objs
}

func (t *treeViewRenderer) AnyObjectNil() bool {
	t.mux.RLock()
	defer t.mux.RUnlock()
	for _, obj := range []fyne.CanvasObject{t.check, t.label, t.deleteBtn} {
		if obj == nil {
			return true
		}
	}
	return false
}

func (t *treeViewRenderer) Refresh() {
	t.mux.Lock()
	defer t.mux.Unlock()
	t.check.Refresh()
	t.label.Refresh()
	t.deleteBtn.Refresh()
}

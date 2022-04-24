package main

import (
	_ "embed"
	"fmt"
	"runtime/debug"
	"todo/data"
	"todo/view"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	osdialog "github.com/sqweek/dialog"
)

//go:embed Icon.png
var icon []byte

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
	w.SetIcon(fyne.NewStaticResource("Icon.png", icon))
	ctx := &view.UiCtx{
		App:        a,
		MainWindow: w,
	}

	ctrlS := desktop.CustomShortcut{
		KeyName:  fyne.KeyS,
		Modifier: desktop.ControlModifier,
	}
	w.Canvas().AddShortcut(&ctrlS, func(_ fyne.Shortcut) {
		ctx.SaveData(model)
	})

	w.SetContent(view.TaskPanel(ctx, model))
	w.SetOnClosed(func() {
		ctx.SaveData(model)
	})
	w.ShowAndRun()
}

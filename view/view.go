package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
)

const (
	bigFloat           float32 = 1_000.0
	descEntryMinHeight float32 = 300
)

var (
	taskNameValidator = validation.NewRegexp(`^[A-Za-z0-9 -_]+$`, "Must be alphanumeric, '-', '_', or spaces")
)

func runningWidth(sizes ...fyne.Size) fyne.Size {
	var maxHeight float32 = 0.0
	var widthSum float32 = 0.0
	for _, sz := range sizes {
		if sz.Height > maxHeight {
			maxHeight = sz.Height
		}
		widthSum += sz.Width
	}
	return fyne.NewSize(widthSum, maxHeight)
}

type minHeightObj struct {
	*widget.Entry

	minHeight float32
}

func minHeightEntry(obj *widget.Entry, minHeight float32) fyne.Widget {
	wid := &minHeightObj{Entry: obj, minHeight: minHeight}
	wid.ExtendBaseWidget(wid)
	return wid
}

func (mho *minHeightObj) MinSize() fyne.Size {
	orig := mho.Entry.MinSize()
	return fyne.NewSize(orig.Width, mho.minHeight)
}

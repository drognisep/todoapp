package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/validation"
)

const (
	BigFloat float32 = 1_000.0
)

var (
	TaskNameValidator = validation.NewRegexp(`^[A-Za-z0-9 -_]+$`, "Must be alphanumeric, '-', '_', or spaces")
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

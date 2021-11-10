package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

var _ fyne.DoubleTappable = (*TappableLabel)(nil)

type TappableLabel struct {
	widget.Label

	OnDoubleTap func()
}

func NewTappableLabel(placeholder string) *TappableLabel {
	lbl := &TappableLabel{
		Label: widget.Label{
			Text:       placeholder,
			Alignment:  fyne.TextAlignLeading,
			Wrapping:   fyne.TextTruncate,
			TextStyle:  fyne.TextStyle{},
		},
	}
	lbl.Text = placeholder
	lbl.ExtendBaseWidget(lbl)
	return lbl
}

func (t *TappableLabel) DoubleTapped(_ *fyne.PointEvent) {
	if fn := t.OnDoubleTap; fn != nil {
		fn()
	}
}

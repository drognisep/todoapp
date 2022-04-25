package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
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
			Text:      placeholder,
			Alignment: fyne.TextAlignLeading,
			Wrapping:  fyne.TextTruncate,
			TextStyle: fyne.TextStyle{},
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

type TappableMarkdown struct {
	widget.RichText

	OnDoubleTap func()
}

func NewTappableMarkdown(placeholder string) *TappableMarkdown {
	mkd := &TappableMarkdown{
		RichText: widget.RichText{},
	}
	mkd.ParseMarkdown(placeholder)
	mkd.ExtendBaseWidget(mkd)
	mkd.Wrapping = fyne.TextWrapWord
	mkd.Scroll = container.ScrollVerticalOnly
	return mkd
}

func (t *TappableMarkdown) DoubleTapped(_ *fyne.PointEvent) {
	if fn := t.OnDoubleTap; fn != nil {
		fn()
	}
}

func (t *TappableMarkdown) MinSize() fyne.Size {
	return fyne.NewSize(0, 300)
}

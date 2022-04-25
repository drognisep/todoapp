package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var _ fyne.DoubleTappable = (*tappableLabel)(nil)

type tappableLabel struct {
	widget.Label

	OnDoubleTap func()
}

func newTappableLabel(placeholder string) *tappableLabel {
	lbl := &tappableLabel{
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

func (t *tappableLabel) DoubleTapped(_ *fyne.PointEvent) {
	if fn := t.OnDoubleTap; fn != nil {
		fn()
	}
}

type tappableMarkdown struct {
	widget.RichText

	OnDoubleTap func()
}

func newTappableMarkdown(placeholder string) *tappableMarkdown {
	mkd := &tappableMarkdown{
		RichText: widget.RichText{},
	}
	mkd.ParseMarkdown(placeholder)
	mkd.ExtendBaseWidget(mkd)
	mkd.Wrapping = fyne.TextWrapWord
	mkd.Scroll = container.ScrollVerticalOnly
	return mkd
}

func (t *tappableMarkdown) DoubleTapped(_ *fyne.PointEvent) {
	if fn := t.OnDoubleTap; fn != nil {
		fn()
	}
}

func (t *tappableMarkdown) MinSize() fyne.Size {
	return fyne.NewSize(0, descEntryMinHeight)
}

package ui

import "fyne.io/fyne/v2"

type FixedSizeLayout struct {
	Width  float32
	Height float32
}

func (f *FixedSizeLayout) getSize(object fyne.CanvasObject) fyne.Size {
	width := object.MinSize().Width
	height := object.MinSize().Height

	if f.Width >= 0 {
		width = f.Width
	}
	if f.Height >= 0 {
		height = f.Height
	}

	return fyne.NewSize(width, height)
}

func (f *FixedSizeLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return f.getSize(objects[0])
}

func (f *FixedSizeLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	for _, o := range objects {
		o.Resize(f.getSize(o))
	}
}

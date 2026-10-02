package customLayout

import "fyne.io/fyne/v2"

type ResponsivePaddingLayout struct {
	padding float32
}

func NewResponsivePaddingLayout(padding float32) *ResponsivePaddingLayout {
	return &ResponsivePaddingLayout{
		padding: padding,
	}
}

func (a *ResponsivePaddingLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	width, height := float32(0), float32(0)

	for _, o := range objects {
		childSize := o.MinSize()

		width += childSize.Width
		height += childSize.Height
	}
	return fyne.NewSize(width, height)
}

func (a *ResponsivePaddingLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	widthPadding := size.Width * a.padding
	heightPadding := size.Height * a.padding

	x := widthPadding
	y := heightPadding
	width := size.Width - (widthPadding * 2)
	height := size.Height - (heightPadding * 2)

	for _, o := range objects {
		o.Move(fyne.NewPos(x, y))
		o.Resize(fyne.NewSize(width, height))
	}
}

package animations

import (
	"github.com/AllenDang/giu"
	"github.com/AllenDang/imgui-go"
)

// ResizeAnimation allows to resize an UI element.
type ResizeAnimation[T giu.Widget] struct {
	widget resizable2D[T]
	sizes  []imgui.Vec2
	id     string
}

// Resize creates ResizeAnimation.
// It requires Widget type parameter (e.g. *giu.ButtonWidget).
func Resize[T giu.Widget](w resizable2D[T], sizes ...imgui.Vec2) *ResizeAnimation[T] {
	return &ResizeAnimation[T]{
		id:     giu.GenAutoID("giu-animations-ResizeAnimation"),
		widget: w,
		sizes:  sizes,
	}
}

// ID allows to set ID manually.
func (r *ResizeAnimation[T]) ID(id string) {
	r.id = id
}

// Init implements Animation interface.
func (r *ResizeAnimation[T]) Init() {
	// noop
}

// Reset implements Animation.
func (r *ResizeAnimation[T]) Reset() {
	// noop
}

// KeyFramesCount implements Animation.
func (r *ResizeAnimation[T]) KeyFramesCount() int {
	return len(r.sizes)
}

// BuildNormal implements Animation.
func (r *ResizeAnimation[T]) BuildNormal(currentKeyFrame KeyFrame, _ StarterFunc) {
	c := imgui.CursorPos()
	imgui.SetCursorPos(imgui.Vec2{
		X: c.X - (r.sizes[currentKeyFrame].X-r.sizes[0].X)/2,
		Y: c.Y - (r.sizes[currentKeyFrame].Y-r.sizes[0].Y)/2,
	})

	r.widget.Size(r.sizes[currentKeyFrame].X, r.sizes[currentKeyFrame].Y).Build()
}

// BuildAnimation implements Animation.
func (r *ResizeAnimation[T]) BuildAnimation(
	animationPercentage, _ float32,
	baseKeyFrame, destinationKeyFrame KeyFrame,
	_ PlayMode, _ StarterFunc,
) {
	delta := imgui.Vec2{
		X: (r.sizes[destinationKeyFrame].X - r.sizes[baseKeyFrame].X) * animationPercentage,
		Y: (r.sizes[destinationKeyFrame].Y - r.sizes[baseKeyFrame].Y) * animationPercentage,
	}

	c := imgui.CursorPos()
	imgui.SetCursorPos(imgui.Vec2{
		X: c.X - (r.sizes[baseKeyFrame].X-r.sizes[0].X+delta.X)/2,
		Y: c.Y - (r.sizes[baseKeyFrame].Y-r.sizes[0].Y+delta.Y)/2,
	})

	r.widget.Size(r.sizes[baseKeyFrame].X+delta.X, r.sizes[baseKeyFrame].Y+delta.Y).Build()
}

type resizable2D[T giu.Widget] interface {
	Size(w, h float32) T
	giu.Widget
}

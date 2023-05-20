package animations

import (
	"github.com/AllenDang/giu"
	"github.com/AllenDang/imgui-go"
)

type TrickCursor byte

const (
	TrickNever TrickCursor = 1 << iota
	TrickCursorBefore
	TrickCursorAfter
)

// ResizeAnimation allows to resize an UI element.
type ResizeAnimation[T giu.Widget] struct {
	widget      resizable2D[T]
	sizes       []imgui.Vec2
	id          string
	trickCursor TrickCursor
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

func (r *ResizeAnimation[T]) TrickCursor(t TrickCursor) *ResizeAnimation[T] {
	r.trickCursor = t
	return r
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
	r.trickCursorBefore(r.sizes[currentKeyFrame], imgui.Vec2{})

	r.widget.Size(r.sizes[currentKeyFrame].X, r.sizes[currentKeyFrame].Y).Build()

	//imgui.SetCursorPos(imgui.Vec2{
	//	X: 0, // This is because cursor is alreay moved to the next line TODO: what if Sameline called?
	//Y: c.Y + r.sizes[0].Y,
	//})
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

	r.trickCursorBefore(r.sizes[baseKeyFrame], delta)

	r.widget.Size(r.sizes[baseKeyFrame].X+delta.X, r.sizes[baseKeyFrame].Y+delta.Y).Build()

	//imgui.SetCursorPos(imgui.Vec2{
	//	X: 0, // This is because cursor is alreay moved to the next line TODO: what if Sameline called?
	//Y: c.Y + r.sizes[0].Y,
	//})
}

func (r *ResizeAnimation[T]) trickCursorBefore(current, delta imgui.Vec2) {
	if r.trickCursor&TrickCursorBefore == 0 {
		return
	}

	c := imgui.CursorPos()
	imgui.SetCursorPos(imgui.Vec2{
		X: c.X - (current.X+delta.X-r.sizes[0].X)/2,
		Y: c.Y - (current.Y+delta.Y-r.sizes[0].Y)/2,
	})
}

type resizable2D[T giu.Widget] interface {
	Size(w, h float32) T
	giu.Widget
}

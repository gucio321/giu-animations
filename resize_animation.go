package animations

import (
	"github.com/AllenDang/giu"
	"github.com/AllenDang/imgui-go"
)

type ResizeAnimation[T giu.Widget] struct {
	widget resizable2D[T]
	sizes  []imgui.Vec2
}

func Resize[T giu.Widget](w resizable2D[T], sizes ...imgui.Vec2) *ResizeAnimation[T] {
	return &ResizeAnimation[T]{
		widget: w,
		sizes:  sizes,
	}
}

func (r ResizeAnimation[T]) Init() {
	// noop
}

func (r ResizeAnimation[T]) Reset() {
	// noop
}

func (r ResizeAnimation[T]) KeyFramesCount() int {
	return len(r.sizes)
}

func (r ResizeAnimation[T]) BuildNormal(currentKeyFrame KeyFrame, starterFunc StarterFunc) {
	r.widget.Size(r.sizes[currentKeyFrame].X, r.sizes[currentKeyFrame].Y).Build()
}

func (r ResizeAnimation[T]) BuildAnimation(animationPercentage, arbitraryPercentage float32, baseKeyFrame, destinationKeyFrame KeyFrame, mode PlayMode, starterFunc StarterFunc) {
	delta := imgui.Vec2{
		X: r.sizes[baseKeyFrame].X + (r.sizes[destinationKeyFrame].X-r.sizes[baseKeyFrame].X)*animationPercentage,
		Y: r.sizes[baseKeyFrame].Y + (r.sizes[destinationKeyFrame].Y-r.sizes[baseKeyFrame].Y)*animationPercentage,
	}

	r.widget.Size(delta.X, delta.Y).Build()
}

type resizable2D[T giu.Widget] interface {
	Size(w, h float32) T
	giu.Widget
}

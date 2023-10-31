package animations

import (
	"github.com/AllenDang/giu"
	"github.com/AllenDang/imgui-go"
)

// TrickCursor allows to enable special Drawing Cursor behavior.
// If enabled, Drawing Cursor is moved before rendering ResizeAnimation.
// It allows to reach an effect of resizing a UI element from its center.
// NOTE: Be careful as this feature (especially TrickCursorAfterX) may break imgui.Sameline() mechanism.
// NOTE: You can test it in _examples/main.go.
type TrickCursor byte

// These constants work as a bitmask. Feel free to use it like TrickCursorBeforeX | TrickCursorAfterY.
const (
	// TrickNever allows to disable TrickCursor.
	TrickNever TrickCursor = 1 << iota

	// TrickCursorBeforeX allows to move Drawing Cursor before rendering ResizeAnimation on X axis.
	TrickCursorBeforeX
	// TrickCursorBeforeY allows to move Drawing Cursor before rendering ResizeAnimation on Y axis.
	TrickCursorBeforeY
	// TrickCursorAfterX allows to move Drawing Cursor after rendering ResizeAnimation on X axis.
	// IMPORTANT: Applying this will interfere with imgui.Sameline() mechanism.
	TrickCursorAfterX
	// TrickCursorAfterY allows to move Drawing Cursor after rendering ResizeAnimation on Y axis.
	TrickCursorAfterY

	// TrickCursorBefore allows to move Drawing Cursor before rendering ResizeAnimation on both axes.
	TrickCursorBefore = TrickCursorBeforeX | TrickCursorBeforeY
	// TrickCursorAfter allows to move Drawing Cursor after rendering ResizeAnimation on both axes.
	TrickCursorAfter = TrickCursorAfterX | TrickCursorAfterY
	// TrickCursorAlways allows to move Drawing Cursor before and after rendering ResizeAnimation on both axes.
	TrickCursorAlways = TrickCursorBefore | TrickCursorAfter
)

// ResizeAnimation allows to resize a UI element.
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
		id:          giu.GenAutoID("giu-animations-ResizeAnimation"),
		widget:      w,
		sizes:       sizes,
		trickCursor: TrickCursorBeforeY,
	}
}

// ID allows to set ID manually.
func (r *ResizeAnimation[T]) ID(id string) {
	r.id = id
}

// TrickCursor allows to set TrickCursor.
// Default: TrickCursorBefore.
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
func (r *ResizeAnimation[T]) BuildNormal(currentKeyFrame KeyFrame, _ StarterFunc, triggerCheck func()) {
	// This may happen if user forgot to pass size vectors. In this case just allow to build unchanged widget.
	if int(currentKeyFrame) > len(r.sizes)-1 {
		r.widget.Build()

		return
	}

	r.trickCursorBefore(r.sizes[currentKeyFrame], imgui.Vec2{})

	r.widget.Size(r.sizes[currentKeyFrame].X, r.sizes[currentKeyFrame].Y).Build()
	triggerCheck()

	r.trickCursorAfter(r.sizes[currentKeyFrame])
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

	newSize := imgui.Vec2{
		X: r.sizes[baseKeyFrame].X + delta.X,
		Y: r.sizes[baseKeyFrame].Y + delta.Y,
	}

	r.widget.Size(newSize.X, newSize.Y).Build()

	r.trickCursorAfter(newSize)
}

func (r *ResizeAnimation[T]) trickCursorBefore(current, delta imgui.Vec2) {
	move := imgui.Vec2{}

	if r.trickCursor&TrickCursorBeforeX != 0 {
		move.X -= (current.X + delta.X - r.sizes[0].X) / 2
	}

	if r.trickCursor&TrickCursorBeforeY != 0 {
		move.Y -= (current.Y + delta.Y - r.sizes[0].Y) / 2
	}

	imgui.Dummy(move)
}

func (r *ResizeAnimation[T]) trickCursorAfter(currentSize imgui.Vec2) {
	move := imgui.Vec2{}

	if r.trickCursor&TrickCursorAfterX != 0 {
		move.X -= (currentSize.X - r.sizes[0].X) / 2
	}

	if r.trickCursor&TrickCursorAfterY != 0 {
		move.Y -= (currentSize.Y - r.sizes[0].Y) / 2
	}

	imgui.Dummy(move)
}

type resizable2D[T giu.Widget] interface {
	Size(w, h float32) T
	giu.Widget
}

package animations

import (
	"fmt"
	"github.com/AllenDang/giu"
	"github.com/AllenDang/imgui-go"
	"log"
	"sync"
)

type resizeAnimationState struct {
	cursorPos imgui.Vec2
	m         *sync.Mutex
}

func (r *resizeAnimationState) Dispose() {
	// noop
}

func (r *ResizeAnimation[T]) newState() *resizeAnimationState {
	return &resizeAnimationState{
		m: &sync.Mutex{},
	}
}

// getState returns animator's state.
// It could not be public, because of concurrency issues.
// There is animation bunch of Animator's methods that allows
// user to obtain certain data.
func (r *ResizeAnimation[T]) getState() *resizeAnimationState {
	if s := giu.Context.GetState(r.id); s != nil {
		state, ok := s.(*resizeAnimationState)
		if !ok {
			log.Panicf("error asserting type of animator state: got %T, wanted *animatorState", s)
		}

		return state
	}

	giu.Context.SetState(r.id, r.newState())

	return r.getState()
}

type ResizeAnimation[T giu.Widget] struct {
	widget resizable2D[T]
	sizes  []imgui.Vec2
	id     string
}

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

func (r *ResizeAnimation[T]) Init() {
	state := r.getState()
	state.m.Lock()
	defer state.m.Unlock()
	state.cursorPos = imgui.CursorPos()
}

func (r *ResizeAnimation[T]) Reset() {
	//r.Init()
}

func (r *ResizeAnimation[T]) KeyFramesCount() int {
	return len(r.sizes)
}

func (r *ResizeAnimation[T]) BuildNormal(currentKeyFrame KeyFrame, starterFunc StarterFunc) {
	state := r.getState()

	state.m.Lock()
	imgui.SetCursorPos(
		imgui.Vec2{
			X: state.cursorPos.X - (r.sizes[currentKeyFrame].X-r.sizes[0].X)/2,
			Y: state.cursorPos.Y - (r.sizes[currentKeyFrame].Y-r.sizes[0].Y)/2,
		})
	state.m.Unlock()

	r.widget.Size(r.sizes[currentKeyFrame].X, r.sizes[currentKeyFrame].Y).Build()
}

func (r *ResizeAnimation[T]) BuildAnimation(animationPercentage, arbitraryPercentage float32, baseKeyFrame, destinationKeyFrame KeyFrame, mode PlayMode, starterFunc StarterFunc) {
	delta := imgui.Vec2{
		X: (r.sizes[destinationKeyFrame].X - r.sizes[baseKeyFrame].X) * animationPercentage,
		Y: (r.sizes[destinationKeyFrame].Y - r.sizes[baseKeyFrame].Y) * animationPercentage,
	}

	state := r.getState()
	state.m.Lock()

	fmt.Println(delta)
	imgui.SetCursorPos(
		imgui.Vec2{
			X: state.cursorPos.X - (r.sizes[baseKeyFrame].X-r.sizes[0].X+delta.X)/2,
			Y: state.cursorPos.Y - (r.sizes[baseKeyFrame].Y-r.sizes[0].Y+delta.Y)/2,
		})

	state.m.Unlock()

	r.widget.Size(r.sizes[baseKeyFrame].X+delta.X, r.sizes[baseKeyFrame].Y+delta.Y).Build()
}

type resizable2D[T giu.Widget] interface {
	Size(w, h float32) T
	giu.Widget
}

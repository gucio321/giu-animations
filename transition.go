package animations

import (
	"log"
	"sync"

	"github.com/AllenDang/giu"
	"github.com/AllenDang/imgui-go"
)

type transitionAnimationState struct {
	layout int
	m      *sync.Mutex
}

func (s *transitionAnimationState) Dispose() {}

var _ Animation = &TransitionAnimation{}

type TransitionAnimation struct {
	id                   string
	renderer1, renderer2 func(starter func())
}

func Transition(renderer1, renderer2 func(starter func())) *TransitionAnimation {
	return &TransitionAnimation{
		id:        giu.GenAutoID("transitionAnimation"),
		renderer1: renderer1,
		renderer2: renderer2,
	}
}

func (t *TransitionAnimation) KeyFrames() int {
	return 2
}

func (t *TransitionAnimation) BuildAnimation(percentage, _ float32, _, _ KeyFrame, starter StarterFunc) {
	state := t.getState()
	// it means the current layout is layout1, so increasing percentage
	if state.layout == 1 {
		percentage = 1 - percentage
	}

	imgui.PushStyleVarFloat(imgui.StyleVarAlpha, percentage)
	t.renderer1(func() { starter(PlayAuto) })
	imgui.PopStyleVar()
	imgui.PushStyleVarFloat(imgui.StyleVarAlpha, 1-percentage)
	t.renderer2(func() { starter(PlayAuto) })
	imgui.PopStyleVar()
}

func (t *TransitionAnimation) Reset() {
	state := t.getState()
	state.m.Lock()

	if state.layout == 0 {
		state.layout = 1
	} else {
		state.layout = 0
	}

	state.m.Unlock()
}

func (t *TransitionAnimation) Init() {
	// noop
}

func (t *TransitionAnimation) BuildNormal(_ KeyFrame, starter StarterFunc) {
	state := t.getState()

	state.m.Lock()
	layout := state.layout
	state.m.Unlock()

	if layout == 0 {
		t.renderer1(func() { starter(PlayAuto) })
	} else {
		t.renderer2(func() { starter(PlayAuto) })
	}
}

func (t *TransitionAnimation) getState() *transitionAnimationState {
	if s := giu.Context.GetState(t.id); s != nil {
		state, ok := s.(*transitionAnimationState)
		if !ok {
			log.Panicf("expected state type *transitionAnimationState, got %T", s)
		}

		return state
	}

	giu.Context.SetState(t.id, t.newState())

	return t.getState()
}

func (t *TransitionAnimation) newState() *transitionAnimationState {
	return &transitionAnimationState{
		m: &sync.Mutex{},
	}
}

package animations

import (
	"log"

	"github.com/AllenDang/giu"
	"github.com/AllenDang/imgui-go"
)

type transitionAnimationState struct {
	layout int
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

func (t *TransitionAnimation) BuildAnimation(percentage float32, starter func()) {
	state := t.getState()
	// it means the current layou is layout1, so increasing procentage
	if state.layout == 1 {
		percentage = 1 - percentage
	}

	imgui.PushStyleVarFloat(imgui.StyleVarAlpha, percentage)
	t.renderer1(starter)
	imgui.PopStyleVar()
	imgui.PushStyleVarFloat(imgui.StyleVarAlpha, 1-percentage)
	t.renderer2(starter)
	imgui.PopStyleVar()
}

func (t *TransitionAnimation) Reset() {
	state := t.getState()
	if state.layout == 0 {
		state.layout = 1
	} else {
		state.layout = 0
	}
}

func (t *TransitionAnimation) Init() {
	// noop
}

func (t *TransitionAnimation) BuildNormal(starter func()) {
	state := t.getState()

	if state.layout == 0 {
		t.renderer1(starter)
	} else {
		t.renderer2(starter)
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
	return &transitionAnimationState{}
}

package animations

import (
	"github.com/AllenDang/giu"
	"github.com/AllenDang/imgui-go"
	"github.com/gucio321/giu-animations/internal/logger"
)

type transitionAnimationState struct {
	layout int
}

func (s *transitionAnimationState) Dispose() {}

var _ Animation = &TransitionWidget{}

type TransitionWidget struct {
	id                   string
	renderer1, renderer2 func(starter func())
}

func Transition(renderer1, renderer2 func(starter func())) *TransitionWidget {
	return &TransitionWidget{
		id:        giu.GenAutoID("transitionAnimation"),
		renderer1: renderer1,
		renderer2: renderer2,
	}
}

func (t *TransitionWidget) BuildAnimation(percentage float32, starter func()) {
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

func (t *TransitionWidget) Reset() {
	state := t.getState()
	if state.layout == 0 {
		state.layout = 1
	} else {
		state.layout = 0
	}
}

func (t *TransitionWidget) Init() {
	// noop
}

func (t *TransitionWidget) BuildNormal(starter func()) {
	state := t.getState()

	if state.layout == 0 {
		t.renderer1(starter)
	} else {
		t.renderer2(starter)
	}

}

func (t *TransitionWidget) getState() *transitionAnimationState {
	if s := giu.Context.GetState(t.id); s != nil {
		state, ok := s.(*transitionAnimationState)
		if !ok {
			logger.Fatalf("expected state type *hoverColorAnimationState, got %T", s)
		}

		return state
	}

	giu.Context.SetState(t.id, t.newState())

	return t.getState()
}

func (t *TransitionWidget) newState() *transitionAnimationState {
	return &transitionAnimationState{}
}

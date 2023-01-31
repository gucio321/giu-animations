package animations

import (
	"github.com/AllenDang/giu"
	"github.com/AllenDang/imgui-go"
	"log"
)

type moveAnimationState struct {
	state    bool
	startPos imgui.Vec2
}

func (m *moveAnimationState) Dispose() {
	// noop
}

func (m *MoveAnimation) getState() *moveAnimationState {
	if s := giu.Context.GetState(m.id); s != nil {
		state, ok := s.(*moveAnimationState)
		if !ok {
			log.Panicf("error asserting type of move animation state: got %T, wanted *moveAnimationState", s)
		}

		return state
	}

	giu.Context.SetState(m.id, m.newState())

	return m.getState()
}

func (m *MoveAnimation) newState() *moveAnimationState {
	return &moveAnimationState{}
}

var _ Animation = &MoveAnimation{}

type MoveAnimation struct {
	id string

	widget   func(starter func()) giu.Widget
	posDelta imgui.Vec2

	alg EasingAlgorithmType
}

func Move(w func(starter func()) giu.Widget, posDelta imgui.Vec2) *MoveAnimation {
	return &MoveAnimation{
		id:       giu.GenAutoID("MoveAnimation"),
		widget:   w,
		posDelta: posDelta,
	}
}

func (m *MoveAnimation) Init() {
	m.getState().startPos = imgui.CursorPos()
}

func (m *MoveAnimation) Reset() {
	state := m.getState()
	state.state = !state.state
}

func (m *MoveAnimation) Algorithm(algorithm EasingAlgorithmType) *MoveAnimation {
	m.alg = algorithm
	return m
}

func (m *MoveAnimation) BuildNormal(starter func()) {
	state := m.getState()
	if state.state {
		p := imgui.Vec2{
			X: m.posDelta.X + state.startPos.X,
			Y: m.posDelta.Y + state.startPos.Y,
		}
		imgui.SetCursorPos(p)
	} else {
		imgui.SetCursorPos(state.startPos)
	}

	m.widget(starter).Build()
}

func (m *MoveAnimation) BuildAnimation(animationPercentage float32, starter func()) {
	state := m.getState()

	switch m.alg {
	case EasingAlgNone:
	// noop
	case EasingAlgInSine:
		animationPercentage = easingAlgInSine(animationPercentage)
	case EasingAlgOutSine:
		animationPercentage = easingAlgOutSine(animationPercentage)
	case EasingAlgInOutSine:
		animationPercentage = easingAlgInOutSine(animationPercentage)
	case EasingAlgInBack:
		animationPercentage = easingAlgInBack(animationPercentage)
	case EasingAlgOutBack:
		animationPercentage = easingAlgOutBack(animationPercentage)
	case EasingAlgInOutBack:
		animationPercentage = easingAlgInOutBack(animationPercentage)
	case EasingAlgInElastic:
		animationPercentage = easingAlgInElastic(animationPercentage)
	case EasingAlgOutElastic:
		animationPercentage = easingAlgOutElastic(animationPercentage)
	case EasingAlgInOutElastic:
		animationPercentage = easingAlgInOutElastic(animationPercentage)
	case EasingAlgInBounce:
		animationPercentage = easingAlgInBounce(animationPercentage)
	case EasingAlgOutBounce:
		animationPercentage = easingAlgOutBounce(animationPercentage)
	case EasingAlgInOutBounce:
		animationPercentage = easingAlgInOutBounce(animationPercentage)
	}

	if !state.state {
		animationPercentage = 1 - animationPercentage
	}

	pos := imgui.Vec2{
		X: state.startPos.X + m.posDelta.X*animationPercentage,
		Y: state.startPos.Y + m.posDelta.Y*animationPercentage,
	}

	imgui.SetCursorScreenPos(pos)
	m.widget(starter).Build()
}

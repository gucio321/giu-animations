package animations

import (
	"log"

	"github.com/AllenDang/giu"
	"github.com/AllenDang/imgui-go"
)

var _ Animation = &MoveAnimation{}

// MoveAnimation moves animation widget from start position to destination.
// You can also specify animation Bézier curve's points.
type MoveAnimation struct {
	id string

	widget func(starter StarterFunc) giu.Widget
	steps  []*MoveStep

	startStep func(startPos imgui.Vec2) *MoveStep
}

// Move creates new *MoveAnimations
// NOTE: You may want to take animation look on StartPos or DefaultStartPos methods to specify animation starting position.
// otherwise the first step specified will be treated as start position.
func Move(w func(starter StarterFunc) giu.Widget, steps ...*MoveStep) *MoveAnimation {
	return &MoveAnimation{
		id:     giu.GenAutoID("MoveAnimation"),
		steps:  steps,
		widget: w,
	}
}

// StartPos allows to specify custom StartPos (item will be moved there immediately).
// argument function will receive cursor position returned by imgui.GetCursorPos while initializing animation.
func (m *MoveAnimation) StartPos(startPosStep func(startPos imgui.Vec2) *MoveStep) *MoveAnimation {
	m.startStep = startPosStep

	return m
}

// DefaultStartPos will set animation default value of MoveStep as animation starting step.
// NOTE: You will lose possibility of setting up any additional properties of MoveStep (like bezier points).
func (m *MoveAnimation) DefaultStartPos() *MoveAnimation {
	return m.StartPos(StepVec)
}

// Init implements Animation.
func (m *MoveAnimation) Init() {
	m.getState().startPos = imgui.CursorPos()
}

// Reset implements Animation.
func (m *MoveAnimation) Reset() {
	// noop
}

// KeyFramesCount implements Animation interface.
func (m *MoveAnimation) KeyFramesCount() int {
	l := len(m.steps)
	if m.startStep != nil {
		l++
	}

	return l
}

// BuildNormal implements Animation.
func (m *MoveAnimation) BuildNormal(currentKF KeyFrame, starter StarterFunc, _ func()) {
	imgui.SetCursorPos(m.getPosition(currentKF))

	m.widget(starter).Build()
}

// BuildAnimation implements Animation.
func (m *MoveAnimation) BuildAnimation(
	animationPercentage, _ float32,
	srcFrame, destFrame KeyFrame,
	mode PlayMode,
	starter StarterFunc,
) {
	startPos := m.getPosition(srcFrame)
	destPos := m.getPosition(destFrame)

	var pos imgui.Vec2

	steps := m.getSteps()

	// srcStep depends on animations play mode
	var (
		srcStep *MoveStep
		srcPos  imgui.Vec2
	)

	switch mode {
	case PlayForward:
		srcStep = steps[srcFrame]
		srcPos = m.getPosition(srcFrame)
	case PlayBackward:
		srcStep = steps[destFrame]
		srcPos = m.getPosition(destFrame)
	}

	if srcStep.useBezier {
		pts := []imgui.Vec2{startPos}
		l := len(srcStep.bezier)

		for i := 0; i < l; i++ {
			var b imgui.Vec2

			switch mode {
			case PlayForward:
				b = srcStep.bezier[i]
			case PlayBackward:
				b = srcStep.bezier[l-i-1]
			}

			pts = append(pts, imgui.Vec2{
				X: b.X + srcPos.X,
				Y: b.Y + srcPos.Y,
			})
		}

		pts = append(pts, destPos)
		pos = bezier(animationPercentage, pts)
	} else {
		pos = vecSum(startPos, vecMul(vecDif(destPos, startPos), animationPercentage))
	}

	imgui.SetCursorPos(pos)
	m.widget(starter).Build()
}

// this will return absolute position.
// If step specifies animation relative position, it will go to the previous step.
func (m *MoveAnimation) getPosition(currentKF KeyFrame) imgui.Vec2 {
	state := m.getState()

	steps := m.getSteps()
	if m.startStep != nil {
		steps = append([]*MoveStep{m.startStep(state.startPos)}, m.steps...)
	}

	pos := imgui.Vec2{}

	for i := int(currentKF); i >= 0; i-- {
		s := steps[i]

		pos = vecSum(pos, s.positionDelta)

		if s.isAbsolute {
			return pos
		}
	}

	return pos
}

// this will return animation list of steps with the first step added if necessary.
func (m *MoveAnimation) getSteps() []*MoveStep {
	state := m.getState()
	steps := m.steps

	if m.startStep != nil {
		steps = append([]*MoveStep{m.startStep(state.startPos)}, m.steps...)
	}

	return steps
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

var _ giu.Disposable = &moveAnimationState{}

type moveAnimationState struct {
	startPos imgui.Vec2
}

// Dispose implements giu.Disposable.
func (m *moveAnimationState) Dispose() {
	// noop
}

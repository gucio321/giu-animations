package animations

import (
	"github.com/AllenDang/giu"
	"github.com/AllenDang/imgui-go"
	"log"
)

var _ giu.Disposable = &moveAnimationState{}

type moveAnimationState struct {
	state    bool
	startPos imgui.Vec2
}

// Dispose implements giu.Disposable
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

// MoveAnimation moves a widget from start position to destination.
// You can also specify a Bézier curve's points.
type MoveAnimation struct {
	id string

	widget func(starter StarterFunc) giu.Widget
	steps  []*MoveStep

	startStep func(startPos imgui.Vec2) *MoveStep
}

// Move creates new *MoveAnimations
// StartPos will be current cursor position
func Move(w func(starter StarterFunc) giu.Widget, steps ...*MoveStep) *MoveAnimation {
	return &MoveAnimation{
		id:     giu.GenAutoID("MoveAnimation"),
		steps:  steps,
		widget: w,
	}
}

// StartPos allows to specify custom StartPos (item will be moved there immediately).
// argument function will receive cursor position returned by imgui.GetCursorPos while initializing animation,
// however the resulting MoveStep does not need to base on it.
func (m *MoveAnimation) StartPos(startPosStep func(startPos imgui.Vec2) *MoveStep) *MoveAnimation {
	m.startStep = startPosStep
	return m
}

func (m *MoveAnimation) DefaultStartPos() *MoveAnimation {
	return m.StartPos(func(p imgui.Vec2) *MoveStep {
		return StepVec(p)
	})
}

// Init implements Animation
func (m *MoveAnimation) Init() {
	m.getState().startPos = imgui.CursorPos()
}

// Reset implements Animation
func (m *MoveAnimation) Reset() {
	state := m.getState()
	state.state = !state.state
}

func (m *MoveAnimation) KeyFrames() int {
	l := len(m.steps)
	if m.startStep != nil {
		l++
	}

	return l
}

// BuildNormal implements Animation
func (m *MoveAnimation) BuildNormal(currentKF KeyFrame, starter StarterFunc) {
	imgui.SetCursorPos(m.getPosition(currentKF))

	m.widget(starter).Build()
}

// BuildAnimation implements Animation
func (m *MoveAnimation) BuildAnimation(animationPercentage, _ float32, srcFrame, destFrame KeyFrame, starter StarterFunc) {
	srcPos := m.getPosition(srcFrame)
	destPos := m.getPosition(destFrame)
	var pos imgui.Vec2

	srcStep := m.steps[srcFrame]
	if srcStep.useBezier {
		pts := []imgui.Vec2{srcPos}
		for _, b := range srcStep.bezier {
			pts = append(pts, imgui.Vec2{
				X: b.X + srcPos.X,
				Y: b.Y + srcPos.Y,
			})
		}
		pts = append(pts, destPos)
		pos = bezier(animationPercentage, pts)
	} else {
		pos = vecSum(srcPos, vecMul(vecDif(destPos, srcPos), animationPercentage))
	}

	imgui.SetCursorScreenPos(pos)
	m.widget(starter).Build()
}

type MoveStep struct {
	positionDelta imgui.Vec2
	isAbsolute    bool

	useBezier bool
	bezier    []imgui.Vec2
}

func Step(x, y float32) *MoveStep {
	return &MoveStep{
		positionDelta: imgui.Vec2{X: x, Y: y},
	}
}

func StepVec(v imgui.Vec2) *MoveStep {
	return &MoveStep{
		positionDelta: v,
	}
}

// Bezier allows to specify Bezier Curve points.
func (m *MoveStep) Bezier(points ...imgui.Vec2) *MoveStep {
	m.useBezier = true
	m.bezier = points

	return m
}

func (m *MoveStep) Absolute() *MoveStep {
	m.isAbsolute = true

	return m
}

func vecSum(vec1, vec2 imgui.Vec2) imgui.Vec2 {
	return imgui.Vec2{
		X: vec1.X + vec2.X,
		Y: vec1.Y + vec2.Y,
	}
}

func vecDif(vec1, vec2 imgui.Vec2) imgui.Vec2 {
	return imgui.Vec2{
		X: vec1.X - vec2.X,
		Y: vec1.Y - vec2.Y,
	}
}

func vecMul(vec1 imgui.Vec2, multiplier float32) imgui.Vec2 {
	return imgui.Vec2{
		X: vec1.X * multiplier,
		Y: vec1.Y * multiplier,
	}
}

func (m *MoveAnimation) getPosition(currentKF KeyFrame) imgui.Vec2 {
	state := m.getState()

	steps := m.steps
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

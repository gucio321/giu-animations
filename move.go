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

	useBezier bool
	bezier    []imgui.Vec2
}

func Move(w func(starter func()) giu.Widget, posDelta imgui.Vec2) *MoveAnimation {
	return &MoveAnimation{
		id:       giu.GenAutoID("MoveAnimation"),
		widget:   w,
		posDelta: posDelta,
	}
}

func (m *MoveAnimation) StartPos(startPos imgui.Vec2) *MoveAnimation {
	m.getState().startPos = startPos
	return m
}

func (m *MoveAnimation) Bezier(points ...imgui.Vec2) *MoveAnimation {
	if len(points) > 2 {
		points = points[:2]
	}

	m.useBezier = true
	m.bezier = points

	return m
}

func (m *MoveAnimation) Init() {
	m.getState().startPos = imgui.CursorPos()
}

func (m *MoveAnimation) Reset() {
	state := m.getState()
	state.state = !state.state
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

func (m *MoveAnimation) BuildAnimation(animationPercentage, _ float32, starter func()) {
	state := m.getState()

	if !state.state {
		animationPercentage = 1 - animationPercentage
	}

	var pos imgui.Vec2
	if m.useBezier {
		pts := []imgui.Vec2{state.startPos}
		for _, b := range m.bezier {
			pts = append(pts, imgui.Vec2{
				X: b.X + state.startPos.X,
				Y: b.Y + state.startPos.Y,
			})
		}
		pts = append(pts, imgui.Vec2{
			X: state.startPos.X + m.posDelta.X,
			Y: state.startPos.Y + m.posDelta.Y,
		})
		pos = bezier(animationPercentage, pts)
	} else {
		pos = imgui.Vec2{
			X: state.startPos.X + m.posDelta.X*animationPercentage,
			Y: state.startPos.Y + m.posDelta.Y*animationPercentage,
		}
	}

	imgui.SetCursorScreenPos(pos)
	m.widget(starter).Build()
}

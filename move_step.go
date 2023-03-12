package animations

import "github.com/AllenDang/imgui-go"

// MoveStep represents a single key frame in context of MoveAnimation.
// If Relative() not set, positionDelta is relative to this in a previous step.
type MoveStep struct {
	positionDelta imgui.Vec2
	isAbsolute    bool

	useBezier bool
	bezier    []imgui.Vec2
}

// Step creates a new instance of MoveStep.
func Step(x, y float32) *MoveStep {
	return &MoveStep{
		positionDelta: imgui.Vec2{X: x, Y: y},
	}
}

// StepVec acts same as Step but takes imgui.Vec2.
func StepVec(v imgui.Vec2) *MoveStep {
	return &MoveStep{
		positionDelta: v,
	}
}

// Bezier allows to specify Bézier curve points.
// Points are relative to position specified in step.
func (m *MoveStep) Bezier(points ...imgui.Vec2) *MoveStep {
	m.useBezier = true
	m.bezier = points

	return m
}

// Absolute tells animation to take position specified in this step as an absolute
// position rather than relative to he previous step.
func (m *MoveStep) Absolute() *MoveStep {
	m.isAbsolute = true

	return m
}

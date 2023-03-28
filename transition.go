package animations

import (
	"github.com/AllenDang/imgui-go"
)

var _ Animation = &TransitionAnimation{}

// TransitionAnimation is a smooth transition between two renderers.
// It may apply to Windows (giu.WindowWidget) as well as to particular widgets/layouts.
type TransitionAnimation struct {
	renderers []func(starterFunc StarterFunc)
}

// Transition creates a new TransitionAnimation.
func Transition(renderers ...func(starter StarterFunc)) *TransitionAnimation {
	return &TransitionAnimation{
		renderers: renderers,
	}
}

// KeyFramesCount implements Animation interface.
func (t *TransitionAnimation) KeyFramesCount() int {
	return len(t.renderers)
}

// Reset implements Animation interface.
func (t *TransitionAnimation) Reset() {
	// noop
}

// Init implements Animation interface.
func (t *TransitionAnimation) Init() {
	// noop
}

// BuildNormal implements Animation interface.
func (t *TransitionAnimation) BuildNormal(f KeyFrame, starter StarterFunc) {
	t.renderers[f](starter)
}

// BuildAnimation implements Animation interface.
func (t *TransitionAnimation) BuildAnimation(
	percentage, _ float32,
	bf, df KeyFrame,
	_ PlayMode,
	starter StarterFunc,
) {
	layout1 := t.renderers[bf]
	layout2 := t.renderers[df]

	imgui.PushStyleVarFloat(imgui.StyleVarAlpha, percentage)
	layout2(starter)
	imgui.PopStyleVar()

	imgui.PushStyleVarFloat(imgui.StyleVarAlpha, 1-percentage)
	layout1(starter)
	imgui.PopStyleVar()
}

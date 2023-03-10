package animations

import (
	"github.com/AllenDang/imgui-go"
)

var _ Animation = &TransitionAnimation{}

type TransitionAnimation struct {
	renderers []func(starterFunc StarterFunc)
}

func Transition(renderers ...func(starter StarterFunc)) *TransitionAnimation {
	return &TransitionAnimation{
		renderers: renderers,
	}
}

func (t *TransitionAnimation) KeyFramesCount() int {
	return len(t.renderers)
}

func (t *TransitionAnimation) Reset() {
	// noop
}

func (t *TransitionAnimation) Init() {
	// noop
}

func (t *TransitionAnimation) BuildNormal(f KeyFrame, starter StarterFunc) {
	t.renderers[f](starter)
}

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

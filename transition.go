package animations

import (
	"github.com/AllenDang/imgui-go"
)

var _ Animation = &TransitionAnimation{}

type TransitionAnimation struct {
	renderers []func(starter func())
}

func Transition(renderers ...func(starter func())) *TransitionAnimation {
	return &TransitionAnimation{
		renderers: renderers,
	}
}

func (t *TransitionAnimation) KeyFrames() int {
	return 2
}

func (t *TransitionAnimation) BuildAnimation(percentage, _ float32, bf, df KeyFrame, starter StarterFunc) {
	layout1 := t.renderers[bf]
	layout2 := t.renderers[df]
	imgui.PushStyleVarFloat(imgui.StyleVarAlpha, percentage)
	layout2(func() { starter(PlayForward) })
	imgui.PopStyleVar()
	imgui.PushStyleVarFloat(imgui.StyleVarAlpha, 1-percentage)
	layout1(func() { starter(PlayForward) })
	imgui.PopStyleVar()
}

func (t *TransitionAnimation) Reset() {
	// noop
}

func (t *TransitionAnimation) Init() {
	// noop
}

func (t *TransitionAnimation) BuildNormal(f KeyFrame, starter StarterFunc) {
	t.renderers[f](func() { starter(PlayAuto) })
}

package animations

import (
	"image/color"
	"log"
	"sync"

	"github.com/AllenDang/giu"
	"github.com/AllenDang/imgui-go"
)

type colorFlowAnimationState struct {
	state bool
	m     *sync.Mutex
}

func (s *colorFlowAnimationState) Dispose() {
	// noop
}

var _ Animation = &ColorFlowAnimation{}

type ColorFlowAnimation struct {
	id string

	giu.Widget
	destinationColor,
	normalColor func() color.RGBA
	applyingStyles []giu.StyleColorID
}

// ColorFlowStyle wraps ColorFlow so that it automatically obtains the color for specified style values.
func ColorFlowStyle(
	widget giu.Widget,
	destiny, normal giu.StyleColorID,
) *ColorFlowAnimation {
	return ColorFlow(
		widget,
		func() color.RGBA {
			return giu.Vec4ToRGBA(imgui.CurrentStyle().GetColor(imgui.StyleColorID(destiny)))
		},
		func() color.RGBA {
			return giu.Vec4ToRGBA(imgui.CurrentStyle().GetColor(imgui.StyleColorID(normal)))
		},
		destiny, normal,
	)
}

func ColorFlow(
	widget giu.Widget,
	destinationColor, normalColor func() color.RGBA,
	applying ...giu.StyleColorID,
) *ColorFlowAnimation {
	return &ColorFlowAnimation{
		id:               giu.GenAutoID("colorFlowAnimation"),
		Widget:           widget,
		destinationColor: destinationColor,
		normalColor:      normalColor,
		applyingStyles:   applying,
	}
}

func (h *ColorFlowAnimation) Reset() {
	data := h.getState()
	data.state = !data.state
}

func (h *ColorFlowAnimation) Init() {
	// noop
}

func (h *ColorFlowAnimation) BuildNormal(starter StarterFunc) {
	data := h.getState()

	var normalColor color.Color

	if data.state {
		normalColor = h.destinationColor()
	} else {
		normalColor = h.normalColor()
	}

	h.build(normalColor)
}

func (h *ColorFlowAnimation) BuildAnimation(percentage, _ float32, _ StarterFunc) {
	// need to call this method here to prevent state from being disposed.
	_ = h.getState()

	normalColor := giu.ToVec4Color(h.normalColor())
	destinationColor := giu.ToVec4Color(h.destinationColor())

	normalColor.X += (destinationColor.X - normalColor.X) * percentage
	normalColor.X = clamp01(normalColor.X)
	normalColor.Y += (destinationColor.Y - normalColor.Y) * percentage
	normalColor.Y = clamp01(normalColor.Y)
	normalColor.Z += (destinationColor.Z - normalColor.Z) * percentage
	normalColor.Z = clamp01(normalColor.Z)

	h.build(giu.Vec4ToRGBA(normalColor))
}

func (h *ColorFlowAnimation) build(c color.Color) {
	for _, s := range h.applyingStyles {
		giu.PushStyleColor(s, c)
	}

	defer giu.PopStyleColorV(len(h.applyingStyles))

	h.Widget.Build()
}

func (h *ColorFlowAnimation) getState() *colorFlowAnimationState {
	if s := giu.Context.GetState(h.id); s != nil {
		state, ok := s.(*colorFlowAnimationState)
		if !ok {
			log.Panicf("expected state type *colorFlowAnimationState, got %T", s)
		}

		return state
	}

	giu.Context.SetState(h.id, h.newState())

	return h.getState()
}

func (h *ColorFlowAnimation) newState() *colorFlowAnimationState {
	return &colorFlowAnimationState{
		m: &sync.Mutex{},
	}
}

func clamp01(val float32) float32 {
	if val <= 0 {
		return 0
	} else if val >= 1 {
		return 1
	}

	return val
}

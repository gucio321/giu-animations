package animations

import (
	"image/color"
	"log"
	"sync"

	"github.com/AllenDang/giu"
	"github.com/AllenDang/imgui-go"
)

type colorFlowAnimationState struct {
	isTriggered bool
	shouldStart bool
	m           *sync.Mutex
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
	// noop
}

func (h *ColorFlowAnimation) Init() {
	// noop
}

func (h *ColorFlowAnimation) BuildNormal(starter func()) {
	data := h.getState()
	data.m.Lock()
	shouldStart := data.shouldStart
	isTriggered := data.isTriggered
	data.m.Unlock()

	if shouldStart {
		data.m.Lock()
		data.shouldStart = false
		data.m.Unlock()
		starter()
	}

	var normalColor color.Color

	if shouldStart {
		isTriggered = !isTriggered
	}

	if isTriggered {
		normalColor = h.destinationColor()
	} else {
		normalColor = h.normalColor()
	}

	if shouldStart {
		isTriggered = !isTriggered
	}

	h.build(normalColor)
	isTriggeredNow := imgui.IsItemHovered()

	data.m.Lock()
	data.shouldStart = isTriggeredNow != isTriggered
	data.isTriggered = isTriggeredNow
	data.m.Unlock()
}

func (h *ColorFlowAnimation) BuildAnimation(percentage, _ float32, _ func()) {
	data := h.getState()
	normalColor := giu.ToVec4Color(h.normalColor())
	destinationColor := giu.ToVec4Color(h.destinationColor())
	data.m.Lock()
	isTriggered := data.isTriggered
	data.m.Unlock()

	data.m.Lock()
	data.m.Unlock()

	if !isTriggered {
		percentage = 1 - percentage
	}

	normalColor.X += (destinationColor.X - normalColor.X) * percentage
	normalColor.Y += (destinationColor.Y - normalColor.Y) * percentage
	normalColor.Z += (destinationColor.Z - normalColor.Z) * percentage

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

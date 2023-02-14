package animations

import (
	"image/color"
	"log"
	"sync"

	"github.com/AllenDang/giu"
	"github.com/AllenDang/imgui-go"
)

type colorFlowAnimationState struct {
	isHovered   bool
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
	hoveredColor,
	normalColor func() color.RGBA
	hoverID, normalID imgui.StyleColorID
}

// ColorFlowStyle wraps ColorFlow so that it automatically obtains the color for specified style values.
func ColorFlowStyle(
	widget giu.Widget,
	hover, normal giu.StyleColorID,
) *ColorFlowAnimation {
	return ColorFlow(
		widget,
		func() color.RGBA {
			return giu.Vec4ToRGBA(imgui.CurrentStyle().GetColor(imgui.StyleColorID(hover)))
		},
		func() color.RGBA {
			return giu.Vec4ToRGBA(imgui.CurrentStyle().GetColor(imgui.StyleColorID(normal)))
		},
		hover, normal,
	)
}

func ColorFlow(
	widget giu.Widget,
	hoverColor, normalColor func() color.RGBA,
	hoverID, normalID giu.StyleColorID,
) *ColorFlowAnimation {
	return &ColorFlowAnimation{
		id:           giu.GenAutoID("hoverColorAnimation"),
		Widget:       widget,
		hoveredColor: hoverColor,
		normalColor:  normalColor,
		hoverID:      imgui.StyleColorID(hoverID),
		normalID:     imgui.StyleColorID(normalID),
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
	isHovered := data.isHovered
	data.m.Unlock()

	if shouldStart {
		data.m.Lock()
		data.shouldStart = false
		data.m.Unlock()
		starter()
	}

	var normalColor imgui.Vec4

	if shouldStart {
		isHovered = !isHovered
	}

	if isHovered {
		normalColor = giu.ToVec4Color(h.hoveredColor())
	} else {
		normalColor = giu.ToVec4Color(h.normalColor())
	}

	if shouldStart {
		isHovered = !isHovered
	}

	h.build(normalColor)
	isHoveredNow := imgui.IsItemHovered()

	data.m.Lock()
	data.shouldStart = isHoveredNow != isHovered
	data.isHovered = isHoveredNow
	data.m.Unlock()
}

func (h *ColorFlowAnimation) BuildAnimation(percentage, _ float32, _ func()) {
	data := h.getState()
	normalColor := giu.ToVec4Color(h.normalColor())
	hoverColor := giu.ToVec4Color(h.hoveredColor())
	data.m.Lock()
	isHovered := data.isHovered
	data.m.Unlock()

	data.m.Lock()
	data.m.Unlock()

	if !isHovered /*&& state.IsRunning()*/ {
		percentage = 1 - percentage
	}

	normalColor.X += (hoverColor.X - normalColor.X) * percentage
	normalColor.Y += (hoverColor.Y - normalColor.Y) * percentage
	normalColor.Z += (hoverColor.Z - normalColor.Z) * percentage

	h.build(normalColor)
}

func (h *ColorFlowAnimation) build(c imgui.Vec4) {
	imgui.PushStyleColor(h.normalID, c)

	if h.hoverID != h.normalID {
		imgui.PushStyleColor(h.hoverID, c)
		defer imgui.PopStyleColor()
	}

	h.Widget.Build()
	imgui.PopStyleColor()
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

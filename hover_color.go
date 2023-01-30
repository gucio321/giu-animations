package animations

import (
	"github.com/AllenDang/giu"
	"github.com/AllenDang/imgui-go"
	"github.com/gucio321/giu-animations/internal/logger"
	"image/color"
	"sync"
)

type hoverColorAnimationState struct {
	isHovered   bool
	shouldStart bool
	m           *sync.Mutex
}

func (s *hoverColorAnimationState) Dispose() {
	// noop
}

var _ Animation = &HoverColorAnimation{}

type HoverColorAnimation struct {
	id string

	giu.Widget
	hoveredColor,
	normalColor func() color.RGBA
	hoverID, normalID imgui.StyleColorID
}

func HoverColorStyle(
	widget giu.Widget,
	hover, normal giu.StyleColorID,
) *HoverColorAnimation {
	return HoverColor(
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

func HoverColor(
	widget giu.Widget,
	hoverColor, normalColor func() color.RGBA,
	hoverID, normalID giu.StyleColorID,
) *HoverColorAnimation {
	return &HoverColorAnimation{
		id:           giu.GenAutoID("hoverColorAnimation"),
		Widget:       widget,
		hoveredColor: hoverColor,
		normalColor:  normalColor,
		hoverID:      imgui.StyleColorID(hoverID),
		normalID:     imgui.StyleColorID(normalID),
	}
}

func (h *HoverColorAnimation) Reset() {
	// noop
}

func (h *HoverColorAnimation) Init() {
	// noop
}

func (h *HoverColorAnimation) BuildNormal(starter func()) {
	logger.Infof("normal build")
	data := h.getState()
	data.m.Lock()
	shouldStart := data.shouldStart
	isHovered := data.isHovered
	data.m.Unlock()

	logger.Infof("should start %v, is hovered %v", shouldStart, isHovered)

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

func (h *HoverColorAnimation) BuildAnimation(percentage float32, starter func()) {
	logger.Infof("animation build; %v%%", percentage)
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

	logger.Infof("percentage after hover check: %v%%", percentage)

	normalColor.X += (hoverColor.X - normalColor.X) * percentage
	normalColor.Y += (hoverColor.Y - normalColor.Y) * percentage
	normalColor.Z += (hoverColor.Z - normalColor.Z) * percentage

	h.build(normalColor)

}

func (h *HoverColorAnimation) build(c imgui.Vec4) {
	imgui.PushStyleColor(h.normalID, c)

	if h.hoverID != h.normalID {
		imgui.PushStyleColor(h.hoverID, c)
		defer imgui.PopStyleColor()
	}

	h.Widget.Build()
	imgui.PopStyleColor()
}

func (a *HoverColorAnimation) getState() *hoverColorAnimationState {
	if s := giu.Context.GetState(a.id); s != nil {
		state, ok := s.(*hoverColorAnimationState)
		if !ok {
			logger.Fatalf("expected state type *hoverColorAnimationState, got %T", s)
		}

		return state
	}

	giu.Context.SetState(a.id, a.newState())

	return a.getState()
}

func (a *HoverColorAnimation) newState() *hoverColorAnimationState {
	return &hoverColorAnimationState{
		m: &sync.Mutex{},
	}
}

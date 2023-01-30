package animations

import (
	"github.com/AllenDang/giu"
	"github.com/AllenDang/imgui-go"
	"github.com/gucio321/giu-animations/internal/logger"
	"image/color"
	"sync"
)

type hoverColorAnimationState struct {
	procentage  float32
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
	d := h.AnimatorWidget.CustomData()
	if d == nil {
		h.AnimatorWidget.SetCustomData(&hoverColorAnimationState{
			m: &sync.Mutex{},
		})

		return
	}

	currentData, ok := d.(*hoverColorAnimationState)
	if !ok {
		logger.Fatalf("expected data type *hoverColorAnimationState, got %T", d)
	}

	currentData.m.Lock()
	currentData.procentage = 0
	currentData.m.Unlock()
}

func (h *HoverColorAnimation) Advance(procentDelta float32) bool {
	d := h.AnimatorWidget.CustomData()
	if d == nil {
		return true
	}

	data, ok := d.(*hoverColorAnimationState)
	if !ok {
		logger.Fatalf("expected data type *hoverColorAnimationState, got %T", d)
	}

	data.m.Lock()
	data.procentage = procentDelta
	data.m.Unlock()

	return true
}

func (h *HoverColorAnimation) Init() {
	h.AnimatorWidget.SetCustomData(&hoverColorAnimationState{
		m: &sync.Mutex{},
	})
}

func (h *HoverColorAnimation) Build() {
	if h.getState().shouldInit {
		h.Init()
		h.getState().shouldInit = false
	}

	d := h.AnimatorWidget.CustomData()
	data, ok := d.(*hoverColorAnimationState)
	if !ok {
		logger.Fatalf("expected data type *hoverColorAnimationState, got %T", d)
	}

	normalColor := giu.ToVec4Color(h.normalColor())
	hoverColor := giu.ToVec4Color(h.hoveredColor())
	data.m.Lock()
	shouldStart := data.shouldStart
	isHovered := data.isHovered
	data.m.Unlock()

	if shouldStart {
		data.m.Lock()
		data.shouldStart = false
		data.m.Unlock()
		h.Start(h.duration, h.fps)
	}

	data.m.Lock()
	procentage := data.procentage
	data.m.Unlock()

	state := h.AnimatorWidget.getState()

	if !isHovered && state.IsRunning() {
		procentage = 1 - procentage
	}

	normalColor.X += (hoverColor.X - normalColor.X) * procentage
	normalColor.Y += (hoverColor.Y - normalColor.Y) * procentage
	normalColor.Z += (hoverColor.Z - normalColor.Z) * procentage

	if !state.IsRunning() {
		if isHovered {
			normalColor = hoverColor
		} else {
			normalColor = giu.ToVec4Color(h.normalColor())
		}
	}

	imgui.PushStyleColor(h.normalID, normalColor)

	if h.hoverID != h.normalID {
		imgui.PushStyleColor(h.hoverID, normalColor)
		defer imgui.PopStyleColor()
	}

	h.Widget.Build()
	imgui.PopStyleColor()
	isHoveredNow := imgui.IsItemHovered()

	data.m.Lock()
	data.shouldStart = isHoveredNow != isHovered && !state.IsRunning()
	data.isHovered = isHoveredNow
	data.m.Unlock()
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
	return &hoverColorAnimationState{}
}

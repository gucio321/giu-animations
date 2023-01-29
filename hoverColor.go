package animations

import (
	"github.com/AllenDang/giu"
	"github.com/AllenDang/imgui-go"
	"github.com/gucio321/giu-animations/internal/logger"
	"image/color"
	"sync"
	"time"
)

type animationData struct {
	procentage  float32
	isHovered   bool
	shouldStart bool
	m           *sync.Mutex
}

var _ Animation = &HoverColorAnimationWidget{}

type HoverColorAnimationWidget struct {
	*AnimatorWidget
	giu.Widget
	hoveredColor,
	normalColor func() color.RGBA
	fps               int
	duration          time.Duration
	hoverID, normalID imgui.StyleColorID
}

func HoverColorAnimationStyle(widget giu.Widget, fps int, duration time.Duration, hover, normal giu.StyleColorID) *HoverColorAnimationWidget {
	return HoverColorAnimation(
		widget,
		fps, duration,
		func() color.RGBA {
			return giu.Vec4ToRGBA(imgui.CurrentStyle().GetColor(imgui.StyleColorID(hover)))
		},
		func() color.RGBA {
			return giu.Vec4ToRGBA(imgui.CurrentStyle().GetColor(imgui.StyleColorID(normal)))
		},
		hover, normal,
	)
}

func HoverColorAnimation(
	widget giu.Widget,
	fps int, duration time.Duration,
	hoverColor, normalColor func() color.RGBA,
	hoverID, normalID giu.StyleColorID,
) *HoverColorAnimationWidget {
	result := &HoverColorAnimationWidget{
		Widget:       widget,
		hoveredColor: hoverColor,
		normalColor:  normalColor,
		fps:          fps,
		duration:     duration,
		hoverID:      imgui.StyleColorID(hoverID),
		normalID:     imgui.StyleColorID(normalID),
	}

	result.AnimatorWidget = Animator(result, nil, nil)

	return result
}

func (h *HoverColorAnimationWidget) Reset() {
	d := h.AnimatorWidget.GetCustomData()
	if d == nil {
		h.AnimatorWidget.SetCustomData(&animationData{
			m: &sync.Mutex{},
		})

		return
	}

	currentData, ok := d.(*animationData)
	if !ok {
		logger.Fatalf("expected data type *animationData, got %T", d)
	}

	currentData.m.Lock()
	currentData.procentage = 0
	currentData.m.Unlock()
}

func (h *HoverColorAnimationWidget) Advance(procentDelta float32) bool {
	d := h.AnimatorWidget.GetCustomData()
	if d == nil {
		return true
	}

	data, ok := d.(*animationData)
	if !ok {
		logger.Fatalf("expected data type *animationData, got %T", d)
	}

	data.m.Lock()
	data.procentage = procentDelta
	data.m.Unlock()

	return true
}

func (h *HoverColorAnimationWidget) Init() {
	h.AnimatorWidget.SetCustomData(&animationData{
		m: &sync.Mutex{},
	})
}

func (h *HoverColorAnimationWidget) Build() {
	if h.GetState().shouldInit {
		h.Init()
		h.GetState().shouldInit = false
	}

	d := h.AnimatorWidget.GetCustomData()
	data, ok := d.(*animationData)
	if !ok {
		logger.Fatalf("expected data type *animationData, got %T", d)
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

	state := h.AnimatorWidget.GetState()

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

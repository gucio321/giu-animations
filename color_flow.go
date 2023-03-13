package animations

import (
	"image/color"

	"github.com/AllenDang/giu"
	"github.com/AllenDang/imgui-go"
)

var _ Animation = &ColorFlowAnimation{}

type ColorFlowAnimation struct {
	id string

	giu.Widget
	applyingStyles []giu.StyleColorID

	color []func() color.RGBA
}

// ColorFlowStyle wraps ColorFlow so that it automatically obtains the color for specified style values.
func ColorFlowStyle(
	widget giu.Widget,
	normal, destiny giu.StyleColorID,
) *ColorFlowAnimation {
	return ColorFlow(
		widget,
		[]giu.StyleColorID{normal, destiny},
		func() color.RGBA {
			return giu.Vec4ToRGBA(imgui.CurrentStyle().GetColor(imgui.StyleColorID(normal)))
		},
		func() color.RGBA {
			return giu.Vec4ToRGBA(imgui.CurrentStyle().GetColor(imgui.StyleColorID(destiny)))
		},
	)
}

func ColorFlow(
	widget giu.Widget,
	applying []giu.StyleColorID,
	colors ...func() color.RGBA,
) *ColorFlowAnimation {
	return &ColorFlowAnimation{
		id:             giu.GenAutoID("colorFlowAnimation"),
		Widget:         widget,
		applyingStyles: applying,
		color:          colors,
	}
}

func (h *ColorFlowAnimation) Reset() {
	// noop
}

func (h *ColorFlowAnimation) Init() {
	// noop
}

func (h *ColorFlowAnimation) KeyFramesCount() int {
	return len(h.color)
}

func (h *ColorFlowAnimation) BuildNormal(currentKeyFrame KeyFrame, _ StarterFunc) {
	normalColor := h.color[currentKeyFrame]()

	h.build(normalColor)
}

func (h *ColorFlowAnimation) BuildAnimation(
	percentage, _ float32,
	sourceKeyFrame, destinyKeyFrame KeyFrame,
	_ PlayMode,
	_ StarterFunc,
) {
	normalColor := giu.ToVec4Color(h.color[sourceKeyFrame]())
	destinationColor := giu.ToVec4Color(h.color[destinyKeyFrame]())

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

func clamp01(val float32) float32 {
	if val <= 0 {
		return 0
	} else if val >= 1 {
		return 1
	}

	return val
}

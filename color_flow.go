package animations

import (
	"image/color"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/AllenDang/giu"
)

var _ Animation = &ColorFlowAnimation{}

// ColorFlowAnimation makes a smooth flow from one color to another
// on all specified StyleColor variables.
type ColorFlowAnimation struct {
	id giu.ID

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
			return giu.Vec4ToRGBA(*imgui.StyleColorVec4(imgui.Col(normal)))
		},
		func() color.RGBA {
			return giu.Vec4ToRGBA(*imgui.StyleColorVec4(imgui.Col(destiny)))
		},
	)
}

// ColorFlowColors takes a colors list instead of list of functions returning colors.
func ColorFlowColors(
	widget giu.Widget,
	applying []giu.StyleColorID,
	colors ...color.Color,
) *ColorFlowAnimation {
	c := make([]func() color.RGBA, len(colors))
	for i := range c {
		c[i] = func() color.RGBA {
			r, g, b, a := colors[i].RGBA()

			return color.RGBA{
				R: byte(r),
				G: byte(g),
				B: byte(b),
				A: byte(a),
			}
		}
	}

	return ColorFlow(widget, applying, c...)
}

// ColorFlow creates a new ColorFlowAnimation.
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

// Reset implements Animation.
func (c *ColorFlowAnimation) Reset() {
	// noop
}

// Init implements Animation.
func (c *ColorFlowAnimation) Init() {
	// noop
}

// KeyFramesCount implements Animation.
func (c *ColorFlowAnimation) KeyFramesCount() int {
	return len(c.color)
}

// BuildNormal builds animation in normal, not-triggered state.
func (c *ColorFlowAnimation) BuildNormal(currentKeyFrame KeyFrame, _ StarterFunc) {
	normalColor := c.color[currentKeyFrame]()

	c.build(normalColor)
}

// BuildAnimation implements Animation.
func (c *ColorFlowAnimation) BuildAnimation(
	percentage, _ float32,
	sourceKeyFrame, destinyKeyFrame KeyFrame,
	_ PlayMode,
	_ StarterFunc,
) {
	normalColor := giu.ToVec4Color(c.color[sourceKeyFrame]())
	destinationColor := giu.ToVec4Color(c.color[destinyKeyFrame]())

	normalColor.X += (destinationColor.X - normalColor.X) * percentage
	normalColor.X = clamp01(normalColor.X)
	normalColor.Y += (destinationColor.Y - normalColor.Y) * percentage
	normalColor.Y = clamp01(normalColor.Y)
	normalColor.Z += (destinationColor.Z - normalColor.Z) * percentage
	normalColor.Z = clamp01(normalColor.Z)

	c.build(giu.Vec4ToRGBA(normalColor))
}

func (c *ColorFlowAnimation) build(col color.Color) {
	for _, s := range c.applyingStyles {
		giu.PushStyleColor(s, col)
	}

	defer giu.PopStyleColorV(len(c.applyingStyles))

	c.Widget.Build()
}

func clamp01(val float32) float32 {
	if val <= 0 {
		return 0
	} else if val >= 1 {
		return 1
	}

	return val
}

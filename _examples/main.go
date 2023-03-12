package main

import (
	"golang.org/x/image/colornames"
	"image/color"
	"time"

	"github.com/AllenDang/giu"
	"github.com/AllenDang/imgui-go"

	animations "github.com/gucio321/giu-animations"
)

var (
	easingAlg   = animations.EasingAlgNone
	playOnHover bool
)

func loop() {
	a := int32(easingAlg)
	animations.Animator(
		animations.Transition(
			func(starter func()) {
				giu.Window("window1").Layout(
					giu.Label("I'm a window 1"),
					animations.Animator(
						animations.ColorFlow(
							giu.Button("start transition").OnClick(func() {
								starter()
							}),
							[]giu.StyleColorID{
								giu.StyleColorButtonHovered,
								giu.StyleColorButton,
							},
							func() color.RGBA {
								return colornames.Blue
							},
							func() color.RGBA {
								return colornames.Red
							},
							func() color.RGBA {
								return colornames.Green
							},
							func() color.RGBA {
								return colornames.Yellow
							},
						),
					).
						Duration(time.Second).
						FPS(60).
						Trigger(animations.TriggerOnChange, animations.PlayForward, imgui.IsItemHovered),
					giu.Checkbox("Play on hover", &playOnHover),
					animations.Animator(
						animations.Move(func(starter animations.StarterFunc) giu.Widget {
							return giu.Child().Layout(
								giu.Row(
									giu.Label("Set easing alg:"),
									giu.SliderInt(&a, 0, int32(animations.EasingAlgMax-1)).Size(100).OnChange(func() {
										easingAlg = animations.EasingAlgorithmType(a)
									}),
								),
								giu.Button("move me!").OnClick(func() {
									starter(animations.PlayForward)
								}),
							).Size(200, 80)
						}, imgui.Vec2{X: 20, Y: 100}).
							Bezier(imgui.Vec2{X: 20, Y: 20}, imgui.Vec2{X: 90}),
					).Duration(time.Second*3).
						FPS(120).
						EasingAlgorithm(easingAlg).
						Trigger(animations.TriggerOnTrue, animations.PlayForward, func() bool {
							return playOnHover && giu.IsItemHovered()
						}),
				)
			},
			func(starter func()) {
				giu.Window("window2").Layout(
					giu.Label("I'm a window 2"),
					animations.Animator(
						animations.ColorFlowStyle(
							giu.Button("start transition").OnClick(func() {
								starter()
							}),
							giu.StyleColorButton, giu.StyleColorButtonHovered,
						),
					).Trigger(animations.TriggerOnChange, animations.PlayForward, imgui.IsItemHovered),
				)
			},
			func(starter func()) {
				giu.Window("window 3").Layout(
					giu.Label("I'm third window!"),
					giu.Button("start transition").OnClick(func() {
						starter()
					}),
				)
			},
		),
	).Build()
}

func main() {
	wnd := giu.NewMasterWindow("Animations presentation [example]", 640, 480, 0)
	wnd.Run(loop)
}

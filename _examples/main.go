package main

import (
	"image/color"
	"time"

	"golang.org/x/image/colornames"

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
			func(starterFunc animations.StarterFunc) {
				giu.Window("window1").Layout(
					giu.Label("I'm a window 1"),
					animations.Animator(
						animations.ColorFlow(
							giu.Button("start transition").OnClick(func() {
								starterFunc.Start(animations.PlayForward)
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
						Trigger(animations.TriggerOnTrue, animations.PlayForward, imgui.IsItemHovered),
					giu.Button("Play whole transition!").OnClick(func() {
						starterFunc.StartWhole(animations.PlayForward)
					}),
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
								giu.Row(
									giu.Button("play backwards").OnClick(func() {
										starter.Start(animations.PlayBackwards)
									}),
									giu.Button("move me!").OnClick(func() {
										starter.Start(animations.PlayForward)
									}),
								),
							).Size(200, 80)
						},
							animations.Step(20, 100).
								Bezier(imgui.Vec2{X: 20, Y: 20}, imgui.Vec2{X: 90}),
						).DefaultStartPos(),
					).Duration(time.Second*3).
						FPS(120).
						EasingAlgorithm(easingAlg).
						Trigger(animations.TriggerOnTrue, animations.PlayForward, func() bool {
							return playOnHover && giu.IsItemHovered()
						}),
				)
			},
			func(starterFunc animations.StarterFunc) {
				giu.Window("window2").Layout(
					giu.Label("I'm a window 2"),
					animations.Animator(
						animations.ColorFlowStyle(
							giu.Button("start transition").OnClick(func() {
								starterFunc.Start(animations.PlayForward)
							}),
							giu.StyleColorButton, giu.StyleColorButtonHovered,
						),
					).Trigger(animations.TriggerOnChange, animations.PlayForward, imgui.IsItemHovered),
				)
			},
			func(starterFunc animations.StarterFunc) {
				giu.Window("window 3").Layout(
					giu.Label("I'm third window!"),
					giu.Row(
						giu.Button("<< Previous Window").OnClick(func() {
							starterFunc.Start(animations.PlayBackwards)
						}),
						giu.Button("Next Window >>").OnClick(func() {
							starterFunc.Start(animations.PlayForward)
						}),
					),
				)
			},
		),
	).Build()
}

func main() {
	wnd := giu.NewMasterWindow("Animations presentation [example]", 640, 480, 0)
	wnd.Run(loop)
}

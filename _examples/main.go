package main

import (
	"time"

	"golang.org/x/image/colornames"

	"github.com/AllenDang/giu"
	"github.com/AllenDang/imgui-go"

	animations "github.com/gucio321/giu-animations/v2"
)

var (
	easingAlg     = animations.EasingAlgNone
	moveOnHover   bool
	resizeOnHover bool
)

func loop() {
	a := int32(easingAlg)
	animations.Animator(
		animations.Transition(
			func(starterFunc animations.StarterFunc) {
				giu.Window("window1").Layout(
					giu.Label("I'm a window 1"),
					animations.Animator(
						animations.ColorFlowColors(
							giu.Button("start transition").OnClick(func() {
								starterFunc.Start(animations.PlayForward)
							}),
							[]giu.StyleColorID{
								giu.StyleColorButtonHovered,
								giu.StyleColorButton,
							},
							colornames.Blue,
							colornames.Red,
							colornames.Green,
							colornames.Yellow,
						),
					).
						ID("colorful-button").
						Duration(time.Second).
						FPS(60).
						Trigger(animations.TriggerOnTrue, animations.PlayForward, imgui.IsItemHovered),
					giu.Row(
						giu.Button("<<").OnClick(func() {
							starterFunc.StartCycle(1, animations.PlayBackward)
						}),
						giu.Label("Play whole transition"),
						giu.Button(">>").OnClick(func() {
							starterFunc.StartCycle(1, animations.PlayForward)
						}),
					),
					giu.Checkbox("Move on hover", &moveOnHover),
					giu.Checkbox("Resize on hover", &resizeOnHover),
					animations.Animator(
						animations.Move(func(starter animations.StarterFunc) giu.Widget {
							return animations.Animator(animations.Resize[*giu.ChildWidget](
								giu.Child().Layout(
									giu.Row(
										giu.Label("Set easing alg:"),
										giu.SliderInt(&a, 0, int32(animations.EasingAlgMax-1)).Size(100).OnChange(func() {
											easingAlg = animations.EasingAlgorithmType(a)
										}),
									),
									giu.Row(
										giu.Button("play backwards").OnClick(func() {
											starter.Start(animations.PlayBackward)
										}),
										giu.Button("move me!").OnClick(func() {
											starter.Start(animations.PlayForward)
										}),
									),
								),
								imgui.Vec2{X: 200, Y: 80},
								imgui.Vec2{X: 250, Y: 130},
							).TrickCursor(animations.TrickNever)).Trigger(animations.TriggerOnChange, animations.PlayForward, func() bool {
								return imgui.IsItemHovered() && resizeOnHover
							})
						},
							animations.Step(20, 100).
								Bezier(imgui.Vec2{X: 20, Y: 20}, imgui.Vec2{X: 90}),
						).DefaultStartPos(),
					).Duration(time.Second*3).
						FPS(120).
						EasingAlgorithm(easingAlg).
						Trigger(animations.TriggerOnTrue, animations.PlayForward, func() bool {
							return moveOnHover && giu.IsItemHovered()
						}),
					animations.Animator(animations.Resize[*giu.ButtonWidget](
						giu.Button("Resize me!"),
						imgui.Vec2{X: 150, Y: 150},
						imgui.Vec2{X: 200, Y: 200},
						imgui.Vec2{X: 250, Y: 250},
						imgui.Vec2{X: 300, Y: 300},
					).TrickCursor(animations.TrickCursorAlways)).Trigger(animations.TriggerOnChange, animations.PlayForward, imgui.IsItemHovered).
						EasingAlgorithm(animations.EasingAlgOutBounce),
					giu.Button("test"),
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
					).
						ID("button-with-animated-default-color").
						Trigger(animations.TriggerOnChange, animations.PlayForward, imgui.IsItemHovered),
					giu.Row(
						giu.Button("<<").OnClick(func() {
							starterFunc.StartCycle(1, animations.PlayBackward)
						}),
						giu.Label("Play whole transition"),
						giu.Button(">>").OnClick(func() {
							starterFunc.StartCycle(1, animations.PlayForward)
						}),
					),
				)
			},
			func(starterFunc animations.StarterFunc) {
				giu.Window("window 3").Layout(
					giu.Label("I'm third window!"),
					giu.Row(
						giu.Button("<< Previous Window").OnClick(func() {
							starterFunc.Start(animations.PlayBackward)
						}),
						giu.Button("Next Window >>").OnClick(func() {
							starterFunc.Start(animations.PlayForward)
						}),
					),
					giu.Row(
						giu.Button("<<").OnClick(func() {
							starterFunc.StartCycle(1, animations.PlayBackward)
						}),
						giu.Label("Play whole transition"),
						giu.Button(">>").OnClick(func() {
							starterFunc.StartCycle(1, animations.PlayForward)
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

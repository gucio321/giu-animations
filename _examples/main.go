package main

import (
	"github.com/AllenDang/giu"
	"github.com/AllenDang/imgui-go"
	animations "github.com/gucio321/giu-animations"
	"golang.org/x/image/colornames"
	"image/color"
	"time"
)

func loop() {
	animations.Animator(
		animations.Transition(
			func(starter func()) {
				giu.Window("window1").Layout(
					giu.Label("I'm a window 1"),
					animations.Animator(
						animations.HoverColor(
							giu.Button("start transition").OnClick(func() {
								starter()
							}),
							func() color.RGBA {
								return colornames.Red
							},
							func() color.RGBA {
								return colornames.Blue
							},
							giu.StyleColorButtonHovered,
							giu.StyleColorButton,
						),
					).Duration(time.Second).FPS(60),
					animations.Animator(
						animations.Move(func(starter func()) giu.Widget {
							return giu.Button("move me!").OnClick(func() {
								starter()
							})
						}, imgui.Vec2{X: 20, Y: 20}),
					).Duration(time.Second).FPS(60),
				)
			},
			func(starter func()) {
				giu.Window("window2").Layout(
					giu.Label("I'm a window 1"),
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

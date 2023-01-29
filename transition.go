/*
 * Copyright (c) 2023 The Greater Heptavirate team (https://github.com/TheGreaterHeptavirate)
 * All Rights Reserved
 *
 * All copies of this software (if not stated otherway) are dedicated
 * ONLY to personal, non-commercial use.
 */

package animations

import (
	"github.com/AllenDang/imgui-go"
	"github.com/gucio321/giu-animations/internal/logger"
	"time"
)

type TransitionWidget struct {
	a *animationWidget
}

func Transition(renderer1, renderer2 func(Animation)) *TransitionWidget {
	result := &TransitionWidget{}
	result.a = newAnimation(result, renderer1, renderer2)
	return result
}

func (t *TransitionWidget) Start(d time.Duration, fps int) {
	t.a.Start(d, fps)
}

func (t *TransitionWidget) Advance(procentDelta float32) bool {
	state := t.a.GetState()
	// it means the current layou is layout1, so increasing procentage
	if state.currentLayout {
		t.a.SetCustomData(procentDelta)
	} else {
		t.a.SetCustomData(1 - procentDelta)
	}

	return true
}

func (t *TransitionWidget) Reset() {
	state := t.a.GetState()
	if state.currentLayout {
		t.a.SetCustomData(float32(0))
	} else {
		t.a.SetCustomData(float32(1))
	}
}

func (t *TransitionWidget) Init() {
	// noop
}

func (t *TransitionWidget) Build() {
	if t.a.BuildNormal(t) {
		return
	}

	d := t.a.GetCustomData()
	layout1ProcentageAlpha, ok := d.(float32)
	if !ok {
		logger.Fatal("invalid custom data type: wanted float32 got %t", d)
	}
	if layout1ProcentageAlpha > 1 {
		logger.Fatalf("animationWidget: procentage alpha is %v (should be in range 0-1)", layout1ProcentageAlpha)
	}

	imgui.PushStyleVarFloat(imgui.StyleVarAlpha, layout1ProcentageAlpha)
	t.a.renderer1(t)
	imgui.PopStyleVar()
	imgui.PushStyleVarFloat(imgui.StyleVarAlpha, 1-layout1ProcentageAlpha)
	t.a.renderer2(t)
	imgui.PopStyleVar()
}

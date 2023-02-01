// Package animations contains my attempt to create a kind of "animations" in imgui.
package animations

import (
	"time"

	"github.com/AllenDang/giu"
)

const (
	DefaultFPS      = 60
	DefaultDuration = time.Second / 4
)

// AnimatorWidget is a manager for animation.
type AnimatorWidget struct {
	id string

	duration time.Duration
	fps      int

	a Animation
}

// Animator creates a new AnimatorWidget.
func Animator(a Animation) *AnimatorWidget {
	result := &AnimatorWidget{
		id:       giu.GenAutoID("Animation"),
		a:        a,
		duration: DefaultDuration,
		fps:      DefaultFPS,
	}

	return result
}

// FPS allows to specify FPS value.
// CAUTION: it will take effect after next call to Start - not applied to currently plaid animation.
func (a *AnimatorWidget) FPS(fps int) *AnimatorWidget {
	a.fps = fps

	return a
}

// Duration allows to specify duration value.
// CAUTION: it will take effect after next call to Start - not applied to currently plaid animation.
func (a *AnimatorWidget) Duration(duration time.Duration) *AnimatorWidget {
	a.duration = duration

	return a
}

// Start starts the animation.
func (a *AnimatorWidget) Start() {
	a.a.Reset()
	state := a.getState()

	state.m.Lock()
	state.isRunning = true
	state.duration = a.duration
	state.m.Unlock()

	go func() {
		tickDuration := time.Second / time.Duration(a.fps)
		for range time.Tick(tickDuration) {
			if state.elapsed > state.duration {
				state.m.Lock()
				state.isRunning = false
				state.elapsed = 0
				state.m.Unlock()

				return
			}

			giu.Update()

			state.m.Lock()
			state.elapsed += tickDuration
			state.m.Unlock()
		}
	}()
}

// Build implements giu.Widget
func (a *AnimatorWidget) Build() {
	if a.shouldInit() {
		a.a.Init()
		s := a.getState()
		s.m.Lock()
		s.shouldInit = false
		s.m.Unlock()
	}

	if a.IsRunning() {
		a.a.BuildAnimation(a.CurrentPercentageProgress(), a.Start)

		return
	}

	a.a.BuildNormal(a.Start)
}

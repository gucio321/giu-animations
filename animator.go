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
// It is a giu.Widget, but you should reather store it in a variable
// in order to be able to call Start.
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
func (t *AnimatorWidget) FPS(fps int) *AnimatorWidget {
	t.fps = fps

	return t
}

// Duration allows to specify duration value.
// CAUTION: it will take effect after next call to Start - not applied to currently plaid animation.
func (t *AnimatorWidget) Duration(duration time.Duration) *AnimatorWidget {
	t.duration = duration

	return t
}

// Start starts the animation.
func (t *AnimatorWidget) Start() {
	t.a.Reset()
	state := t.getState()

	state.m.Lock()
	state.isRunning = true
	state.duration = t.duration
	state.m.Unlock()

	go func() {
		tickDuration := time.Second / time.Duration(t.fps)
		for range time.Tick(tickDuration) {
			if state.elapsed > state.duration {
				state.m.Lock()
				state.isRunning = false
				state.elapsed = 0
				state.currentLayout = !state.currentLayout
				state.m.Unlock()

				return
			}

			giu.Update()
			state.elapsed += tickDuration
		}
	}()
}

// Build implements giu.Widget
func (t *AnimatorWidget) Build() {
	if t.shouldInit() {
		t.a.Init()
		s := t.getState()
		s.m.Lock()
		s.shouldInit = false
		s.m.Unlock()
	}

	if t.IsRunning() {
		t.a.BuildAnimation(t.CurrentPercentageProgress(), t.Start)

		return
	}

	t.a.BuildNormal(t.Start)
}

// Package animations contains my attempt to create a kind of "animations" in imgui.
package animations

import (
	"time"

	"github.com/AllenDang/giu"
)

const (
	// DefaultFPS is FPS value that should suit most use-cases.
	// Animator takes this value by default and it could be changed by (*Animator).FPS()
	DefaultFPS = 60
	// DefaultDuration is animation's duration set by default.
	// You can change this by (*Animator).Durations()
	DefaultDuration = time.Second / 4
)

// AnimatorWidget is a manager for animation.
type AnimatorWidget struct {
	id string

	duration        time.Duration
	fps             int
	easingAlgorithm EasingAlgorithmType
	triggerType     TriggerType
	triggerFunc     TriggerFunc

	a Animation
}

// Animator creates a new AnimatorWidget.
func Animator(a Animation) *AnimatorWidget {
	result := &AnimatorWidget{
		id:              giu.GenAutoID("Animation"),
		a:               a,
		duration:        DefaultDuration,
		fps:             DefaultFPS,
		easingAlgorithm: EasingAlgNone,
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

// EasingAlgorithm allows to specify easing algorithm
func (a *AnimatorWidget) EasingAlgorithm(alg EasingAlgorithmType) *AnimatorWidget {
	a.easingAlgorithm = alg

	return a
}

func (a *AnimatorWidget) Trigger(triggerType TriggerType, f TriggerFunc) *AnimatorWidget {
	a.triggerType = triggerType
	a.triggerFunc = f

	return a
}

// Start starts the animation.
func (a *AnimatorWidget) Start(playMode PlayMode) {
	a.a.Reset()
	state := a.getState()

	state.m.Lock()
	state.isRunning = true
	state.duration = a.duration
	state.m.Unlock()

	go a.playAnimation(playMode)
}

func (a *AnimatorWidget) playAnimation(playMode PlayMode) {
	state := a.getState()
	{
		state.m.Lock()
		switch playMode {
		case PlayForward:
			state.elapsed = 0
		case PlayBackwards:
			state.elapsed = state.duration
		}
		state.m.Unlock()
	}

	tickDuration := time.Second / time.Duration(a.fps)
	for range time.Tick(tickDuration) {
		giu.Update()
		state.m.Lock()
		switch playMode {
		case PlayForward:
			if state.elapsed >= state.duration {
				state.isRunning = false
				state.elapsed = 0
				state.m.Unlock()

				return
			}

			state.elapsed += tickDuration
		case PlayBackwards:
			if state.elapsed <= 0 {
				state.isRunning = false
				state.elapsed = 0
				state.m.Unlock()

				return
			}

			state.elapsed -= tickDuration
		}

		state.m.Unlock()
	}
}

// Build implements giu.Widget
func (a *AnimatorWidget) Build() {
	s := a.getState()
	if a.shouldInit() {
		a.a.Init()
		s.m.Lock()
		s.shouldInit = false
		s.m.Unlock()
	}

	if a.IsRunning() {
		p := a.CurrentPercentageProgress()
		a.a.BuildAnimation(Ease(a.easingAlgorithm, p), p, a.Start)

		return
	}

	a.a.BuildNormal(a.Start)

	if a.triggerFunc != nil {
		triggerValue := a.triggerFunc()
		switch a.triggerType {
		case TriggerNever:
			// noop
		case TriggerOnTrue:
			if triggerValue {
				a.Start(PlayForward)
			}
		case TriggerOnChange:
			s.m.Lock()
			if s.triggerStatus != triggerValue {
				if triggerValue {
					a.Start(PlayForward)
				} else {
					a.Start(PlayBackwards)
				}
			}

			s.triggerStatus = triggerValue
			s.m.Unlock()
		}
	}
}

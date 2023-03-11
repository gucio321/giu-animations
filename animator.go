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

	duration time.Duration
	fps      int

	easingAlgorithm EasingAlgorithmType

	triggerType TriggerType
	triggerFunc TriggerFunc

	numKeyFrames int

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
		numKeyFrames:    a.KeyFrames(),
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
	state := a.getState()
	state.m.Lock()
	cf := state.currentKeyFrame
	state.m.Unlock()
	delta := 1
	if playMode == PlayBackwards {
		delta = -1
	}

	a.StartKeyFrames(cf, getWithDelta(cf, a.numKeyFrames, delta), playMode)
}

func (a *AnimatorWidget) StartKeyFrames(beginKF, destinationKF KeyFrame, playMode PlayMode) {
	state := a.getState()
	state.m.Lock()
	state.currentKeyFrame = beginKF
	state.longTimeDestinationKeyFrame = destinationKF
	switch playMode {
	case PlayForward:
		state.destinationKeyFrame = getWithDelta(beginKF, a.numKeyFrames, 1)
	case PlayBackwards:
		state.destinationKeyFrame = getWithDelta(beginKF, a.numKeyFrames, -1)
	}
	state.m.Unlock()

	a.start(playMode)
}

func (a *AnimatorWidget) StartWhole(playMode PlayMode) {
	// TODO: set state here
	var begin, end = 0, a.numKeyFrames - 1
	if playMode == PlayBackwards {
		begin, end = end, begin
	}

	a.StartKeyFrames(KeyFrame(begin), KeyFrame(end), playMode)
}

func (a *AnimatorWidget) start(playMode PlayMode) {
	a.a.Reset()
	state := a.getState()

	state.m.Lock()

	if state.isRunning {
		state.reset <- true
	}

	state.isRunning = true
	state.duration = a.duration

	state.m.Unlock()

	go a.playAnimation(playMode)
}

func (a *AnimatorWidget) playAnimation(playMode PlayMode) {
	state := a.getState()
	state.m.Lock()

	switch playMode {
	case PlayForward:
		state.elapsed = 0
	case PlayBackwards:
		state.elapsed = state.duration
	}

	resetChan := state.reset
	state.m.Unlock()

	tickDuration := time.Second / time.Duration(a.fps)
	for {
		state.m.Lock()
		if state.currentKeyFrame == state.longTimeDestinationKeyFrame {
			state.m.Unlock()
			break
		}

		state.m.Unlock()

	AnimationLoop:
		for {
			select {
			case <-time.Tick(tickDuration):
				giu.Update()
				state.m.Lock()
				switch playMode {
				case PlayForward:
					if state.elapsed >= state.duration {
						a.stopAnimation(state)
						state.currentKeyFrame = getWithDelta(state.currentKeyFrame, a.numKeyFrames, 1)
						state.destinationKeyFrame = getWithDelta(state.currentKeyFrame, a.numKeyFrames, 1)
						state.m.Unlock()

						break AnimationLoop
					}

					state.elapsed += tickDuration
				case PlayBackwards:
					if state.elapsed <= 0 {
						state.currentKeyFrame = getWithDelta(state.currentKeyFrame, a.numKeyFrames, -1)
						state.destinationKeyFrame = getWithDelta(state.currentKeyFrame, a.numKeyFrames, -1)
						a.stopAnimation(state)
						state.m.Unlock()

						break AnimationLoop
					}

					state.elapsed -= tickDuration
				}

				state.m.Unlock()
			case <-resetChan:
				return
			}
		}
	}
}

func (a *AnimatorWidget) stopAnimation(state *animatorState) {
	state.isRunning = false
	state.elapsed = 0

	// call update last time to build animation normally at least once (before Power Saving Mechanism freezes updating)
	// This is important mainly because of triggers that might have to be run.
	giu.Update()
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

	s.m.Lock()
	cf, df := s.currentKeyFrame, s.destinationKeyFrame
	s.m.Unlock()

	if a.IsRunning() {
		p := a.CurrentPercentageProgress()
		a.a.BuildAnimation(Ease(a.easingAlgorithm, p), p, cf, df, a.Start)

		return
	}

	a.a.BuildNormal(cf, a.Start)

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
			triggerStatus := s.triggerStatus
			s.m.Unlock()

			if triggerStatus != triggerValue {
				if triggerValue {
					a.Start(PlayForward)
				} else {
					a.Start(PlayBackwards)
				}
			}

			s.m.Lock()
			s.triggerStatus = triggerValue
			s.m.Unlock()
		}
	}
}

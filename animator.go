// Package animations contains my attempt to create animation kind of "animations" in imgui.
package animations

import (
	"time"

	"github.com/AllenDang/giu"
)

const (
	// DefaultFPS is FPS value that should suit most use-cases.
	// Animator takes this value by default and it could be changed by (*Animator).FPS().
	DefaultFPS = 60
	// DefaultDuration is animation's duration set by default.
	// You can change this by (*Animator).Durations().
	DefaultDuration = time.Second / 4
)

var (
	_ giu.Widget  = &AnimatorWidget{}
	_ StarterFunc = &AnimatorWidget{}
)

// AnimatorWidget is animation manager for Animation.
// It is a giu.Widget (so you can use it in any giu.Layout).
// This type provides a wide API that allows you to manage your animation
// such as Start* functions or parameters like FPS or Duration.
// It is actually responsible for advancement of animation.
// NOTE: This type should be concurrent-safe. If you find any data race
// please open an issue in source repository.
type AnimatorWidget struct {
	// giu ID used for storing internal state
	id string

	// playback properties
	duration time.Duration
	fps      int

	easingAlgorithm EasingAlgorithmType

	// triggers
	triggerType    TriggerType
	triggerPlyMode PlayMode
	triggerFunc    TriggerFunc

	animation    Animation
	numKeyFrames int // <- Filled in in Animator call. SHOULD NOT CHANGE after it.
}

// Animator creates animation new AnimatorWidget.
func Animator(a Animation) *AnimatorWidget {
	result := &AnimatorWidget{
		id:              giu.GenAutoID("Animation"),
		animation:       a,
		duration:        DefaultDuration,
		fps:             DefaultFPS,
		easingAlgorithm: EasingAlgNone,
		numKeyFrames:    a.KeyFramesCount(),
	}

	return result
}

// ID sets animation custom ID to this AnimatorWidget
// It may be really important when using TransitionAnimation, because
// sometimes when using sub-animators inside of Transition, it may happen
// that the second AnimatorWidget will receive the same ID as the previous one.
// It may cause unexpected behaviors.
func (a *AnimatorWidget) ID(newID string) *AnimatorWidget {
	a.id = newID

	return a
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

// EasingAlgorithm allows to specify easing algorithm.
func (a *AnimatorWidget) EasingAlgorithm(alg EasingAlgorithmType) *AnimatorWidget {
	a.easingAlgorithm = alg

	return a
}

// Trigger sets automatic triggering of animation.
//
//	Example: (*AnimatorWidget).Trigger(TriggerOnChange, imgui.IsItemHovered)
func (a *AnimatorWidget) Trigger(triggerType TriggerType, playMode PlayMode, f TriggerFunc) *AnimatorWidget {
	a.triggerType = triggerType
	a.triggerPlyMode = playMode
	a.triggerFunc = f

	return a
}

// Start starts the animation.
// It plays one single frame forwards/backwards (depending on playMode).
func (a *AnimatorWidget) Start(playMode PlayMode) {
	state := a.getState()
	state.m.Lock()
	cf := state.currentKeyFrame
	state.m.Unlock()

	delta := 1

	if playMode == PlayBackward {
		delta = -1
	}

	destinationFrame := getWithDelta(cf, a.numKeyFrames, delta)
	a.StartKeyFrames(cf, destinationFrame, 0, playMode)
}

// StartKeyFrames initializes animation playback from beginKF to destination KF in direction
// specified by playMode.
func (a *AnimatorWidget) StartKeyFrames(beginKF, destinationKF KeyFrame, cyclesCount int, playMode PlayMode) {
	state := a.getState()

	state.m.Lock()
	state.currentKeyFrame = beginKF
	state.longTimeDestinationKeyFrame = destinationKF

	switch playMode {
	case PlayForward:
		state.destinationKeyFrame = getWithDelta(beginKF, a.numKeyFrames, 1)
	case PlayBackward:
		state.destinationKeyFrame = getWithDelta(beginKF, a.numKeyFrames, -1)
	}

	state.numberOfCycles = cyclesCount

	state.m.Unlock()

	a.start(playMode)
}

// StartCycle plays an animation from start to end (optionally from end to start).
func (a *AnimatorWidget) StartCycle(numberOfCycles int, playMode PlayMode) {
	state := a.getState()
	b := state.currentKeyFrame
	a.StartKeyFrames(b, b, numberOfCycles, playMode)
}

// internal start method. Stops animator if running and re-initializes it.
// It will call playAnimation in a new goroutine.
func (a *AnimatorWidget) start(playMode PlayMode) {
	a.animation.Reset()
	state := a.getState()

	state.m.Lock()

	if state.isRunning {
		state.reset <- true
	}

	state.isRunning = true
	state.duration = a.duration

	state.playMode = playMode

	state.m.Unlock()

	go a.playAnimation(playMode)
}

// playAnimation is where the animation is plaid.
// It runs a for loop through all the frames that it should go.
// It will exit if any message received on state.reset.
func (a *AnimatorWidget) playAnimation(playMode PlayMode) {
	state := a.getState()
	state.m.Lock()

	resetChan := state.reset
	state.m.Unlock()

	for {
		state.m.Lock()
		state.elapsed = 0

		if state.currentKeyFrame == state.longTimeDestinationKeyFrame {
			if state.numberOfCycles == 0 {
				state.m.Unlock()

				break
			}

			state.numberOfCycles--
		}

		state.m.Unlock()

		tickDuration := time.Second / time.Duration(a.fps)
		ticker := time.NewTicker(tickDuration)
	AnimationLoop:
		for {
			select {
			case <-ticker.C:
				giu.Update()
				state.m.Lock()
				if state.elapsed >= state.duration {
					ticker.Stop()

					state.elapsed = 0

					// call update last time to build animation normally at least once (before Power Saving Mechanism freezes updating)
					// This is important mainly because of triggers that might have to be run.
					giu.Update()

					delta := 1
					if playMode == PlayBackward {
						delta = -1
					}

					state.currentKeyFrame = getWithDelta(state.currentKeyFrame, a.numKeyFrames, delta)
					state.destinationKeyFrame = getWithDelta(state.currentKeyFrame, a.numKeyFrames, delta)

					state.m.Unlock()

					break AnimationLoop
				}

				state.elapsed += tickDuration

				state.m.Unlock()
			case <-resetChan:
				return
			}
		}
	}

	state.m.Lock()
	state.isRunning = false
	state.m.Unlock()
}

// Build implements giu.Widget.
func (a *AnimatorWidget) Build() {
	s := a.getState()

	if a.shouldInit() {
		a.animation.Init()
		s.m.Lock()
		s.shouldInit = false
		s.m.Unlock()
	}

	s.m.Lock()
	cf, df := s.currentKeyFrame, s.destinationKeyFrame
	playMode := s.playMode
	s.m.Unlock()

	if a.IsRunning() {
		p := a.CurrentPercentageProgress()
		a.animation.BuildAnimation(
			Ease(a.easingAlgorithm, p), p,
			cf, df,
			playMode,
			a,
		)

		return
	}

	a.animation.BuildNormal(cf, a)

	if a.triggerFunc != nil {
		triggerValue := a.triggerFunc()

		switch a.triggerType {
		case TriggerNever:
			// noop
		case TriggerOnTrue:
			if triggerValue {
				a.Start(a.triggerPlyMode)
			}
		case TriggerOnChange:
			s.m.Lock()
			triggerStatus := s.triggerStatus
			s.m.Unlock()

			if triggerStatus != triggerValue {
				a.Start(a.triggerPlyMode)
			}

			s.m.Lock()
			s.triggerStatus = triggerValue
			s.m.Unlock()
		}
	}
}

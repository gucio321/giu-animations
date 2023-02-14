package animations

import (
	"log"
	"sync"
	"time"

	"github.com/AllenDang/giu"
)

var _ giu.Disposable = &animatorState{}

type animatorState struct {
	shouldInit bool
	isRunning  bool

	elapsed  time.Duration
	duration time.Duration

	triggerStatus bool

	m *sync.Mutex
}

// Dispose implements giu.Disposable
func (s *animatorState) Dispose() {
	// noop
}

func (a *AnimatorWidget) newState() *animatorState {
	return &animatorState{
		shouldInit: true,
		m:          &sync.Mutex{},
	}
}

// getState returns animator's state.
// It could not be public, because of concurrency issues.
// There is a bunch of Animator's methods that allows
// user to obtain certain data.
func (a *AnimatorWidget) getState() *animatorState {
	if s := giu.Context.GetState(a.id); s != nil {
		state, ok := s.(*animatorState)
		if !ok {
			log.Panicf("error asserting type of animator state: got %T, wanted *animatorState", s)
		}

		return state
	}

	giu.Context.SetState(a.id, a.newState())

	return a.getState()
}

// IsRunning returns true if the animation is already running.
func (a *AnimatorWidget) IsRunning() bool {
	s := a.getState()
	s.m.Lock()
	defer s.m.Unlock()
	return s.isRunning
}

func (a *AnimatorWidget) shouldInit() bool {
	s := a.getState()
	s.m.Lock()
	defer s.m.Unlock()

	return s.shouldInit
}

// CurrentPercentageProgress returns a float value from range <0, 1>
// representing current progress of an animation.
// If animation is not running, it will return 0.
func (a *AnimatorWidget) CurrentPercentageProgress() float32 {
	if !a.IsRunning() {
		return 0
	}

	s := a.getState()
	s.m.Lock()
	defer s.m.Unlock()

	result := float32(s.elapsed) / float32(s.duration)
	if result > 1 {
		return 1
	}

	return result
}

package animations

import (
	"github.com/gucio321/giu-animations/internal/logger"
	"sync"
	"time"

	"github.com/AllenDang/giu"
)

var _ giu.Disposable = &animatorState{}

type animatorState struct {
	isRunning     bool
	currentLayout bool
	elapsed       time.Duration
	duration      time.Duration
	customData    any
	shouldInit    bool
	m             *sync.Mutex
}

func (s *animatorState) Dispose() {
	// noop
}

func (t *AnimatorWidget) newState() *animatorState {
	return &animatorState{
		shouldInit: true,
		m:          &sync.Mutex{},
	}
}

// getState returns animator's state.
// It could not be public, because of concurrency issues.
// There is a bunch of Animator's methods that allows
// user to obtain certain data.
func (t *AnimatorWidget) getState() *animatorState {
	if s := giu.Context.GetState(t.id); s != nil {
		state, ok := s.(*animatorState)
		if !ok {
			logger.Fatalf("error asserting type of ttransition state: got %T", s)
		}

		return state
	}

	giu.Context.SetState(t.id, t.newState())

	return t.getState()
}

// CustomData returns previously saved custom data or nil.
// it is concurrency-safe
func (t *AnimatorWidget) CustomData() any {
	state := t.getState()

	state.m.Lock()
	defer state.m.Unlock()

	return state.customData
}

// CustomData is a generic custom data getter.
func CustomData[T any](a *AnimatorWidget) (out T, err error) {
	d := a.CustomData()
	if d == nil {
		return out, nil
	}

	out, ok := d.(T)
	if !ok {
		return out, ErrInvalidDataType
	}

	return out, nil
}

// SetCustomData sets custom data.
func (t *AnimatorWidget) SetCustomData(d any) {
	state := t.getState()
	state.m.Lock()
	state.customData = d
	state.m.Unlock()
}

// IsRunning returns true if the animation is already running.
func (s *animatorState) IsRunning() bool {
	s.m.Lock()
	defer s.m.Unlock()
	return s.isRunning
}

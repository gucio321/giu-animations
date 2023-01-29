// Package animations contains my attempt to create a kind of "animations" in imgui.
package animations

import (
	"github.com/gucio321/giu-animations/internal/logger"
	"sync"
	"time"

	"github.com/AllenDang/giu"
)

func (t *AnimatorWidget) GetState() *animatorState {
	if s := giu.Context.GetState(t.id); s != nil {
		state, ok := s.(*animatorState)
		if !ok {
			logger.Fatalf("error asserting type of ttransition state: got %T", s)
		}

		return state
	}

	giu.Context.SetState(t.id, t.newState())

	return t.GetState()
}

func (t *AnimatorWidget) newState() *animatorState {
	return &animatorState{
		shouldInit: true,
		m:          &sync.Mutex{},
	}
}

type AnimatorWidget struct {
	id                   string
	renderer1, renderer2 func(this Animation)
	Animation
}

func Animator(a Animation, renderer1, renderer2 func(this Animation)) *AnimatorWidget {
	result := &AnimatorWidget{
		id:        giu.GenAutoID("Animation"),
		renderer1: renderer1,
		renderer2: renderer2,
		Animation: a,
	}

	return result
}

func (t *AnimatorWidget) Start(duration time.Duration, fps int) {
	t.Animation.Reset()
	state := t.GetState()

	if state.isRunning {
		logger.Fatal("AnimatorWidget: StartTransition called, but transition is already running")
	}

	state.isRunning = true
	state.duration = duration

	go func() {
		tickDuration := time.Second / time.Duration(fps)
		for range time.Tick(tickDuration) {
			if state.elapsed > state.duration {
				state.m.Lock()
				state.isRunning = false
				state.elapsed = 0
				state.currentLayout = !state.currentLayout
				state.m.Unlock()

				return
			}

			procentDelta := float32(state.elapsed) / float32(state.duration)

			if !t.Advance(procentDelta) {
				return
			}

			giu.Update()
			state.elapsed += tickDuration
		}
	}()
}

func (t *AnimatorWidget) GetCustomData() any {
	state := t.GetState()
	state.m.Lock()
	defer state.m.Unlock()
	return state.customData
}

func (t *AnimatorWidget) SetCustomData(d any) {
	state := t.GetState()
	state.m.Lock()
	state.customData = d
	state.m.Unlock()
}

func (t *AnimatorWidget) BuildNormal(a Animation) (proceeded bool) {
	state := t.GetState()

	if !state.IsRunning() {
		if !state.currentLayout {
			t.renderer1(a)
		} else {
			t.renderer2(a)
		}

		return true
	}

	return false
}

func (s *animatorState) IsRunning() bool {
	s.m.Lock()
	defer s.m.Unlock()
	return s.isRunning
}

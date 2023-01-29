/*
 * Copyright (c) 2023 The Greater Heptavirate team (https://github.com/TheGreaterHeptavirate)
 * All Rights Reserved
 *
 * All copies of this software (if not stated otherway) are dedicated
 * ONLY to personal, non-commercial use.
 */

// Package animations contains my attempt to create a kind of "animations" in imgui.
package animations

import (
	"github.com/gucio321/giu-animations/internal/logger"
	"sync"
	"time"

	"github.com/AllenDang/giu"
)

// ANimation is an interface implemented by each animation.
type Animation interface {
	Advance(procentageDelta float32) (shouldContinue bool)
	Start(duration time.Duration, fps int)
	Reset()
	Init()
	giu.Widget
}

type transitionState struct {
	isRunning     bool
	currentLayout bool
	elapsed       time.Duration
	duration      time.Duration
	customData    any
	shouldInit    bool
	m             *sync.Mutex
}

func (s *transitionState) IsRunning() bool {
	s.m.Lock()
	defer s.m.Unlock()
	return s.isRunning
}

func (s *transitionState) Dispose() {
	// noop
}

func (t *animationWidget) GetState() *transitionState {
	if s := giu.Context.GetState(t.id); s != nil {
		state, ok := s.(*transitionState)
		if !ok {
			logger.Fatalf("error asserting type of ttransition state: got %T", s)
		}

		return state
	}

	giu.Context.SetState(t.id, t.newState())

	return t.GetState()
}

func (t *animationWidget) newState() *transitionState {
	return &transitionState{
		shouldInit: true,
		m:          &sync.Mutex{},
	}
}

type animationWidget struct {
	id                   string
	renderer1, renderer2 func(this Animation)
	Animation
}

func newAnimation(a Animation, renderer1, renderer2 func(this Animation)) *animationWidget {
	result := &animationWidget{
		id:        giu.GenAutoID("Animation"),
		renderer1: renderer1,
		renderer2: renderer2,
		Animation: a,
	}

	return result
}

func (t *animationWidget) Start(duration time.Duration, fps int) {
	t.Animation.Reset()
	state := t.GetState()

	if state.isRunning {
		logger.Fatal("animationWidget: StartTransition called, but transition is already running")
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

func (t *animationWidget) GetCustomData() any {
	state := t.GetState()
	state.m.Lock()
	defer state.m.Unlock()
	return state.customData
}

func (t *animationWidget) SetCustomData(d any) {
	state := t.GetState()
	state.m.Lock()
	state.customData = d
	state.m.Unlock()
}

func (t *animationWidget) BuildNormal(a Animation) (proceeded bool) {
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

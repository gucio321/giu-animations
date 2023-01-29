package animations

import (
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

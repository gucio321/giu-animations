package animations

import (
	"sync"
	"time"
)

type animatorState struct {
	isRunning     bool
	currentLayout bool
	elapsed       time.Duration
	duration      time.Duration
	customData    any
	shouldInit    bool
	m             *sync.Mutex
}

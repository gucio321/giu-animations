package animations

import (
	"github.com/AllenDang/giu"
)

// Animation is an interface implemented by each animation.
type Animation interface {
	// Init is called once, immediately on start.
	Init()
	// Reset is called whenever needs to restart animation.
	Reset()
	// Advance is called every frame
	Advance(procentageDelta float32) (shouldContinue bool)
	giu.Widget
}

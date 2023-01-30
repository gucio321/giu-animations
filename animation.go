package animations

// Animation is an interface implemented by each animation.
type Animation interface {
	// Init is called once, immediately on start.
	Init()
	// Reset is called whenever needs to restart animation.
	Reset()
	// Advance is called every frame
	//Advance(procentageDelta float32) (shouldContinue bool)

	// BuildNormal is called every frame when animation is not running
	BuildNormal()
	// BuildAnimation is called when running an animation.
	// It receives the current animation progress as a float, where
	// 0 >= animationPercentage <= 1
	BuildAnimation(animationPercentage float32)
}

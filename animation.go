package animations

// Animation is an interface implemented by each animation.
// Every type that implements Animation interface is liable to
// be used as an argument of Animator method.
type Animation interface {
	// Init is called once, immediately on start.
	Init()
	// Reset is called whenever needs to restart animation.
	Reset()

	// KeyFramesCount is used mainly by the AnimatorWidget.
	// It returns animation number of key frames.
	KeyFramesCount() int

	// BuildNormal is called every frame when animation is not running
	// starter is animation link to Animator.Start
	// triggerCheck should be placed where the animator should check for trigger. If not called - Animator will do that
	// after rendering the widget (after calling BuildNormal).
	BuildNormal(currentKeyFrame KeyFrame, starterFunc StarterFunc, triggerCheck func())
	// BuildAnimation is called when running an animation.
	// It receives several important arguments:
	// - animationPercentage after applying specified by Animator
	//   easing algorithm.
	//   ATTENTION: this value may be less than 0 or greater than 1
	// - pure percentage status of animation before applying
	//   any algorithm.
	//   it is always in range <0, 1> (0 <= arbitraryPercentage <= 1)
	//   NOTE: this value should NOT be used in most cases, because it will
	//   disable user from specifying Easing Algorithm and most use-cases
	//   does not want this, however you may want to use for comparing something.
	// - base and destination Key Frames - your animation should be played from first to the second.
	// - animation's Play Mode - use it if it is important to know what is the exact play direction.
	// - starter functions set (see StarterFunc)
	// starter is animation link to (*Animator).Start() method.
	BuildAnimation(
		animationPercentage, arbitraryPercentage float32,
		baseKeyFrame, destinationKeyFrame KeyFrame,
		mode PlayMode,
		starterFunc StarterFunc,
	)
}

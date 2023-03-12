package animations

type StarterFunc func(playMode PlayMode)

// Animation is an interface implemented by each animation.
type Animation interface {
	// Init is called once, immediately on start.
	Init()
	// Reset is called whenever needs to restart animation.
	Reset()

	// KeyFramesCount is used mainly by the AnimatorWidget.
	// It returns a number of key frames.
	KeyFramesCount() int

	// BuildNormal is called every frame when animation is not running
	// starter is a link to Animator.Start
	BuildNormal(currentKeyFrame KeyFrame, starterFunc StarterFunc)
	// BuildAnimation is called when running an animation.
	// It receives two values:
	// - first one is animationPercentage after applying specified by Animator
	//   easing algorithm.
	//   ATTENTION: this value may be less than 0 or larger than 1
	// - second value is arbitrary percentage progress before applying
	//   easing algorithm.
	//   it is always in range <0, 1> (0 <= arbitraryPercentage <= 1)
	//   NOTE: this value should NOT be used in most cases, because it will
	//   disable user from specifying Easing Algorithm and most use-cases
	//   does not want this. Exceptions, I see for now may be:
	// 	 - your animation does not accept negative (or larger than 1) progress values
	// starter is a link to (*Animator).Start() method.
	BuildAnimation(animationPercentage, arbitraryPercentage float32, baseKeyFrame, destinationKeyFrame KeyFrame, starterFunc StarterFunc)
}

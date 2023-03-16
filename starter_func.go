package animations

// StarterFunc contains animation reference to all Starters of AnimatorWidget
type StarterFunc interface {
	Start(mode PlayMode)
	StartKeyFrames(beginKF, destinyKF KeyFrame, cyclesCount int, mode PlayMode)
	StartCycle(cyclesCount int, mode PlayMode)
}

package animations

// StarterFunc contains a reference to all Starters of AnimatorWidget
type StarterFunc interface {
	Start(mode PlayMode)
	StartKeyFrames(beginKF, destinyKF KeyFrame, mode PlayMode)
	StartWhole(mode PlayMode)
}

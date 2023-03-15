package animations

import "log"

// KeyFrame represents the most important states of an animation.
// Each animation declares an algorithm of calculating states between
// its KeyFrames.
type KeyFrame byte

// getWithDelta returns clamp(0, count-1, current + delta)
// TODO(gucio321): I believe it could be done simpler.
func getWithDelta(current KeyFrame, count, delta int) KeyFrame {
	// special case
	if count == 1 && (delta == 1 || delta == -1) {
		return 0
	}

	if delta >= count {
		log.Panicf("multiple-cycles not supported yet (delta=%v >= max=%v)", delta, count)
	}

	result := int(current) + delta
	if result < 0 {
		result = count - result
	}

	if result >= count {
		result -= count
	}

	return KeyFrame(result)
}

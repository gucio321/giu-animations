package animations

import "log"

type KeyFrame byte

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

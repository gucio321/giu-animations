package animations

type KeyFrame byte

func getWithDelta(current KeyFrame, count, delta int) KeyFrame {
	if delta >= count {
		panic("multiple-cycles not supported yet (delta >= max")
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

package animations

import "math"

type EasingAlgorithm func(plainPercentage float32) (percentage float32)

// EasingAlgorithmType represents a type of easing algorithm used for animation.
// Refer https://easings.net/
type EasingAlgorithmType byte

const (
	EasingAlgNone EasingAlgorithmType = iota
	EasingAlgInSine
	EasingAlgOutSine
	EasingAlgInOutSine
	EasingAlgInElastic
	EasingAlgOutElastic
	EasingAlgInOutElastic
	EasingAlgInBounce
	EasingAlgOutBounce
	EasingAlgInOutBounce
	EasingAlgMax
)

func easingAlgInSine(plainPercentage float32) float32 {
	return 1 - float32(math.Cos((float64(plainPercentage)*math.Pi)/2))
}

func easingAlgOutSine(p float32) float32 {
	return float32(math.Sin((float64(p) * math.Pi) / 2))
}

func easingAlgInOutSine(p float32) float32 {
	return -(float32(math.Cos(math.Pi*float64(p))) - 1) / 2
}

func easingAlgInElastic(p float32) float32 {
	if p == 0 {
		return 0
	} else if p == 1 {
		return 1
	}

	return -float32(math.Pow(2, 10*float64(p)-10) * math.Sin((float64(p)*10-10.75)*(2*math.Pi)/3))
}

func easingAlgOutElastic(p float32) float32 {
	if p == 0 {
		return 0
	} else if p == 1 {
		return 1
	}

	return float32(math.Pow(2, -10*float64(p))*math.Sin((float64(p)*10-0.75)*(2*math.Pi)/3)) + 1
}

func easingAlgInOutElastic(p float32) float32 {
	c5 := (2 * math.Pi) / 4.5
	switch {
	case p == 0:
		return 0
	case p < 0.5:
		return -float32((math.Pow(2, 20*float64(p)-10) * math.Sin((20*float64(p)-11.125)*c5)) / 2)
	case p >= 0.5:
		return float32((math.Pow(2, -20*float64(p)+10)*math.Sin((20*float64(p)-11.125)*c5))/2) + 1
	case p == 1:
		return 1
	default:
		return 0
	}
}

func easingAlgInBounce(p float32) float32 {
	return 1 - easingAlgOutBounce(1-p)
}

func easingAlgOutBounce(p float32) float32 {
	const n1 = 7.5625
	const d1 = 2.75

	if p < 1/d1 {
		return n1 * p * p
	} else if p < 2/d1 {
		p -= 1.5 / d1
		return n1*p*p + 0.75
	} else if p < 2.5/d1 {
		p -= 2.25 / d1
		return n1*p*p + 0.9375
	} else {
		p -= 2.625 / d1
		return n1*p*p + 0.984375
	}
}

func easingAlgInOutBounce(p float32) float32 {
	if p < 0.5 {
		return easingAlgInBounce(p*2) * 0.5
	}

	return easingAlgOutBounce(p*2-1)*0.5 + 0.5
}

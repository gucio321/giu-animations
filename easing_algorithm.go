package animations

import (
	"log"
	"math"
)

type EasingAlgorithm func(plainPercentage float32) (percentage float32)

// EasingAlgorithmType represents a type of easing algorithm used for animation.
// Refer https://easings.net/
type EasingAlgorithmType byte

const (
	EasingAlgNone EasingAlgorithmType = iota
	EasingAlgInSine
	EasingAlgOutSine
	EasingAlgInOutSine
	EasingAlgInBack
	EasingAlgOutBack
	EasingAlgInOutBack
	EasingAlgInElastic
	EasingAlgOutElastic
	EasingAlgInOutElastic
	EasingAlgInBounce
	EasingAlgOutBounce
	EasingAlgInOutBounce
	EasingAlgMax
)

// Ease takes EasingAlgorithmType and plain percentage value t and returns eased value.
// The following condition is expected to be met, however they are not restricted anyhow:
// 0 <= t <= 1
func Ease(alg EasingAlgorithmType, t float32) float32 {
	switch alg {
	case EasingAlgNone:
		return t
	case EasingAlgInSine:
		return easingAlgInSine(t)
	case EasingAlgOutSine:
		return easingAlgOutSine(t)
	case EasingAlgInOutSine:
		return easingAlgInOutSine(t)
	case EasingAlgInBack:
		return easingAlgInBack(t)
	case EasingAlgOutBack:
		return easingAlgOutBack(t)
	case EasingAlgInOutBack:
		return easingAlgInOutBack(t)
	case EasingAlgInElastic:
		return easingAlgInElastic(t)
	case EasingAlgOutElastic:
		return easingAlgOutElastic(t)
	case EasingAlgInOutElastic:
		return easingAlgInOutElastic(t)
	case EasingAlgInBounce:
		return easingAlgInBounce(t)
	case EasingAlgOutBounce:
		return easingAlgOutBounce(t)
	case EasingAlgInOutBounce:
		return easingAlgInOutBounce(t)
	}

	log.Panicf("Unknown easing type %v", alg)
	return -1 // unreachable
}

func easingAlgInSine(plainPercentage float32) float32 {
	return 1 - float32(math.Cos((float64(plainPercentage)*math.Pi)/2))
}

func easingAlgOutSine(p float32) float32 {
	return float32(math.Sin((float64(p) * math.Pi) / 2))
}

func easingAlgInOutSine(p float32) float32 {
	return -(float32(math.Cos(math.Pi*float64(p))) - 1) / 2
}

func easingAlgInBack(p float32) float32 {
	s := float32(1.70158)
	return p * p * ((s+1)*p - s)
}

func easingAlgOutBack(p float32) float32 {
	s := float32(1.70158)
	p -= 1
	return p*p*((s+1)*p+s) + 1
}

func easingAlgInOutBack(p float32) float32 {
	s := float32(1.70158) * 1.525
	p *= 2
	if p < 1 {
		return 0.5 * (p * p * ((s+1)*p - s))
	}
	p -= 2
	return 0.5 * (p*p*((s+1)*p+s) + 2)
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

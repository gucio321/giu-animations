package animations

import (
	"log"
	"math"
)

// EasingAlgorithm describes what exactly an Easing Function is.
type EasingAlgorithm func(plainPercentage float32) (percentage float32)

// EasingAlgorithmType represents animation type of easing algorithm used for animation.
// Refer https://easings.net/
type EasingAlgorithmType byte

// Easing Algorithm types.
const (
	EasingAlgNone EasingAlgorithmType = iota

	EasingAlgInSine
	EasingAlgOutSine
	EasingAlgInOutSine

	EasingAlgInQuad
	EasingAlgOutQuad
	EasingAlgInOutQuad

	EasingAlgInCubic
	EasingAlgOutCubic
	EasingAlgInOutCubic

	EasingAlgInQuart
	EasingAlgOutQuart
	EasingAlgInOutQuart

	EasingAlgInQuint
	EasingAlgOutQuint
	EasingAlgInOutQuint

	EasingAlgInExpo
	EasingAlgOutExpo
	EasingAlgInOutExpo

	EasingAlgInCirc
	EasingAlgOutCirc
	EasingAlgInOutCirc

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
// 0 <= t <= 1.
func Ease(alg EasingAlgorithmType, t float32) float32 {
	algs := map[EasingAlgorithmType]EasingAlgorithm{
		EasingAlgNone: func(t float32) float32 { return t },

		// sine
		EasingAlgInSine:    easingAlgInSine,
		EasingAlgOutSine:   easingAlgOutSine,
		EasingAlgInOutSine: easingAlgInOutSine,

		// quad
		EasingAlgInQuad:    easingAlgInQuad,
		EasingAlgOutQuad:   easingAlgOutQuad,
		EasingAlgInOutQuad: easingAlgInOutQuad,

		// cubic
		EasingAlgInCubic:    easingAlgInCubic,
		EasingAlgOutCubic:   easingAlgOutCubic,
		EasingAlgInOutCubic: easingAlgInOutCubic,

		// quart
		EasingAlgInQuart:    easingAlgInQuart,
		EasingAlgOutQuart:   easingAlgOutQuart,
		EasingAlgInOutQuart: easingAlgInOutQuart,

		// quint
		EasingAlgInQuint:    easingAlgInQuint,
		EasingAlgOutQuint:   easingAlgOutQuint,
		EasingAlgInOutQuint: easingAlgInOutQuint,

		// expo
		EasingAlgInExpo:    easingAlgInExpo,
		EasingAlgOutExpo:   easingAlgOutExpo,
		EasingAlgInOutExpo: easingAlgInOutExpo,

		// circ
		EasingAlgInCirc:    easingAlgInCirc,
		EasingAlgOutCirc:   easingAlgOutCirc,
		EasingAlgInOutCirc: easingAlgInOutCirc,

		// back
		EasingAlgInBack:    easingAlgInBack,
		EasingAlgOutBack:   easingAlgOutBack,
		EasingAlgInOutBack: easingAlgInOutBack,

		// elastic
		EasingAlgInElastic:    easingAlgInElastic,
		EasingAlgOutElastic:   easingAlgOutElastic,
		EasingAlgInOutElastic: easingAlgInOutElastic,

		// bounce
		EasingAlgInBounce:    easingAlgInBounce,
		EasingAlgOutBounce:   easingAlgOutBounce,
		EasingAlgInOutBounce: easingAlgInOutBounce,
	}

	a, found := algs[alg]
	if !found {
		log.Panicf("Unknown easing type %v", alg)
	}

	return a(t)
}

// === Easing Algorithm Implementations ===

// - Sine:

func easingAlgInSine(plainPercentage float32) float32 {
	return 1 - float32(math.Cos((float64(plainPercentage)*math.Pi)/2))
}

func easingAlgOutSine(p float32) float32 {
	return float32(math.Sin((float64(p) * math.Pi) / 2))
}

func easingAlgInOutSine(p float32) float32 {
	return -(float32(math.Cos(math.Pi*float64(p))) - 1) / 2
}

//- quad

func easingAlgInQuad(p float32) float32 {
	return p * p
}

func easingAlgOutQuad(p float32) float32 {
	return 1 - (1-p)*(1-p)
}

func easingAlgInOutQuad(p float32) float32 {
	if p < .5 {
		return 2 * p * p
	}

	return float32(1 - math.Pow(float64(-2*p+2), 2)/2)
}

// - cubic

func easingAlgInCubic(p float32) float32 {
	return p * p * p
}

func easingAlgOutCubic(p float32) float32 {
	return float32(1 - math.Pow(float64(1-p), 3))
}

func easingAlgInOutCubic(p float32) float32 {
	if p < .5 {
		return 4 * p * p * p
	}

	return float32(1 - math.Pow(float64(-2*p+2), 3)/2)
}

// - quart.
func easingAlgInQuart(p float32) float32 {
	return p * p * p * p
}

func easingAlgOutQuart(p float32) float32 {
	return float32(1 - math.Pow(float64(1-p), 4))
}

func easingAlgInOutQuart(p float32) float32 {
	if p < .5 {
		return 8 * p * p * p * p
	}

	return float32(1 - math.Pow(float64(-2*p+2), 4)/2)
}

// - quint

func easingAlgInQuint(p float32) float32 {
	return float32(math.Pow(float64(p), 5))
}

func easingAlgOutQuint(p float32) float32 {
	return float32(1 - math.Pow(float64(1-p), 5))
}

func easingAlgInOutQuint(p float32) float32 {
	if p < .5 {
		return float32(16 * math.Pow(float64(p), 5))
	}

	return float32(1 - math.Pow(float64(-2*p+2), 5)/2)
}

// - expo

func easingAlgInExpo(p float32) float32 {
	if p == 0 {
		return 0
	}

	return float32(math.Pow(2, float64(10*p-10)))
}

func easingAlgOutExpo(p float32) float32 {
	if p == 0 {
		return 1
	}

	return float32(1 - math.Pow(2, float64(-10*p)))
}

func easingAlgInOutExpo(p float32) float32 {
	switch {
	case p == 0:
		return 0
	case p == 1:
		return 1
	case p < .5:
		return float32(math.Pow(2, float64(20*p-10)) / 2)
	default:
		return float32((2 - math.Pow(2, float64(-20*p+10))) / 2)
	}
}

// - circ

func easingAlgInCirc(p float32) float32 {
	return float32(1 - math.Sqrt(1-math.Pow(float64(p), 2)))
}

func easingAlgOutCirc(p float32) float32 {
	return float32(math.Sqrt(1 - math.Pow(float64(p-1), 2)))
}

func easingAlgInOutCirc(p float32) float32 {
	if p < .5 {
		return float32((1 - math.Sqrt(1-math.Pow(float64(2*p), 2))) / 2)
	}

	return float32((math.Sqrt(1-math.Pow(float64(-2*p+2), 2)) + 1) / 2)
}

// - back

func easingAlgInBack(p float32) float32 {
	const s = float32(1.70158)

	return p * p * ((s+1)*p - s)
}

func easingAlgOutBack(p float32) float32 {
	const s = float32(1.70158)

	p--

	return p*p*((s+1)*p+s) + 1
}

func easingAlgInOutBack(p float32) float32 {
	const s = float32(1.70158) * 1.525

	p *= 2

	if p < 1 {
		return 0.5 * (p * p * ((s+1)*p - s))
	}

	p -= 2

	return 0.5 * (p*p*((s+1)*p+s) + 2)
}

// - elastic

func easingAlgInElastic(p float32) float32 {
	switch p {
	case 0:
		return 0
	case 1:
		return 1
	default:
		return -float32(math.Pow(2, 10*float64(p)-10) * math.Sin((float64(p)*10-10.75)*(2*math.Pi)/3))
	}
}

func easingAlgOutElastic(p float32) float32 {
	switch p {
	case 0:
		return 0
	case 1:
		return 1
	default:
		return float32(math.Pow(2, -10*float64(p))*math.Sin((float64(p)*10-0.75)*(2*math.Pi)/3)) + 1
	}
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
	}

	return 0
}

// - bounce

func easingAlgInBounce(p float32) float32 {
	return 1 - easingAlgOutBounce(1-p)
}

func easingAlgOutBounce(p float32) float32 {
	const (
		n1 = 7.5625
		d1 = 2.75
	)

	switch {
	case p < 1/d1:
		return n1 * p * p
	case p < 2/d1:
		p -= 1.5 / d1

		return n1*p*p + 0.75
	case p < 2.5/d1:
		p -= 2.25 / d1

		return n1*p*p + 0.9375
	}

	p -= 2.625 / d1

	return n1*p*p + 0.984375
}

func easingAlgInOutBounce(p float32) float32 {
	if p < 0.5 {
		return easingAlgInBounce(p*2) * 0.5
	}

	return easingAlgOutBounce(p*2-1)*0.5 + 0.5
}

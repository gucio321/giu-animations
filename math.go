package animations

import (
	"math"

	"github.com/AllenDang/imgui-go"
)

func factorial(n int) int {
	if n == 0 {
		return 1
	}

	return n * factorial(n-1)
}

// refer: http://zobaczycmatematyke.krk.pl/025-Zolkos-Krakow/bezier.html
func bezier(t float32, points []imgui.Vec2) imgui.Vec2 {
	var result imgui.Vec2
	for i := 0; i < len(points); i++ {
		d := float32(factorial(len(points)-1)) /
			float32(factorial(i)*factorial(len(points)-1-i)) *
			(float32(math.Pow(float64(t), float64(i))) *
				(float32(math.Pow(float64(1-t), float64(len(points)-1-i)))))
		result.X += points[i].X * d
		result.Y += points[i].Y * d
	}

	return result
}

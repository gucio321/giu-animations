package animations

import "github.com/AllenDang/cimgui-go/imgui"

func vecSum(vec1, vec2 imgui.Vec2) imgui.Vec2 {
	return imgui.Vec2{
		X: vec1.X + vec2.X,
		Y: vec1.Y + vec2.Y,
	}
}

func vecDif(vec1, vec2 imgui.Vec2) imgui.Vec2 {
	return imgui.Vec2{
		X: vec1.X - vec2.X,
		Y: vec1.Y - vec2.Y,
	}
}

func vecMul(vec1 imgui.Vec2, multiplier float32) imgui.Vec2 {
	return imgui.Vec2{
		X: vec1.X * multiplier,
		Y: vec1.Y * multiplier,
	}
}

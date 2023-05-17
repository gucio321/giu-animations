package animations

import (
	"github.com/AllenDang/giu"
	"log"
)

type ResizeAnimation struct {
	widget resizable2D
}

func Resize(w giu.Widget) *ResizeAnimation {
	return &ResizeAnimation{
		widget: getResizer(w),
	}
}

func (r ResizeAnimation) Init() {
	//TODO implement me
	panic("implement me")
}

func (r ResizeAnimation) Reset() {
	//TODO implement me
	panic("implement me")
}

func (r ResizeAnimation) KeyFramesCount() int {
	//TODO implement me
	panic("implement me")
}

func (r ResizeAnimation) BuildNormal(currentKeyFrame KeyFrame, starterFunc StarterFunc) {
	//TODO implement me
	panic("implement me")
}

func (r ResizeAnimation) BuildAnimation(animationPercentage, arbitraryPercentage float32, baseKeyFrame, destinationKeyFrame KeyFrame, mode PlayMode, starterFunc StarterFunc) {
	//TODO implement me
	panic("implement me")
}

type resizable2D interface {
	Size(w, h float32)
	giu.Widget
}

type resizable1D interface {
	Size(w float32)
	giu.Widget
}

type resizable1DTo2D struct {
	resizable1D
}

func convert1DResizableTo2D(w resizable1D) *resizable1DTo2D {
	return &resizable1DTo2D{w}
}

func (r *resizable1DTo2D) Size(w, _ float32) {
	r.resizable1D.Size(w)
}

func getResizer(w giu.Widget) resizable2D {
	if r, is1D := w.(resizable1D); is1D {
		return convert1DResizableTo2D(r)
	}

	if r, is2D := w.(resizable2D); is2D {
		return r
	}

	log.Panicf("Widget of type %T implements neither Size(width, height float32) nor Size(width float32) method. - it is not supported")

	return nil
}

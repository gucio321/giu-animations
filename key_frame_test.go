package animations

import "testing"

func Test_getWithDelta(t *testing.T) {
	type args struct {
		current KeyFrame
		count   int
		delta   int
	}
	tests := []struct {
		name string
		args args
		want KeyFrame
	}{
		{"nothing 0/0", args{0, 1, 0}, 0},
		{"simply next frame 0->1", args{0, 2, 1}, 1},
		{"simply previous frame 0<-1", args{1, 2, -1}, 0},
		{"next 1->0", args{1, 2, 1}, 0},
		{"previous 1<-0", args{0, 2, -1}, 1},
		{"previous 2<-0", args{0, 3, -1}, 2},
		{"Only one frame", args{0, 1, 1}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getWithDelta(tt.args.current, tt.args.count, tt.args.delta); got != tt.want {
				t.Errorf("getWithDelta() = %v, want %v", got, tt.want)
			}
		})
	}
}

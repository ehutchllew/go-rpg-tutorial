package animations

type Animation struct {
	First        int
	frame        int
	frameCounter float32
	Last         int
	SpeedInTps   float32 // how many ticks before next frame
	Step         int     // how many indices to move per frame
}

func (a *Animation) Frame() int {
	return a.frame
}

func (a *Animation) Update() {
	a.frameCounter -= 1.0
	if a.frameCounter < 0.0 {
		a.frameCounter = a.SpeedInTps
		a.frame += a.Step
		if a.frame > a.Last {
			a.frame = a.First
		}
	}
}

func NewAnimation(first, last, step int, speed float32) *Animation {
	return &Animation{
		First:        first,
		frame:        first,
		frameCounter: speed,
		Last:         last,
		SpeedInTps:   speed,
		Step:         step,
	}
}

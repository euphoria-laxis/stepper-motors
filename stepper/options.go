package stepper

type Options struct {
	running    bool
	threshold  uint
	currentPos int
	nSteps     uint
	sequence   [8][4]uint8
	pins       [4]int
}

// StepAngle = (Step angle / gear reduction ratio) = (5.625 / 63.68395)
const StepAngle float64 = 0.0883268076179

var (
	Sequence28BYJ48 = [8][4]uint8{
		{1, 0, 0, 0},
		{1, 1, 0, 0},
		{0, 1, 0, 0},
		{0, 1, 1, 0},
		{0, 0, 1, 0},
		{0, 0, 1, 1},
		{0, 0, 0, 1},
		{1, 0, 0, 1},
	}
	defaultOptions = Options{
		running:    false,
		threshold:  0,
		currentPos: 0,
		nSteps:     0,
		// Switching sequence for the 28BYJ48 (clockwise)
		sequence: Sequence28BYJ48,
		pins:     [4]int{0, 1, 2, 3},
	}
)

type OptFunc func(options *Options)

func SetGPIOs(pins [4]int) OptFunc {
	return func(options *Options) {
		options.pins = pins
	}
}

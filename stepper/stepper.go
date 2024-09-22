package stepper

import (
	"log/slog"
	"math"
	"os"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

type StepperMotor struct {
	*Options
}

// NewStepperMotor StepperMotor constructor
func NewStepperMotor(opts ...OptFunc) *StepperMotor {
	o := defaultOptions
	for _, fn := range opts {
		fn(&o)
	}

	return &StepperMotor{&o}
}

// IsRunning returns s.running value
func (s *StepperMotor) IsRunning() bool {
	return s.running
}

// GetThreshold returns s.threshold value
func (s *StepperMotor) GetThreshold() uint {
	return s.threshold
}

// GetCurrentPosition returns s.currentPos value
func (s *StepperMotor) GetCurrentPosition() int {
	return s.currentPos
}

// GetSteps returns steps count
func GetSteps(angle uint) uint {
	return uint(math.Round(float64(angle) / StepAngle))
}

// GetNumOfSteps returns s.nSteps value
func (s *StepperMotor) GetNumOfSteps() uint {
	return s.nSteps
}

// SetThreshold returns s.threshold value
func (s *StepperMotor) SetThreshold(threshold uint) {
	s.threshold = threshold
}

// Run the stepper motor using given params (direction, angle and speed)
func (s *StepperMotor) Run(direction Direction, angle uint, speed Speed) {
	s.running = true
	// Verify if direction value is valid
	if direction != DirectionClock && direction != DirectionCounterClock {
		slog.Error("bad direction given to stepper motor (must be 1 or -1)")
		os.Exit(2)
	}
	// Delay between each step of the switching sequence (in microseconds)
	td := uint((5 * 100 / speed.Float64()) * 1000)
	ang := direction.Int() * int(angle)
	n := float64(s.GetCurrentPosition() + ang)
	// Set the right number of steps to execute, taking in account of the m_threshold
	var degrees uint
	if math.Abs(n) > float64(s.threshold) && s.threshold != 0 {
		degrees = s.threshold - uint(direction.Int()*s.currentPos)
	} else {
		degrees = angle
	}
	steps := GetSteps(degrees)
	// To go counterclockwise we need to reverse the switching sequence
	if direction == DirectionCounterClock {
		s.sequence = reverseSequence(s.sequence)
	}
	var i, count uint
	count = 0
	var pins [4]rpio.Pin
	for idx, p := range s.pins {
		pin := rpio.Pin(p)
		pin.Output()
		pins[idx] = pin
	}
	// Run the sequence
	for i = 0; i < steps; i++ {
		// Stop stepper motor if current position is greater than threshold
		if s.threshold != 0 && math.Abs(float64(s.currentPos)) >= float64(s.threshold) {
			break
		}
		// Reset count if it's greater than sequence length
		if int(count) >= len(s.sequence) {
			count = 0
		}
		// Update state
		s.nSteps = steps - (i + 1)
		if s.nSteps > 0 {
			currentAngle := ((float64(steps) - float64(s.nSteps)) / float64(steps)) * float64(degrees)
			s.currentPos = direction.Int() * int(currentAngle)
		} else {
			s.currentPos = 0
		}
		// Set GPIOs value according to sequence
		for j, pin := range pins {
			if s.sequence[count][j] == 1 {
				pin.High()
			} else {
				pin.Low()
			}
		}
		count++
		time.Sleep(time.Duration(td) * time.Microsecond) // minimum delay 5ms (speed 100%)
	}
	// Cleanup (recommended in order to prevent stepper motor overheating)
	for _, pin := range pins {
		pin.Low()
	}
	// Restore the original sequence for the next operations
	if direction == DirectionCounterClock {
		s.sequence = Sequence28BYJ48
	}
	// Update the state
	s.running = false
}

// Wait for given time.Duration
func Wait(t time.Duration) {
	time.Sleep(t)
}

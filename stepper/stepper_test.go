package stepper

import (
	"sync"
	"testing"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

/**
											### RPI 4B GPIO MAP ###
						 ## Map to plug stepper motor ULN2003 driver board to the device ##

						+-----+---------+------+---+---Pi 4B--+---+------+---------+-----+
						| BCM |   Name  | Mode | V | Physical | V | Mode | Name    | BCM |
						+-----+---------+------+---+----++----+---+------+---------+-----+
						|     |    3.3v |      |   |  1 || 2  |   |      | 5v      |     |
						|   2 |   SDA.1 | ALT0 | 1 |  3 || 4  |   |      | 5v      |     |
						|   3 |   SCL.1 | ALT0 | 1 |  5 || 6  |   |      | 0v      |     |
						|   4 | GPIO. 7 |   IN | 1 |  7 || 8  | 1 | ALT5 | TxD     | 14  |
						|     |      0v |      |   |  9 || 10 | 1 | ALT5 | RxD     | 15  |
Stepper In 1	->		|  17 | GPIO. 0 |   IN | 0 | 11 || 12 | 0 | IN   | GPIO. 1 | 18  | 	<-  Stepper In 2
Stepper In 3  	->		|  27 | GPIO. 2 |   IN | 0 | 13 || 14 |   |      | 0v      |     |
Stepper In 4	->		|  22 | GPIO. 3 |   IN | 0 | 15 || 16 | 0 | IN   | GPIO. 4 | 23  |
						|     |    3.3v |      |   | 17 || 18 | 0 | IN   | GPIO. 5 | 24  |
						|  10 |    MOSI | ALT0 | 0 | 19 || 20 |   |      | 0v      |     |
						|   9 |    MISO | ALT0 | 0 | 21 || 22 | 0 | IN   | GPIO. 6 | 25  |
						|  11 |    SCLK | ALT0 | 0 | 23 || 24 | 1 | OUT  | CE0     | 8   |
						|     |      0v |      |   | 25 || 26 | 1 | OUT  | CE1     | 7   |
						|   0 |   SDA.0 |   IN | 1 | 27 || 28 | 1 | IN   | SCL.0   | 1   |
						|   5 | GPIO.21 |   IN | 0 | 29 || 30 |   |      | 0v      |     |
						|   6 | GPIO.22 |   IN | 0 | 31 || 32 | 0 | IN   | GPIO.26 | 12  |
						|  13 | GPIO.23 |   IN | 0 | 33 || 34 |   |      | 0v      |     |
						|  19 | GPIO.24 |   IN | 0 | 35 || 36 | 0 | IN   | GPIO.27 | 16  |
						|  26 | GPIO.25 |   IN | 0 | 37 || 38 | 0 | IN   | GPIO.28 | 20  |
						|     |      0v |      |   | 39 || 40 | 0 | IN   | GPIO.29 | 21  |
						+-----+---------+------+---+----++----+---+------+---------+-----+
						| BCM |   Name  | Mode | V | Physical | V | Mode | Name    | BCM |
						+-----+---------+------+---+---Pi 4B--+---+------+---------+-----+
*/

/**
 * RPi GPIO | Physical
 * ---------+---------
 * GPIO 17  |    11
 * GPIO 18  |    12
 * GPIO 27  |    13
 * GPIO 22  |    15
 */
var smgpios = [4]int{17, 18, 27, 22}
var sm *StepperMotor

func TestStepper(t *testing.T) {
	t.Log("stepper/stepper.go tests")
	// go-rpio initialization
	if err := rpio.Open(); err != nil {
		t.Error("rpio.Open() failed")
		t.Fail()
	}
	t.Run("Test StepperMotor constructor", testNewStepperMotor())
	t.Run("Test StepperMotor getters", testStepperMotorGetters(sm))
	t.Run("Test StepperMotor getters", testStepperMotorGetters(sm))
	var p params
	p.speed = Speed100
	p.direction = DirectionClock
	p.angle = 90
	p.threshold = 45
	t.Run("Test StepperMotor run", testStepperMotorRun(sm, &p))
}

func testNewStepperMotor() func(t *testing.T) {
	return func(t *testing.T) {
		sm = NewStepperMotor(SetGPIOs(smgpios))
		if sm == nil {
			t.Error("NewStepperMotor() returned nil")
		}
	}
}

func testStepperMotorGetters(sm *StepperMotor) func(t *testing.T) {
	return func(t *testing.T) {
		if sm.IsRunning() {
			t.Error("StepperMotor() should not be running")
		}
		if sm.GetCurrentPosition() != 0 {
			t.Error("StepperMotor() current position should be 0")
		}
		if sm.GetNumOfSteps() != 0 {
			t.Error("StepperMotor() num of steps should be 0")
		}
		if sm.GetThreshold() != 0 {
			t.Error("StepperMotor() threshold should be 0")
		}
	}
}

type params struct {
	speed     Speed
	direction Direction
	angle     uint
	threshold uint
}

// testStepperMotorRun test if the stepper motor runs correctly
func testStepperMotorRun(sm *StepperMotor, p *params) func(t *testing.T) {
	return func(t *testing.T) {
		// Set threshold
		sm.SetThreshold(p.threshold)
		// Run stepper motor in a routine
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			sm.Run(p.direction, p.angle, p.speed)
		}()
		steps := GetSteps(p.angle)
		stepDuration := uint((5 * 100 / p.speed.Float64()) * 1000) // equal to 5ms for 100% speed
		wait := time.Duration((steps/2)*stepDuration) * time.Microsecond
		Wait(wait) // wait until half of the routine execution time
		if !sm.IsRunning() {
			t.Error("StepperMotor should be running")
		}
		if sm.GetNumOfSteps() <= 0 || sm.GetNumOfSteps() >= steps {
			t.Errorf("StepperMotor invalid num of steps : %d", steps)
		}
		if sm.GetCurrentPosition() == 0 || uint(sm.GetCurrentPosition()) > p.angle {
			t.Error("StepperMotor current position shouldn't be 0")
		}
		if sm.GetThreshold() != p.threshold {
			t.Errorf("StepperMotor threshold should be %d", p.threshold)
		}
		wg.Wait() // wait for routine to be done
		cleanup()
	}
}

// cleanup set pins to LOW to prevent stepper motor overheating
func cleanup() {
	for _, gpio := range smgpios {
		pin := rpio.Pin(gpio)
		pin.Output()
		pin.Low()
	}
}

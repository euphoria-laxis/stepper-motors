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

func TestStepperMotor(t *testing.T) {
	// go-rpio initialization
	if err := rpio.Open(); err != nil {
		t.Error("rpio.Open() failed")
		t.Fail()
	}
	t.Run("Test StepperMotor constructor", testNewStepperMotor())
	t.Run("Test StepperMotor getters", testStepperMotorGetters(sm))
	t.Run("Test StepperMotor getters", testStepperMotorGetters(sm))
	t.Run("Test StepperMotor run", testStepperMotorRun(sm))
}

func testNewStepperMotor() func(t *testing.T) {
	return func(t *testing.T) {
		sm = NewStepperMotor(SetGPIOs(smgpios))
		if sm == nil {
			t.Error("NewStepperMotor() returned nil")
			t.Fail()
		}
	}
}

func testStepperMotorGetters(sm *StepperMotor) func(t *testing.T) {
	return func(t *testing.T) {
		if sm.IsRunning() {
			t.Error("StepperMotor() should not be running")
			t.Fail()
		}
		if sm.GetCurrentPosition() != 0 {
			t.Error("StepperMotor() current position should be 0")
			t.Fail()
		}
		if sm.GetNumOfSteps() != 0 {
			t.Error("StepperMotor() num of steps should be 0")
			t.Fail()
		}
		if sm.GetThreshold() != 0 {
			t.Error("StepperMotor() threshold should be 0")
			t.Fail()
		}
	}
}

func testStepperMotorRun(sm *StepperMotor) func(t *testing.T) {
	return func(t *testing.T) {
		// Run stepper motor in a routine
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			sm.Run(DirectionClock, 180, Speed40)
		}()
		Wait(200 * time.Millisecond)
		if !sm.IsRunning() {
			t.Error("StepperMotor() should be running")
			t.Fail()
		}
		Wait(2 * time.Second)
		wg.Wait() // wait for routine to be done
	}
}

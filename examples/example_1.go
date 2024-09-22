package main

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/euphoria-laxis/stepper-motors/stepper"
	"github.com/stianeikeland/go-rpio/v4"
)

/**
						+-----+---------+------+---+---Pi 4B--+---+------+---------+-----+
						| BCM |   Name  | Mode | V | Physical | V | Mode | Name    | BCM |
						+-----+---------+------+---+----++----+---+------+---------+-----+
						|     |    3.3v |      |   |  1 || 2  |   |      | 5v      |     |
						|   2 |   SDA.1 | ALT0 | 1 |  3 || 4  |   |      | 5v      |     |
						|   3 |   SCL.1 | ALT0 | 1 |  5 || 6  |   |      | 0v      |     |
						|   4 | GPIO. 7 |   IN | 1 |  7 || 8  | 1 | ALT5 | TxD     | 14  |
						|     |      0v |      |   |  9 || 10 | 1 | ALT5 | RxD     | 15  |
stepper 1 in 1	->		|  17 | GPIO. 0 |   IN | 0 | 11 || 12 | 0 | IN   | GPIO. 1 | 18  | 	<-  stepper 1 in 2
stepper 1 in 3  ->		|  27 | GPIO. 2 |   IN | 0 | 13 || 14 |   |      | 0v      |     |
stepper 1 in 4	->		|  22 | GPIO. 3 |   IN | 0 | 15 || 16 | 0 | IN   | GPIO. 4 | 23  |
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

const banner = `
> Example: 		Running stepper motor
> Author: 		Euphoria Laxis
> Version:		v0.1.0
> Repository: 	https://github.com/euphoria-laxis/stepper-motors.git
`

// RPi GPIO | Physical
// -------------------
// GPIO 17  |    11
// GPIO 18  |    12
// GPIO 27  |    13
// GPIO 22  |    15
var smgpios = [4]int{17, 18, 27, 22} // Stepper motor GPIOs position (BCM)

func main() {
	fmt.Print(banner) // Display a banner when running program

	// go-rpio initialization
	if err := rpio.Open(); err != nil {
		slog.Error("rpio.Open() failed")
		os.Exit(1)
	}

	// Create stepper.StepperMotor instance
	sm := stepper.NewStepperMotor(
		stepper.SetGPIOs(smgpios),
	)

	// Params that will be used here
	smParams := params{
		direction: stepper.DirectionClock,  // direction in which the stepper motor will rotate
		speed:     stepper.Speed60,         // stepper motor speed (valid speeds are 20%, 40%, 60%, 80% and 100%)
		angle:     270,                     // will rotate to 270° from initial position
		duration:  2000 * time.Millisecond, // time to wait after running
	}

	// Run stepper motor in a routine
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		runStepperMotor(sm, &smParams)
	}()
	wg.Wait() // wait for routine to be done

	// Rerun stepper motor in opposite direction
	smParams.direction = stepper.DirectionCounterClock
	wg.Add(1)
	go func() {
		defer wg.Done()
		runStepperMotor(sm, &smParams)
	}()
	wg.Wait()

	// Close GPIOs
	err := rpio.Close()
	if err != nil {
		slog.Error("rpio.Close() failed")
	}
}

// params contains the required information to run the stepper motor
type params struct {
	direction stepper.Direction
	speed     stepper.Speed
	angle     uint
	duration  time.Duration
}

// runStepperMotor run stepper motor
func runStepperMotor(sm *stepper.StepperMotor, p *params) {
	sm.Run(p.direction, p.angle, p.speed)
	stepper.Wait(p.duration)
}

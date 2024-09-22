package stepper

// Direction is used to give stepper motor direction (clockwise or counterclockwise)
type Direction int8

const (
	DirectionClock        Direction = 1
	DirectionCounterClock Direction = -1
)

// Int return value as int
func (d Direction) Int() int {
	return int(d)
}

// Float64 return value as float64
func (d Direction) Float64() float64 {
	return float64(d)
}

// Speed is the stepper motor speed, it ensures to use only allowed speed values
type Speed uint8

const (
	Speed20  = 20
	Speed40  = 40
	Speed60  = 60
	Speed80  = 80
	Speed100 = 100
)

// Uint returns the speed value as an uint
func (s Speed) Uint() uint {
	return uint(s)
}

// Float64 returns the speed value as a float64
func (s Speed) Float64() float64 {
	return float64(s)
}

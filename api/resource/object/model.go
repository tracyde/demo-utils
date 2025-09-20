package object

import "time"

// Friendly
// Operator says:
//  “This is Viper One-One, F-16C. My position 34.123 North, 132.456 East, altitude twenty-four thousand five hundred feet,
// speed four-two-zero knots, heading two-seven-five.
// Mission status: on task.
// Time 2025-08-29T12:45:30Z.”
//
// {
//   name: "Viper 1-1",
//   type: "friendly",
//   platform: "F-16C",
// 	 position: {
// 	 	latitude: 34.123,
// 	 	longitude: 132.456,
// 	 	altitude: 24500,
// 	 	speed: 420,
// 	 	heading: 275
// 	 },
// 	 status: "on task",
//   time: "2025-08-29T12:45:30Z",
// }
//

// Surveillance Track (air target)
// Operator says:
//  “One unknown aircraft detected. Position 35.678 North, 129.432 East, altitude fifteen thousand two hundred feet, speed five-two-zero knots, heading zero-nine-zero. Track quality seven. No IFF. Sensor: AESA radar. Time 2025-08-29T12:47:10Z.”

// EW / Threat Emitter
// Operator says:
//  “SA-21 radar emission assessed. Position 36.100 North, 129.900 East. Confidence zero point eight five. Activity: tracking. Threat radius forty kilometers. Time 2025-08-29T12:49:00Z.”

type Object struct {
	Name         string
	Type         string
	Platform     string
	Position     Position
	Sensor       string
	Iff          bool
	TrackQuality int
	Confidence   float64
	Activity     string
	Status       string
	Time         time.Time
}

type Position struct {
	Latitude  float64
	Longitude float64
	Altitude  float64
	Speed     float64
	Heading   float64
}

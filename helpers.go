package satellite

import (
	"math"
)

const TWOPI float64 = math.Pi * 2.0
const DEG2RAD float64 = math.Pi / 180.0
const RAD2DEG float64 = 180.0 / math.Pi
const XPDOTP float64 = 1440.0 / (2.0 * math.Pi)

// Holds latitude and Longitude in either degrees or radians
type LatLong struct {
	Latitude, Longitude float64
}

type Vector3 struct {
	X, Y, Z float64
}

type LookAngles struct {
	Azimuth, Elevation, Range float64
}

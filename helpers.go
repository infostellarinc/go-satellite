package satellite

import (
	"math"
)

const TWOPI float64 = math.Pi * 2.0
const JULIAN_DAY_JAN_1_2000 float64 = 2451545.0
const JULIAN_CENTURY float64 = 36525.0
const SECONDS_IN_DAY float64 = 86400.0
const DEG2RAD float64 = math.Pi / 180.0
const RAD2DEG float64 = 180.0 / math.Pi
const XPDOTP float64 = 1440.0 / TWOPI
const GRAVITY_EARTH float64 = 398600.4418
const EQUATOR_RADIUS float64 = 6378.137
const POLAR_RADIUS float64 = 6356.7523142

// Holds latitude and Longitude in either degrees or radians
type LatLong struct {
	Latitude, Longitude float64
}

type Vector3 struct {
	X, Y, Z float64
}

func (v Vector3) Equals(v2 Vector3) bool {
	return closeFloat(v.X, v2.X) && closeFloat(v.Y, v2.Y) && closeFloat(v.Z, v2.Z)
}

func closeFloat(a, b float64) bool {
	return math.Abs(a-b) < 1e-4
}

type LookAngles struct {
	Azimuth, Elevation, Range float64
}

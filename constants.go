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

package satellite

import (
	"errors"
	"math"
)

// this procedure converts the day of the year, epochDays, to the equivalent month day, hour, minute and second.
func days2mdhms(year int64, epochDays float64) (float64, float64, float64, float64, float64) {
	lmonth := [12]int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

	if year%4 == 0 {
		lmonth = [12]int{31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	}

	dayofyr := math.Floor(epochDays)

	i := 1.0
	inttemp := 0.0

	for dayofyr > inttemp+float64(lmonth[int(i-1)]) && i < 22 {
		inttemp += float64(lmonth[int(i-1)])
		i += 1
	}

	month := i
	day := dayofyr - inttemp

	temp := (epochDays - dayofyr) * 24.0
	hour := math.Floor(temp)

	temp = (temp - hour) * 60.0
	minute := math.Floor(temp)

	second := (temp - minute) * 60.0

	return month, day, hour, minute, second
}

// Calc julian date given year, month, day, hour, minute and second
// the julian date is defined by each elapsed day since noon, jan 1, 4713 bc.
func JDay(year, month, day, hour, minute, second int) float64 {
	return (367.0*float64(year) - math.Floor((7*(float64(year)+math.Floor((float64(month)+9)/12.0)))*0.25) + math.Floor(275*float64(month)/9.0) + float64(day) + 1721013.5 + ((float64(second)/60.0+float64(minute))/60.0+float64(hour))/24.0)
}

// this function finds the greenwich sidereal time (iau-82).
func gstime(jdut1 float64) float64 {
	tut1 := (jdut1 - JULIAN_DAY_JAN_1_2000) / JULIAN_CENTURY
	result := -6.2e-6*tut1*tut1*tut1 + 0.093104*tut1*tut1 + (876600.0*3600+8640184.812866)*tut1 + 67310.54841
	result = math.Mod((result * DEG2RAD / 240.0), TWOPI)

	if result < 0.0 {
		result += TWOPI
	}

	return result
}

// Calc GST given year, month, day, hour, minute and second.
func GSTimeFromDate(year, month, day, hour, minute, second int) float64 {
	jDay := JDay(year, month, day, hour, minute, second)
	return gstime(jDay)
}

// Convert Earth Centered Inertial coordinated into equivalent latitude, longitude, altitude and velocity.
// Reference: http://celestrak.com/columns/v02n03/
func ECIToLLA(eciCoords Vector3, gmst float64) (altitude, velocity float64, ret LatLong) {
	a := EQUATOR_RADIUS // Semi-major Axis
	b := POLAR_RADIUS   // Semi-minor Axis
	f := (a - b) / a    // Flattening
	e2 := ((2 * f) - math.Pow(f, 2))

	sqx2y2 := math.Sqrt(math.Pow(eciCoords.X, 2) + math.Pow(eciCoords.Y, 2))

	// Spherical Earth Calculations
	longitude := math.Atan2(eciCoords.Y, eciCoords.X) - gmst
	latitude := math.Atan2(eciCoords.Z, sqx2y2)

	// Oblate Earth Fix
	c := 0.0
	for _ = range 20 {
		c = 1 / math.Sqrt(1-e2*(math.Sin(latitude)*math.Sin(latitude)))
		latitude = math.Atan2(eciCoords.Z+(a*c*e2*math.Sin(latitude)), sqx2y2)
	}

	// Calc Alt
	altitude = (sqx2y2 / math.Cos(latitude)) - (a * c)

	// Orbital Speed ≈ sqrt(μ / r) where μ = std. gravitaional parameter
	velocity = math.Sqrt(GRAVITY_EARTH / (altitude + EQUATOR_RADIUS))

	ret.Latitude = latitude
	ret.Longitude = longitude

	return
}

var ErrInvalidLatitude = errors.New("latitude not within bounds -pi/2 to +pi/2")

// Convert LatLong in radians to LatLong in degrees.
func LatLongDeg(rad LatLong) (LatLong, error) {
	var result LatLong
	result.Longitude = math.Mod(rad.Longitude/math.Pi*180, 360)
	if result.Longitude > 180 {
		result.Longitude = 360 - result.Longitude
	} else if result.Longitude < -180 {
		result.Longitude = 360 + result.Longitude
	}

	if rad.Latitude < (-math.Pi/2) || rad.Latitude > math.Pi/2 {
		return LatLong{}, ErrInvalidLatitude
	}
	result.Latitude = (rad.Latitude / math.Pi * 180)
	return result, nil
}

// Calculate GMST from Julian date.
// Reference: The 1992 Astronomical Almanac, page B6.
func ThetaGJD(jday float64) float64 {
	_, ut := math.Modf(jday + 0.5)
	jday -= ut
	tu := (jday - JULIAN_DAY_JAN_1_2000) / JULIAN_CENTURY
	gmst := 24110.54841 + tu*(8640184.812866+tu*(0.093104-tu*6.2e-6))
	gmst = math.Mod(gmst+86400.0*1.00273790934*ut, SECONDS_IN_DAY)
	result := TWOPI * gmst / SECONDS_IN_DAY
	return result
}

// Convert latitude, longitude and altitude(km) into equivalent Earth Centered Intertial coordinates(km)
// Reference: The 1992 Astronomical Almanac, page K11.
func LLAToECI(obsCoords LatLong, alt, jday float64) Vector3 {
	re := EQUATOR_RADIUS
	theta := math.Mod(ThetaGJD(jday)+obsCoords.Longitude, TWOPI)
	r := (re + alt) * math.Cos(obsCoords.Latitude)
	var eciObs Vector3
	eciObs.X = r * math.Cos(theta)
	eciObs.Y = r * math.Sin(theta)
	eciObs.Z = (re + alt) * math.Sin(obsCoords.Latitude)
	return eciObs
}

// Convert Earth Centered Intertial coordinates into Earth Cenetered Earth Final coordinates
// Reference: http://ccar.colorado.edu/ASEN5070/handouts/coordsys.doc
func ECIToECEF(eciCoords Vector3, gmst float64) Vector3 {
	var ecfCoords Vector3
	ecfCoords.X = eciCoords.X*math.Cos(gmst) + eciCoords.Y*math.Sin(gmst)
	ecfCoords.Y = eciCoords.X*-math.Sin(gmst) + eciCoords.Y*math.Cos(gmst)
	ecfCoords.Z = eciCoords.Z
	return ecfCoords
}

// Calculate look angles for given satellite position and observer position
// obsAlt in km
// Reference: http://celestrak.com/columns/v02n02/
func ECIToLookAngles(eciSat Vector3, obsCoords LatLong, obsAlt, jday float64) LookAngles {
	theta := math.Mod(ThetaGJD(jday)+obsCoords.Longitude, TWOPI)
	obsPos := LLAToECI(obsCoords, obsAlt, jday)

	rx := eciSat.X - obsPos.X
	ry := eciSat.Y - obsPos.Y
	rz := eciSat.Z - obsPos.Z

	topS := math.Sin(obsCoords.Latitude)*math.Cos(theta)*rx + math.Sin(obsCoords.Latitude)*math.Sin(theta)*ry - math.Cos(obsCoords.Latitude)*rz
	topE := -math.Sin(theta)*rx + math.Cos(theta)*ry
	topZ := math.Cos(obsCoords.Latitude)*math.Cos(theta)*rx + math.Cos(obsCoords.Latitude)*math.Sin(theta)*ry + math.Sin(obsCoords.Latitude)*rz

	var lookAngles LookAngles
	lookAngles.Azimuth = math.Atan(-topE / topS)
	if topS > 0 {
		lookAngles.Azimuth += math.Pi
	}
	if lookAngles.Azimuth < 0 {
		lookAngles.Azimuth += TWOPI
	}
	lookAngles.Range = math.Sqrt(rx*rx + ry*ry + rz*rz)
	lookAngles.Elevation = math.Asin(topZ / lookAngles.Range)

	return lookAngles
}

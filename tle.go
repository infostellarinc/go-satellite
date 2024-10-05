package satellite

import (
	"fmt"
	"strconv"
	"strings"
)

type TLE struct {
	Line1 string `json:"LINE1"`
	Line2 string `json:"LINE2"`

	CatalogNumber string
	EpochYear     int64
	EpochDay      float64

	// aka ndot
	FirstTimeDerivativeOfMeanMotion float64
	// aka nddot
	SecondTimeDerivativeOfMeanMotion float64
	BStar                            float64

	Inclination                   float64
	RightAscensionOfAscendingNode float64
	Eccentricity                  float64
	ArgumentOfPerigee             float64
	MeanAnomaly                   float64
	MeanMotion                    float64

	OrbitNumberAtEpoch int64
}

// Parses a two line element dataset into a Satellite struct
func ParseTLE(line1, line2 string, gravConst Gravity) (TLE, error) {
	var tle TLE
	tle.Line1 = line1
	tle.Line2 = line2

	var err error
	// sat.gravity, err = getGravConst(gravConst)
	// if err != nil {
	// 	return Satellite{}, fmt.Errorf("getGravConst: %w", err)
	// }

	// LINE 1 BEGIN
	tle.CatalogNumber = strings.TrimSpace(line1[2:7])
	tle.EpochYear, err = strconv.ParseInt(line1[18:20], 10, 0)
	if err != nil {
		return TLE{}, fmt.Errorf("epoch year: %w", err)
	}
	tle.EpochDay, err = strconv.ParseFloat(line1[20:32], 64)
	if err != nil {
		return TLE{}, fmt.Errorf("epoch days: %w", err)
	}

	// These three can be negative / positive
	tle.FirstTimeDerivativeOfMeanMotion, err = strconv.ParseFloat(strings.Replace(line1[33:43], " ", "", 2), 64)
	if err != nil {
		return TLE{}, fmt.Errorf("first time derivative of mean motion: %w", err)
	}
	tle.SecondTimeDerivativeOfMeanMotion, err = strconv.ParseFloat(strings.Replace(line1[44:45]+"."+line1[45:50]+"e"+line1[50:52], " ", "", 2), 64)
	if err != nil {
		return TLE{}, fmt.Errorf("second time derivative of mean motion: %w", err)
	}
	tle.BStar, err = strconv.ParseFloat(strings.Replace(line1[53:54]+"."+line1[54:59]+"e"+line1[59:61], " ", "", 2), 64)
	if err != nil {
		return TLE{}, fmt.Errorf("b star: %w", err)
	}
	// Note: skips ephemeris type, element number, checksum
	// LINE 1 END

	// LINE 2 BEGIN
	tle.Inclination, err = strconv.ParseFloat(strings.Replace(line2[8:16], " ", "", 2), 64)
	if err != nil {
		return TLE{}, fmt.Errorf("inclincation: %w", err)
	}
	tle.RightAscensionOfAscendingNode, err = strconv.ParseFloat(strings.Replace(line2[17:25], " ", "", 2), 64)
	if err != nil {
		return TLE{}, fmt.Errorf("right ascension of ascending node: %w", err)
	}
	tle.Eccentricity, err = strconv.ParseFloat("."+line2[26:33], 64)
	if err != nil {
		return TLE{}, fmt.Errorf("eccentricity: %w", err)
	}
	tle.ArgumentOfPerigee, err = strconv.ParseFloat(strings.Replace(line2[34:42], " ", "", 2), 64)
	if err != nil {
		return TLE{}, fmt.Errorf("argument of perigee: %w", err)
	}
	tle.MeanAnomaly, err = strconv.ParseFloat(strings.Replace(line2[43:51], " ", "", 2), 64)
	if err != nil {
		return TLE{}, fmt.Errorf("mean anomoly: %w", err)
	}
	tle.MeanMotion, err = strconv.ParseFloat(strings.Replace(line2[52:63], " ", "", 2), 64)
	if err != nil {
		return TLE{}, fmt.Errorf("mean motion: %w", err)
	}

	tle.OrbitNumberAtEpoch, err = strconv.ParseInt(strings.TrimSpace(line2[63:68]), 10, 0)
	if err != nil {
		return TLE{}, fmt.Errorf("orbit number at epoch: %w", err)
	}
	// Note: skips checksum
	// LINE 2 END
	return tle, nil
}

// Converts a two line element data set into a Satellite struct and runs sgp4init
func TLEToSat(line1, line2 string, gravConst Gravity) (Satellite, error) {
	tle, err := ParseTLE(line1, line2, gravConst)
	if err != nil {
		return Satellite{}, fmt.Errorf("could not parse tle: %w", err)
	}

	var sat Satellite
	sat.gravity, err = getGravConst(gravConst)
	if err != nil {
		return Satellite{}, fmt.Errorf("getGravConst: %w", err)
	}
	sat.ndot = tle.FirstTimeDerivativeOfMeanMotion / (XPDOTP * 1440.0)
	sat.nddot = tle.SecondTimeDerivativeOfMeanMotion / (XPDOTP * 1440.0 * 1440)
	sat.bstar = tle.BStar
	sat.inclo = tle.Inclination * DEG2RAD
	sat.nodeo = tle.RightAscensionOfAscendingNode * DEG2RAD
	sat.ecco = tle.Eccentricity
	sat.argpo = tle.ArgumentOfPerigee * DEG2RAD
	sat.mo = tle.MeanAnomaly * DEG2RAD
	sat.no = tle.MeanMotion / XPDOTP

	var year int64
	if tle.EpochYear < 57 {
		year = tle.EpochYear + 2000
	} else {
		year = tle.EpochYear + 1900
	}

	month, day, hour, minute, second := days2mdhms(year, tle.EpochDay)

	sat.jdsatepoch = JDay(int(year), int(month), int(day), int(hour), int(minute), second)

	_, _, err = sgp4init(sat.jdsatepoch-2433281.5, &sat)
	if err != nil {
		return Satellite{}, fmt.Errorf("sgp4init: %w", err)
	}

	return sat, nil
}

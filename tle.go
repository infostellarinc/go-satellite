package satellite

import (
	"fmt"
	"strconv"
	"strings"
)

// Parses a two line element dataset into a Satellite struct
func ParseTLE(line1, line2 string, gravConst Gravity) (Satellite, error) {
	var sat Satellite
	sat.Line1 = line1
	sat.Line2 = line2

	sat.Error = 0
	var err error
	sat.whichconst, err = getGravConst(gravConst)
	if err != nil {
		return Satellite{}, fmt.Errorf("getGravConst: %w", err)
	}

	// LINE 1 BEGIN
	sat.satnum, err = strconv.ParseInt(strings.TrimSpace(line1[2:7]), 10, 0)
	if err != nil {
		return Satellite{}, fmt.Errorf("satellite number: %w", err)
	}
	sat.epochyr, err = strconv.ParseInt(line1[18:20], 10, 0)
	if err != nil {
		return Satellite{}, fmt.Errorf("epoch year: %w", err)
	}
	sat.epochdays, err = strconv.ParseFloat(line1[20:32], 64)
	if err != nil {
		return Satellite{}, fmt.Errorf("epoch days: %w", err)
	}

	// These three can be negative / positive
	sat.ndot, err = strconv.ParseFloat(strings.Replace(line1[33:43], " ", "", 2), 64)
	if err != nil {
		return Satellite{}, fmt.Errorf("first time derivative of mean motion: %w", err)
	}
	sat.nddot, err = strconv.ParseFloat(strings.Replace(line1[44:45]+"."+line1[45:50]+"e"+line1[50:52], " ", "", 2), 64)
	if err != nil {
		return Satellite{}, fmt.Errorf("second time derivative of mean motion: %w", err)
	}
	sat.bstar, err = strconv.ParseFloat(strings.Replace(line1[53:54]+"."+line1[54:59]+"e"+line1[59:61], " ", "", 2), 64)
	if err != nil {
		return Satellite{}, fmt.Errorf("b star: %w", err)
	}
	// Note: skips ephemeris type, element number, checksum
	// LINE 1 END

	// LINE 2 BEGIN
	sat.inclo, err = strconv.ParseFloat(strings.Replace(line2[8:16], " ", "", 2), 64)
	if err != nil {
		return Satellite{}, fmt.Errorf("inclincation: %w", err)
	}
	sat.nodeo, err = strconv.ParseFloat(strings.Replace(line2[17:25], " ", "", 2), 64)
	if err != nil {
		return Satellite{}, fmt.Errorf("right ascension of ascending node: %w", err)
	}
	sat.ecco, err = strconv.ParseFloat("."+line2[26:33], 64)
	if err != nil {
		return Satellite{}, fmt.Errorf("eccentricity: %w", err)
	}
	sat.argpo, err = strconv.ParseFloat(strings.Replace(line2[34:42], " ", "", 2), 64)
	if err != nil {
		return Satellite{}, fmt.Errorf("argument of perigee: %w", err)
	}
	sat.mo, err = strconv.ParseFloat(strings.Replace(line2[43:51], " ", "", 2), 64)
	if err != nil {
		return Satellite{}, fmt.Errorf("mean anomoly: %w", err)
	}
	sat.no, err = strconv.ParseFloat(strings.Replace(line2[52:63], " ", "", 2), 64)
	if err != nil {
		return Satellite{}, fmt.Errorf("mean motion: %w", err)
	}
	// Note: skips orbit number at epoch, checksum
	// LINE 2 END
	return sat, nil
}

// Converts a two line element data set into a Satellite struct and runs sgp4init
func TLEToSat(line1, line2 string, gravConst Gravity) (Satellite, error) {
	//sat := Satellite{Line1: line1, Line2: line2}
	sat, err := ParseTLE(line1, line2, gravConst)
	if err != nil {
		return Satellite{}, fmt.Errorf("could not parse tle: %w", err)
	}

	opsmode := "i"

	sat.no = sat.no / XPDOTP
	sat.ndot = sat.ndot / (XPDOTP * 1440.0)
	sat.nddot = sat.nddot / (XPDOTP * 1440.0 * 1440)

	sat.inclo = sat.inclo * DEG2RAD
	sat.nodeo = sat.nodeo * DEG2RAD
	sat.argpo = sat.argpo * DEG2RAD
	sat.mo = sat.mo * DEG2RAD

	var year int64
	if sat.epochyr < 57 {
		year = sat.epochyr + 2000
	} else {
		year = sat.epochyr + 1900
	}

	mon, day, hr, min, sec := days2mdhms(year, sat.epochdays)

	sat.jdsatepoch = JDay(int(year), int(mon), int(day), int(hr), int(min), int(sec))

	sgp4init(&opsmode, sat.jdsatepoch-2433281.5, &sat)

	return sat, nil
}

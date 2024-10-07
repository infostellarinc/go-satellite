package satellite

import (
	"testing"
	"time"
)

func TestECItoLookAngles(t *testing.T) {
	tests := []struct {
		name                string
		line1               string
		line2               string
		gravConst           Gravity
		time                time.Time
		latitudeDegree      float64
		longitudeDegree     float64
		altitude            float64
		wantAzimuthDegree   float64
		wantElevationDegree float64
	}{
		{
			name:                "ISS#22825 at 2020-05-23T20:23:37",
			line1:               "1 25544U 98067A   20140.34419374 -.00000374  00000-0  13653-5 0  9990",
			line2:               "2 25544  51.6433 131.2277 0001338 330.3524 173.1622 15.49372617227549",
			gravConst:           GravityWGS72,
			time:                time.Date(2020, 5, 23, 20, 23, 37, 0, time.UTC),
			latitudeDegree:      55.6167,
			longitudeDegree:     12.6500,
			altitude:            0.005,
			wantAzimuthDegree:   181.2902281625632,
			wantElevationDegree: 42.06164214709452,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sat, err := TLEToSat(tt.line1, tt.line2, tt.gravConst)
			if err != nil {
				t.Fatalf("ParseTLE() error = %v", err)
			}

			jday := JDayTime(tt.time)

			pos, _, err := Propagate(sat, tt.time)
			if err != nil {
				t.Errorf("sat.Propagate() error = %v", err)
			}

			coordinates := Coordinates{
				Latitude:  tt.latitudeDegree * DEG2RAD,
				Longitude: tt.longitudeDegree * DEG2RAD,
				Altitude:  tt.altitude,
			}

			lookAngles := ECIToLookAngles(pos, coordinates, jday, sat.GravityConst)

			if !closeFloat(lookAngles.Azimuth*RAD2DEG, tt.wantAzimuthDegree) {
				t.Errorf("ECItoLookAngles() Azimuth = %v, want %v", lookAngles.Azimuth*RAD2DEG, tt.wantAzimuthDegree)
			}

			if !closeFloat(lookAngles.Elevation*RAD2DEG, tt.wantElevationDegree) {
				t.Errorf("ECItoLookAngles() Elevation = %v, want %v", lookAngles.Elevation*RAD2DEG, tt.wantElevationDegree)
			}
		})
	}
}

package satellite

import (
	"errors"
	"testing"
)

func TestParseTLE(t *testing.T) {
	tests := []struct {
		name        string
		line1       string
		line2       string
		gravConst   Gravity
		expectedErr error
		satNum      string
		epochyr     int64
		epochdays   float64
		ndot        float64
		nddot       float64
		bstar       float64
		inclo       float64
		nodeo       float64
		ecco        float64
		argpo       float64
		mo          float64
		no          float64
	}{
		{
			name:        "ISS#25544",
			line1:       "1 25544U 98067A   08264.51782528 -.00002182  00000-0 -11606-4 0  2927",
			line2:       "2 25544  51.6416 247.4627 0006703 130.5360 325.0288 15.72125391563537",
			expectedErr: nil,
			satNum:      "25544",
			epochyr:     8,
			epochdays:   264.51782528,
			ndot:        -2.182e-05,
			nddot:       0,
			bstar:       -1.1606e-05,
			inclo:       51.6416,
			nodeo:       247.4627,
			ecco:        0.0006703,
			argpo:       130.536,
			mo:          325.0288,
			no:          15.72125391,
		},
		{
			name:        "NOAA 19#33591",
			line1:       "1 33591U 09005A   16163.48990228  .00000077  00000-0  66998-4 0  9990",
			line2:       "2 33591  99.0394 120.2160 0013054 232.8317 127.1662 14.12079902378332",
			expectedErr: nil,
			satNum:      "33591",
			epochyr:     16,
			epochdays:   163.48990228,
			ndot:        7.7e-07,
			nddot:       0,
			bstar:       .66998e-4,
			inclo:       99.0394,
			nodeo:       120.216,
			ecco:        0.0013054,
			argpo:       232.8317,
			mo:          127.1662,
			no:          14.12079902,
		},
		{
			name:        "TITAN 3C#04632",
			line1:       "1 04632U 70093B   04031.91070959 -.00000084  00000-0  10000-3 0  9955",
			line2:       "2 04632  11.4628 273.1101 1450506 207.6000 143.9350  1.20231981 44145",
			expectedErr: nil,
			satNum:      "04632",
			epochyr:     4,
			epochdays:   31.91070959,
			ndot:        -8.4e-07,
			nddot:       0,
			bstar:       .1e-3,
			inclo:       11.4628,
			nodeo:       273.1101,
			ecco:        0.1450506,
			argpo:       207.6,
			mo:          143.935,
			no:          1.20231981,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tle, err := ParseTLE(test.line1, test.line2)
			if err == nil && test.expectedErr != nil {
				t.Fatalf("expected error %v, got nil", test.expectedErr)
			}
			if err != nil && test.expectedErr == nil {
				t.Fatalf("expected nil, got error %v", err)
			}
			if err != nil && test.expectedErr != nil && !errors.Is(err, test.expectedErr) {
				t.Fatalf("expected error %v, got %v", test.expectedErr, err)
			}
			if tle.CatalogNumber != test.satNum {
				t.Fatalf("expected satnum %s, got %s", test.satNum, tle.CatalogNumber)
			}
			if tle.EpochYear != test.epochyr {
				t.Fatalf("expected epochyr %d, got %d", test.epochyr, tle.EpochYear)
			}
			if tle.EpochDay != test.epochdays {
				t.Fatalf("expected epochdays %f, got %f", test.epochdays, tle.EpochDay)
			}
			if tle.FirstTimeDerivativeOfMeanMotion != test.ndot {
				t.Fatalf("expected ndot %f, got %f", test.ndot, tle.FirstTimeDerivativeOfMeanMotion)
			}
			if tle.SecondTimeDerivativeOfMeanMotion != test.nddot {
				t.Fatalf("expected nddot %f, got %f", test.nddot, tle.SecondTimeDerivativeOfMeanMotion)
			}
			if tle.BStar != test.bstar {
				t.Fatalf("expected bstar %f, got %f", test.bstar, tle.BStar)
			}
			if tle.Inclination != test.inclo {
				t.Fatalf("expected inclo %f, got %f", test.inclo, tle.Inclination)
			}
			if tle.RightAscensionOfAscendingNode != test.nodeo {
				t.Fatalf("expected nodeo %f, got %f", test.nodeo, tle.RightAscensionOfAscendingNode)
			}
			if tle.Eccentricity != test.ecco {
				t.Fatalf("expected ecco %f, got %f", test.ecco, tle.Eccentricity)
			}
			if tle.ArgumentOfPerigee != test.argpo {
				t.Fatalf("expected argpo %f, got %f", test.argpo, tle.ArgumentOfPerigee)
			}
			if tle.MeanAnomaly != test.mo {
				t.Fatalf("expected mo %f, got %f", test.mo, tle.MeanAnomaly)
			}
			if tle.MeanMotion != test.no {
				t.Fatalf("expected no %f, got %f", test.no, tle.MeanMotion)
			}
		})
	}
}

package satellite

import (
	"fmt"
	"math"
)

// GravConst holds variables that are dependent upon selected gravity model.
type GravConst struct {
	mu            float64
	radiusearthkm float64
	xke           float64
	tumin         float64
	j2            float64
	j3            float64
	j4            float64
	j3oj2         float64
	flattening    float64
}

type Gravity string

const (
	GravityWGS72Old Gravity = "wgs72old"
	GravityWGS72    Gravity = "wgs72"
	GravityWGS84    Gravity = "wgs84"
)

// Returns a GravConst with correct information on requested model provided through the name parameter
func getGravConst(name Gravity) (GravConst, error) {
	var grav GravConst
	switch name {
	case GravityWGS72Old:
		grav.mu = 398600.79964
		grav.radiusearthkm = 6378.135
		grav.xke = 0.0743669161
		grav.tumin = 1.0 / grav.xke
		grav.j2 = 0.001082616
		grav.j3 = -0.00000253881
		grav.j4 = -0.00000165597
		grav.j3oj2 = grav.j3 / grav.j2
		grav.flattening = 1 / 298.26
	case GravityWGS72:
		grav.mu = 398600.8
		grav.radiusearthkm = 6378.135
		grav.xke = 60.0 / math.Sqrt(grav.radiusearthkm*grav.radiusearthkm*grav.radiusearthkm/grav.mu)
		grav.tumin = 1.0 / grav.xke
		grav.j2 = 0.001082616
		grav.j3 = -0.00000253881
		grav.j4 = -0.00000165597
		grav.j3oj2 = grav.j3 / grav.j2
		grav.flattening = 1 / 298.26
	case GravityWGS84:
		grav.mu = 398600.5
		grav.radiusearthkm = EQUATOR_RADIUS
		grav.xke = 60.0 / math.Sqrt(grav.radiusearthkm*grav.radiusearthkm*grav.radiusearthkm/grav.mu)
		grav.tumin = 1.0 / grav.xke
		grav.j2 = 0.00108262998905
		grav.j3 = -0.00000253215306
		grav.j4 = -0.00000161098761
		grav.j3oj2 = grav.j3 / grav.j2
		grav.flattening = 1 / 298.257223563
	default:
		return grav, fmt.Errorf("'%s' is not a valid gravity model", name)
	}

	return grav, nil
}

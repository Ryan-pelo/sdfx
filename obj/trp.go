//-----------------------------------------------------------------------------
/*

Truncated Rectangular Pyramid

This a rectangular base pyramid that has rounded edges and has been truncated.

It's an attractive object in its own right, but it's particularly useful for
sand-casting patterns because the slope implements a pattern draft and the
rounded edges minimise sand crumbling.

*/
//-----------------------------------------------------------------------------

package obj

import (
	"errors"
	"math"

	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

// TruncRectPyramidParms defines the parameters for a truncated rectangular pyramid.
type TruncRectPyramidParms struct {
	Size        sdf.V3  // size of truncated pyramid
	BaseAngle   float64 // base angle of pyramid (radians)
	BaseRadius  float64 // base corner radius
	RoundRadius float64 // edge rounding radius
}

// TruncRectPyramid3D returns a truncated rectangular pyramid with rounded edges.
func TruncRectPyramid3D(k *TruncRectPyramidParms) (sdf.SDF3, error) {
	if k.Size.LessThanZero() {
		return nil, errors.New("size vector components < 0")
	}
	if k.BaseAngle <= 0 || k.BaseAngle > sdf.DtoR(90) {
		return nil, errors.New("base angle must be (0,90] degrees")
	}
	if k.BaseRadius < 0 {
		return nil, errors.New("base radius < 0")
	}
	if k.RoundRadius < 0 {
		return nil, errors.New("round radius < 0")
	}
	h := k.Size.Z
	dr := h / math.Tan(k.BaseAngle)
	rb := k.BaseRadius + dr
	rt := sdf.Max(k.BaseRadius-dr, 0)
	round := sdf.Min(0.5*rt, k.RoundRadius)
	s := sdf.Cone3D(2.0*h, rb, rt, round)
	wx := sdf.Max(k.Size.X-2.0*k.BaseRadius, 0)
	wy := sdf.Max(k.Size.Y-2.0*k.BaseRadius, 0)
	s = sdf.Elongate3D(s, sdf.V3{wx, wy, 0})
	s = sdf.Cut3D(s, sdf.V3{0, 0, 0}, sdf.V3{0, 0, 1})
	return s, nil
}

//-----------------------------------------------------------------------------

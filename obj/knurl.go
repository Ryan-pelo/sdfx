//-----------------------------------------------------------------------------
/*

Knurled Cylinders

See: https://en.wikipedia.org/wiki/Knurling

This code builds a knurl with the intersection of left and right hand
multistart screw "threads".

*/
//-----------------------------------------------------------------------------

package obj

import (
	"errors"
	"math"

	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

// KnurlParms specifies the knurl parameters.
type KnurlParms struct {
	Length float64 // length of cylinder
	Radius float64 // radius of cylinder
	Pitch  float64 // knurl pitch
	Height float64 // knurl height
	Theta  float64 // knurl helix angle
}

// knurlProfile returns a 2D knurl profile.
func knurlProfile(k *KnurlParms) sdf.SDF2 {
	knurl := sdf.NewPolygon()
	knurl.Add(k.Pitch/2, 0)
	knurl.Add(k.Pitch/2, k.Radius)
	knurl.Add(0, k.Radius+k.Height)
	knurl.Add(-k.Pitch/2, k.Radius)
	knurl.Add(-k.Pitch/2, 0)
	//knurl.Render("knurl.dxf")
	return sdf.Polygon2D(knurl.Vertices())
}

// Knurl3D returns a knurled cylinder.
func Knurl3D(k *KnurlParms) (sdf.SDF3, error) {
	if k.Length <= 0 {
		return nil, errors.New("length <= 0")
	}
	if k.Radius <= 0 {
		return nil, errors.New("radius <= 0")
	}
	if k.Pitch <= 0 {
		return nil, errors.New("pitch <= 0")
	}
	if k.Height <= 0 {
		return nil, errors.New("height <= 0")
	}
	if k.Theta < 0 {
		return nil, errors.New("theta < 0")
	}
	if k.Theta >= sdf.DtoR(90) {
		return nil, errors.New("theta >= 90")
	}
	// Work out the number of starts using the desired helix angle.
	n := int(sdf.Tau * k.Radius * math.Tan(k.Theta) / k.Pitch)
	// build the knurl profile.
	knurl2d := knurlProfile(k)
	// create the left/right hand spirals
	knurl0_3d := sdf.Screw3D(knurl2d, k.Length, k.Pitch, n)
	knurl1_3d := sdf.Screw3D(knurl2d, k.Length, k.Pitch, -n)
	return sdf.Intersect3D(knurl0_3d, knurl1_3d), nil
}

// KnurledHead3D returns a generic cylindrical knurled head.
func KnurledHead3D(
	r float64, // radius
	h float64, // height
	pitch float64, // knurl pitch
) (sdf.SDF3, error) {
	cylinderRound := r * 0.05
	knurlLength := pitch * math.Floor((h-cylinderRound)/pitch)
	k := KnurlParms{
		Length: knurlLength,
		Radius: r,
		Pitch:  pitch,
		Height: pitch * 0.3,
		Theta:  sdf.DtoR(45),
	}
	knurl3d, err := Knurl3D(&k)
	if err != nil {
		return nil, err
	}
	return sdf.Union3D(sdf.Cylinder3D(h, r, cylinderRound), knurl3d), nil
}

//-----------------------------------------------------------------------------

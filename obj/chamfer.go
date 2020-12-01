//-----------------------------------------------------------------------------
/*

Chamfered Cylinder

*/
//-----------------------------------------------------------------------------

package obj

import "github.com/deadsy/sdfx/sdf"

//-----------------------------------------------------------------------------

// ChamferedCylinder intersects a chamfered cylinder with an SDF3.
func ChamferedCylinder(s sdf.SDF3, kb, kt float64) sdf.SDF3 {
	// get the length and radius from the bounding box
	l := s.BoundingBox().Max.Z
	r := s.BoundingBox().Max.X
	p := sdf.NewPolygon()
	p.Add(0, -l)
	p.Add(r, -l).Chamfer(r * kb)
	p.Add(r, l).Chamfer(r * kt)
	p.Add(0, l)
	return sdf.Intersect3D(s, sdf.Revolve3D(sdf.Polygon2D(p.Vertices())))
}

//-----------------------------------------------------------------------------

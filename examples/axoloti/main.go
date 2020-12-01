//-----------------------------------------------------------------------------
/*

Axoloti Board Mounting Kit

*/
//-----------------------------------------------------------------------------

package main

import (
	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

// material shrinkage
var shrink = 1.0 / 0.999 // PLA ~0.1%
//var shrink = 1.0/0.995; // ABS ~0.5%

//-----------------------------------------------------------------------------

var frontPanelThickness = 3.0
var frontPanelLength = 170.0
var frontPanelHeight = 50.0
var frontPanelYOffset = 15.0

var baseWidth = 50.0
var baseLength = 170.0
var baseThickness = 3.0

var baseFootWidth = 10.0
var baseFootCornerRadius = 3.0

var pcbWidth = 50.0
var pcbLength = 160.0

var pillarHeight = 16.8

//-----------------------------------------------------------------------------

// multiple standoffs
func standoffs() sdf.SDF3 {

	k := &obj.StandoffParms{
		PillarHeight:   pillarHeight,
		PillarDiameter: 6.0,
		HoleDepth:      10.0,
		HoleDiameter:   2.4,
	}

	zOfs := 0.5 * (pillarHeight + baseThickness)

	// from the board mechanicals
	positions := sdf.V3Set{
		{3.5, 10.0, zOfs},   // H1
		{3.5, 40.0, zOfs},   // H2
		{54.0, 40.0, zOfs},  // H3
		{156.5, 10.0, zOfs}, // H4
		//{54.0, 10.0, zOfs},  // H5
		{156.5, 40.0, zOfs}, // H6
		{44.0, 10.0, zOfs},  // H7
		{116.0, 10.0, zOfs}, // H8
	}

	return sdf.Multi3D(obj.Standoff3D(k), positions)
}

//-----------------------------------------------------------------------------

// base returns the base mount.
func base() sdf.SDF3 {
	// base
	pp := &obj.PanelParms{
		Size:         sdf.V2{baseLength, baseWidth},
		CornerRadius: 5.0,
		HoleDiameter: 3.5,
		HoleMargin:   [4]float64{7.0, 20.0, 7.0, 20.0},
		HolePattern:  [4]string{"xx", "x", "xx", "x"},
	}
	s0 := obj.Panel2D(pp)

	// cutout
	l := baseLength - (2.0 * baseFootWidth)
	w := 18.0
	s1 := sdf.Box2D(sdf.V2{l, w}, baseFootCornerRadius)
	yOfs := 0.5 * (baseWidth - pcbWidth)
	s1 = sdf.Transform2D(s1, sdf.Translate2d(sdf.V2{0, yOfs}))

	s2 := sdf.Extrude3D(sdf.Difference2D(s0, s1), baseThickness)
	xOfs := 0.5 * pcbLength
	yOfs = pcbWidth - (0.5 * baseWidth)
	s2 = sdf.Transform3D(s2, sdf.Translate3d(sdf.V3{xOfs, yOfs, 0}))

	// standoffs
	s3 := standoffs()

	s4 := sdf.Union3D(s2, s3)
	s4.(*sdf.UnionSDF3).SetMin(sdf.PolyMin(3.0))

	return s4
}

//-----------------------------------------------------------------------------
// front panel cutouts

type panelHole struct {
	center sdf.V2   // center of hole
	hole   sdf.SDF2 // 2d hole
}

// button positions
var pbX = 53.0
var pb0 = sdf.V2{pbX, 0.8}
var pb1 = sdf.V2{pbX + 5.334, 0.8}

// panelCutouts returns the 2D front panel cutouts
func panelCutouts() sdf.SDF2 {

	sMidi := sdf.Circle2D(0.5 * 17.0)
	sJack := sdf.Circle2D(0.5 * 11.5)
	sLed := sdf.Box2D(sdf.V2{1.6, 1.6}, 0)

	fb := &obj.FingerButtonParms{
		Width:  4.0,
		Gap:    0.6,
		Length: 20.0,
	}
	sButton := sdf.Transform2D(obj.FingerButton2D(fb), sdf.Rotate2d(sdf.DtoR(-90)))

	jackX := 123.0
	midiX := 18.8
	ledX := 62.9

	holes := []panelHole{
		{sdf.V2{midiX, 10.2}, sMidi},                         // MIDI DIN Jack
		{sdf.V2{midiX + 20.32, 10.2}, sMidi},                 // MIDI DIN Jack
		{sdf.V2{jackX, 8.14}, sJack},                         // 1/4" Stereo Jack
		{sdf.V2{jackX + 19.5, 8.14}, sJack},                  // 1/4" Stereo Jack
		{sdf.V2{107.6, 2.3}, sdf.Circle2D(0.5 * 5.5)},        // 3.5 mm Headphone Jack
		{sdf.V2{ledX, 0.5}, sLed},                            // LED
		{sdf.V2{ledX + 3.635, 0.5}, sLed},                    // LED
		{pb0, sButton},                                       // Push Button
		{pb1, sButton},                                       // Push Button
		{sdf.V2{84.1, 1.0}, sdf.Box2D(sdf.V2{16.0, 7.5}, 0)}, // micro SD card
		{sdf.V2{96.7, 1.0}, sdf.Box2D(sdf.V2{11.0, 7.5}, 0)}, // micro USB connector
		{sdf.V2{73.1, 7.1}, sdf.Box2D(sdf.V2{7.5, 15.0}, 0)}, // fullsize USB connector
	}

	s := make([]sdf.SDF2, len(holes))
	for i, k := range holes {
		s[i] = sdf.Transform2D(k.hole, sdf.Translate2d(k.center))
	}

	return sdf.Union2D(s...)
}

//-----------------------------------------------------------------------------

// frontPanel returns the front panel mount.
func frontPanel() sdf.SDF3 {

	// overall panel
	pp := &obj.PanelParms{
		Size:         sdf.V2{frontPanelLength, frontPanelHeight},
		CornerRadius: 5.0,
		HoleDiameter: 3.5,
		HoleMargin:   [4]float64{5.0, 5.0, 5.0, 5.0},
		HolePattern:  [4]string{"xx", "x", "xx", "x"},
	}
	panel := obj.Panel2D(pp)

	xOfs := 0.5 * pcbLength
	yOfs := (0.5 * frontPanelHeight) - frontPanelYOffset
	panel = sdf.Transform2D(panel, sdf.Translate2d(sdf.V2{xOfs, yOfs}))

	// extrude to 3d
	fp := sdf.Extrude3D(sdf.Difference2D(panel, panelCutouts()), frontPanelThickness)

	// Add buttons to the finger button
	bHeight := 4.0
	b := sdf.Cylinder3D(bHeight, 1.4, 0)
	b0 := sdf.Transform3D(b, sdf.Translate3d(pb0.ToV3(-0.5*bHeight)))
	b1 := sdf.Transform3D(b, sdf.Translate3d(pb1.ToV3(-0.5*bHeight)))

	return sdf.Union3D(fp, b0, b1)
}

//-----------------------------------------------------------------------------

// mountingKit creates the STLs for the axoloti mount kit
func mountingKit() {

	// front panel
	s0 := frontPanel()
	sx := sdf.Transform3D(s0, sdf.RotateY(sdf.DtoR(180.0)))
	sdf.RenderSTL(sdf.ScaleUniform3D(sx, shrink), 400, "panel.stl")

	// base
	s1 := base()
	sdf.RenderSTL(sdf.ScaleUniform3D(s1, shrink), 400, "base.stl")

	// both together
	s0 = sdf.Transform3D(s0, sdf.Translate3d(sdf.V3{0, 80, 0}))
	s3 := sdf.Union3D(s0, s1)
	sdf.RenderSTL(sdf.ScaleUniform3D(s3, shrink), 400, "panel_and_base.stl")
}

//-----------------------------------------------------------------------------

func main() {
	mountingKit()
}

//-----------------------------------------------------------------------------

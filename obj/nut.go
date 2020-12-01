//-----------------------------------------------------------------------------
/*

Nuts: Simple Nut for 3d printing.

*/
//-----------------------------------------------------------------------------

package obj

import (
	"errors"
	"fmt"

	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

// NutParms defines the parameters for a nut.
type NutParms struct {
	Thread    string  // name of thread
	Style     string  // head style "hex" or "knurl"
	Tolerance float64 // add to internal thread radius
}

// Nut returns a simple nut suitable for 3d printing.
func Nut(k *NutParms) (sdf.SDF3, error) {
	// validate parameters
	t, err := sdf.ThreadLookup(k.Thread)
	if err != nil {
		return nil, err
	}
	if k.Tolerance < 0 {
		return nil, errors.New("tolerance < 0")
	}

	// nut body
	var nut sdf.SDF3
	nr := t.HexRadius()
	nh := t.HexHeight()
	switch k.Style {
	case "hex":
		nut, err = HexHead3D(nr, nh, "tb")
	case "knurl":
		nut, err = KnurledHead3D(nr, nh, nr*0.25)
	default:
		return nil, fmt.Errorf("unknown style \"%s\"", k.Style)
	}
	if err != nil {
		return nil, err
	}

	// internal thread
	isoThread := sdf.ISOThread(t.Radius+k.Tolerance, t.Pitch, false)
	thread := sdf.Screw3D(isoThread, nh, t.Pitch, 1)

	return sdf.Difference3D(nut, thread), nil
}

//-----------------------------------------------------------------------------

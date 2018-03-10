//-----------------------------------------------------------------------------
/*

Text Example

*/
//-----------------------------------------------------------------------------

package main

import (
	"fmt"
	"os"

	. "github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

func main() {

	f, err := LoadFont("cmr10.ttf")
	//f, err := LoadFont("Times_New_Roman.ttf")

	if err != nil {
		fmt.Printf("can't read font file %s\n", err)
		os.Exit(1)
	}

	t := NewText("SDFX! Hello World!")

	s2d, err := TextSDF2(f, t)
	if err != nil {
		fmt.Printf("can't generate text sdf2 %s\n", err)
		os.Exit(1)
	}

	RenderDXF(s2d, 600, "shape.dxf")

	s3d := ExtrudeRounded3D(s2d, 400, 20)
	RenderSTL(s3d, 600, "shape.stl")
}

//-----------------------------------------------------------------------------
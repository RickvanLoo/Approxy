package VHDL

import (
	"math"
	"testing"
)

func TestRecursive4(t *testing.T) {
	o1 := New2DUnsignedAcc(2)
	o2 := New2DUnsignedAcc(2)
	o3 := New2DUnsignedAcc(2)
	o4 := New2DUnsignedAcc(2)
	RecLutArray := [4]*LUT2D{o1, o2, o3, o4}
	rec4 := NewRecursive4(RecLutArray)

	maxval := int(math.Exp2(4))

	for a := 0; a < maxval; a++ {
		for b := 0; b < maxval; b++ {

			out := rec4.ReturnVal(uint(a), uint(b))
			test := a * b

			if out != uint(test) {
				t.Errorf("!!ERROR!!: %d * %d != %d!!", a, b, out)
			} else {
				t.Logf("PASS: %d * %d == %d", a, b, out)
			}

		}
	}
}

package VHDL

import (
	"math"
	"testing"
)

func TestRecursive8(t *testing.T) {
	o1 := New2DUnsignedAcc("o1", 2)
	o2 := New2DUnsignedAcc("o2", 2)
	o3 := New2DUnsignedAcc("o3", 2)
	o4 := New2DUnsignedAcc("o4", 2)
	RecLutArray := [4]*LUT2D{o1, o2, o3, o4}
	rec4 := NewRecursive4("testRec4", RecLutArray)
	rec4array := [4]*Recursive4{rec4, rec4, rec4, rec4}
	rec8 := NewRecursive8("testRec8", rec4array)
	pass := true

	maxval := int(math.Exp2(8))

	for a := 0; a < maxval; a++ {
		for b := 0; b < maxval; b++ {

			out := rec8.ReturnVal(uint(a), uint(b))
			test := a * b

			if out != uint(test) {
				t.Errorf("!!ERROR!!: %d * %d != %d!!", a, b, out)
				pass = false
			}

		}
	}

	if pass {
		t.Logf("PASS ALL: No Errors detected for Rec8")
	}
}

func TestRecursive8Overflow(t *testing.T) {
	m3 := M3().LUT2D
	RecLutArray := [4]*LUT2D{m3, m3, m3, m3} //Known Overflow Configuration
	rec4 := NewRecursive4("testRec4OverFlow", RecLutArray)
	rec4LutArray := [4]*Recursive4{rec4, rec4, rec4, rec4}
	rec8 := NewRecursive8("testRec8Overflow", rec4LutArray)

	maxval := int(math.Exp2(8))

	for a := 0; a < maxval; a++ {
		for b := 0; b < maxval; b++ {
			rec8.ReturnVal(uint(a), uint(b))
		}
	}

	if !rec8.OverflowError {
		t.Errorf("!!ERROR!! Overflow malfunction for rec4")
	} else {
		t.Logf("PASS: Overflow detected for known Overflow Configuration")
	}
}

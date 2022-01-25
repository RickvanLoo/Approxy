package VHDL

import (
	"math"
	"testing"
)

func TestRecursive4(t *testing.T) {
	o1 := New2DUnsignedAcc("o1", 2)
	o2 := New2DUnsignedAcc("o2", 2)
	o3 := New2DUnsignedAcc("o3", 2)
	o4 := New2DUnsignedAcc("o4", 2)
	RecLutArray := [4]*LUT2D{o1, o2, o3, o4}
	rec4 := NewRecursive4("testRec4", RecLutArray)
	pass := true

	maxval := int(math.Exp2(4))

	for a := 0; a < maxval; a++ {
		for b := 0; b < maxval; b++ {

			out := rec4.ReturnVal(uint(a), uint(b))
			test := a * b

			if out != uint(test) {
				t.Errorf("!!ERROR!!: %d * %d != %d!!", a, b, out)
				pass = false
			}
		}
	}

	if pass {
		t.Logf("PASS ALL: No Errors detected for Rec4")
	}
}

func TestRecursive4Overflow(t *testing.T) {
	m3 := M3().LUT2D
	RecLutArray := [4]*LUT2D{m3, m3, m3, m3} //Known Overflow Configuration
	rec4 := NewRecursive4("testRec4OverFlow", RecLutArray)

	maxval := int(math.Exp2(4))

	for a := 0; a < maxval; a++ {
		for b := 0; b < maxval; b++ {
			rec4.ReturnVal(uint(a), uint(b))
		}
	}

	if !rec4.OverflowError {
		t.Errorf("!!ERROR!! Overflow malfunction for rec4")
	} else {
		t.Logf("PASS: Overflow detected for known Overflow Configuration")
	}
}

func TestOutputArrayRec4(t *testing.T) {
	o1 := New2DUnsignedAcc("o1", 2)
	o2 := New2DUnsignedAcc("o2", 2)
	o3 := New2DUnsignedAcc("o3", 2)
	o4 := New2DUnsignedAcc("o4", 2)
	RecLutArray := [4]*LUT2D{o1, o2, o3, o4}
	rec4 := NewRecursive4("testRec4", RecLutArray)

	output := rec4.GenerateVHDLEntityArray()

	top := output[0].ReturnData().EntityName
	if top != "testRec4" {
		t.Errorf("Expected for 0: testRec4")
	}
}

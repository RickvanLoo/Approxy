package VHDL

import (
	"testing"
)

func Benchmark4BitMac(b *testing.B) {
	m1 := M1().LUT2D
	m2 := M2().LUT2D
	m3 := M3().LUT2D
	m4 := M4().LUT2D
	rec4 := NewRecursive4("ApproxRec4", [4]*LUT2D{m1, m2, m3, m4})

	for n := 0; n < b.N; n++ {
		NewMAC(rec4, 100)
	}
}

func Benchmark8BitMac(b *testing.B) {
	m1 := M1().LUT2D
	m2 := M2().LUT2D
	m3 := M3().LUT2D
	m4 := M4().LUT2D
	rec4 := NewRecursive4("ApproxRec4", [4]*LUT2D{m1, m2, m3, m4})
	rec8 := NewRecursive8("AproxRec8", [4]*Recursive4{rec4, rec4, rec4, rec4})

	for n := 0; n < b.N; n++ {
		NewMAC(rec8, 100)
	}
}

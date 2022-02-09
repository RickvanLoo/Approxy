package VHDL

import (
	"testing"
)

func Benchmark4BitMac(b *testing.B) {
	m1 := M1().LUT2D
	m2 := M2().LUT2D
	m3 := M3().LUT2D
	m4 := M4().LUT2D
	rec4 := NewRecursive4("ApproxRec4", [4]VHDLEntityMultiplier{m1, m2, m3, m4})

	for n := 0; n < b.N; n++ {
		NewMAC(rec4, 100)
	}
}

func Benchmark8BitMac(b *testing.B) {
	m1 := M1().LUT2D
	m2 := M2().LUT2D
	m3 := M3().LUT2D
	m4 := M4().LUT2D
	rec4 := NewRecursive4("ApproxRec4", [4]VHDLEntityMultiplier{m1, m2, m3, m4})
	rec8 := NewRecursive8("AproxRec8", [4]VHDLEntityMultiplier{rec4, rec4, rec4, rec4})

	for n := 0; n < b.N; n++ {
		NewMAC(rec8, 100)
	}
}

func TestSimpleMac(t *testing.T) {
	acc := New2DUnsignedAcc("Acc", 4)
	mac := NewMAC(acc, 1)

	if mac.CurrentValue != 0 {
		t.Errorf("Init not 0!")
	}

	if mac.ReturnVal(2, 2) != 4 {
		t.Errorf("First val is not 4!")
	}

	if mac.ReturnVal(2, 2) != 8 {
		t.Errorf("Second val is not 8!")
	}

	if mac.ReturnVal(2, 2) != 12 {
		t.Errorf("Third val is not 12!")
	}

	mac.ResetVal()

	if mac.CurrentValue != 0 {
		t.Errorf("Reset not 0!")
	}
}

func TestOverflowMac(t *testing.T) {
	acc := New2DUnsignedAcc("Acc", 4)
	mac := NewMAC(acc, 4)

	if mac.CurrentValue != 0 {
		t.Errorf("Init not 0!")
	}

	val := mac.ReturnVal(15, 15)
	if val != 225 {
		t.Errorf("First val is not 225! rec:%d", val)
	}

	val = mac.ReturnVal(15, 15)
	if val != 450 {
		t.Errorf("First val is not 450! rec:%d", val)
	}

	val = mac.ReturnVal(15, 15)
	if val != 675 {
		t.Errorf("First val is not 675! rec:%d", val)
	}

	val = mac.ReturnVal(15, 15)
	if val != 900 {
		t.Errorf("First val is not 900! rec:%d", val)
	}

	val = mac.ReturnVal(15, 15)
	if val == 1125 {
		t.Errorf("MAC Should overflow, not return 1125! rec:%d", val)
	}

	overflowcheck, _ := OverflowCheckGeneric(val, mac.OutputSize)
	if val != overflowcheck {
		t.Errorf("Overflow Expected %d, rec %d", overflowcheck, val)
	}

	mac.ResetVal()

	if mac.CurrentValue != 0 {
		t.Errorf("Reset not 0!")
	}
}

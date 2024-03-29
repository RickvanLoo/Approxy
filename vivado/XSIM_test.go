package vivado

import (
	"testing"

	VHDL "github.com/RickvanLoo/Approxy/vhdl"
)

func TestXSIMDuplicate2(t *testing.T) {
	m1 := VHDL.M1().LUT2D
	rec4 := VHDL.NewRecursive4("ApproxRec4", [4]VHDL.VHDLEntityMultiplier{m1, m1, m1, m1})

	xsim := CreateXSIM("Output", "TestName", rec4.GenerateVHDLEntityArray())
	lengthvhdlentityXSIM := len(xsim.VHDLEntities)

	if lengthvhdlentityXSIM != 2 {
		t.Errorf("Duplication Issue!, l=%d, expect=2", lengthvhdlentityXSIM)
	}

	topEntityName := xsim.VHDLEntities[0].ReturnData().EntityName
	if topEntityName != "ApproxRec4" {
		t.Errorf("Expected ApproxRec4, got %s", topEntityName)
	}

}

func TestXSIMDuplicate3(t *testing.T) {
	m1 := VHDL.M1().LUT2D
	m2 := VHDL.M2().LUT2D

	rec4 := VHDL.NewRecursive4("ApproxRec4", [4]VHDL.VHDLEntityMultiplier{m1, m2, m1, m2})

	xsim := CreateXSIM("Output", "TestName", rec4.GenerateVHDLEntityArray())
	lengthvhdlentityXSIM := len(xsim.VHDLEntities)

	if lengthvhdlentityXSIM != 3 {
		t.Errorf("Duplication Issue!, l=%d, expect=3", lengthvhdlentityXSIM)
	}

	topEntityName := xsim.VHDLEntities[0].ReturnData().EntityName
	if topEntityName != "ApproxRec4" {
		t.Errorf("Expected ApproxRec4, got %s", topEntityName)
	}
}

func TestXSIMDuplicate4(t *testing.T) {
	m1 := VHDL.M1().LUT2D
	m2 := VHDL.M2().LUT2D
	m3 := VHDL.M3().LUT2D

	rec4 := VHDL.NewRecursive4("ApproxRec4", [4]VHDL.VHDLEntityMultiplier{m1, m2, m3, m2})

	xsim := CreateXSIM("Output", "TestName", rec4.GenerateVHDLEntityArray())
	lengthvhdlentityXSIM := len(xsim.VHDLEntities)

	if lengthvhdlentityXSIM != 4 {
		t.Errorf("Duplication Issue!, l=%d, expect=4", lengthvhdlentityXSIM)
	}

	topEntityName := xsim.VHDLEntities[0].ReturnData().EntityName
	if topEntityName != "ApproxRec4" {
		t.Errorf("Expected ApproxRec4, got %s", topEntityName)
	}
}

func TestXSIMDuplicate5(t *testing.T) {
	m1 := VHDL.M1().LUT2D
	m2 := VHDL.M2().LUT2D
	m3 := VHDL.M3().LUT2D
	m4 := VHDL.M4().LUT2D

	rec4 := VHDL.NewRecursive4("ApproxRec4", [4]VHDL.VHDLEntityMultiplier{m1, m2, m3, m4})

	xsim := CreateXSIM("Output", "TestName", rec4.GenerateVHDLEntityArray())
	lengthvhdlentityXSIM := len(xsim.VHDLEntities)

	if lengthvhdlentityXSIM != 5 {
		t.Errorf("Duplication Issue!, l=%d, expect=5", lengthvhdlentityXSIM)
	}

	topEntityName := xsim.VHDLEntities[0].ReturnData().EntityName
	if topEntityName != "ApproxRec4" {
		t.Errorf("Expected ApproxRec4, got %s", topEntityName)
	}
}

func TestXSIMDuplicate7(t *testing.T) {
	m1 := VHDL.M1().LUT2D
	m2 := VHDL.M2().LUT2D
	m3 := VHDL.M3().LUT2D
	m4 := VHDL.M4().LUT2D

	rec4_1 := VHDL.NewRecursive4("ApproxRec4_1", [4]VHDL.VHDLEntityMultiplier{m1, m2, m3, m4})
	rec4_2 := VHDL.NewRecursive4("ApproxRec4_2", [4]VHDL.VHDLEntityMultiplier{m4, m3, m2, m1})
	rec8 := VHDL.NewRecursive8("ApproxRec8", [4]VHDL.VHDLEntityMultiplier{rec4_1, rec4_2, rec4_1, rec4_2})

	xsim := CreateXSIM("Output", "TestName", rec8.GenerateVHDLEntityArray())
	lengthvhdlentityXSIM := len(xsim.VHDLEntities)

	if lengthvhdlentityXSIM != 7 {
		t.Errorf("Duplication Issue!, l=%d, expect=7", lengthvhdlentityXSIM)
	}

	topEntityName := xsim.VHDLEntities[0].ReturnData().EntityName
	if topEntityName != "ApproxRec8" {
		t.Errorf("Expected ApproxRec8, got %s", topEntityName)
	}

	secEntityName := xsim.VHDLEntities[1].ReturnData().EntityName
	if secEntityName != "ApproxRec4_1" && secEntityName != "ApproxRec4_2" {
		t.Errorf("Expected ApproxRec4_1 or ApproxRec4_2, got %s", secEntityName)
	}

	trdEntityName := xsim.VHDLEntities[2].ReturnData().EntityName
	if trdEntityName != "ApproxRec4_1" && trdEntityName != "ApproxRec4_2" {
		t.Errorf("Expected ApproxRec4_1 or ApproxRec4_2, got %s", trdEntityName)
	}
}

func TestXSIMRec8(t *testing.T) {
	m1 := VHDL.M1().LUT2D
	m2 := VHDL.M1().LUT2D
	m3 := VHDL.M1().LUT2D
	m4 := VHDL.M1().LUT2D

	rec4_1 := VHDL.NewRecursive4("ApproxRec4_1", [4]VHDL.VHDLEntityMultiplier{m1, m2, m3, m4})
	rec4_2 := VHDL.NewRecursive4("ApproxRec4_1", [4]VHDL.VHDLEntityMultiplier{m4, m3, m2, m1})
	rec8 := VHDL.NewRecursive8("ApproxRec8", [4]VHDL.VHDLEntityMultiplier{rec4_1, rec4_2, rec4_1, rec4_2})

	xsim := CreateXSIM("Output", "TestName", rec8.GenerateVHDLEntityArray())
	lengthvhdlentityXSIM := len(xsim.VHDLEntities)

	if lengthvhdlentityXSIM != 3 {
		t.Errorf("Duplication Issue!, l=%d, expect=3", lengthvhdlentityXSIM)
	}

	topEntityName := xsim.VHDLEntities[0].ReturnData().EntityName
	if topEntityName != "ApproxRec8" {
		t.Errorf("Expected ApproxRec8, got %s", topEntityName)
	}

	secEntityName := xsim.VHDLEntities[1].ReturnData().EntityName
	if secEntityName != "ApproxRec4_1" {
		t.Errorf("Expected ApproxRec4_1 got %s", secEntityName)
	}

	trdEntityName := xsim.VHDLEntities[2].ReturnData().EntityName
	if trdEntityName != "M1" {
		t.Errorf("Expected M1, got %s", trdEntityName)
	}
}

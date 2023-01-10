package vhdl

// /M1-M4 are 2-bit Unsigned Approximate Multiplayers defined by Gillani et al.

// M1 creates a 2-bit UnsignedApproxMultiplyer on basis of LUT2D
// Added Modifications:
// 3x3=7
func M1() *UnsignedApproxMultiplyer {
	m := NewUnsignedApprox("M1", 2)
	e7 := Modification{3, 3, 7}
	m.AddModification(e7)
	m.ExecModifications()
	return m
}

// M2 creates a 2-bit UnsignedApproxMultiplyer on basis of LUT2D
// Added Modifications:
// 1x1=0
// 1x3=2
// 3x1=2
func M2() *UnsignedApproxMultiplyer {
	m := NewUnsignedApprox("M2", 2)
	e0 := Modification{1, 1, 0}
	e1 := Modification{1, 3, 2}
	e2 := Modification{3, 1, 2}
	m.AddModification(e0)
	m.AddModification(e1)
	m.AddModification(e2)
	m.ExecModifications()
	return m
}

// M3 creates a 2-bit UnsignedApproxMultiplyer on basis of LUT2D
// 3x3=11
func M3() *UnsignedApproxMultiplyer {
	m := NewUnsignedApprox("M3", 2)
	e11 := Modification{3, 3, 11}
	m.AddModification(e11)
	m.ExecModifications()
	return m
}

// M4 creates a 2-bit UnsignedApproxMultiplyer on basis of LUT2D
// 3x3=5
func M4() *UnsignedApproxMultiplyer {
	m := NewUnsignedApprox("M4", 2)
	e5 := Modification{3, 3, 5}
	m.AddModification(e5)
	m.ExecModifications()
	return m
}

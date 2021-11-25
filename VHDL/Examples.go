package VHDL

///M1-M4 are 2-bit Unsigned Approximate Multiplayers defined by Gillani et al.
func M1() *UnsignedApproxMultiplyer {
	m := NewUnsignedApprox(2)
	m.LUT2D.EntityName = "M1"
	e7 := Modification{3, 3, 7}
	m.AddModication(e7)
	m.ExecModifications()
	return m
}

func M2() *UnsignedApproxMultiplyer {
	m := NewUnsignedApprox(2)
	m.LUT2D.EntityName = "M2"
	e0 := Modification{1, 1, 0}
	e1 := Modification{1, 3, 2}
	e2 := Modification{3, 1, 2}
	m.AddModication(e0)
	m.AddModication(e1)
	m.AddModication(e2)
	m.ExecModifications()
	return m
}

func M3() *UnsignedApproxMultiplyer {
	m := NewUnsignedApprox(2)
	m.LUT2D.EntityName = "M3"
	e11 := Modification{3, 3, 11}
	m.AddModication(e11)
	m.ExecModifications()
	return m
}

func M4() *UnsignedApproxMultiplyer {
	m := NewUnsignedApprox(2)
	m.LUT2D.EntityName = "M4"
	e5 := Modification{3, 3, 5}
	m.AddModication(e5)
	m.ExecModifications()
	return m
}

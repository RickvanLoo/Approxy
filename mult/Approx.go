package mult

import "strconv"

type UnsignedApproxMultiplyer struct {
	Mult *UnsignedAccurateMultiplyer
	Mods []Modification
}

type Modification struct {
	A uint
	B uint
	O uint
}

func NewUnsignedApprox(size uint) *UnsignedApproxMultiplyer {
	m := new(UnsignedApproxMultiplyer)
	m.Mult = NewUnsignedAcc(size)
	m.Mult.Name = "uApprox" + strconv.Itoa(int(size)) + "bitMult"
	return m
}

func (m *UnsignedApproxMultiplyer) AddModication(mod Modification) {
	m.Mods = append(m.Mods, mod)
}

func (m *UnsignedApproxMultiplyer) ExecModifications() {
	for _, mod := range m.Mods {
		m.Mult.changeVal(mod.A, mod.B, mod.O)
	}
}

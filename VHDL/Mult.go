package VHDL

import (
	"math"
	"strconv"
)

type UnsignedApproxMultiplyer struct {
	LUT2D *LUT2D
	Mods  []Modification
}

type Modification struct {
	A uint
	B uint
	O uint
}

func NewUnsignedAcc(BitSize uint) *LUT2D {
	m := new(LUT2D)
	m.BitSize = BitSize
	m.EntityName = "uAcc" + strconv.Itoa(int(BitSize)) + "bitMult"
	m.OutputSize = 2 * m.BitSize
	LUTSize := int(math.Pow(2, float64(BitSize)))

	var LUT [][]uint
	for x := 0; x < LUTSize; x++ {
		var row []uint
		for y := 0; y < LUTSize; y++ {
			row = append(row, uint(x*y))
		}
		LUT = append(LUT, row)
	}

	m.LUT = LUT

	return m
}

func NewUnsignedApprox(BitSize uint) *UnsignedApproxMultiplyer {
	m := new(UnsignedApproxMultiplyer)
	m.LUT2D = NewUnsignedAcc(BitSize)
	m.LUT2D.EntityName = "uApprox" + strconv.Itoa(int(BitSize)) + "bitMult"
	return m
}

func (m *UnsignedApproxMultiplyer) AddModication(mod Modification) {
	m.Mods = append(m.Mods, mod)
}

func (m *UnsignedApproxMultiplyer) ExecModifications() {
	for _, mod := range m.Mods {
		m.LUT2D.changeVal(mod.A, mod.B, mod.O)
	}
}

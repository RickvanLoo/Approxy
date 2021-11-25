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

//New2DUnsignedAcc creates a 2D LUT structure, containing the LUT of an accurate Bitsize multiplyer.
func New2DUnsignedAcc(BitSize uint) *LUT2D {
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

//NewUnsignedApprox creates a UnsignedApproxMultiplayer based upon a 2D LUT structure
//Using New2DUnsignedAcc to create an accurate multiplyer, but adds the option to add modification
//Modification are used to change values within the 2D LUT
func NewUnsignedApprox(BitSize uint) *UnsignedApproxMultiplyer {
	m := new(UnsignedApproxMultiplyer)
	m.LUT2D = New2DUnsignedAcc(BitSize)
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

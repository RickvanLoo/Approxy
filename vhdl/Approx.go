package vhdl

import (
	"math"
)

// UnsignedApproxMultiplyer creates an Approximate Multiplier on basis of the Accurate LUT2D struct, by enabling to add modification to it.
// This struct is simply an extension of the original LUT2D to make it easier to create approximate multiplier structures
// Template is "multbehav.vhd" as defined in LUT2D.go
type UnsignedApproxMultiplyer struct {
	LUT2D *LUT2D
	Mods  []Modification
}

// Modification is struct containing a single modification to a multiplier.
// For Example, we want 3x3=7 instead of 3x3=9, this makes the modification A=3, B=3, O=7
type Modification struct {
	A uint
	B uint
	O uint
}

// New2DUnsignedAcc creates a 2D LUT structure, containing the LUT of an accurate Bitsize multiplyer.
func New2DUnsignedAcc(EntityName string, BitSize uint) *LUT2D {
	m := New2DLUT(EntityName, BitSize)

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

// NewUnsignedApprox creates a UnsignedApproxMultiplayer based upon a 2D LUT structure
// Using New2DUnsignedAcc to create an accurate multiplyer, but adds the option to add modification
// Modification are used to change values within the 2D LUT
func NewUnsignedApprox(EntityName string, BitSize uint) *UnsignedApproxMultiplyer {
	m := new(UnsignedApproxMultiplyer)
	m.LUT2D = New2DUnsignedAcc(EntityName, BitSize)
	return m
}

// AddModification adds a modification to UnsignedApproxMultiplyer
func (m *UnsignedApproxMultiplyer) AddModification(mod Modification) {
	m.Mods = append(m.Mods, mod)
}

// ExecModifications executes all added modifications to the internal LUT2D structure
// Warning: this is final
func (m *UnsignedApproxMultiplyer) ExecModifications() {
	for _, mod := range m.Mods {
		m.LUT2D.changeVal(mod.A, mod.B, mod.O)
	}
}

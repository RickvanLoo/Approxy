//go:build exclude

package vhdl

// BOILERPLATE CODE
// REMOVE EXCLUDE BEFORE USING

// Implements:

// type VHDLEntityMultiplier interface {
// 	VHDLEntity
// 	Multiplier
// }

// // VHDLEntity describes an interface for testable and synthesizable VHDL structures
// type VHDLEntity interface {
// 	ReturnData() *EntityData
// 	GenerateVHDL(string)
// 	GenerateTestData(string)
// 	GenerateVHDLEntityArray() []VHDLEntity
// 	String() string //MSB -> LSB
// }

// // Multiplier describes an interface for generic Multipliers
// type Multiplier interface {
// 	ReturnVal(uint, uint) uint
// 	Overflow() bool
// 	MeanAbsoluteError() float64
// }

// Boilerplate is an empty Approxy model without functionality
type Boilerplate struct {
	EntityName string
	BitSize    uint
	VHDLFile   string
	TestFile   string
}

// NewBoilerplate is a creator
func NewBoilerplate(EntityName string, BitSize uint) *Boilerplate {
	m := new(Boilerplate)
	m.EntityName = EntityName
	m.BitSize = BitSize
	m.VHDLFile, m.TestFile = FileNameGen(m.EntityName)

	return m
}

// ReturnData returns metadata
func (m *Boilerplate) ReturnData() *EntityData {
	d := new(EntityData)
	d.EntityName = m.EntityName
	d.BitSize = m.BitSize
	d.OutputSize = 2 * m.BitSize
	d.VHDLFile = m.VHDLFile
	d.TestFile = m.TestFile
	return d
}

// GenerateVHDL generates VHDL on basis of template.
func (m *Boilerplate) GenerateVHDL(FolderPath string) {
	template := "boilerplate.vhd" //FIXME
	CreateFile(FolderPath, m.VHDLFile, template, m)
}

// GenerateTestData should export ALL behaviour to a textfile. Inputs and output seperated by \t, binary
// Example:
// 0000 0000 0000
// 0000 0001 0000
// 0001 0001 0001
// 0010 0010 1000
// etc.
func (m *Boilerplate) GenerateTestData(FolderPath string) {
	// CreateTestData
}

// GenerateVHDLEntityArray creates an array of potentially multiple VHDLEntities, sorted by priority for synthesizing
// For example: Multiplier A uses a VHDL portmap for the smaller Multiplier B & C. B & C need to be synthesized first, hence A will be last in the array
func (m *Boilerplate) GenerateVHDLEntityArray() []VHDLEntity {
	var out []VHDLEntity
	out = append(out, m)
	return out
}

// For implementing the fmt.Stringer interface
func (m *Boilerplate) String() string {
	return m.EntityName
}

// ReturnVal returns the output of the multiplier
func (m *Boilerplate) ReturnVal(a uint, b uint) uint {
	return a * b
}

// Overflow returns a boolean if any internal overflow has occured
func (m *Boilerplate) Overflow() bool {
	return false
}

// MeanAbsoluteError returns the MeanAbsoluteError of the multiplier in float64
func (m *Boilerplate) MeanAbsoluteError() float64 {
	return 0
}

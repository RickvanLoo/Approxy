package vhdl

// UnsignedNumericAccurateMultiplyer defines a multiplyer based upon 'accuratebehav.vhd', this is a BitSize A*B=prod multiplyer, using the IEEE Numeric lib
// The struct is implemented as a VHDLEntity Interface
type UnsignedNumericAccurateMultiplyer struct {
	EntityName string
	BitSize    uint
	VHDLFile   string
	TestFile   string
}

// NewAccurateNumMultiplyer creates the struct UnsignedNumericAccurateMultiplyer
// EntityName: Name of Entity
// BitSize: Bit-length of Multiplier inputs
func NewAccurateNumMultiplyer(EntityName string, BitSize uint) *UnsignedNumericAccurateMultiplyer {
	m := new(UnsignedNumericAccurateMultiplyer)
	m.EntityName = EntityName
	m.BitSize = BitSize
	m.VHDLFile, m.TestFile = FileNameGen(m.EntityName)

	return m
}

// GenerateVHDL creates the VHDL file in FolderPath
func (m *UnsignedNumericAccurateMultiplyer) GenerateVHDL(FolderPath string) {
	CreateFile(FolderPath, m.VHDLFile, "accuratebehav.vhd", m)
}

// GenerateTestData creates a plaintext testdata file containing both inputs and the output in binary seperated by \t
// GenerateTestData uses the function from New2DUnsignedAcc, since their behaviour is identical.
func (m *UnsignedNumericAccurateMultiplyer) GenerateTestData(FolderPath string) {
	Accurate2D := New2DUnsignedAcc(m.EntityName, m.BitSize)
	Accurate2D.TestFile = m.TestFile
	Accurate2D.GenerateTestData(FolderPath)
}

// ReturnData returns a struct containing metadata of the multiplier
func (m *UnsignedNumericAccurateMultiplyer) ReturnData() *EntityData {
	// EntityName string
	// BitSize    uint
	// VHDLFile   string
	// TestFile   string
	d := new(EntityData)
	d.EntityName = m.EntityName
	d.BitSize = m.BitSize
	d.OutputSize = 2 * m.BitSize
	d.VHDLFile = m.VHDLFile
	d.TestFile = m.TestFile
	return d
}

func (m *UnsignedNumericAccurateMultiplyer) String() string {
	return m.EntityName
}

// GenerateVHDLEntityArray creates an array of potentially multiple VHDLEntities, sorted by priority for synthesizing
// For example: Multiplier A uses a VHDL portmap for the smaller Multiplier B & C. B & C need to be synthesized first, hence A will be last in the array
func (m *UnsignedNumericAccurateMultiplyer) GenerateVHDLEntityArray() []VHDLEntity {
	var out []VHDLEntity
	out = append(out, m)
	return out
}

// ReturnVal returns the output of the multiplier
// This multiplier is accurate
func (m *UnsignedNumericAccurateMultiplyer) ReturnVal(a uint, b uint) uint {
	return a * b
}

// Overflow returns a boolean if any internal overflow has occured
// This Multiplier is accurate, hence no overflow
func (m *UnsignedNumericAccurateMultiplyer) Overflow() bool {
	return false
}

// MeanAbsoluteError returns the MeanAbsoluteError of the multiplier in float64
func (m *UnsignedNumericAccurateMultiplyer) MeanAbsoluteError() float64 {
	return 0
}

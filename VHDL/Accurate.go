package VHDL

//UnsignedNumericAccurateMultiplyer defines a multiplyer based upon 'accuratebehav.vhd', this is a BitSize A*B=prod multiplyer, using the IEEE Numeric lib
//The struct is implemented as a VHDLEntity Interface
type UnsignedNumericAccurateMultiplyer struct {
	EntityName string
	BitSize    uint
	VHDLFile   string
	TestFile   string
}

//Creates a VHDL file of an UnsignedNumericAccurateMultiplyer.
func NewAccurateNumMultiplyer(EntityName string, BitSize uint) *UnsignedNumericAccurateMultiplyer {
	m := new(UnsignedNumericAccurateMultiplyer)
	m.EntityName = EntityName
	m.BitSize = BitSize
	m.VHDLFile, m.TestFile = FileNameGen(m.EntityName)

	return m
}

func (m *UnsignedNumericAccurateMultiplyer) GenerateVHDL(FolderPath string) {
	CreateFile(FolderPath, m.VHDLFile, "accuratebehav.vhd", m)
}

//GenerateTestData uses the function from New2DUnsignedAcc, since their behaviour is identical.
func (m *UnsignedNumericAccurateMultiplyer) GenerateTestData(FolderPath string) {
	Accurate2D := New2DUnsignedAcc(m.EntityName, m.BitSize)
	Accurate2D.TestFile = m.TestFile
	Accurate2D.GenerateTestData(FolderPath)
}

func (m *UnsignedNumericAccurateMultiplyer) ReturnData() *EntityData {
	// EntityName string
	// BitSize    uint
	// VHDLFile   string
	// TestFile   string
	d := new(EntityData)
	d.EntityName = m.EntityName
	d.BitSize = m.BitSize
	d.VHDLFile = m.VHDLFile
	d.TestFile = m.TestFile
	return d
}

func (m *UnsignedNumericAccurateMultiplyer) String() string {
	return m.EntityName
}

func (m *UnsignedNumericAccurateMultiplyer) GenerateVHDLEntityArray() []VHDLEntity {
	var out []VHDLEntity
	out = append(out, m)
	return out
}

package VHDL

//UnsingedNumericAccurateMultiplyer defines a multiplyer based upon 'accuratebehav.vhd', this is a BitSize A*B=prod multiplyer, using the IEEE Numeric lib
type UnsignedNumericAccurateMultiplyer struct {
	EntityName string
	BitSize    uint
	VHDLFile   string
	TestFile   string
}

//Creates a VHDL file of an UnsignedNumericAccurateMultiplyer.
func AccurateMultToFile(EntityName string, BitSize uint, FolderPath string) *UnsignedNumericAccurateMultiplyer {
	m := new(UnsignedNumericAccurateMultiplyer)
	m.EntityName = EntityName
	m.BitSize = BitSize
	m.VHDLFile, m.TestFile = FileNameGen(m.EntityName)

	CreateVHDLFile(FolderPath, m.VHDLFile, "accuratebehav.vhd", m)

	return m
}

//GenerateTestData uses the function from New2DUnsignedAcc, since their behaviour is identical.
func (m *UnsignedNumericAccurateMultiplyer) GenerateTestData(FolderPath string) {
	Accurate2D := New2DUnsignedAcc(m.EntityName, m.BitSize)
	Accurate2D.TestFile = m.TestFile
	Accurate2D.GenerateTestData(FolderPath)
}

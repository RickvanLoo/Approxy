package VHDL

import "log"

type MAC struct {
	EntityName string
	BitSize    uint
	OutputSize uint
	Multiplier *EntityData
	VHDLFile   string
	TestFile   string
}

func NewMAC(Multiplier VHDLEntity) *MAC {
	mac := new(MAC)
	mac.Multiplier = Multiplier.ReturnData()
	mac.EntityName = "MAC_" + mac.Multiplier.EntityName
	mac.BitSize = mac.Multiplier.BitSize
	mac.OutputSize = mac.BitSize * 2
	mac.VHDLFile, mac.TestFile = FileNameGen(mac.EntityName)

	mac.TestFile = mac.Multiplier.TestFile // Fix because MAC does not have testing

	return mac
}

func (mac *MAC) ReturnData() *EntityData {
	d := new(EntityData)
	d.BitSize = mac.BitSize
	d.EntityName = mac.EntityName
	d.TestFile = mac.TestFile
	d.VHDLFile = mac.VHDLFile

	return d
}

func (mac *MAC) GenerateVHDL(FolderPath string) {
	CreateFile(FolderPath, mac.VHDLFile, "macbehav.vhd", mac)
}

func (mac *MAC) GenerateTestData(FolderPath string) {
	log.Println("ERROR: Generating Test Data for MAC not (yet) supported. Continuing...")
}

func (mac *MAC) String() string {
	return "MAC -> " + mac.Multiplier.EntityName
}

// ReturnData() *EntityData
// GenerateVHDL(string)
// GenerateTestData(string)
// String() string //MSB -> LSB

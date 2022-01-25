package VHDL

import (
	"log"
	"math"
)

type MAC struct {
	EntityName string
	BitSize    uint
	OutputSize uint
	Multiplier VHDLEntityMultiplier
	VHDLFile   string
	TestFile   string
}

//TODO : Make it the same as Rec4, where GenerateVHDL generates the preceding VHDL files as well.

func NewMAC(Multiplier VHDLEntityMultiplier, T int) *MAC {
	mac := new(MAC)
	mac.Multiplier = Multiplier
	MultiplierData := mac.Multiplier.ReturnData()
	mac.EntityName = "MAC_" + MultiplierData.EntityName
	mac.BitSize = MultiplierData.BitSize
	mac.VHDLFile, mac.TestFile = FileNameGen(mac.EntityName)

	mac.TestFile = MultiplierData.TestFile // Fix because MAC does not have testing

	maxinput := int(math.Exp2(float64(mac.BitSize)))
	maxval := 0

	for a := 0; a < maxinput; a++ {
		for b := 0; b < maxinput; b++ {
			val := Multiplier.ReturnVal(uint(a), uint(b))
			if val > uint(maxval) {
				maxval = int(val)
			}
		}
	}

	bitvalue := math.Log2(float64(maxval * T))
	mac.OutputSize = uint(math.Ceil(bitvalue))

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
	mac.Multiplier.GenerateVHDL(FolderPath)
	CreateFile(FolderPath, mac.VHDLFile, "macbehav.vhd", mac)
}

func (mac *MAC) GenerateTestData(FolderPath string) {
	log.Println("ERROR: Generating Test Data for MAC not (yet) supported. Continuing...")
}

func (mac *MAC) String() string {
	return "MAC -> " + mac.Multiplier.ReturnData().EntityName
}

// ReturnData() *EntityData
// GenerateVHDL(string)
// GenerateTestData(string)
// String() string //MSB -> LSB

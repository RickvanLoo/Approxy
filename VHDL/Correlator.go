package VHDL

import "log"

//MACArray[0] = ReYReS
//MACArray[1] = ImYImS
//MACArray[2] = ReYImS
//MACArray[3] = ImYReS
type Correlator struct {
	EntityName string
	BitSize    uint
	OutputSize uint
	MACArray   [4]*MAC //Size of 4
	VHDLFile   string
	TestFile   string
}

//TODO : Make it the same as Rec4, where GenerateVHDL generates the preceding VHDL files as well.

func NewCorrelator(EntityName string, MACArray [4]*MAC) *Correlator {
	corr := new(Correlator)
	corr.EntityName = EntityName
	corr.MACArray = MACArray
	MACData := corr.MACArray[0].ReturnData()
	corr.BitSize = MACData.BitSize
	corr.OutputSize = corr.BitSize * 3
	corr.VHDLFile, corr.TestFile = FileNameGen(corr.EntityName)
	return corr
}

func (corr *Correlator) ReturnData() *EntityData {
	d := new(EntityData)
	d.BitSize = corr.BitSize
	d.EntityName = corr.EntityName
	d.TestFile = corr.TestFile
	d.VHDLFile = corr.VHDLFile
	return d
}

func (corr *Correlator) GenerateVHDL(FolderPath string) {
	for _, mac := range corr.MACArray {
		mac.GenerateVHDL(FolderPath)
	}
	CreateFile(FolderPath, corr.VHDLFile, "corrbehav.vhd", corr)

}

func (corr *Correlator) GenerateTestData(FolderPath string) {
	log.Println("ERROR: Generating Test Data for Correlator not (yet) supported. Continuing...")
}

func (corr *Correlator) String() string {
	str := corr.EntityName + " -> [" + corr.MACArray[0].EntityName + ","
	str += corr.MACArray[1].EntityName + ","
	str += corr.MACArray[2].EntityName + ","
	str += corr.MACArray[3].EntityName + "]"
	return str
}

func (corr *Correlator) GenerateVHDLEntityArray() []VHDLEntity {

	var out []VHDLEntity

	out = append(out, corr)

	for _, mac := range corr.MACArray {
		out = append(out, mac)
	}

	return out
}

// ReturnData() *EntityData
// GenerateVHDL(string)
// GenerateTestData(string)
// String() string //MSB -> LSB

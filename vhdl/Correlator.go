package vhdl

import "log"

//TODO: Not-finished/Not-working

// Correlator is a non-finished struct implementing template "corrbehav.vhd" within the Approxy software
type Correlator struct {
	EntityName string
	BitSize    uint
	OutputSize uint
	MACArray   [4]VHDLEntity //Size of 4
	VHDLFile   string
	TestFile   string
}

//TODO : Make it the same as Rec4, where GenerateVHDL generates the preceding VHDL files as well.

// NewCorrelator returns a new correlator on basis of 4 Multiply-Accumulators as defined
// MACArray[0] = ReYReS
// MACArray[1] = ImYImS
// MACArray[2] = ReYImS
// MACArray[3] = ImYReS
func NewCorrelator(EntityName string, MACArray [4]VHDLEntity) *Correlator {
	corr := new(Correlator)
	corr.EntityName = EntityName
	corr.MACArray = MACArray
	MACData := corr.MACArray[0].ReturnData()
	corr.BitSize = MACData.BitSize
	corr.OutputSize = corr.MACArray[0].ReturnData().OutputSize //Different Outputsize for MACArray not yet supported
	corr.VHDLFile, corr.TestFile = FileNameGen(corr.EntityName)
	return corr
}

// ReturnData returns a struct containing metadata of the correlator
func (corr *Correlator) ReturnData() *EntityData {
	d := new(EntityData)
	d.BitSize = corr.BitSize
	d.EntityName = corr.EntityName
	d.TestFile = corr.TestFile
	d.VHDLFile = corr.VHDLFile
	return d
}

// GenerateVHDL creates the VHDL file in FolderPath
func (corr *Correlator) GenerateVHDL(FolderPath string) {
	for _, mac := range corr.MACArray {
		mac.GenerateVHDL(FolderPath)
	}
	CreateFile(FolderPath, corr.VHDLFile, "corrbehav.vhd", corr)

}

// GenerateTestData creates a plaintext testdata file containing both inputs and the output in binary seperated by \t
// Not Implemented for CORRELATOR!
func (corr *Correlator) GenerateTestData(FolderPath string) {
	log.Println("ERROR: Generating Test Data for Correlator not (yet) supported. Continuing...")
}

func (corr *Correlator) String() string {
	str := corr.EntityName + " -> [" + corr.MACArray[0].ReturnData().EntityName + ";"
	str += corr.MACArray[1].ReturnData().EntityName + ";"
	str += corr.MACArray[2].ReturnData().EntityName + ";"
	str += corr.MACArray[3].ReturnData().EntityName + "]"
	return str
}

// GenerateVHDLEntityArray creates an array of potentially multiple VHDLEntities, sorted by priority for synthesizing
// For example: Multiplier A uses a VHDL portmap for the smaller Multiplier B & C. B & C need to be synthesized first, hence A will be last in the array
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

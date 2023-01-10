package vhdl

import (
	"log"
	"strconv"
)

// Scaler creates a linear scaling VHDL topfile for multipliers on basis of template  "scaler.vhd"
// Note that synthesis has to be done OOC (Out-of-Context) or all multipliers besides one get optimised away
type Scaler struct {
	Entity     VHDLEntityMultiplier
	LUTName    string
	EntityName string
	BitSize    uint
	ScaleN     uint
	VHDLFile   string
	TestFile   string
	MAC        bool
	OutputSize uint
}

// New2DScaler creates a new Scaler struct
func New2DScaler(Entity VHDLEntityMultiplier, N uint) *Scaler {
	scl := new(Scaler)
	scl.Entity = Entity
	scl.LUTName = scl.Entity.ReturnData().EntityName
	scl.BitSize = scl.Entity.ReturnData().BitSize
	scl.EntityName = scl.LUTName + "_scaler"
	scl.ScaleN = N
	scl.VHDLFile = scl.LUTName + "_scaler.vhd"
	scl.MAC = false
	scl.OutputSize = scl.BitSize * 2

	scl.VHDLFile, scl.TestFile = FileNameGen(scl.EntityName)
	scl.TestFile = scl.Entity.ReturnData().TestFile //hacky

	return scl
}

// SetMAC changes Scaling template for Multiply-Accumulators
func (scl *Scaler) SetMAC(b bool, OutputSize uint) {
	scl.MAC = b
	scl.OutputSize = OutputSize
}

// GenerateTestData creates a plaintext testdata file containing both inputs and the output in binary seperated by \t
// Outputs simply the testdata of Multiplier that is being scaled
func (scl *Scaler) GenerateTestData(FolderPath string) {
	//log.Println("ERROR: Generating Test Data for Scaler not supported. Continuing...")
	scl.Entity.GenerateTestData(FolderPath)
}

// GenerateVHDL creates the VHDL files in FolderPath
func (scl *Scaler) GenerateVHDL(FolderPath string) {
	CreateFile(FolderPath, scl.VHDLFile, "scaler.vhd", scl)
	scl.Entity.GenerateVHDL(FolderPath)
}

// ReturnData returns a struct containing metadata of the multiplier
func (scl *Scaler) ReturnData() *EntityData {
	// EntityName string
	// BitSize    uint
	// VHDLFile   string
	// TestFile   string
	d := new(EntityData)
	d.EntityName = scl.EntityName
	d.BitSize = scl.BitSize
	d.OutputSize = scl.Entity.ReturnData().OutputSize
	d.VHDLFile = scl.VHDLFile
	d.TestFile = scl.TestFile
	return d
}

func (scl *Scaler) String() string {
	str := scl.EntityName + " N=" + strconv.Itoa(int(scl.ScaleN)) + " -> " + scl.Entity.String()
	return str
}

// GenerateVHDLEntityArray creates an array of potentially multiple VHDLEntities, sorted by priority for synthesizing
// For example: Multiplier A uses a VHDL portmap for the smaller Multiplier B & C. B & C need to be synthesized first, hence A will be last in the array
func (scl *Scaler) GenerateVHDLEntityArray() []VHDLEntity {
	//Not Correct, but Scaler is not XSIM'able anyway
	//log.Println("ERROR: Generating VHDLEntityArray for Scaler not supported. Continuing...")
	var out []VHDLEntity
	out = append(out, scl)
	out = append(out, scl.Entity.GenerateVHDLEntityArray()...)
	return out
}

// MeanAbsoluteError is not supported
func (scl *Scaler) MeanAbsoluteError() float64 {
	log.Println("ERROR: MeanAbsoluteError for Scaler not supported. Try non-scaled Entity! Continuing...")
	return 0
}

// Overflow is not supported
func (scl *Scaler) Overflow() bool {
	log.Println("ERROR: Overflow() for Scaler not supported. Try non-scaled Entity! Continuing...")
	return false
}

// ReturnVal simply returns value of scaled multiplier
func (scl *Scaler) ReturnVal(a uint, b uint) uint {
	//log.Println("ERROR: ReturnVal() for Scaler not supported. Try non-scaled Entity! Continuing...")
	return scl.Entity.ReturnVal(a, b)
}

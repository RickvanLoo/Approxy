package VHDL

import (
	"log"
	"strconv"
)

type Scaler struct {
	Entity     VHDLEntity
	LUTName    string
	EntityName string
	BitSize    uint
	ScaleN     uint
	VHDLFile   string
}

func New2DScaler(Entity VHDLEntity, N uint) *Scaler {
	scl := new(Scaler)
	scl.Entity = Entity
	scl.LUTName = scl.Entity.ReturnData().EntityName
	scl.BitSize = scl.Entity.ReturnData().BitSize
	scl.EntityName = scl.LUTName + "_scaler"
	scl.ScaleN = N
	scl.VHDLFile = scl.LUTName + "_scaler.vhd"

	return scl
}

func (scl *Scaler) GenerateTestData(FolderPath string) {
	log.Println("ERROR: Generating Test Data for Scaler not supported. Continuing...")
}

func (scl *Scaler) GenerateVHDL(FolderPath string) {
	CreateFile(FolderPath, scl.VHDLFile, "scaler.vhd", scl)
}

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
	d.TestFile = "" //Not Supported!
	return d
}

func (scl *Scaler) String() string {
	str := scl.EntityName + " N=" + strconv.Itoa(int(scl.ScaleN)) + " -> " + scl.Entity.String()
	return str
}

func (scl *Scaler) GenerateVHDLEntityArray() []VHDLEntity {
	//Not Correct, but Scaler is not XSIM'able anyway
	log.Println("ERROR: Generating VHDLEntityArray for Scaler not supported. Continuing...")
	var out []VHDLEntity
	out = append(out, scl)
	return out
}

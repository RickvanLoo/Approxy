package VHDL

type Scaler struct {
	LUT        *LUT2D
	LUTName    string
	EntityName string
	BitSize    uint
	ScaleN     uint
}

func CreateScalerVHDL(m *LUT2D, N uint, FolderPath string) *Scaler {
	scl := new(Scaler)
	scl.LUT = m
	scl.LUTName = scl.LUT.EntityName
	scl.BitSize = scl.LUT.BitSize
	scl.EntityName = scl.LUTName + "_scaler"
	scl.ScaleN = N

	name := scl.LUTName + "_scaler.vhd"

	CreateVHDLFile(FolderPath, name, "scaler.vhd", scl)

	return scl
}

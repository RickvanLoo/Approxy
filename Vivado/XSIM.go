package Vivado

import (
	"badmath/VHDL"
	"log"
	"os"
	"os/exec"
	"strings"
)

type XSIM struct {
	SimFile       string
	TopFile       string
	TestFile      string
	TopEntityName string
	FolderPath    string
	VHDLEntities  []VHDL.VHDLEntity
	BitSize       uint
}

//Creates an XSIM on bases of an array of VHDLEntities
//TopEntity at index[0]!
func CreateXSIM(FolderPath string, SimName string, VHDLEntities []VHDL.VHDLEntity) *XSIM {
	XSIM := new(XSIM)
	XSIM.SimFile = SimName + ".vhd"
	XSIM.TopFile = VHDLEntities[0].ReturnData().VHDLFile
	XSIM.TestFile = VHDLEntities[0].ReturnData().TestFile
	XSIM.TopEntityName = VHDLEntities[0].ReturnData().EntityName
	XSIM.BitSize = VHDLEntities[0].ReturnData().BitSize
	XSIM.FolderPath = FolderPath
	XSIM.VHDLEntities = VHDLEntities
	//XSIM.removeDuplicate()

	return XSIM
}

func (x *XSIM) Exec() {
	VHDL.CreateFile(x.FolderPath, x.SimFile, "xsim.vhd", x)

	//This is ugly as hell, but it works, and is readable

	for i := len(x.VHDLEntities) - 1; i >= 0; i-- {
		loadBehav := exec.Command("xvhdl", x.VHDLEntities[i].ReturnData().VHDLFile)
		loadBehav.Dir = x.FolderPath
		loadBehav.Stdout = os.Stdout
		loadBehav.Stderr = os.Stderr
		err := loadBehav.Run()
		if err != nil {
			log.Println(err)
		}
	}

	loadSim := exec.Command("xvhdl", x.SimFile)
	loadSim.Dir = x.FolderPath
	loadSim.Stdout = os.Stdout
	loadSim.Stderr = os.Stderr

	xelab := exec.Command("xelab", "-debug", "typical", "sim", "-s", x.TopEntityName+"top_sim")
	xelab.Dir = x.FolderPath
	xelab.Stdout = os.Stdout
	xelab.Stderr = os.Stderr

	xsim := exec.Command("xsim", x.TopEntityName+"top_sim", "--log", x.TopEntityName+"_xsimlog")
	xsim.Dir = x.FolderPath
	xsim.Stderr = os.Stderr
	xsim.Stdout = os.Stdout
	xsim.Stdin = strings.NewReader("run all\n")

	err := loadSim.Run()
	if err != nil {
		log.Println(err)
	}

	err = xelab.Run()
	if err != nil {
		log.Println(err)
	}
	err = xsim.Run()
	if err != nil {
		log.Println(err)
	}
}

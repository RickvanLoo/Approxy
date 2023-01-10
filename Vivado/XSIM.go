package Vivado

import (
	"log"
	"os"
	"os/exec"
	"strings"

	VHDL "github.com/RickvanLoo/Approxy/vhdl"
)

type XSIM struct {
	SimFile       string
	TopFile       string
	TestFile      string
	TopEntityName string
	FolderPath    string
	TemplateFile  string
	VHDLEntities  []VHDL.VHDLEntity
	BitSize       uint
	OutputSize    uint
	ScaleN        uint
	PostSim       bool //For special functionality made for PostSimming
}

// Creates an XSIM on bases of an array of VHDLEntities
// TopEntity at index[0]!
func CreateXSIM(FolderPath string, SimName string, VHDLEntities []VHDL.VHDLEntity) *XSIM {
	XSIM := new(XSIM)
	XSIM.SimFile = SimName + ".vhd"
	XSIM.TopFile = VHDLEntities[0].ReturnData().VHDLFile
	XSIM.TestFile = VHDLEntities[0].ReturnData().TestFile
	XSIM.TopEntityName = VHDLEntities[0].ReturnData().EntityName
	XSIM.BitSize = VHDLEntities[0].ReturnData().BitSize
	XSIM.FolderPath = FolderPath
	XSIM.TemplateFile = "xsim_mult.vhd" //Default option
	XSIM.VHDLEntities = VHDLEntities
	XSIM.OutputSize = 2 * XSIM.BitSize
	XSIM.ScaleN = 1
	XSIM.PostSim = false

	return XSIM
}

// Set TB to Multiplier sweep for single multipliers
func (x *XSIM) SetTemplateMultiplier() {
	x.TemplateFile = "xsim_mult.vhd"
}

// Set TB to Sequential Multiplier testing for MAC
func (x *XSIM) SetTemplateSequential(OutputSize uint) {
	x.TemplateFile = "xsim_seq.vhd"
	x.OutputSize = OutputSize
}

// Set TB to Multiplier sweep for Scaler
func (x *XSIM) SetTemplateScaler(N uint) {
	x.TemplateFile = "xsim_mult_scaler.vhd"
	x.ScaleN = N
}

func (x *XSIM) SetTemplateSequentialScaler(N uint, OutputSize uint) {
	x.TemplateFile = "xsim_seq_scaler.vhd"
	x.OutputSize = OutputSize
	x.ScaleN = N
}

// Set TB to not functionally verify model, but retreive data instead
func (x *XSIM) SetTemplateReverse() {
	x.TemplateFile = "xsim_reverse.vhd"
}

// Recreate TB File, neccesary for PostPlacement sim after behavioural analysis
// PostSim here is a bool that switches defined pre/post blocks within TB templates
func (x *XSIM) CreateFile(PostSim bool) {
	x.PostSim = PostSim
	VHDL.CreateFile(x.FolderPath, x.SimFile, x.TemplateFile, x)
}

// Exec() creates the TB file and runs a behavioural analysis
func (x *XSIM) Exec() {
	x.PostSim = false
	VHDL.CreateFile(x.FolderPath, x.SimFile, x.TemplateFile, x)

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

// FuncSim() does not(!) create a TB, but runs a PostPlacement analysis on the current available TB and dumps SAIF
func (x *XSIM) Funcsim() {
	loadTop := exec.Command("xvhdl", x.TopEntityName+"_funcsim.vhd")
	loadTop.Dir = x.FolderPath
	loadTop.Stdout = os.Stdout
	loadTop.Stderr = os.Stderr

	loadSim := exec.Command("xvhdl", x.SimFile)
	loadSim.Dir = x.FolderPath
	loadSim.Stdout = os.Stdout
	loadSim.Stderr = os.Stderr

	xelab := exec.Command("xelab", "-debug", "typical", "-L", "secureip", "-L", "unisims_ver", "sim", "-s", "top_funcsim")
	xelab.Dir = x.FolderPath
	xelab.Stdout = os.Stdout
	xelab.Stderr = os.Stderr

	xsim := exec.Command("xsim", "top_funcsim", "--log", x.TopEntityName+"_funcsimlog")
	xsim.Dir = x.FolderPath
	xsim.Stderr = os.Stderr
	xsim.Stdout = os.Stdout
	xsimCommand := "open_saif " + x.TopEntityName + "_dump.saif\n"
	xsimCommand = xsimCommand + "log_saif [ get_objects -r ]\n" //default name of simulation is sim, default name of module is testmod
	xsimCommand = xsimCommand + "run all\n"
	xsimCommand = xsimCommand + "close_saif\n"
	xsim.Stdin = strings.NewReader(xsimCommand)

	err := loadTop.Run()
	if err != nil {
		log.Println(err)
	}

	err = loadSim.Run()
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

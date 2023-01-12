package vivado

import (
	"log"
	"os"
	"os/exec"

	VHDL "github.com/RickvanLoo/Approxy/vhdl"
)

// VivadoTCL is a struct that contains the name and location of a generated Vivado TCL file
// See templates/vivado.tcl
type VivadoTCL struct {
	FolderPath string
	TopName    string
	FileName   string
	Settings   *VivadoTCLSettings
}

// VivadoTCLSettings contains all settings used in the TCL template for enabling/disabling functionality
type VivadoTCLSettings struct {
	PartName        string
	OOC             bool // Enables Out-of-Context synthesis
	NODSP           bool // When TRUE disables the use of DSP slices for the synthesis of the design
	WriteCheckpoint bool // Enables exporting Checkpoints after synthesis and place+route
	Placement       bool // Enables Placement
	Route           bool // Enables Routing
	Funcsim         bool // Note: Set to false to disable funcsim VHDL file for post-processing.
	Utilization     bool // Enables Exporting of Utilization Report
	Hierarchical    bool // Enables Hierachical Utilization Report //Note: Setting this to false, does break utilization report parsing
	Clk             bool // Set this flag to TRUE if the design has a clock, else potentially breaking
	Timing          bool
}

// CreateVivadoTCL generates a Vivado TCL file named FileName within FolderPath, based on the top-level Entity and Settings
func CreateVivadoTCL(FolderPath string, FileName string, Entity VHDL.VHDLEntity, Settings *VivadoTCLSettings) *VivadoTCL {
	TCL := new(VivadoTCL)
	TCL.FolderPath = FolderPath
	TCL.TopName = Entity.ReturnData().EntityName
	TCL.FileName = FileName

	TCL.Settings = Settings

	VHDL.CreateFile(FolderPath, TCL.FileName, "vivado.tcl", TCL)

	return TCL
}

// Exec starts Vivado in Batch mode to execute the TCL file
func (tcl *VivadoTCL) Exec() {
	cmd := exec.Command("vivado", "-mode", "batch", "-source", tcl.FileName)
	cmd.Dir = tcl.FolderPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
}

// PowerPostPlacementGeneration creates and executes a TCL on basis of template reportpower.tcl
// Requires exported SAIF and post-PR checkpoint file
func (tcl *VivadoTCL) PowerPostPlacementGeneration() {
	VHDL.CreateFile(tcl.FolderPath, "power_"+tcl.FileName, "reportpower.tcl", tcl)
	cmd := exec.Command("vivado", "-mode", "batch", "-source", "power_"+tcl.FileName)
	cmd.Dir = tcl.FolderPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
}

// ClearLogs is used to clear all Vivado logs within the main folder ending on *.jou and *.log
func ClearLogs() {
	cmd1 := exec.Command("/bin/sh", "-c", "rm *.jou")
	cmd2 := exec.Command("/bin/sh", "-c", "rm *.log")
	cmd1.Stdout = os.Stdout
	cmd1.Stderr = os.Stderr
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr
	err := cmd1.Run()
	if err != nil {
		log.Println(err)
	}
	err = cmd2.Run()
	if err != nil {
		log.Println(err)
	}
}

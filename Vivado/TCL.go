package Vivado

import (
	"approxy/VHDL"
	"log"
	"os"
	"os/exec"
)

type VivadoTCL struct {
	FolderPath string
	TopName    string
	FileName   string
	Settings   *VivadoTCLSettings
}

type VivadoTCLSettings struct {
	PartName        string
	OOC             bool
	NO_DSP          bool
	WriteCheckpoint bool
	Placement       bool
	Route           bool
	Funcsim         bool //Note: Set to false to disable funcsim VHDL file for post-processing.
	Utilization     bool
	Hierarchical    bool //Note: Setting this to false, does break utilization report parsing
	Clk             bool //WARNING: Set only to TRUE if the design has a clock.
	Timing          bool
}

func CreateVivadoTCL(FolderPath string, FileName string, Entity VHDL.VHDLEntity, Settings *VivadoTCLSettings) *VivadoTCL {
	TCL := new(VivadoTCL)
	TCL.FolderPath = FolderPath
	TCL.TopName = Entity.ReturnData().EntityName
	TCL.FileName = FileName

	TCL.Settings = Settings

	VHDL.CreateFile(FolderPath, TCL.FileName, "vivado.tcl", TCL)

	return TCL
}

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

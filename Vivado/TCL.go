package Vivado

import (
	"log"
	"os"
	"os/exec"
	"text/template"
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
	Utilization     bool
}

func CreateVivadoTCL(FolderPath string, FileName string, TopName string, Settings *VivadoTCLSettings) *VivadoTCL {
	TCL := new(VivadoTCL)
	TCL.FolderPath = FolderPath
	TCL.TopName = TopName
	TCL.FileName = FileName

	TCL.Settings = Settings

	templatepath := "template/vivado.tcl"
	templatename := "vivado.tcl"

	t, err := template.New(templatename).ParseFiles(templatepath)
	if err != nil {
		log.Print(err)
		return TCL
	}

	f, err := os.Create(FolderPath + "/" + FileName)

	if err != nil {
		log.Println("create file: ", err)
		return TCL
	}

	err = t.ExecuteTemplate(f, templatename, TCL)
	if err != nil {
		log.Print("execute: ", err)
		return TCL
	}

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

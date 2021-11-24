package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/RickvanLoo/badmath/mult"
)

type VivadoTCL struct {
	PartName        string
	Folder          string
	Top             string
	OOC             bool
	WriteCheckpoint bool
	Placement       bool
	Utilization     bool
}

type XSIM struct {
	SimFile       string
	TopFile       string
	TestData      string
	Behav         *mult.UnsignedAccurateMultiplyer
	TopEntityName string
	Bitsize       uint
}

func CreateVivadoTCL(folder string, name string, top string) {
	TCL := new(VivadoTCL)
	TCL.PartName = "Xc7z030fbg676-3"
	TCL.OOC = true
	TCL.Folder = folder
	TCL.Top = top
	TCL.WriteCheckpoint = true
	TCL.Utilization = true
	TCL.Placement = true

	templatepath := "template/vivado.tcl"
	templatename := "vivado.tcl"

	t, err := template.New(templatename).ParseFiles(templatepath)
	if err != nil {
		log.Print(err)
		return
	}

	f, err := os.Create(folder + "/" + name)

	if err != nil {
		log.Println("create file: ", err)
		return
	}

	err = t.ExecuteTemplate(f, templatename, TCL)
	if err != nil {
		log.Print("execute: ", err)
		return
	}
}

func ExecuteVivadoTCL(folder string, name string) {
	cmd := exec.Command("vivado", "-mode", "batch", "-source", name)
	cmd.Dir = folder
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

func CreateXSIM(folder string, top string, datafile string, newsim string, behav *mult.UnsignedAccurateMultiplyer) *XSIM {
	XSIM := new(XSIM)
	XSIM.SimFile = newsim
	XSIM.TopFile = top
	XSIM.TestData = datafile
	XSIM.Behav = behav
	XSIM.TopEntityName = XSIM.Behav.Name
	XSIM.Bitsize = XSIM.Behav.Bitsize

	templatepath := "template/xsim.vhd"
	templatename := "xsim.vhd"

	t, err := template.New(templatename).ParseFiles(templatepath)
	if err != nil {
		log.Print(err)
		return XSIM
	}

	f, err := os.Create(folder + "/" + XSIM.SimFile)

	if err != nil {
		log.Println("create file: ", err)
		return XSIM
	}

	err = t.ExecuteTemplate(f, templatename, XSIM)
	if err != nil {
		log.Print("execute: ", err)
		return XSIM
	}

	return XSIM
}

func (x *XSIM) Exec(folder string) {
	//This is ugly as hell, but it works, and is readable
	loadBehav := exec.Command("xvhdl", x.TopFile)
	loadBehav.Dir = folder
	loadBehav.Stdout = os.Stdout
	loadBehav.Stderr = os.Stderr

	loadSim := exec.Command("xvhdl", x.SimFile)
	loadSim.Dir = folder
	loadSim.Stdout = os.Stdout
	loadSim.Stderr = os.Stderr

	xelab := exec.Command("xelab", "-debug", "typical", "sim", "-s", x.Behav.Name+"top_sim")
	xelab.Dir = folder
	xelab.Stdout = os.Stdout
	xelab.Stderr = os.Stderr

	xsim := exec.Command("xsim", x.Behav.Name+"top_sim", "--log", x.Behav.Name+"_xsimlog")
	xsim.Dir = folder
	xsim.Stderr = os.Stderr
	xsim.Stdout = os.Stdout
	xsim.Stdin = strings.NewReader("run all\n")

	err := loadBehav.Run()
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

package Vivado

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

type XSIM struct {
	SimFile       string
	TopFile       string
	TestData      string
	TopEntityName string
	Bitsize       uint
}

func CreateXSIM(FolderPath string, TopFile string, DataFile string, SimName string, TopEntityName string, BitSize uint) *XSIM {
	XSIM := new(XSIM)
	XSIM.SimFile = SimName
	XSIM.TopFile = TopFile
	XSIM.TestData = DataFile
	XSIM.TopEntityName = TopEntityName
	XSIM.Bitsize = BitSize

	templatepath := "template/xsim.vhd"
	templatename := "xsim.vhd"

	t, err := template.New(templatename).ParseFiles(templatepath)
	if err != nil {
		log.Print(err)
		return XSIM
	}

	f, err := os.Create(FolderPath + "/" + XSIM.SimFile)

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

	xelab := exec.Command("xelab", "-debug", "typical", "sim", "-s", x.TopEntityName+"top_sim")
	xelab.Dir = folder
	xelab.Stdout = os.Stdout
	xelab.Stderr = os.Stderr

	xsim := exec.Command("xsim", x.TopEntityName+"top_sim", "--log", x.TopEntityName+"_xsimlog")
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

package main

import (
	"log"
	"os"
	"os/exec"
	"text/template"
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
	path := folder + "/" + name

	cmd := exec.Command("vivado", "-mode", "batch", "-source", path)
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

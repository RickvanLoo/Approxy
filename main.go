package main

import (
	"fmt"
	"log"
	"os"

	"github.com/RickvanLoo/badmath/VHDL"
	Viv "github.com/RickvanLoo/badmath/Vivado"
)

var OutputPath string
var VivadoSettings *Viv.VivadoTCLSettings

func ClearPath(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		log.Println(err)
	}
}

func CreatePath(path string) {

	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	//ClearLogs()
	fmt.Println("badmath...")

	OutputPath = "output"
	ClearPath(OutputPath)
	CreatePath(OutputPath)

	VivadoSettings = new(Viv.VivadoTCLSettings)
	VivadoSettings.NO_DSP = true
	VivadoSettings.OOC = true
	VivadoSettings.PartName = "Xc7z030fbg676-3"
	VivadoSettings.Placement = true
	VivadoSettings.Utilization = true
	VivadoSettings.WriteCheckpoint = true

	acc16 := VHDL.UAM_To_VHDL("Acc16", 16, OutputPath, "Acc16.vhd")
	tcl := Viv.CreateVivadoTCL(OutputPath, "main1.tcl", acc16.EntityName, VivadoSettings)
	tcl.Exec()
}

func ScaleM1() {
	M1 := VHDL.M1()
	M1.LUT2D.Print()
	M1.LUT2D.VHDLtoFile(OutputPath, "M1.vhd")
	Scaler := VHDL.CreateScaler(M1.LUT2D, 100, OutputPath)
	tcl := Viv.CreateVivadoTCL(OutputPath, "main.tcl", Scaler.EntityName, VivadoSettings)
	tcl.Exec()
}

func M1M2M3M4() {
	M1 := VHDL.M1()
	M1.LUT2D.Print()
	M2 := VHDL.M2()
	M2.LUT2D.Print()
	M3 := VHDL.M3()
	M3.LUT2D.Print()
	M4 := VHDL.M4()
	M4.LUT2D.Print()

	M1.LUT2D.VHDLtoFile(OutputPath, "m1.vhd")
	M2.LUT2D.VHDLtoFile(OutputPath, "m2.vhd")
	M3.LUT2D.VHDLtoFile(OutputPath, "m3.vhd")
	M4.LUT2D.VHDLtoFile(OutputPath, "m4.vhd")

	M1.LUT2D.GenerateOutput(OutputPath, "testb1.txt")
	M2.LUT2D.GenerateOutput(OutputPath, "testb2.txt")
	M3.LUT2D.GenerateOutput(OutputPath, "testb3.txt")
	M4.LUT2D.GenerateOutput(OutputPath, "testb4.txt")

	XSIM1 := Viv.CreateXSIM(OutputPath, "m1.vhd", "testb1.txt", "topsim1.vhd", M1.LUT2D.EntityName, M1.LUT2D.BitSize)
	XSIM2 := Viv.CreateXSIM(OutputPath, "m2.vhd", "testb2.txt", "topsim2.vhd", M2.LUT2D.EntityName, M2.LUT2D.BitSize)
	XSIM3 := Viv.CreateXSIM(OutputPath, "m3.vhd", "testb3.txt", "topsim3.vhd", M3.LUT2D.EntityName, M3.LUT2D.BitSize)
	XSIM4 := Viv.CreateXSIM(OutputPath, "m4.vhd", "testb4.txt", "topsim4.vhd", M4.LUT2D.EntityName, M4.LUT2D.BitSize)

	XSIM1.Exec(OutputPath)
	XSIM2.Exec(OutputPath)
	XSIM3.Exec(OutputPath)
	XSIM4.Exec(OutputPath)

	tcl1 := Viv.CreateVivadoTCL(OutputPath, "main1.tcl", M1.LUT2D.EntityName, VivadoSettings)
	tcl2 := Viv.CreateVivadoTCL(OutputPath, "main2.tcl", M2.LUT2D.EntityName, VivadoSettings)
	tcl3 := Viv.CreateVivadoTCL(OutputPath, "main3.tcl", M3.LUT2D.EntityName, VivadoSettings)
	tcl4 := Viv.CreateVivadoTCL(OutputPath, "main4.tcl", M4.LUT2D.EntityName, VivadoSettings)

	tcl1.Exec()
	tcl2.Exec()
	tcl3.Exec()
	tcl4.Exec()

}

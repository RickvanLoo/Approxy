package main

import (
	"fmt"
	"log"
	"os"

	"github.com/RickvanLoo/badmath/VHDL"
	Viv "github.com/RickvanLoo/badmath/Vivado"
)

var OutputPath string

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

	M1M2M3M4()

}

func ScaleM1() {
	M1 := VHDL.M1()
	M1.LUT2D.Print()
	M1.LUT2D.VHDLtoFile(OutputPath, "M1.vhd")
	Scaler := VHDL.CreateScaler(M1.LUT2D, 100, OutputPath)
	Viv.CreateVivadoTCL(OutputPath, "main.tcl", Scaler.EntityName)
	Viv.ExecuteVivadoTCL(OutputPath, "main.tcl")
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

	Viv.CreateVivadoTCL(OutputPath, "main1.tcl", M1.LUT2D.EntityName)
	Viv.CreateVivadoTCL(OutputPath, "main2.tcl", M2.LUT2D.EntityName)
	Viv.CreateVivadoTCL(OutputPath, "main3.tcl", M3.LUT2D.EntityName)
	Viv.CreateVivadoTCL(OutputPath, "main4.tcl", M4.LUT2D.EntityName)

	Viv.ExecuteVivadoTCL(OutputPath, "main1.tcl")
	Viv.ExecuteVivadoTCL(OutputPath, "main2.tcl")
	Viv.ExecuteVivadoTCL(OutputPath, "main3.tcl")
	Viv.ExecuteVivadoTCL(OutputPath, "main4.tcl")

}

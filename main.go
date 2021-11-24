package main

import (
	"fmt"
	"log"
	"os"

	"github.com/RickvanLoo/badmath/mult"
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

	M1 := mult.M1()
	M1.Mult.Print()
	M1.Mult.VHDLtoFile(OutputPath, "M1.vhd")
	Scaler := mult.CreateScaler(M1.Mult, 31440, OutputPath)
	CreateVivadoTCL(OutputPath, "main.tcl", Scaler.EntityName)
	ExecuteVivadoTCL(OutputPath, "main.tcl")
}

func M1M2M3M4() {
	M1 := mult.M1()
	M1.Mult.Print()
	M2 := mult.M2()
	M2.Mult.Print()
	M3 := mult.M3()
	M3.Mult.Print()
	M4 := mult.M4()
	M4.Mult.Print()

	M1.Mult.VHDLtoFile(OutputPath, "m1.vhd")
	M2.Mult.VHDLtoFile(OutputPath, "m2.vhd")
	M3.Mult.VHDLtoFile(OutputPath, "m3.vhd")
	M4.Mult.VHDLtoFile(OutputPath, "m4.vhd")

	M1.Mult.GenerateOutput(OutputPath, "testb1.txt")
	M2.Mult.GenerateOutput(OutputPath, "testb2.txt")
	M3.Mult.GenerateOutput(OutputPath, "testb3.txt")
	M4.Mult.GenerateOutput(OutputPath, "testb4.txt")

	XSIM1 := CreateXSIM(OutputPath, "m1.vhd", "testb1.txt", "topsim1.vhd", M1.Mult)
	XSIM2 := CreateXSIM(OutputPath, "m2.vhd", "testb2.txt", "topsim2.vhd", M2.Mult)
	XSIM3 := CreateXSIM(OutputPath, "m3.vhd", "testb3.txt", "topsim3.vhd", M3.Mult)
	XSIM4 := CreateXSIM(OutputPath, "m4.vhd", "testb4.txt", "topsim4.vhd", M4.Mult)

	XSIM1.Exec(OutputPath)
	XSIM2.Exec(OutputPath)
	XSIM3.Exec(OutputPath)
	XSIM4.Exec(OutputPath)

	CreateVivadoTCL(OutputPath, "main1.tcl", M1.Mult.Name)
	CreateVivadoTCL(OutputPath, "main2.tcl", M2.Mult.Name)
	CreateVivadoTCL(OutputPath, "main3.tcl", M3.Mult.Name)
	CreateVivadoTCL(OutputPath, "main4.tcl", M4.Mult.Name)

	ExecuteVivadoTCL(OutputPath, "main1.tcl")
	ExecuteVivadoTCL(OutputPath, "main2.tcl")
	ExecuteVivadoTCL(OutputPath, "main3.tcl")
}

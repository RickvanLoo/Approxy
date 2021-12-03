package main

import (
	"fmt"
	"log"
	"os"

	VHDL "badmath/VHDL"

	Viv "badmath/Vivado"
)

var OutputPath string
var VivadoSettings *Viv.VivadoTCLSettings
var Results []*Result

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

	//o1 = AH*BH
	//o2 = AH*BL
	//o3 = AL*BH
	//o4 = AL*BL

	// o1 := VHDL.New2DUnsignedAcc("Acc", 2)
	// o2 := VHDL.New2DUnsignedAcc("Acc", 2)
	// o3 := VHDL.New2DUnsignedAcc("Acc", 2)
	// o4 := VHDL.New2DUnsignedAcc("Acc", 2)

	// o1 := VHDL.M4().LUT2D
	// o2 := VHDL.M4().LUT2D
	// o3 := VHDL.M4().LUT2D
	// o4 := VHDL.M4().LUT2D

	o1 := VHDL.New2DUnsignedAcc("Acc", 2)
	o2 := VHDL.M4().LUT2D
	o3 := VHDL.M1().LUT2D
	o4 := VHDL.New2DUnsignedAcc("Acc", 2)

	RecLutArray := [4]*VHDL.LUT2D{o1, o2, o3, o4}
	rec4 := VHDL.NewRecursive4("rec4", RecLutArray)

	rec4.GenerateTestData(OutputPath)
	rec4.GenerateVHDL(OutputPath)
	xsim := Viv.CreateXSIM(OutputPath, "SimRec4", rec4.GenerateVHDLEntityArray())
	xsim.Exec()
	err := Viv.ParseXSIMReport(OutputPath, rec4)
	if err != nil {
		log.Fatalln(err)
	}
	rec4scaler := VHDL.New2DScaler(rec4, 100)
	rec4scaler.GenerateVHDL(OutputPath)
	tcl := Viv.CreateVivadoTCL(OutputPath, "main1.tcl", rec4scaler, VivadoSettings)
	tcl.Exec()
	util := Viv.ParseUtilizationReport(OutputPath, rec4scaler)

	Result := NewResult(rec4scaler, util, rec4.Overflow(), rec4.MeanAbsoluteError())
	Result.PrettyPrint()
	Results = append(Results, Result)

}

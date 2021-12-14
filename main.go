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
var M1 *VHDL.LUT2D
var M2 *VHDL.LUT2D
var M3 *VHDL.LUT2D
var M4 *VHDL.LUT2D
var Acc *VHDL.LUT2D

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

	m1 := VHDL.M1().LUT2D
	m2 := VHDL.M2().LUT2D
	//3 := VHDL.M3().LUT2D
	//m4 := VHDL.M4().LUT2D
	acc := VHDL.New2DUnsignedAcc("Acc", 2)

	//rec4multirun()
	//Accurate4bit()

	//[M1,Acc,M2,M1] BIG
	//[M1,M1,M1,M1] SMALL

	bigrec4 := VHDL.NewRecursive4("BigRec4", [4]*VHDL.LUT2D{m1, acc, m2, m1})
	smallrec4 := VHDL.NewRecursive4("SmallRec4", [4]*VHDL.LUT2D{m1, m1, m1, m1})

	bigrec4.GenerateTestData(OutputPath)
	bigrec4.GenerateVHDL(OutputPath)
	smallrec4.GenerateTestData(OutputPath)
	smallrec4.GenerateVHDL(OutputPath)

	xsimbigrec := Viv.CreateXSIM(OutputPath, "BigSim", bigrec4.GenerateVHDLEntityArray())
	xsimsmallrec := Viv.CreateXSIM(OutputPath, "SmallSim", smallrec4.GenerateVHDLEntityArray())
	xsimbigrec.Exec()
	xsimsmallrec.Exec()

	err := Viv.ParseXSIMReport(OutputPath, smallrec4)
	if err != nil {
		log.Fatalln(err)
	}

	err = Viv.ParseXSIMReport(OutputPath, bigrec4)
	if err != nil {
		log.Fatalln(err)
	}

	//Scale

	bigscale := VHDL.New2DScaler(bigrec4, 100)
	smallscale := VHDL.New2DScaler(smallrec4, 100)

	bigscale.GenerateVHDL(OutputPath)
	smallscale.GenerateVHDL(OutputPath)

	tclbig := Viv.CreateVivadoTCL(OutputPath, "big.tcl", bigscale, VivadoSettings)
	tclsmall := Viv.CreateVivadoTCL(OutputPath, "small.tcl", smallscale, VivadoSettings)
	tclbig.Exec()
	tclsmall.Exec()

}

func Accurate4bit() {
	acc4 := VHDL.NewAccurateNumMultiplyer("Acc4", 4, OutputPath)
	acc4.GenerateVHDL(OutputPath)
	acc4.GenerateTestData(OutputPath)

	accscaler := VHDL.New2DScaler(acc4, 100)
	accscaler.GenerateVHDL(OutputPath)
	tcl := Viv.CreateVivadoTCL(OutputPath, "main1.tcl", accscaler, VivadoSettings)
	tcl.Exec()
	util := Viv.ParseUtilizationReport(OutputPath, accscaler)

	fmt.Printf("%+v\n", util)
}

func rec4multirun() {
	M1 = VHDL.M1().LUT2D                  //1
	M2 = VHDL.M2().LUT2D                  //2
	M3 = VHDL.M3().LUT2D                  //3
	M4 = VHDL.M4().LUT2D                  //4
	Acc = VHDL.New2DUnsignedAcc("Acc", 2) //5

	options := []int{1, 2, 3, 4, 5}
	Cartesian4 := cartN(options, options, options, options)

	m := make(map[int]*VHDL.LUT2D)
	m[1] = M1
	m[2] = M2
	m[3] = M3
	m[4] = M4
	m[5] = Acc

	var resultStrings []string

	file, err := os.Create("output/FinalResults.txt")
	if err != nil {
		log.Fatalln(err)
	}

	for i := 400; i < len(Cartesian4); i++ {
		array := [4]*VHDL.LUT2D{m[Cartesian4[i][0]], m[Cartesian4[i][1]], m[Cartesian4[i][2]], m[Cartesian4[i][3]]}
		id_array := fmt.Sprintf("%d,", i)
		resultstring := id_array + rec4scalerRun(array) + "\n"
		file.WriteString(resultstring)
		resultStrings = append(resultStrings, resultstring)
	}

	for _, value := range resultStrings {
		fmt.Println(value)
	}

	file.Close()
}

func rec4scalerRun(array [4]*VHDL.LUT2D) string {
	rec4 := VHDL.NewRecursive4("rec4", array)

	rec4.GenerateTestData(OutputPath)
	rec4.GenerateVHDL(OutputPath)
	xsim := Viv.CreateXSIM(OutputPath, "SimRec4", rec4.GenerateVHDLEntityArray())
	xsim.Exec()
	err := Viv.ParseXSIMReport(OutputPath, rec4)
	if err != nil {
		log.Println(rec4.Overflow())
		log.Fatalln(err)
	}
	rec4scaler := VHDL.New2DScaler(rec4, 100)
	rec4scaler.GenerateVHDL(OutputPath)
	tcl := Viv.CreateVivadoTCL(OutputPath, "main1.tcl", rec4scaler, VivadoSettings)
	tcl.Exec()
	util := Viv.ParseUtilizationReport(OutputPath, rec4scaler)

	Result := NewResult(rec4scaler, util, rec4.Overflow(), rec4.MeanAbsoluteError())
	Result.PrettyPrint()
	log.Println(Result.String())
	Results = append(Results, Result)

	return Result.String()
}

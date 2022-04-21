package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

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
	VivadoSettings.Hierarchical = false
	VivadoSettings.Route = true
	VivadoSettings.Funcsim = false

	// M1 = VHDL.M1().LUT2D
	// M2 = VHDL.M2().LUT2D
	// M3 = VHDL.M3().LUT2D
	// M4 = VHDL.M4().LUT2D
	Acc = VHDL.New2DUnsignedAcc("Acc", 2)

	// for i := 1; i < 65; i++ {
	// 	name := "Acc" + strconv.Itoa(i)
	// 	AccM := VHDL.NewAccurateNumMultiplyer(name, uint(i))
	// 	AccM.GenerateVHDL(OutputPath)
	// 	tcl := Viv.CreateVivadoTCL(OutputPath, name, AccM, VivadoSettings)
	// 	tcl.Exec()

	// 	AccMDSP := VHDL.NewAccurateNumMultiplyer(name+"DSP", uint(i))
	// 	AccMDSP.GenerateVHDL(OutputPath)

	// 	tcldsp := Viv.CreateVivadoTCL(OutputPath, name+"DSP", AccMDSP, VivadoDSPSettings)
	// 	tcldsp.Exec()
	// }

	rec4 := VHDL.NewRecursive4("rec4", [4]VHDL.VHDLEntityMultiplier{Acc, Acc, Acc, Acc})
	rec8 := VHDL.NewRecursive8("rec8", [4]VHDL.VHDLEntityMultiplier{rec4, rec4, rec4, rec4})
	AccM := VHDL.NewMAC(rec8, 128)
	AccM.GenerateVHDL(OutputPath)
	AccM.GenerateTestData(OutputPath)
	VivadoSettings.Funcsim = true

	sim := Viv.CreateXSIM(OutputPath, "acctest", AccM.GenerateVHDLEntityArray())
	sim.SetTemplateSequential(AccM.OutputSize)
	sim.Exec()

	err := Viv.ParseXSIMReport(OutputPath, AccM)
	if err != nil {
		log.Fatalln(err)
	}

	syn := Viv.CreateVivadoTCL(OutputPath, "main.tcl", AccM, VivadoSettings)
	syn.Exec()
	sim.Funcsim()
	syn.ReportPowerPostPlacement()

}

func Rec8total() {
	_, length := ReturnRec8Run(0)

	log.Println(length)

	file, err := os.Create("output/ResultsAll.txt")
	if err != nil {
		log.Fatalln(err)
	}

	defer file.Close()

	var TotalTime time.Duration

	TotalTime = 0

	var Reset = "\033[0m"

	var Yellow = "\033[33m"

	for i := 2245; i < length; i++ {
		start := time.Now()
		testRec8, _ := ReturnRec8Run(i)
		result := Rec8Run(testRec8)
		file.WriteString(result + "\n")
		elapsed := time.Since(start)
		TotalTime = TotalTime + elapsed
		log.Printf(Yellow+"Last run took %s\n"+Reset, elapsed)
		log.Printf(Yellow+"Total time: %s\n"+Reset, TotalTime)
	}

}

func ReturnRec8Run(i int) (*VHDL.Recursive8, int) {
	M1 = VHDL.M1().LUT2D
	M2 = VHDL.M2().LUT2D
	M3 = VHDL.M3().LUT2D
	M4 = VHDL.M4().LUT2D
	Acc = VHDL.New2DUnsignedAcc("Acc", 2)

	//[Acc,M4,M2,M1],1399,1.875000 //1
	//[Acc,M1,M4,M1],1481,1.625000 //2
	//[Acc,M1,M1,M1],1494,1.125000 //3
	//[Acc,M1,M3,M1],1496,1.015625 //4
	//[Acc,Acc,M1,M1],1499,0.625000 //5
	//[Acc,M1,Acc,Acc],1598,0.500000 //6

	options := []int{1, 2, 3, 4, 5, 6, 7}
	Cartesian4 := cartN(options, options, options, options)

	rec := make(map[int]*VHDL.Recursive4)
	rec[1] = VHDL.NewRecursive4("RecA421", [4]VHDL.VHDLEntityMultiplier{Acc, M4, M2, M1})
	rec[2] = VHDL.NewRecursive4("RecA141", [4]VHDL.VHDLEntityMultiplier{Acc, M1, M4, M1})
	rec[3] = VHDL.NewRecursive4("RecA111", [4]VHDL.VHDLEntityMultiplier{Acc, M1, M1, M1})
	rec[4] = VHDL.NewRecursive4("RecA131", [4]VHDL.VHDLEntityMultiplier{Acc, M1, M3, M1})
	rec[5] = VHDL.NewRecursive4("RecAA11", [4]VHDL.VHDLEntityMultiplier{Acc, Acc, M1, M1})
	rec[6] = VHDL.NewRecursive4("RecA1AA", [4]VHDL.VHDLEntityMultiplier{Acc, M1, Acc, Acc})
	rec[7] = VHDL.NewRecursive4("RecAAAA", [4]VHDL.VHDLEntityMultiplier{Acc, Acc, Acc, Acc})

	rec4array := [4]VHDL.VHDLEntityMultiplier{rec[Cartesian4[i][0]], rec[Cartesian4[i][1]], rec[Cartesian4[i][2]], rec[Cartesian4[i][3]]}

	valuesText := []string{}

	for j := range Cartesian4[i] {
		number := Cartesian4[i][j]
		text := strconv.Itoa(number)
		valuesText = append(valuesText, text)
	}

	result := strings.Join(valuesText, "")

	rec8 := VHDL.NewRecursive8("Rec"+result, rec4array)

	return rec8, len(Cartesian4)
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
		array := [4]VHDL.VHDLEntityMultiplier{m[Cartesian4[i][0]], m[Cartesian4[i][1]], m[Cartesian4[i][2]], m[Cartesian4[i][3]]}
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

func rec4scalerRun(array [4]VHDL.VHDLEntityMultiplier) string {
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

func Rec8Run(rec8 *VHDL.Recursive8) string {

	rec8.GenerateTestData(OutputPath)
	rec8.GenerateVHDL(OutputPath)
	xsim := Viv.CreateXSIM(OutputPath, "Sim_"+rec8.EntityName, rec8.GenerateVHDLEntityArray())
	xsim.Exec()
	err := Viv.ParseXSIMReport(OutputPath, rec8)
	if err != nil {
		log.Println(rec8.Overflow())
		log.Fatalln(err)
	}
	rec8scaler := VHDL.New2DScaler(rec8, 100)
	rec8scaler.GenerateVHDL(OutputPath)
	tcl := Viv.CreateVivadoTCL(OutputPath, "main1.tcl", rec8scaler, VivadoSettings)
	tcl.Exec()
	util := Viv.ParseUtilizationReport(OutputPath, rec8scaler)

	Result := NewResult(rec8scaler, util, rec8.Overflow(), rec8.MeanAbsoluteError())
	Result.PrettyPrint()
	log.Println(Result.String())
	Results = append(Results, Result)

	return Result.String()
}

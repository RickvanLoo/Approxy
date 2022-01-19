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
	VivadoSettings.Hierarchical = false

	//o1 = AH*BH
	//o2 = AH*BL
	//o3 = AL*BH
	//o4 = AL*BL

	m1 := VHDL.M1().LUT2D
	m2 := VHDL.M2().LUT2D
	m3 := VHDL.M3().LUT2D
	m4 := VHDL.M4().LUT2D

	rec4 := VHDL.NewRecursive4("ApproxRec4", [4]*VHDL.LUT2D{m1, m2, m3, m4})
	rec8 := VHDL.NewRecursive8("ApproxRec8", [4]*VHDL.Recursive4{rec4, rec4, rec4, rec4})
	rec8.GenerateTestData(OutputPath)
	rec8.GenerateVHDL(OutputPath)

	test := Viv.CreateXSIM(OutputPath, "rec8Test", rec8.GenerateVHDLEntityArray())
	test.Exec()

	rec8scaler := VHDL.New2DScaler(rec8, 100)
	rec8scaler.GenerateVHDL(OutputPath)

	synth := Viv.CreateVivadoTCL(OutputPath, "main.tcl", rec8scaler, VivadoSettings)
	synth.Exec()

}

func Improvcheck() {
	m1 := VHDL.M1().LUT2D
	m2 := VHDL.M2().LUT2D
	// m3 := VHDL.M3().LUT2D
	m4 := VHDL.M4().LUT2D
	acc := VHDL.New2DUnsignedAcc("Acc", 2)

	acc_perfect := VHDL.NewAccurateNumMultiplyer("AccPerfect", 4, OutputPath)
	acc_perfect.GenerateVHDL(OutputPath)
	mac_acc_perfect := VHDL.NewMAC(acc_perfect)
	mac_acc_perfect.GenerateVHDL(OutputPath)

	corrAcc := VHDL.NewCorrelator("CorrAccurate", [4]*VHDL.MAC{mac_acc_perfect, mac_acc_perfect, mac_acc_perfect, mac_acc_perfect})
	corrAcc.GenerateVHDL(OutputPath)

	//[Acc,M4,M2,M1] = small

	rec4 := VHDL.NewRecursive4("ApproxRec4", [4]*VHDL.LUT2D{acc, m4, m2, m1})
	rec4.GenerateVHDL(OutputPath)

	mac_rec4 := VHDL.NewMAC(rec4)

	corr := VHDL.NewCorrelator("CorrApprox", [4]*VHDL.MAC{mac_rec4, mac_rec4, mac_rec4, mac_rec4})
	corr.GenerateVHDL(OutputPath)

	VivadoProjectAccurate := Viv.CreateVivadoTCL(OutputPath, "approx.tcl", corrAcc, VivadoSettings)
	VivadoProjectAccurate.Exec()

	VivadoProject := Viv.CreateVivadoTCL(OutputPath, "main.tcl", corr, VivadoSettings)
	VivadoProject.Exec()

	///CHECK IF IMPROVEMENT

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

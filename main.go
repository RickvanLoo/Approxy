package main

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	math_rand "math/rand"
	"os"
	"strconv"
	"time"

	VHDL "badmath/VHDL"
	Viv "badmath/Vivado"
)

var OutputPath string
var ReportPath string
var Reset string
var Yellow string

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

func init() {
	//Setup random seed for RNG
	//https://stackoverflow.com/questions/12321133/how-to-properly-seed-random-number-generator
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	math_rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))

	fmt.Println("badmath...")

	OutputPath = "output"
	ReportPath = "report"
	ClearPath(OutputPath)
	CreatePath(OutputPath)
	CreatePath(ReportPath)

	VivadoSettings = new(Viv.VivadoTCLSettings)
	VivadoSettings.NO_DSP = true
	VivadoSettings.OOC = true
	VivadoSettings.PartName = "Xc7z030fbg676-3"
	VivadoSettings.Placement = true
	VivadoSettings.Utilization = true
	VivadoSettings.WriteCheckpoint = true
	VivadoSettings.Hierarchical = true
	VivadoSettings.Route = true
	VivadoSettings.Funcsim = true
	VivadoSettings.Clk = false //IMPORTANT
	VivadoSettings.Timing = true

	M1 = VHDL.M1().LUT2D
	M2 = VHDL.M2().LUT2D
	M3 = VHDL.M3().LUT2D
	M4 = VHDL.M4().LUT2D
	Acc = VHDL.New2DUnsignedAcc("Acc", 2)

	Reset = "\033[0m"
	Yellow = "\033[33m"
}

func main() {

	ErrorRun(500, 1000)
	//Rec4Run(100)
}

func ErrorRun(ScaleN int, Nval int) {
	CurrentRun := Viv.StartRun(ReportPath, OutputPath, "ErrorRun_"+strconv.Itoa(ScaleN)+"_"+strconv.Itoa(Nval))
	CurrentRun.ClearData()
	CurrentRun.AddData("Disc", "Running "+strconv.Itoa(ScaleN)+" accurate 8-bit Multipliers to determine power error, i="+strconv.Itoa(Nval))

	rec8 := VHDL.NewAccurateNumMultiplyer("recacc8", 8)
	AccM := VHDL.New2DScaler(rec8, uint(ScaleN))

	AccM.GenerateVHDL(OutputPath)
	AccM.GenerateTestData(OutputPath)

	sim := Viv.CreateXSIM(OutputPath, AccM.EntityName+"_test", AccM.GenerateVHDLEntityArray())
	sim.SetTemplateScaler(uint(ScaleN))
	sim.Exec()

	err := Viv.ParseXSIMReport(OutputPath, AccM)
	if err != nil {
		log.Fatalln(err)
	}

	syn := Viv.CreateVivadoTCL(OutputPath, "main.tcl", AccM, VivadoSettings)
	syn.Exec()

	for i := 0; i < 50; i++ {
		if CurrentRun.Exists("Run" + strconv.Itoa(i)) {
			log.Printf(Yellow + "Warning, skipping Run: " + AccM.EntityName + "\n" + Reset)
			continue
		}
		start := time.Now()

		sim.CreateFile(true)
		VHDL.NormalTestData(AccM, OutputPath, uint(Nval))
		sim.Funcsim()
		syn.PowerPostPlacementGeneration()

		elapsed := time.Since(start)
		log.Printf(Yellow+"Last run took %s\n"+Reset, elapsed)

		Report := Viv.CreateReport(OutputPath, AccM)
		Report.EntityName = "Run" + strconv.Itoa(i)
		Report.AddData("Error", "0")
		Report.AddData("ElapsedTime", elapsed.String())
		CurrentRun.AddReport(*Report)
	}
}

func Rec4Run(ScaleN int, Nval int) {
	CurrentRun := Viv.StartRun(ReportPath, OutputPath, "Rec4Run_"+strconv.Itoa(ScaleN)+"_"+strconv.Itoa(Nval))
	CurrentRun.ClearData()
	CurrentRun.AddData("Disc", "Full Recursive 4-bit run using M1,M2,M3,M4,Acc, N="+strconv.Itoa(ScaleN)+" i="+strconv.Itoa(Nval))

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

	for i := 0; i < len(Cartesian4); i++ {
		start := time.Now()

		array := [4]VHDL.VHDLEntityMultiplier{m[Cartesian4[i][0]], m[Cartesian4[i][1]], m[Cartesian4[i][2]], m[Cartesian4[i][3]]}
		Name := "Rec_" + strconv.Itoa(Cartesian4[i][0]) + strconv.Itoa(Cartesian4[i][1]) + strconv.Itoa(Cartesian4[i][2]) + strconv.Itoa(Cartesian4[i][3])
		rec4 := VHDL.NewRecursive4(Name, array)
		rec4scaler := VHDL.New2DScaler(rec4, uint(ScaleN))

		if CurrentRun.Exists(rec4scaler.EntityName) {
			log.Printf(Yellow + "Warning, skipping Entity: " + rec4scaler.EntityName + "\n" + Reset)
			continue
		}

		rec4scaler.GenerateVHDL(OutputPath)
		rec4scaler.GenerateTestData(OutputPath)

		sim := Viv.CreateXSIM(OutputPath, rec4scaler.EntityName+"_test", rec4scaler.GenerateVHDLEntityArray())
		sim.SetTemplateScaler(uint(ScaleN))
		sim.Exec()

		err := Viv.ParseXSIMReport(OutputPath, rec4scaler)
		if err != nil {
			log.Fatalln(err)
		}

		syn := Viv.CreateVivadoTCL(OutputPath, "main.tcl", rec4scaler, VivadoSettings)
		syn.Exec()
		sim.CreateFile(true)
		VHDL.NormalTestData(rec4scaler, OutputPath, uint(Nval))
		sim.Funcsim()
		syn.PowerPostPlacementGeneration()

		elapsed := time.Since(start)
		log.Printf(Yellow+"Last run took %s\n"+Reset, elapsed)

		Report := Viv.CreateReport(OutputPath, rec4scaler)
		Report.AddData("MeanAbsoluteError", strconv.FormatFloat(rec4.MeanAbsoluteError(), 'E', -1, 64))
		Report.AddData("AverageRelativeError", strconv.FormatFloat(rec4.AverageRelativeError(), 'E', -1, 64))
		Report.AddData("Overflow", strconv.FormatBool(rec4.Overflow()))
		Report.AddData("ElapsedTime", elapsed.String())
		CurrentRun.AddReport(*Report)

		ClearPath(OutputPath)
		CreatePath(OutputPath)
	}

}

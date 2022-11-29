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

	VHDL "approxy/VHDL"
	Viv "approxy/Vivado"
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

// Be careful with this one
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

// Run once at startup
func init() {
	//Setup random seed for RNG
	//https://stackoverflow.com/questions/12321133/how-to-properly-seed-random-number-generator
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	math_rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))

	fmt.Println("approxy...")

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
	PaperExample()
}

// Listing A.1 of Appendix within thesis + added timing reporting
func PaperExample() {
	//Start of Approxy Run
	start := time.Now()
	CurrentRun := Viv.StartRun(ReportPath, OutputPath, "Rec_1311")
	CurrentRun.ClearData()

	//Generation of Recursive Multipler
	Rec_1234 := VHDL.NewRecursive4("Rec1234", [4]VHDL.VHDLEntityMultiplier{M1, M2, M3, M4})
	Rec_1234.GenerateTestData(OutputPath)
	Rec_1234.GenerateVHDL(OutputPath)

	//Optional Verification of Recursive4
	verify := Viv.CreateXSIM(OutputPath, "prePR", Rec_1234.GenerateVHDLEntityArray())
	verify.SetTemplateMultiplier()
	time_till_analysis := time.Since(start)
	start = time.Now()

	verify.Exec()
	time_analysis := time.Since(start)
	start = time.Now()

	err := Viv.ParseXSIMReport(OutputPath, Rec_1234)
	if err != nil {
		log.Fatalln(err)
	}

	//Generation of N=1000 Scaler of Rec1234
	Rec_1234_scaler := VHDL.New2DScaler(Rec_1234, 1000)
	Rec_1234_scaler.GenerateTestData(OutputPath)
	Rec_1234_scaler.GenerateVHDL(OutputPath)

	//Synth + Place + Route
	viv := Viv.CreateVivadoTCL(OutputPath, "main.tcl", Rec_1234_scaler, VivadoSettings)
	time_till_synth := time.Since(start)
	start = time.Now()

	viv.Exec()
	time_spr := time.Since(start)
	start = time.Now()

	//PostSynthesisAnalysis
	post_analysis := Viv.CreateXSIM(OutputPath, "postPR", Rec_1234_scaler.GenerateVHDLEntityArray())
	post_analysis.SetTemplateScaler(1000)
	post_analysis.CreateFile(true)                         //Create PostPR Testbench
	VHDL.NormalTestData(Rec_1234_scaler, OutputPath, 1000) //Create i=1000 Normal Test Data for 4-bit
	time_beforefunc := time.Since(start)
	start = time.Now()

	post_analysis.Funcsim()
	time_func := time.Since(start) //Funcsim
	start = time.Now()

	viv.PowerPostPlacementGeneration() //Export PostPR data

	//Create Report
	Report := Viv.CreateReport(OutputPath, Rec_1234)
	Report.AddData("MAE_Uniform", strconv.FormatFloat(Rec_1234.MeanAbsoluteError(), 'E', -1, 64))
	Report.AddData("MAE_Normal_1000", strconv.FormatFloat(Rec_1234.MeanAbsoluteErrorNormalDist(1000), 'E', -1, 64))
	Report.AddData("Overflow", strconv.FormatBool(Rec_1234.Overflow()))
	Report.AddData("time_till_analysis", time_till_analysis.String())
	Report.AddData("time_analysis", time_analysis.String())
	Report.AddData("time_till_synth", time_till_synth.String())
	Report.AddData("time_synthpr", time_spr.String())
	Report.AddData("time_beforefunc", time_beforefunc.String())
	Report.AddData("time_func", time_func.String())

	CurrentRun.AddReport(*Report)

	time_reportgen := time.Since(start)
	log.Println(time_reportgen)
}

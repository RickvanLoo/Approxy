package main

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	math_rand "math/rand"
	"os"
	"strconv"

	VHDL "github.com/RickvanLoo/Approxy/vhdl"
	Viv "github.com/RickvanLoo/Approxy/vivado"
)

// OutputPath points to the path that is used as temp folder for exported VHDL, testdata, checkpoints and Vivado report
var OutputPath string

// ReportPath point to the path that is used for Approxy Reports
var ReportPath string

// Reset clears stdout formatting
var Reset string

// Yellow adds yellow colour formatting to stdout
var Yellow string

// VivadoSettings used globally
var VivadoSettings *Viv.VivadoTCLSettings

// M1 Multiplier available globally
var M1 *VHDL.LUT2D

// M2 Multiplier available globally
var M2 *VHDL.LUT2D

// M3 Multiplier available globally
var M3 *VHDL.LUT2D

// M4 Multiplier available globally
var M4 *VHDL.LUT2D

// Acc Multiplier available globally
var Acc *VHDL.LUT2D

// ClearPath completely removes a path folder with all its contents
// Be careful with this one
func ClearPath(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		log.Println(err)
	}
}

// CreatePath creates the path folder described in path
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
	// ClearPath(OutputPath)
	// CreatePath(OutputPath)
	// CreatePath(ReportPath)

	VivadoSettings = new(Viv.VivadoTCLSettings)
	VivadoSettings.NODSP = true
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

// PaperExample starts an Approxy Run
// First a 4-bit Recursive Multiplier, Rec_1234, is created on basis of {M1,M2,M3,M4}
// This Multiplier is verified using XSIM, the report is parsed, and execution of this function will fail if the design does not verify
// A Scaler is created to be able to design and simulate 1000 times Rec_1234
// This scaled design is synthesized, placed and routed using Vivado
// A new XSIM environment is created for the Scaled design for PostPR analysis
// The simulation is executed and fed i=1000 normally distributed input values
// A report is created for this N=1000 scaled Rec_1234 and added to the Run
// Listing A.1 of Appendix within thesis
func PaperExample() {
	//Start of Approxy Run
	CurrentRun := Viv.StartRun(ReportPath, OutputPath, "Rec_1234")
	CurrentRun.ClearData()

	//Generation of Recursive Multipler
	Rec1234 := VHDL.NewRecursive4("Rec1234", [4]VHDL.VHDLEntityMultiplier{M1, M2, M3, M4})
	Rec1234.GenerateTestData(OutputPath)
	Rec1234.GenerateVHDL(OutputPath)

	//Optional Verification of Recursive4
	verify := Viv.CreateXSIM(OutputPath, "prePR", Rec1234.GenerateVHDLEntityArray())
	verify.SetTemplateMultiplier()

	verify.Exec()

	err := Viv.ParseXSIMReport(OutputPath, Rec1234)
	if err != nil {
		log.Fatalln(err)
	}

	//Generation of N=1000 Scaler of Rec1234
	Rec1234scaler := VHDL.New2DScaler(Rec1234, 1000)
	Rec1234scaler.GenerateTestData(OutputPath)
	Rec1234scaler.GenerateVHDL(OutputPath)

	//Synth + Place + Route
	viv := Viv.CreateVivadoTCL(OutputPath, "main.tcl", Rec1234scaler, VivadoSettings)

	viv.Exec()

	//PostSynthesisAnalysis
	postanalysis := Viv.CreateXSIM(OutputPath, "postPR", Rec1234scaler.GenerateVHDLEntityArray())
	postanalysis.SetTemplateScaler(1000)
	postanalysis.CreateFile(true)                        //Create PostPR Testbench
	VHDL.NormalTestData(Rec1234scaler, OutputPath, 1000) //Create i=1000 Normal Test Data for 4-bit

	postanalysis.Funcsim()

	viv.PowerPostPlacementGeneration() //Export PostPR data

	//Create Report
	Report := Viv.CreateReport(OutputPath, Rec1234scaler)
	Report.AddData("MAE_Uniform", strconv.FormatFloat(Rec1234.MeanAbsoluteError(), 'E', -1, 64))
	Report.AddData("MAE_Normal_1000", strconv.FormatFloat(Rec1234.MeanAbsoluteErrorNormalDist(1000), 'E', -1, 64))
	Report.AddData("Overflow", strconv.FormatBool(Rec1234.Overflow()))

	CurrentRun.AddReport(*Report)
}

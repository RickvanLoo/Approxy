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
	//CreateSingleVHDL()
	//RedoExternal()
	PaperExample()
}

func PaperExample() {
	//Start of BadMath Run
	start := time.Now()
	CurrentRun := Viv.StartRun(ReportPath, OutputPath, "Rec_1234")
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
	Report := Viv.CreateReport(OutputPath, Rec_1234_scaler)
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

func CreateSingleVHDL() {
	// Rec1311 := VHDL.NewRecursive4("Rec1311", [4]VHDL.VHDLEntityMultiplier{M1, M3, M1, M1})
	// Rec1311.GenerateTestData(OutputPath)
	// Rec1311.GenerateVHDL(OutputPath)

	// MAC_1311 := VHDL.NewMAC(Rec1311, 64)
	// MAC_1311.GenerateTestData(OutputPath)
	// MAC_1311.GenerateVHDL(OutputPath)
	// MAC_1311 := VHDL.NewUnsignedAccurateMAC(4, 64)
	// MAC_1311.GenerateTestData(OutputPath)
	// MAC_1311.GenerateVHDL(OutputPath)

	approx1 := VHDL.NewExternalMult("approx1", 4, "mult_approx_a4.vhd")
	approx1.GenerateVHDL(OutputPath)
	approx1.GenerateTestData(OutputPath)
	sim_behav := Viv.CreateXSIM(OutputPath, "behavcheck", approx1.GenerateVHDLEntityArray())
	sim_behav.SetTemplateReverse()
	sim_behav.Exec()
	approx1.ParseXSIMOutput(OutputPath)

	MAC_approx1 := VHDL.NewMAC(approx1, 64)
	MAC_approx1.GenerateTestData(OutputPath)
	MAC_approx1.GenerateVHDL(OutputPath)

	approx1_scaler := VHDL.New2DScaler(MAC_approx1, 1000)
	approx1_scaler.SetMAC(true, MAC_approx1.OutputSize)
	approx1_scaler.GenerateTestData(OutputPath)
	approx1_scaler.GenerateVHDL(OutputPath)

	sim_scaler := Viv.CreateXSIM(OutputPath, approx1_scaler.EntityName+"_test", approx1_scaler.GenerateVHDLEntityArray())
	sim_scaler.SetTemplateSequentialScaler(1000, MAC_approx1.OutputSize)

	syn := Viv.CreateVivadoTCL(OutputPath, "main.tcl", approx1_scaler, VivadoSettings)
	syn.Exec()
	sim_scaler.CreateFile(true)
	MAC_approx1.ResetVal()
	VHDL.NormalTestData(approx1_scaler, OutputPath, 1000)
	sim_scaler.Funcsim()
	syn.PowerPostPlacementGeneration()

}

func RedoExternal() {
	CurrentRun := Viv.StartRun(ReportPath, OutputPath, "REDO_External_XMAA")
	CurrentRun.ClearData()

	ExternalMult := VHDL.NewExternalMult("XA_Config_V1_First_0000", 4, "lpACLib/ConfigMultV1_XA/Config4x4MultV1First.vhd")
	ExternalMult.GenerateVHDL(OutputPath)
	ExternalMult.GenerateTestData(OutputPath)

	ExtraFile0 := VHDL.NewExternalMult("dontcare0", 4, "lpACLib/ConfigMultV1_XA/Config2x2MultV1.vhd")
	ExtraFile0.GenerateVHDL(OutputPath)
	ExtraFile0.GenerateTestData(OutputPath)
	ExternalMult.AddVHDLEntity(ExtraFile0)

	ExtraFile1 := VHDL.NewExternalMult("dontcare1", 4, "lpACLib/ConfigMultV1_XA/Approx2x2MultV1.vhd")
	ExtraFile1.GenerateVHDL(OutputPath)
	ExtraFile1.GenerateTestData(OutputPath)
	ExternalMult.AddVHDLEntity(ExtraFile1)

	ExtraFile2 := VHDL.NewExternalMult("dontcare2", 4, "lpACLib/ConfigMultV1_XA/AdderIMPACTFirstApproxMultiBit.vhd")
	ExtraFile2.GenerateVHDL(OutputPath)
	ExtraFile2.GenerateTestData(OutputPath)
	ExternalMult.AddVHDLEntity(ExtraFile2)

	ExtraFile3 := VHDL.NewExternalMult("dontcare3", 4, "lpACLib/ConfigMultV1_XA/AdderAccurateOneBit.vhd")
	ExtraFile3.GenerateVHDL(OutputPath)
	ExtraFile3.GenerateTestData(OutputPath)
	ExternalMult.AddVHDLEntity(ExtraFile3)

	ExtraFile4 := VHDL.NewExternalMult("dontcare4", 4, "lpACLib/ConfigMultV1_XA/AdderIMPACTFirstApproxOneBit.vhd")
	ExtraFile4.GenerateVHDL(OutputPath)
	ExtraFile4.GenerateTestData(OutputPath)
	ExternalMult.AddVHDLEntity(ExtraFile4)

	sim_single := Viv.CreateXSIM(OutputPath, ExternalMult.EntityName+"_test", ExternalMult.GenerateVHDLEntityArray())
	sim_single.SetTemplateReverse()
	sim_single.Exec()
	ExternalMult.ParseXSIMOutput(OutputPath)

	Ext_scaler := VHDL.New2DScaler(ExternalMult, 1000)
	Ext_scaler.GenerateVHDL(OutputPath)

	sim_scaler := Viv.CreateXSIM(OutputPath, Ext_scaler.EntityName+"_test", Ext_scaler.GenerateVHDLEntityArray())
	sim_scaler.SetTemplateScaler(1000)

	syn := Viv.CreateVivadoTCL(OutputPath, "main.tcl", Ext_scaler, VivadoSettings)
	syn.Exec()
	sim_scaler.CreateFile(true)
	VHDL.NormalTestData(Ext_scaler, OutputPath, 1000)
	sim_scaler.Funcsim()
	syn.PowerPostPlacementGeneration()

	Report := Viv.CreateReport(OutputPath, Ext_scaler)
	Report.AddData("MAE_Uniform", strconv.FormatFloat(ExternalMult.MeanAbsoluteError(), 'E', -1, 64))
	Report.AddData("MAE_Normal_1000", strconv.FormatFloat(ExternalMult.MeanAbsoluteErrorNormalDist(1000), 'E', -1, 64))
	Report.AddData("ARE", strconv.FormatFloat(ExternalMult.AverageRelativeError(), 'E', -1, 64))
	CurrentRun.AddReport(*Report)
}

func RedoRec4Error(ScaleN int, Nval int) {
	CurrentRun := Viv.StartRun(ReportPath, OutputPath, "REDO_Rec4Run_"+strconv.Itoa(ScaleN)+"_"+strconv.Itoa(Nval))
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

		array := [4]VHDL.VHDLEntityMultiplier{m[Cartesian4[i][0]], m[Cartesian4[i][1]], m[Cartesian4[i][2]], m[Cartesian4[i][3]]}
		Name := "Rec_" + strconv.Itoa(Cartesian4[i][0]) + strconv.Itoa(Cartesian4[i][1]) + strconv.Itoa(Cartesian4[i][2]) + strconv.Itoa(Cartesian4[i][3])
		rec4 := VHDL.NewRecursive4(Name, array)
		rec4scaler := VHDL.New2DScaler(rec4, uint(ScaleN))

		rec4scaler.GenerateVHDL(OutputPath)
		rec4scaler.GenerateTestData(OutputPath)

		Report := Viv.CreateReport(OutputPath, rec4scaler)
		Report.AddData("MeanAbsoluteErrorNORM", strconv.FormatFloat(rec4.MeanAbsoluteErrorNormalDist(Nval), 'E', -1, 64))
		Report.AddData("Overflow", strconv.FormatBool(rec4.Overflow()))
		CurrentRun.AddReport(*Report)

		ClearPath(OutputPath)
		CreatePath(OutputPath)
	}
}

func AccurateRun() {
	CurrentRun := Viv.StartRun(ReportPath, OutputPath, "AccurateRunFINAL")
	CurrentRun.ClearData()
	CurrentRun.AddData("Disc", "AccurateRunFinal 4-bit N=1000, I=1000")

	Acc4 := VHDL.NewAccurateNumMultiplyer("Acc4", 4)
	Acc4Scale := VHDL.New2DScaler(Acc4, 1000)

	Acc4Scale.GenerateVHDL(OutputPath)
	Acc4Scale.GenerateTestData(OutputPath)

	sim := Viv.CreateXSIM(OutputPath, Acc4Scale.EntityName+"_test", Acc4Scale.GenerateVHDLEntityArray())
	sim.SetTemplateScaler(1000)

	syn := Viv.CreateVivadoTCL(OutputPath, "main.tcl", Acc4Scale, VivadoSettings)
	syn.Exec()
	sim.CreateFile(true)
	VHDL.NormalTestData(Acc4Scale, OutputPath, 1000)
	sim.Funcsim()
	syn.PowerPostPlacementGeneration()

	Report := Viv.CreateReport(OutputPath, Acc4Scale)
	CurrentRun.AddReport(*Report)
}

func ErrorRun(ScaleN int, Nval int) {
	CurrentRun := Viv.StartRun(ReportPath, OutputPath, "ErrorRun4_"+strconv.Itoa(ScaleN)+"_"+strconv.Itoa(Nval))
	CurrentRun.ClearData()
	CurrentRun.AddData("Disc", "Running "+strconv.Itoa(ScaleN)+" accurate 4-bit Multipliers to determine power error, i="+strconv.Itoa(Nval))

	rec8 := VHDL.NewAccurateNumMultiplyer("recacc4", 4)
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
		timeleft := time.Duration((len(Cartesian4)-i)*10) * time.Minute
		finishedat := time.Now().Add(timeleft)

		log.Printf(Yellow+"%d Rec4 Simulations left!\n"+Reset, len(Cartesian4)-i)
		log.Println(Yellow + "Finished at: " + finishedat.Format("02/01/2006 15:04:05") + Reset)
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

func PowerEst() {
	CurrentRun := Viv.StartRun(ReportPath, OutputPath, "PowerEst")
	CurrentRun.ClearData()
	CurrentRun.AddData("Disc", "Running PowerEst")

	rec4 := VHDL.NewRecursive4("RecM3", [4]VHDL.VHDLEntityMultiplier{M3, M3, M3, M3})
	rec4_scaler := VHDL.New2DScaler(rec4, 1000)

	rec4_scaler.GenerateVHDL(OutputPath)
	rec4_scaler.GenerateTestData(OutputPath)

	sim := Viv.CreateXSIM(OutputPath, rec4_scaler.EntityName+"_test", rec4_scaler.GenerateVHDLEntityArray())
	sim.SetTemplateScaler(1000)
	sim.Exec()

	err := Viv.ParseXSIMReport(OutputPath, rec4_scaler)
	if err != nil {
		log.Fatalln(err)
	}

	syn := Viv.CreateVivadoTCL(OutputPath, "main.tcl", rec4_scaler, VivadoSettings)
	syn.Exec()

	sim.CreateFile(true)
	VHDL.NormalTestData(rec4_scaler, OutputPath, 1000)
	sim.Funcsim()
	syn.PowerPostPlacementGeneration()

	Report := Viv.CreateReport(OutputPath, rec4_scaler)
	Report.AddData("MeanAbsoluteError", strconv.FormatFloat(rec4.MeanAbsoluteError(), 'E', -1, 64))
	Report.AddData("AverageRelativeError", strconv.FormatFloat(rec4.AverageRelativeError(), 'E', -1, 64))
	Report.AddData("Overflow", strconv.FormatBool(rec4.Overflow()))
	CurrentRun.AddReport(*Report)
	ClearPath(OutputPath)
	CreatePath(OutputPath)
}

func SingleRun(name string, Entity VHDL.VHDLEntityMultiplier, scaleN int, testi int) {
	ClearPath(OutputPath)
	CreatePath(OutputPath)

	CurrentRun := Viv.StartRun(ReportPath, OutputPath, name+strconv.Itoa(scaleN)+"_"+strconv.Itoa(testi))
	CurrentRun.ClearData()
	CurrentRun.AddData("Disc", "Running "+name+strconv.Itoa(scaleN)+"_"+strconv.Itoa(testi))

	Entity_Scaler := VHDL.New2DScaler(Entity, uint(scaleN))

	Entity_Scaler.GenerateVHDL(OutputPath)
	Entity_Scaler.GenerateTestData(OutputPath)

	sim := Viv.CreateXSIM(OutputPath, Entity_Scaler.EntityName+"_test", Entity_Scaler.GenerateVHDLEntityArray())
	sim.SetTemplateScaler(uint(scaleN))
	sim.Exec()

	err := Viv.ParseXSIMReport(OutputPath, Entity_Scaler)
	if err != nil {
		log.Fatalln(err)
	}

	syn := Viv.CreateVivadoTCL(OutputPath, "main.tcl", Entity_Scaler, VivadoSettings)
	syn.Exec()

	sim.CreateFile(true)
	VHDL.NormalTestData(Entity_Scaler, OutputPath, uint(testi))
	sim.Funcsim()
	syn.PowerPostPlacementGeneration()

	Report := Viv.CreateReport(OutputPath, Entity_Scaler)
	Report.AddData("MeanAbsoluteError", strconv.FormatFloat(Entity.MeanAbsoluteError(), 'E', -1, 64))
	Report.AddData("Overflow", strconv.FormatBool(Entity.Overflow()))
	CurrentRun.AddReport(*Report)

}

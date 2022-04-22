package main

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	math_rand "math/rand"
	"os"
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

func init() {
	//Setup random seed for RNG
	//https://stackoverflow.com/questions/12321133/how-to-properly-seed-random-number-generator
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	math_rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
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
	VivadoSettings.Clk = false //IMPORTANT

	M1 = VHDL.M1().LUT2D
	M2 = VHDL.M2().LUT2D
	M3 = VHDL.M3().LUT2D
	M4 = VHDL.M4().LUT2D
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

	var Reset = "\033[0m"

	var Yellow = "\033[33m"

	start := time.Now()

	// rec6 := VHDL.NewRecursive4("rec4_6", [4]VHDL.VHDLEntityMultiplier{Acc, M1, Acc, Acc})
	// rec1 := VHDL.NewRecursive4("rec4_1", [4]VHDL.VHDLEntityMultiplier{Acc, M4, M2, M1})

	// rec8 := VHDL.NewRecursive8("rec6111", [4]VHDL.VHDLEntityMultiplier{rec6, rec1, rec1, rec1})
	rec8 := VHDL.NewAccurateNumMultiplyer("recacc8", 8)
	AccM := VHDL.New2DScaler(rec8, 500)
	AccM.GenerateVHDL(OutputPath)
	AccM.GenerateTestData(OutputPath)
	VivadoSettings.Funcsim = true

	sim := Viv.CreateXSIM(OutputPath, "acctest", AccM.GenerateVHDLEntityArray())
	sim.SetTemplateScaler(500)
	sim.Exec()

	err := Viv.ParseXSIMReport(OutputPath, AccM)
	if err != nil {
		log.Fatalln(err)
	}

	syn := Viv.CreateVivadoTCL(OutputPath, "main.tcl", AccM, VivadoSettings)
	syn.Exec()
	sim.CreateFile(true)
	VHDL.RandomizeTestData(AccM, OutputPath, 1000)
	sim.Funcsim()
	syn.ReportPowerPostPlacement()

	elapsed := time.Since(start)
	log.Printf(Yellow+"Last run took %s\n"+Reset, elapsed)

}

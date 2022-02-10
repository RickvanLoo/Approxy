package main

import (
	VHDL "badmath/VHDL"
	Viv "badmath/Vivado"
	"testing"
)

func startBench() {
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
	VivadoSettings.Hierarchical = true

	M1 = VHDL.M1().LUT2D
	M2 = VHDL.M2().LUT2D
	M3 = VHDL.M3().LUT2D
	M4 = VHDL.M4().LUT2D
	Acc = VHDL.New2DUnsignedAcc("Acc", 2)
}

func bit2MultN(N uint, b *testing.B) {
	startBench()
	for i := 0; i < b.N; i++ {
		M1.GenerateVHDL(OutputPath)
		M1.GenerateTestData(OutputPath)
		M1_Scaler := VHDL.New2DScaler(M1, N)
		M1_Scaler.GenerateVHDL(OutputPath)
		test := Viv.CreateXSIM(OutputPath, "SimBenchMarkM1", M1.GenerateVHDLEntityArray())
		test.Exec()
		err := Viv.ParseXSIMReport(OutputPath, M1)
		if err != nil {
			b.Fail()
		}
		synthplace := Viv.CreateVivadoTCL(OutputPath, "ExecBenchM1", M1_Scaler, VivadoSettings)
		synthplace.Exec()
		util := Viv.ParseUtilizationReport(OutputPath, M1_Scaler)
		if util.TotalLUT == 0 {
			b.Fail()
		}
	}
}

func bit4MultN(N uint, b *testing.B) {
	startBench()
	for i := 0; i < b.N; i++ {
		rec4test := VHDL.NewRecursive4("rec4bench", [4]VHDL.VHDLEntityMultiplier{M1, M1, M1, M1})
		rec4test.GenerateVHDL(OutputPath)
		rec4test.GenerateTestData(OutputPath)

		rec4testscaler := VHDL.New2DScaler(rec4test, N)
		rec4testscaler.GenerateVHDL(OutputPath)

		test := Viv.CreateXSIM(OutputPath, "SimBenchMarkRec4", rec4test.GenerateVHDLEntityArray())
		test.Exec()
		err := Viv.ParseXSIMReport(OutputPath, rec4test)
		if err != nil {
			b.Fail()
		}
		synthplace := Viv.CreateVivadoTCL(OutputPath, "SimBenchMarkRec4", rec4testscaler, VivadoSettings)
		synthplace.Exec()
		util := Viv.ParseUtilizationReport(OutputPath, rec4testscaler)
		if util.TotalLUT == 0 {
			b.Fail()
		}
	}
}

func Benchmark2bit_1(b *testing.B) {
	startBench()
	for i := 0; i < b.N; i++ {
		M1.GenerateVHDL(OutputPath)
		M1.GenerateTestData(OutputPath)
		test := Viv.CreateXSIM(OutputPath, "SimBenchMarkM1", M1.GenerateVHDLEntityArray())
		test.Exec()
		err := Viv.ParseXSIMReport(OutputPath, M1)
		if err != nil {
			b.Fail()
		}
		synthplace := Viv.CreateVivadoTCL(OutputPath, "ExecBenchM1", M1, VivadoSettings)
		synthplace.Exec()
		util := Viv.ParseUtilizationReport(OutputPath, M1)
		if util.TotalLUT == 0 {
			b.Fail()
		}
	}
}

func Benchmark2bit_10(b *testing.B) {
	bit2MultN(10, b)
}

func Benchmark2bit_100(b *testing.B) {
	bit2MultN(100, b)
}

func Benchmark4bit_1(b *testing.B) {
	startBench()
	for i := 0; i < b.N; i++ {
		rec4test := VHDL.NewRecursive4("rec4bench", [4]VHDL.VHDLEntityMultiplier{M1, M1, M1, M1})
		rec4test.GenerateVHDL(OutputPath)
		rec4test.GenerateTestData(OutputPath)

		test := Viv.CreateXSIM(OutputPath, "SimBenchMarkRec4", rec4test.GenerateVHDLEntityArray())
		test.Exec()
		err := Viv.ParseXSIMReport(OutputPath, rec4test)
		if err != nil {
			b.Fail()
		}
		synthplace := Viv.CreateVivadoTCL(OutputPath, "SimBenchMarkRec4", rec4test, VivadoSettings)
		synthplace.Exec()
		util := Viv.ParseUtilizationReport(OutputPath, rec4test)
		if util.TotalLUT == 0 {
			b.Fail()
		}
	}
}

func Benchmark4bit_10(b *testing.B) {
	bit4MultN(10, b)
}

func Benchmark4bit_100(b *testing.B) {
	bit4MultN(100, b)
}

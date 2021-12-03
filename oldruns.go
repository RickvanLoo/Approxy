package main

import (
	VHDL "badmath/VHDL"
	Viv "badmath/Vivado"
)

func Accurate() {
	acc8 := VHDL.NewAccurateNumMultiplyer("Acc8", 8, OutputPath)
	acc8.GenerateVHDL(OutputPath)
	acc8.GenerateTestData(OutputPath)
	// xsim := Viv.CreateXSIM(OutputPath, "Acc8.vhd", "testAcc8.txt", "topsim.vhd", acc8.EntityName, acc8.BitSize)
	// xsim.Exec()
	tcl := Viv.CreateVivadoTCL(OutputPath, "main1.tcl", acc8, VivadoSettings)
	tcl.Exec()
}

func ScaleM1() {
	M1 := VHDL.M1()
	M1.LUT2D.Print()
	M1.LUT2D.GenerateVHDL(OutputPath)
	Scaler := VHDL.New2DScaler(M1.LUT2D, 100)
	tcl := Viv.CreateVivadoTCL(OutputPath, "main.tcl", Scaler, VivadoSettings)
	tcl.Exec()
}

func M1M2M3M4() {
	M1 := VHDL.M1()
	M1.LUT2D.Print()
	M2 := VHDL.M2()
	M2.LUT2D.Print()
	M3 := VHDL.M3()
	M3.LUT2D.Print()
	M4 := VHDL.M4()
	M4.LUT2D.Print()

	M1.LUT2D.GenerateVHDL(OutputPath)
	M2.LUT2D.GenerateVHDL(OutputPath)
	M3.LUT2D.GenerateVHDL(OutputPath)
	M4.LUT2D.GenerateVHDL(OutputPath)

	M1.LUT2D.GenerateTestData(OutputPath)
	M2.LUT2D.GenerateTestData(OutputPath)
	M3.LUT2D.GenerateTestData(OutputPath)
	M4.LUT2D.GenerateTestData(OutputPath)

	// XSIM1 := Viv.CreateXSIM(OutputPath, "m1.vhd", "testb1.txt", "topsim1.vhd", M1.LUT2D.EntityName, M1.LUT2D.BitSize)
	// XSIM2 := Viv.CreateXSIM(OutputPath, "m2.vhd", "testb2.txt", "topsim2.vhd", M2.LUT2D.EntityName, M2.LUT2D.BitSize)
	// XSIM3 := Viv.CreateXSIM(OutputPath, "m3.vhd", "testb3.txt", "topsim3.vhd", M3.LUT2D.EntityName, M3.LUT2D.BitSize)
	// XSIM4 := Viv.CreateXSIM(OutputPath, "m4.vhd", "testb4.txt", "topsim4.vhd", M4.LUT2D.EntityName, M4.LUT2D.BitSize)

	// XSIM1.Exec()
	// XSIM2.Exec()
	// XSIM3.Exec()
	// XSIM4.Exec()

	tcl1 := Viv.CreateVivadoTCL(OutputPath, "main1.tcl", M1.LUT2D, VivadoSettings)
	tcl2 := Viv.CreateVivadoTCL(OutputPath, "main2.tcl", M2.LUT2D, VivadoSettings)
	tcl3 := Viv.CreateVivadoTCL(OutputPath, "main3.tcl", M3.LUT2D, VivadoSettings)
	tcl4 := Viv.CreateVivadoTCL(OutputPath, "main4.tcl", M4.LUT2D, VivadoSettings)

	tcl1.Exec()
	tcl2.Exec()
	tcl3.Exec()
	tcl4.Exec()

}

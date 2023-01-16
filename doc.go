/*
#vhdl
approxy/vhdl module encapsulates the design of multipliers for digital hardware

Most .go files implement a single model centered around a VHDL template, which leaves behavioural analysis to Approxy.
The Approxy models export synthesizable VHDL code and testdata for futher analysis in tooling such as Xilinx Vivado

All models either implement VHDLEntityMultiplier or VHDLEntity only.

type VHDLEntityMultiplier interface {
	VHDLEntity
	Multiplier
}

type VHDLEntity interface {
	ReturnData() *EntityData
	GenerateVHDL(string)
	GenerateTestData(string)
	GenerateVHDLEntityArray() []VHDLEntity
	String() string //MSB -> LSB
}

type Multiplier interface {
	ReturnVal(uint, uint) uint
	Overflow() bool
	MeanAbsoluteError() float64
}

#File Overview:

File: Accurate.go
Struct: UnsignedNumericAccurateMultiplyer
Template: accuratebehav.vhd
Description: A variable sized accurate multiplier on basis of the IEEE Numeric library, gives full freedom to synthesis tool.

File: AccurateMac.go
Struct: UnsignedNumericAccurateMAC
Template: macaccurate.vhd
Description: A variable sized accurate MAC on basis of the IEEE Numeric library, gives full freedom to synthesis tool.

File: Approx.go
Struct: UnsignedApproxMultiplyer
Template: multbehav.vhd
Description: An extension of the LUT2D model. Implements an accurate variable sized multiplier, using a nested case structure
Adds the possibility of adding modifications to make the multiplier approximate.

File: Correlator.go
Struct: Correlator
Template: corrbehav.vhd
Description: Not finished, implements a correlator on basis of a [4]VHDLEntity, needs to be MAC

File: doc.go
Struct: --
Template: --
Description: This documentation file

File: Examples.go
Struct: UnsignedApproxMultiplyer
Template: --
Description: Implements M1, M2, M3 and M4 within functions using "[][]uint" from Approx.go

File: External.go
Struct: ExternalMult
Template: Any within templates/external
Description: Uses XSIM simulation to retrieve multiplier functionality from VHDL file, encapsulates it within an VHDLEntityMultiplier like other models.

File: Generic.
Struct: --
Template: --
Description: Approxy Interfaces, helper functions

File: LUT2D.go
Struct: LUT2D
Template: multbehav.go
Description: Implements a multiplier model around [][]uint, implements within VHDL using nested case structure. Does not implement behaviour: Use Approx.go/UnsignedApproxMultiplyer

File: Mac.go
Struct: MAC
Template: macbehav.vhd
Description: Creates an MAC, using another Approxy model/VHDLEntityMultiplier as it's multiplier basis. Addition is accurate.

File: Recursive4.go
Struct: Recursive4
Template: rec4behav.vhd
Description: Creates a 4-bit Recursive Multiplier using an array of four 2-bit VHDLEntityMultipliers

File: Recursive8.go
Struct: Recursive8
Template: rec8behav.vhd
Description: Creates a 8-bit Recursive Multiplier using an array of four 4-bit VHDLEntityMultipliers

File: Scaler.go
Struct: Scaler
Template: scaler.vhd
Description: Creates a scaler top-level entity, implements N amount of VHDLEntityMultipliers

#vivado
approxy/vivado module encapsulates all automated interaction with the Vivado and XSIM tooling.

This package automates the Vivado workflow for models within Approxy/vhdl

Requirements:
- Vivado within $PATH (recommended: Vivado v.2021.1)
- XSIM within $PATH

Implements:
1. Verification using XSIM and Approxy/vhdl testdata
2. Generation of Vivado TCL scripts and execution (Synthesis->Place->Route->Report Generation)
3. (Limited) parsing of Vivado reports
4. Post-PR simulation for power estimation

#Overview

File: ParsePower.go
Description: Functions for parsing XML Vivado Power Reports

File: ParseTiming.go
Description: Functions for parsing plaintext Vivado Timing Reports

File: ParseUtilization.go
Description: Functions for parsing plaintext Vivado Utilization Reports

File: Reporting.go
Description: Functionality to combine parsed reports within a Report struct, multiple reports can be combined within a Run to automate the analysis of various designs.

File: TCL.go
Template: vivado.tcl and reportpower.tcl
Description: Generation of Vivado TCL scripts for synthesis/P+R and reporting power, support for Vivado settings/flags

File: XSIM.go
Template: xsim_mult_scalder.vhd, xsim_mult.vhd, xsim_reverse.vhd, xsim_seq_scaler.vhd, xsim_seq.vhd
Description: Automated behavioural verification of approxy/vhdl models, execution of functional simulation for SAIF generation

File: doc.go
Description: This documentation file.
*/

package main

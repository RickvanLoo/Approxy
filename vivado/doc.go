/*
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

package vivado

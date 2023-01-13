# Approxy
Golang tool generating VHDL multipliers using approximate computing. 

![Go Build](https://github.com/RickvanLoo/approxy/actions/workflows/go.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/RickvanLoo/Approxy?style=flat-square)](https://goreportcard.com/report/github.com/RickvanLoo/Approxy)
[![Go Reference](https://pkg.go.dev/badge/github.com/RickvanLoo/Approxy.svg)](https://pkg.go.dev/github.com/RickvanLoo/Approxy)
[![Coverage Status](https://coveralls.io/repos/github/RickvanLoo/Approxy/badge.svg)](https://coveralls.io/github/RickvanLoo/Approxy)
![Approxy Workflow](approxy.png)

## Requirements
- Go 1.17
- "vivado", "xvhdl", "xelab" and "xsim" within $PATH
- Tested for:
    - Vivado v2021.1
    - Linux (not tested for Windows, it probably work)
    - Requires "vivado", "xvhdl", "xelab" and "xsim" within $PATH


## How to use
- Clone repository to local computer
- Edit main.go
- Build by running in main folder:
- >go install approxy
- Run by using:
- >approxy

### By default Approxy will make two folders:
- OutputPath or /output, is used as a temp folder to output VHDL and TCL files, and is used to output Vivado logs
- ReportPath or /report, is used to output JSON reports with final result data

## :warning: Warning: 
By default Approxy clears the OutputPath folder before executing the main() function. Make sure to back up useful files (VHDL, TCL logs, Vivado reports, etc.) from the previous execution, before running 'Approxy' again. Clearing functionality can naturally be removed in the init() function, but might interfere with Vivado synthesis.

## Recreating VHDL models
If simply wanting to recreate the VHDL multipliers as discussed in the thesis the following code example is enough, and will output a synthesizable VHDL file and file containing all input/output combinations. 

E.g. a Recursive 4-bit multiplier using M1, M2, M3, and M4 (Within thesis described as R<sub>1234</sub>)

    Rec_1234 := VHDL.NewRecursive4("Rec1234", [4]VHDL.VHDLEntityMultiplier{M1, M2, M3, M4})
    Rec_1234.GenerateTestData(OutputPath)
    Rec_1234.GenerateVHDL(OutputPath)


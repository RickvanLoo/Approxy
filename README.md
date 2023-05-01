# Approxy
Golang tool generating VHDL multipliers using approximate computing. 

This tooling can be used to investigate various approximate computing design for (Xilinx) FPGAs, by generating VHDL and automating its analysis. 
The final result is the best found 4-bit multiplier, using a set of referenced techniques for my Master's Thesis, but the tool is not limited to this. 

![Go Build](https://github.com/RickvanLoo/approxy/actions/workflows/go.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/RickvanLoo/Approxy?style=flat-square)](https://goreportcard.com/report/github.com/RickvanLoo/Approxy)
[![Go Reference](https://pkg.go.dev/badge/github.com/RickvanLoo/Approxy.svg)](https://pkg.go.dev/github.com/RickvanLoo/Approxy)
[![Coverage Status](https://coveralls.io/repos/github/RickvanLoo/Approxy/badge.svg)](https://coveralls.io/github/RickvanLoo/Approxy)
![Approxy Workflow](approxy.png)

View the accompanied MSc Thesis here: http://essay.utwente.nl/93756/1/van%20Loo_MA_EEMSC.pdf

## Final Result
- R1311: 4-bit, optimized for power consumption & MAE for normally distributed inputs
- LUT/multiplier: 22.792
- Power/multiplier: 129uW
- MaxFreq: 817.322MHz
- CARRY4: 0
- Mean Absolute Error (Uniform Distributed Input): 0.01104736328125
- Error Rate (Uniform): 49
- Error Magnitude (Uniform): 40
- Mean Absolute Error (Normally Distributed Input): 0.000515625

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
- >go build
- Run by using:
- >./Approxy

### By default Approxy will make two folders:
- OutputPath or /output, is used as a temp folder to output VHDL and TCL files, and is used to output Vivado logs
- ReportPath or /report, is used to output JSON reports with final result data

## :warning: Warning: 
By default Approxy clears the OutputPath folder before executing the main() function. Make sure to back up useful files (VHDL, TCL logs, Vivado reports, etc.) from the previous execution, before running 'Approxy' again. Clearing functionality can naturally be removed in the init() function, but might interfere with Vivado synthesis.

## Recreating VHDL models
Read Approxy Github Wiki: https://github.com/RickvanLoo/Approxy/wiki

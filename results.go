package main

import (
	"badmath/VHDL"
	"badmath/Vivado"
	"fmt"
	"os"
	"text/tabwriter"
)

type Result struct {
	EntityString string
	Utilization  *Vivado.Utilization
	Overflow     bool
	MAE          float64
}

//TODO: Redo VHDL.VHDLEntity.String(), make it part of UnsignedMultiplication!!!
func NewResult(Entity VHDL.VHDLEntity, Utilization *Vivado.Utilization, Overflow bool, MAE float64) *Result {
	R := new(Result)
	R.EntityString = Entity.String()
	R.Utilization = Utilization
	R.Overflow = Overflow
	R.MAE = MAE
	return R
}

func (r *Result) PrettyPrint() {
	fmt.Printf("\n")

	var topbottomstr string
	for i := 0; i < 82; i++ {
		topbottomstr += "-"
	}
	writer := tabwriter.NewWriter(os.Stdout, 0, 10, 1, '\t', tabwriter.AlignRight)
	fmt.Fprintln(writer, topbottomstr)
	fmt.Fprintln(writer, "BADMATH - RESULT REPORT")
	fmt.Fprintln(writer, topbottomstr)
	fmt.Fprintln(writer, "TopEntity:\t"+r.EntityString)
	fmt.Fprintf(writer, "Overflow: %t\tMAE: %f\n", r.Overflow, r.MAE)
	fmt.Fprintf(writer, "Total LUTs\tLogic LUTs\tLUTRAMs\tSRLs\tFFs\tRAMB36\tRAMB18\tDSP Blocks\n")
	fmt.Fprintf(writer, "%d\t", r.Utilization.TotalLUT)
	fmt.Fprintf(writer, "%d\t", r.Utilization.LogicLUT)
	fmt.Fprintf(writer, "%d\t", r.Utilization.LUTRAMs)
	fmt.Fprintf(writer, "%d\t", r.Utilization.SRLs)
	fmt.Fprintf(writer, "%d\t", r.Utilization.FFs)
	fmt.Fprintf(writer, "%d\t", r.Utilization.RAMB36)
	fmt.Fprintf(writer, "%d\t", r.Utilization.RAMB18)
	fmt.Fprintf(writer, "%d\n", r.Utilization.DSP)
	fmt.Fprintln(writer, topbottomstr)
	writer.Flush()
}

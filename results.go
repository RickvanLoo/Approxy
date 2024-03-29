//go:build exclude

//THIS FILE IS NOT BEING BUILD

package main

//This file is depricated, please use reporting/result functionality within package Approxy/vivado

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	VHDL "github.com/RickvanLoo/Approxy/vhdl"
	"github.com/RickvanLoo/Approxy/vivado"
)

type Result struct {
	EntityString string
	EntityHash   string
	Utilization  *vivado.Utilization
	Overflow     bool
	MAE          float64
}

// TODO: Redo VHDL.VHDLEntity.String(), make it part of UnsignedMultiplication!!!
func NewResult(Entity VHDL.VHDLEntity, Utilization *vivado.Utilization, Overflow bool, MAE float64) *Result {
	R := new(Result)
	R.EntityString = Entity.String()
	R.EntityHash = toHash(R.EntityString)
	R.Utilization = Utilization
	R.Overflow = Overflow
	R.MAE = MAE
	return R
}

// Only for identifier purposes, do not use this for crypto purposes
func toHash(EntityString string) string {
	hash := md5.Sum([]byte(EntityString))
	return hex.EncodeToString(hash[:])
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
	fmt.Fprintln(writer, "MD5 Hash:\t"+r.EntityHash)
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

func (r *Result) String() string {
	return r.EntityHash + "," + r.EntityString + "," + strconv.Itoa(r.Utilization.TotalLUT) + "," + fmt.Sprintf("%f", r.MAE) + "," + fmt.Sprintf("%t", r.Overflow)
}

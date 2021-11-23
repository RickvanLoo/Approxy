package main

import (
	"fmt"
	"log"
	"os"

	"github.com/RickvanLoo/badmath/mult"
)

var OutputPath string

func main() {
	ClearLogs()
	fmt.Println("Hello, world.")

	OutputPath = "output"
	err := os.RemoveAll(OutputPath)
	if err != nil {
		log.Println(err)
	}

	M1 := mult.M1()
	M1.Mult.Print()
	err = os.MkdirAll(OutputPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	M1.Mult.VHDLtoFile(OutputPath, "m1.vhd")

	CreateVivadoTCL(OutputPath, "main.tcl", M1.Mult.Name)
	ExecuteVivadoTCL(OutputPath, "main.tcl")
}

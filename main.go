package main

import (
	"fmt"

	"github.com/RickvanLoo/badmath/mult"
)

func main() {
	fmt.Println("Hello, world.")

	mult := mult.New_uAM(2)
	mult.print()
	mult.VHDLtoFile("output.vhd")
}

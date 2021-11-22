package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, world.")

	mult := New_uAM(8)
	mult.print()
	mult.VHDLtoFile("output.vhd")
}

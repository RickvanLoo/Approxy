package mult

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"text/template"
)

type UnsignedAccurateMultiplyer struct {
	Name    string
	Bitsize uint
	Osize   uint
	LUT     [][]uint
}

func NewUnsignedAcc(size uint) *UnsignedAccurateMultiplyer {
	m := new(UnsignedAccurateMultiplyer)
	m.Bitsize = size
	m.Name = "uAcc" + strconv.Itoa(int(size)) + "bitMult"
	m.Osize = 2 * m.Bitsize
	LUTSize := int(math.Pow(2, float64(size)))

	var LUT [][]uint
	for x := 0; x < LUTSize; x++ {
		var row []uint
		for y := 0; y < LUTSize; y++ {
			row = append(row, uint(x*y))
		}
		LUT = append(LUT, row)
	}

	m.LUT = LUT

	return m
}

func (m *UnsignedAccurateMultiplyer) Print() {
	fmt.Printf("%+v\n", m)
}

//returnVal, returns output for two inputs
func (m *UnsignedAccurateMultiplyer) ReturnVal(x, y uint) uint {
	return m.LUT[x][y]
}

func (m *UnsignedAccurateMultiplyer) changeVal(x, y, val uint) {
	m.LUT[x][y] = val
}

//convertindex converts integer to binary, adds trailing zeroes for input-vectors, used for VHDL Template
func (m *UnsignedAccurateMultiplyer) convertindex(value interface{}) string {
	var format string
	format = "%0" + strconv.Itoa(int(m.Bitsize)) + "b"
	return fmt.Sprintf(format, value)
}

//convertval converts integer to binary, adds trailing zeroes for output-vectors, used for VHDL Template
func (m *UnsignedAccurateMultiplyer) convertval(value interface{}) string {
	var format string
	format = "%0" + strconv.Itoa(int(m.Osize)) + "b"
	return fmt.Sprintf(format, value)
}

//Outputs VHDL code for generated multiplier to path file
func (m *UnsignedAccurateMultiplyer) VHDLtoFile(folder string, name string) {
	funcMap := template.FuncMap{"indexconv": m.convertindex, "valconv": m.convertval}
	templatepath := "template/multbehav.vhd"
	templatename := "multbehav.vhd"

	path := folder + "/" + name

	t, err := template.New(templatename).Funcs(funcMap).ParseFiles(templatepath)
	if err != nil {
		log.Print(err)
		return
	}

	f, err := os.Create(path)

	if err != nil {
		log.Println("create file: ", err)
		return
	}

	err = t.ExecuteTemplate(f, templatename, m)
	if err != nil {
		log.Print("execute: ", err)
		return
	}
}

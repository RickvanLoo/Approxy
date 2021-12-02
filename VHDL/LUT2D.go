package VHDL

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"text/template"
)

type LUT2D struct {
	EntityName string
	BitSize    uint
	OutputSize uint
	LUT        [][]uint
	VHDLFile   string
	TestFile   string
}

func (m *LUT2D) Print() {
	fmt.Printf("%+v\n", m)
}

//returnVal, returns output for two inputs
func (m *LUT2D) ReturnVal(a, b uint) uint {
	return m.LUT[a][b]
}

func (m *LUT2D) changeVal(a, b, prod uint) {
	m.LUT[a][b] = prod
}

//convertindex converts integer to binary, adds trailing zeroes for input-vectors, used for VHDL Template
func (m *LUT2D) convertindex(value interface{}) string {
	format := "%0" + strconv.Itoa(int(m.BitSize)) + "b"
	return fmt.Sprintf(format, value)
}

//convertval converts integer to binary, adds trailing zeroes for output-vectors, used for VHDL Template
func (m *LUT2D) convertval(value interface{}) string {
	format := "%0" + strconv.Itoa(int(m.OutputSize)) + "b"
	return fmt.Sprintf(format, value)
}

func New2DLUT(EntityName string, BitSize uint) *LUT2D {
	m := new(LUT2D)
	m.BitSize = BitSize
	m.EntityName = EntityName
	m.OutputSize = 2 * m.BitSize
	m.VHDLFile, m.TestFile = FileNameGen(m.EntityName)
	return m
}

//Outputs VHDL code for generated multiplier to path file
//TODO: Make VHDLtoFile call the Generic function
func (m *LUT2D) VHDLtoFile(FolderPath string) {
	funcMap := template.FuncMap{"indexconv": m.convertindex, "valconv": m.convertval}
	templatepath := "template/multbehav.vhd"
	templatename := "multbehav.vhd"

	path := FolderPath + "/" + m.VHDLFile

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

func (m *LUT2D) GenerateTestData(FolderPath string) {
	fmtstr := "%0" + strconv.Itoa(int(m.BitSize)) + "b %0" + strconv.Itoa(int(m.BitSize)) + "b %0" + strconv.Itoa(int(m.OutputSize)) + "b\n"
	path := FolderPath + "/" + m.TestFile

	file, err := os.Create(path)
	if err != nil {
		log.Println(err)
	}

	writer := bufio.NewWriter(file)

	for x, row := range m.LUT {
		for y, z := range row {
			//Remove newline from last output to not mess up XSIM
			if (x == (len(m.LUT) - 1)) && (y == (len(row) - 1)) {
				fmtstr = strings.TrimSuffix(fmtstr, "\n")
			}

			_, err = fmt.Fprintf(writer, fmtstr, x, y, z)
			if err != nil {
				log.Println(err)
			}
		}
	}

	writer.Flush()

}

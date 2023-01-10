package vhdl

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"text/template"
)

// LUT2D defines a 2D LUT structure.
// Interface == VHDLEntity
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

// returnVal, returns output for two inputs
func (m *LUT2D) ReturnVal(a, b uint) uint {
	return m.LUT[a][b]
}

func (m *LUT2D) changeVal(a, b, prod uint) {
	m.LUT[a][b] = prod
}

// convertindex converts integer to binary, adds trailing zeroes for input-vectors, used for VHDL Template
func (m *LUT2D) convertindex(value interface{}) string {
	format := "%0" + strconv.Itoa(int(m.BitSize)) + "b"
	return fmt.Sprintf(format, value)
}

// convertval converts integer to binary, adds trailing zeroes for output-vectors, used for VHDL Template
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

// Outputs VHDL code for generated multiplier to path file
// TODO: Make VHDLtoFile call the Generic function
func (m *LUT2D) GenerateVHDL(FolderPath string) {
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

func (m *LUT2D) ReturnData() *EntityData {
	// EntityName string
	// BitSize    uint
	// VHDLFile   string
	// TestFile   string
	d := new(EntityData)
	d.EntityName = m.EntityName
	d.BitSize = m.BitSize
	d.OutputSize = m.OutputSize
	d.VHDLFile = m.VHDLFile
	d.TestFile = m.TestFile
	return d
}

func (m *LUT2D) Overflow() bool {
	maxval := int(math.Exp2(float64(m.BitSize)))
	overflowval := int(math.Exp2(float64(m.OutputSize)))

	for a := 0; a < maxval; a++ {
		for b := 0; b < maxval; b++ {
			if m.ReturnVal(uint(a), uint(b)) > uint(overflowval)-1 { //double check this with _test.go
				log.Printf("WARNING: Overflow for %d*%d=%d>%d\n", a, b, m.ReturnVal(uint(a), uint(b)), overflowval-1)
				log.Printf("Accurate: %d*%d=%d\n", a, b, a*b)
				return true
			}
		}
	}

	return false
}

func (m *LUT2D) MeanAbsoluteError() float64 {
	maxval := int(math.Exp2(float64(m.BitSize)))
	accum := float64(0)
	for a := 1; a < maxval; a++ {
		for b := 1; b < maxval; b++ {
			accResult := float64(a * b)
			r4Result := m.ReturnVal(uint(a), uint(b))
			accum += math.Abs(float64(r4Result) - accResult)
		}
	}

	return float64(1.0/math.Exp2(float64(m.OutputSize))) * accum

}

func (m *LUT2D) String() string {
	return m.EntityName
}

func (m *LUT2D) GenerateVHDLEntityArray() []VHDLEntity {
	var out []VHDLEntity
	out = append(out, m)
	return out
}

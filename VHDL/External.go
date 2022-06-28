package VHDL

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type ExternalMult struct {
	EntityName   string
	BitSize      uint
	OutputSize   uint
	VHDLFile     string
	TestFile     string
	TemplateFile string
	Behaviour    [][]uint
}

func NewExternalMult(EntityName string, BitSize uint, Template string) *ExternalMult {
	ext := new(ExternalMult)
	ext.EntityName = EntityName
	ext.BitSize = BitSize
	ext.OutputSize = 2 * ext.BitSize
	ext.TemplateFile = Template
	ext.VHDLFile, ext.TestFile = FileNameGen(ext.EntityName)
	maxval := int(math.Exp2(float64(ext.BitSize)))

	ext.Behaviour = make([][]uint, maxval)
	for i := range ext.Behaviour {
		ext.Behaviour[i] = make([]uint, maxval)
	}
	return ext
}

// ReturnData() *EntityData
// GenerateVHDL(string)
// GenerateTestData(string)
// GenerateVHDLEntityArray() []VHDLEntity
// String() string //MSB -> LSB
func (ext *ExternalMult) ReturnData() *EntityData {
	d := new(EntityData)
	d.BitSize = ext.BitSize
	d.EntityName = ext.EntityName
	d.OutputSize = ext.OutputSize
	d.TestFile = ext.TestFile
	d.VHDLFile = ext.VHDLFile
	return d
}

func (ext *ExternalMult) GenerateVHDL(FolderPath string) {
	CreateFile(FolderPath, ext.VHDLFile, "external/"+ext.TemplateFile, ext)
}

func (ext *ExternalMult) GenerateTestData(FolderPath string) {
	fmtstr := "%0" + strconv.Itoa(int(ext.BitSize)) + "b %0" + strconv.Itoa(int(ext.BitSize)) + "b\n"
	path := FolderPath + "/" + ext.TestFile

	file, err := os.Create(path)
	if err != nil {
		log.Println(err)
	}

	writer := bufio.NewWriter(file)

	maxval := int(math.Exp2(float64(ext.BitSize)))

	for a := 0; a < maxval; a++ {
		for b := 0; b < maxval; b++ {

			if (a == (maxval - 1)) && (b == (maxval - 1)) {
				fmtstr = strings.TrimSuffix(fmtstr, "\n")
			}

			_, err = fmt.Fprintf(writer, fmtstr, a, b)
			if err != nil {
				log.Println(err)
			}

		}
	}

	writer.Flush()
}

func (ext *ExternalMult) GenerateVHDLEntityArray() []VHDLEntity {
	var out []VHDLEntity
	out = append(out, ext)
	return out
}

func (ext *ExternalMult) String() string {
	return ext.EntityName
}

func (ext *ExternalMult) ParseXSIMOutput(FolderPath string) {
	file, err := os.Open(FolderPath + "/" + "out_" + ext.TestFile)
	if err != nil {
		log.Println(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text []string

	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	file.Close()

	maxval := int(math.Exp2(float64(ext.BitSize)))

	textindex := 0

	for a := 0; a < maxval; a++ {
		for b := 0; b < maxval; b++ {

			uintout, _ := strconv.ParseUint(text[textindex], 10, 64)
			ext.Behaviour[a][b] = uint(uintout)
			textindex = textindex + 1
		}
	}
}

func (ext *ExternalMult) ReturnVal(a, b uint) uint {
	return ext.Behaviour[a][b]
}

func (ext *ExternalMult) MeanAbsoluteError() float64 {
	maxval := int(math.Exp2(float64(ext.BitSize)))
	accum := float64(0)
	for a := 1; a < maxval; a++ {
		for b := 1; b < maxval; b++ {
			accResult := float64(a * b)
			r4Result := ext.ReturnVal(uint(a), uint(b))
			accum += math.Abs(float64(r4Result) - accResult)
		}
	}

	return float64(1.0/math.Exp2(float64(ext.OutputSize))) * accum
}

func (ext *ExternalMult) AverageRelativeError() float64 {
	maxval := int(math.Exp2(float64(ext.BitSize)))
	accum := float64(0)
	for a := 1; a < maxval; a++ {
		for b := 1; b < maxval; b++ {
			accResult := float64(a * b)
			r4Result := ext.ReturnVal(uint(a), uint(b))
			accum += math.Abs((float64(r4Result) - accResult) / accResult)
		}
	}

	return float64(1.0/math.Exp2(float64(ext.OutputSize))) * accum

}
func (ext *ExternalMult) Overflow() bool {
	log.Println("Overflow not supported for External Multipliers")
	return false
}

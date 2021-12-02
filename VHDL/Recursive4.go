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

//LUTArray[0] = AH*BH
//LUTArray[1] = AH*BL
//LUTArray[2] = AL*BH
//LUTArray[3] = AL*BL
type Recursive4 struct {
	EntityName string
	BitSize    uint      //Default to 4
	OutputSize uint      //Default to 8
	LUTArray   [4]*LUT2D //Size of 4
	VHDLFile   string
	TestFile   string
}

func NewRecursive4(EntityName string, LUTArray [4]*LUT2D) *Recursive4 {
	r4 := new(Recursive4)
	r4.BitSize = 4
	r4.OutputSize = 8
	r4.EntityName = EntityName
	r4.VHDLFile, r4.TestFile = FileNameGen(r4.EntityName)
	r4.LUTArray = LUTArray

	return r4
}

func (r4 *Recursive4) ReturnVal(a uint, b uint) uint {
	AHBH_LUT := r4.LUTArray[0]
	AHBL_LUT := r4.LUTArray[1]
	ALBH_LUT := r4.LUTArray[2]
	ALBL_LUT := r4.LUTArray[3]

	bin_input := make([]byte, 2)
	bin_input[0] = byte(a)
	bin_input[1] = byte(b)

	maskH := byte(0b00001100)
	maskL := byte(0b00000011)

	AHALBHBL := make([]byte, 4)
	AHALBHBL[0] = (bin_input[0] & maskH) >> 2 //AH
	AHALBHBL[1] = bin_input[0] & maskL        //AL
	AHALBHBL[2] = (bin_input[1] & maskH) >> 2 //BH
	AHALBHBL[3] = bin_input[1] & maskL        //BL

	AH := uint(AHALBHBL[0])
	AL := uint(AHALBHBL[1])
	BH := uint(AHALBHBL[2])
	BL := uint(AHALBHBL[3])

	AHBH := AHBH_LUT.ReturnVal(AH, BH)
	AHBL := AHBL_LUT.ReturnVal(AH, BL)
	ALBH := ALBH_LUT.ReturnVal(AL, BH)
	ALBL := ALBL_LUT.ReturnVal(AL, BL)

	output := ALBL + (ALBH << 2) + (AHBL << 2) + (AHBH << 4)

	return output
}

func (r4 *Recursive4) GenerateTestData(FolderPath string) {
	fmtstr := "%0" + strconv.Itoa(int(r4.BitSize)) + "b %0" + strconv.Itoa(int(r4.BitSize)) + "b %0" + strconv.Itoa(int(r4.OutputSize)) + "b\n"
	path := FolderPath + "/" + r4.TestFile

	file, err := os.Create(path)
	if err != nil {
		log.Println(err)
	}

	writer := bufio.NewWriter(file)

	maxval := int(math.Exp2(4))

	for a := 0; a < maxval; a++ {
		for b := 0; b < maxval; b++ {

			if (a == 15) && (b == 15) {
				fmtstr = strings.TrimSuffix(fmtstr, "\n")
			}

			out := r4.ReturnVal(uint(a), uint(b))

			_, err = fmt.Fprintf(writer, fmtstr, a, b, out)
			if err != nil {
				log.Println(err)
			}

		}
	}

	writer.Flush()

}

func (r4 *Recursive4) GenerateVHDL(FolderPath string) {
	for _, mult := range r4.LUTArray {
		mult.GenerateVHDL(FolderPath)
	}

	CreateFile(FolderPath, r4.VHDLFile, "rec4behav.vhd", r4)
}

func (r4 *Recursive4) ReturnData() *EntityData {
	// EntityName string
	// BitSize    uint
	// VHDLFile   string
	// TestFile   string
	d := new(EntityData)
	d.EntityName = r4.EntityName
	d.BitSize = r4.BitSize
	d.VHDLFile = r4.VHDLFile
	d.TestFile = r4.TestFile
	return d
}

func (r4 *Recursive4) GenerateVHDLEntityArray() []VHDLEntity {

	var out []VHDLEntity

	out = append(out, r4)

	for _, mult := range r4.LUTArray {
		out = append(out, mult)
	}

	return out
}

func (r4 *Recursive4) Overflow() bool {
	maxval := int(math.Exp2(4))
	overflowval := int(math.Exp2(8))

	for a := 0; a < maxval; a++ {
		for b := 0; b < maxval; b++ {
			if r4.ReturnVal(uint(a), uint(b)) > uint(overflowval)-1 { //double check this with _test.go
				log.Printf("WARNING: Overflow for %d*%d=%d>%d\n", a, b, r4.ReturnVal(uint(a), uint(b)), overflowval-1)
				log.Printf("Accurate: %d*%d=%d\n", a, b, a*b)
				return true
			}
		}
	}

	return false
}

func (r4 *Recursive4) MeanAbsoluteError() float64 {
	maxval := int(math.Exp2(4))
	accum := float64(0)
	for a := 1; a < maxval; a++ {
		for b := 1; b < maxval; b++ {
			accResult := float64(a * b)
			r4Result := r4.ReturnVal(uint(a), uint(b))
			accum += math.Abs(float64(r4Result) - accResult)
		}
	}

	return float64(1.0/256.0) * accum

}

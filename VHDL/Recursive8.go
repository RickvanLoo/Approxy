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

// type VHDLEntity interface {
// 	ReturnData() *EntityData -- Done
// 	GenerateVHDL(string) -- Done
// 	GenerateTestData(string) -- Done
// 	String() string //MSB -> LSB -- Done
// }

// type UnsignedMultiplyer interface {
// 	ReturnVal(uint, uint) uint  --- Done
// 	Overflow() bool
// 	MeanAbsoluteError() float64
// }

type Recursive8 struct {
	EntityName    string
	BitSize       uint                    //Default to 8
	OutputSize    uint                    //Default to 16
	Rec4Array     [4]VHDLEntityMultiplier //Size of 4
	VHDLFile      string
	TestFile      string
	OverflowError bool
}

func NewRecursive8(EntityName string, Rec4Array [4]VHDLEntityMultiplier) *Recursive8 {
	r8 := new(Recursive8)
	r8.BitSize = 8
	r8.OutputSize = 16
	r8.EntityName = EntityName
	r8.VHDLFile, r8.TestFile = FileNameGen(r8.EntityName)
	r8.Rec4Array = Rec4Array
	r8.OverflowError = false

	for _, mult := range r8.Rec4Array {
		if mult.ReturnData().BitSize != 4 {
			log.Println("Created Recursive8 found LUT where Bitsize is not 4")
		}

		if mult.ReturnData().OutputSize != 8 {
			log.Println("Created Recursive8 found LUT where OutputSize is not 8")
		}
	}

	return r8
}

func (r8 *Recursive8) ReturnVal(a uint, b uint) uint {
	AHBH_LUT := r8.Rec4Array[0]
	AHBL_LUT := r8.Rec4Array[1]
	ALBH_LUT := r8.Rec4Array[2]
	ALBL_LUT := r8.Rec4Array[3]

	bin_input := make([]uint8, 2)
	bin_input[0] = uint8(a)
	bin_input[1] = uint8(b)

	maskH := byte(0b11110000)
	maskL := byte(0b00001111)

	AHALBHBL := make([]uint8, 4)
	AHALBHBL[0] = (bin_input[0] & maskH) >> 4 //AH
	AHALBHBL[1] = bin_input[0] & maskL        //AL
	AHALBHBL[2] = (bin_input[1] & maskH) >> 4 //BH
	AHALBHBL[3] = bin_input[1] & maskL        //BL

	AH := uint(AHALBHBL[0])
	AL := uint(AHALBHBL[1])
	BH := uint(AHALBHBL[2])
	BL := uint(AHALBHBL[3])

	AHBH := AHBH_LUT.ReturnVal(AH, BH)
	AHBL := AHBL_LUT.ReturnVal(AH, BL)
	ALBH := ALBH_LUT.ReturnVal(AL, BH)
	ALBL := ALBL_LUT.ReturnVal(AL, BL)

	output := ALBL + (ALBH << 4) + (AHBL << 4) + (AHBH << 8)
	//Next function masks the output in 8-bit, like VHDL/Vivado would do.
	//Overflow check is best effort, but if generating the whole output-space, we can fully determine overflow
	output, overflowcheck := OverflowCheck16bit(output)
	r8.flagOverflow(overflowcheck)
	return output
}

func (r8 *Recursive8) flagOverflow(input bool) {
	if input {
		r8.OverflowError = true
	}
}

func (r8 *Recursive8) ReturnData() *EntityData {
	d := new(EntityData)
	d.BitSize = r8.BitSize
	d.EntityName = r8.EntityName
	d.OutputSize = r8.OutputSize
	d.TestFile = r8.TestFile
	d.VHDLFile = r8.VHDLFile
	return d
}

func (r8 *Recursive8) GenerateVHDL(FolderPath string) {
	for _, rec4 := range r8.Rec4Array {
		rec4.GenerateVHDL(FolderPath)
	}

	CreateFile(FolderPath, r8.VHDLFile, "rec8behav.vhd", r8)
}

func (r8 *Recursive8) GenerateTestData(FolderPath string) {
	fmtstr := "%0" + strconv.Itoa(int(r8.BitSize)) + "b %0" + strconv.Itoa(int(r8.BitSize)) + "b %0" + strconv.Itoa(int(r8.OutputSize)) + "b\n"
	path := FolderPath + "/" + r8.TestFile

	file, err := os.Create(path)
	if err != nil {
		log.Println(err)
	}

	writer := bufio.NewWriter(file)

	maxval := int(math.Exp2(8))

	for a := 0; a < maxval; a++ {
		for b := 0; b < maxval; b++ {

			if (a == 255) && (b == 255) {
				fmtstr = strings.TrimSuffix(fmtstr, "\n")
			}

			out := r8.ReturnVal(uint(a), uint(b))

			_, err = fmt.Fprintf(writer, fmtstr, a, b, out)
			if err != nil {
				log.Println(err)
			}

		}
	}

	writer.Flush()

}

func (r8 *Recursive8) String() string {
	//AHBH -> AHBL -> ALBH -> ALAL
	str := r8.EntityName + " -> [" + r8.Rec4Array[0].ReturnData().EntityName + ";"
	str += r8.Rec4Array[1].ReturnData().EntityName + ";"
	str += r8.Rec4Array[2].ReturnData().EntityName + ";"
	str += r8.Rec4Array[3].ReturnData().EntityName + "]"
	return str
}

func (r8 *Recursive8) GenerateVHDLEntityArray() []VHDLEntity {

	var out []VHDLEntity
	var LUT2DArray []VHDLEntity
	var REC4Array []VHDLEntity

	for _, rec4 := range r8.Rec4Array {
		REC4Array = append(REC4Array, rec4.GenerateVHDLEntityArray()[0])
		LUT2DArray = append(LUT2DArray, rec4.GenerateVHDLEntityArray()[1:]...)
	}

	REC4Array = RemoveDuplicate(REC4Array)
	LUT2DArray = RemoveDuplicate(LUT2DArray)

	out = append(out, r8)
	out = append(out, REC4Array...)
	out = append(out, LUT2DArray...)

	return out
}

func (r8 *Recursive8) Overflow() bool {
	return r8.OverflowError
}

func (r8 *Recursive8) MeanAbsoluteError() float64 {
	maxval := int(math.Exp2(8))
	accum := float64(0)
	for a := 1; a < maxval; a++ {
		for b := 1; b < maxval; b++ {
			accResult := float64(a * b)
			r4Result := r8.ReturnVal(uint(a), uint(b))
			accum += math.Abs(float64(r4Result) - accResult)
		}
	}

	return float64(1.0/65536.0) * accum
}

package vhdl

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

//TODO: Recursive4 and Recursive8 show identical behaviour for different bitsizes => MERGE

// Recursive4 creates a 4-bit recursive multiplier on basis of four 2-bit VHDLEntityMultipliers, addition is accurate.
type Recursive4 struct {
	EntityName    string
	BitSize       uint                    //Default to 4
	OutputSize    uint                    //Default to 8
	LUTArray      [4]VHDLEntityMultiplier //Size of 4
	VHDLFile      string
	TestFile      string
	OverflowError bool
}

// NewRecursive4 creates a Recursive4 struct
// LUTArray[0] = AH*BH
// LUTArray[1] = AH*BL
// LUTArray[2] = AL*BH
// LUTArray[3] = AL*BL
func NewRecursive4(EntityName string, LUTArray [4]VHDLEntityMultiplier) *Recursive4 {
	r4 := new(Recursive4)
	r4.BitSize = 4
	r4.OutputSize = 8
	r4.EntityName = EntityName
	r4.VHDLFile, r4.TestFile = FileNameGen(r4.EntityName)
	r4.LUTArray = LUTArray

	for _, mult := range r4.LUTArray {
		if mult.ReturnData().BitSize != 2 {
			log.Println("Created Recursive4 found LUT where Bitsize is not 2")
		}

		if mult.ReturnData().OutputSize != 4 {
			log.Println("Created Recursive4 found LUT where OutputSize is not 4")
		}
	}

	r4.OverflowError = false

	return r4
}

// ReturnVal returns the output of the multiplier
func (r4 *Recursive4) ReturnVal(a uint, b uint) uint {
	AHBHLUT := r4.LUTArray[0]
	AHBLLUT := r4.LUTArray[1]
	ALBHLUT := r4.LUTArray[2]
	ALBLLUT := r4.LUTArray[3]

	bininput := make([]byte, 2)
	bininput[0] = byte(a)
	bininput[1] = byte(b)

	maskH := byte(0b00001100)
	maskL := byte(0b00000011)

	AHALBHBL := make([]byte, 4)
	AHALBHBL[0] = (bininput[0] & maskH) >> 2 //AH
	AHALBHBL[1] = bininput[0] & maskL        //AL
	AHALBHBL[2] = (bininput[1] & maskH) >> 2 //BH
	AHALBHBL[3] = bininput[1] & maskL        //BL

	AH := uint(AHALBHBL[0])
	AL := uint(AHALBHBL[1])
	BH := uint(AHALBHBL[2])
	BL := uint(AHALBHBL[3])

	AHBH := AHBHLUT.ReturnVal(AH, BH)
	AHBL := AHBLLUT.ReturnVal(AH, BL)
	ALBH := ALBHLUT.ReturnVal(AL, BH)
	ALBL := ALBLLUT.ReturnVal(AL, BL)

	output := ALBL + (ALBH << 2) + (AHBL << 2) + (AHBH << 4)
	//Next function masks the output in 8-bit, like VHDL/Vivado would do.
	//Overflow check is best effort, but if generating the whole output-space, we can fully determine overflow
	output, overflowcheck := OverflowCheck8bit(output)
	r4.flagOverflow(overflowcheck)
	return output
}

func (r4 *Recursive4) flagOverflow(input bool) {
	if input {
		r4.OverflowError = true
	}
}

// GenerateTestData creates a plaintext testdata file containing both inputs and the output in binary seperated by \t
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

// GenerateVHDL creates the VHDL files in FolderPath
func (r4 *Recursive4) GenerateVHDL(FolderPath string) {
	for _, mult := range r4.LUTArray {
		mult.GenerateVHDL(FolderPath)
	}

	CreateFile(FolderPath, r4.VHDLFile, "rec4behav.vhd", r4)
}

// ReturnData returns a struct containing metadata of the multiplier
func (r4 *Recursive4) ReturnData() *EntityData {
	// EntityName string
	// BitSize    uint
	// VHDLFile   string
	// TestFile   string
	d := new(EntityData)
	d.EntityName = r4.EntityName
	d.BitSize = r4.BitSize
	d.OutputSize = r4.OutputSize
	d.VHDLFile = r4.VHDLFile
	d.TestFile = r4.TestFile
	return d
}

// GenerateVHDLEntityArray creates an array of potentially multiple VHDLEntities, sorted by priority for synthesizing
// For example: Multiplier A uses a VHDL portmap for the smaller Multiplier B & C. B & C need to be synthesized first, hence A will be last in the array
func (r4 *Recursive4) GenerateVHDLEntityArray() []VHDLEntity {

	var out []VHDLEntity

	for _, mult := range r4.LUTArray {
		out = append(out, mult)
	}

	out = RemoveDuplicate(out)

	out = append([]VHDLEntity{r4}, out...)

	return out
}

// Overflow returns a boolean if any internal overflow has occured
func (r4 *Recursive4) Overflow() bool {
	return r4.OverflowError
}

// MeanAbsoluteError returns the MeanAbsoluteError of the multiplier in float64
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

// MeanAbsoluteErrorNormalDist returns the MeanAbsoluteError of the multiplier in float64, using N random normal distributed inputs
func (r4 *Recursive4) MeanAbsoluteErrorNormalDist(N int) float64 {
	accum := float64(0)

	for i := 0; i < N; i++ {
		a := RandomNormalInput(4)
		b := RandomNormalInput(4)
		accResult := float64(a * b)
		r4Result := r4.ReturnVal(uint(a), uint(b))
		accum += math.Abs(float64(r4Result) - accResult)
	}

	return float64(1.0/float64(N)) * accum
}

// AverageRelativeError returns the AverageRelativeError of the multiplier in float64
func (r4 *Recursive4) AverageRelativeError() float64 {
	maxval := int(math.Exp2(4))
	accum := float64(0)
	for a := 1; a < maxval; a++ {
		for b := 1; b < maxval; b++ {
			accResult := float64(a * b)
			r4Result := r4.ReturnVal(uint(a), uint(b))
			accum += math.Abs((float64(r4Result) - accResult) / accResult)
		}
	}

	return float64(1.0/256.0) * accum

}

func (r4 *Recursive4) String() string {
	//AHBH -> AHBL -> ALBH -> ALAL
	str := r4.EntityName + " -> [" + r4.LUTArray[0].ReturnData().EntityName + ";"
	str += r4.LUTArray[1].ReturnData().EntityName + ";"
	str += r4.LUTArray[2].ReturnData().EntityName + ";"
	str += r4.LUTArray[3].ReturnData().EntityName + "]"
	return str
}

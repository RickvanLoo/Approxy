package VHDL

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"math/big"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"text/template"
)

// Generic functions for badmath/VHDL

type VHDLEntityMultiplier interface {
	VHDLEntity
	Multiplier
}

//VHDLEntity describes an interface for testable and synthesizable VHDL structures
type VHDLEntity interface {
	ReturnData() *EntityData
	GenerateVHDL(string)
	GenerateTestData(string)
	GenerateVHDLEntityArray() []VHDLEntity
	String() string //MSB -> LSB
}

//Multiplier describes an interface for generic Multipliers
type Multiplier interface {
	ReturnVal(uint, uint) uint
	Overflow() bool
	MeanAbsoluteError() float64
}

//EntityData encapsulates basic data for a VHDL structure
type EntityData struct {
	EntityName string
	BitSize    uint
	OutputSize uint
	VHDLFile   string
	TestFile   string
}

//https://stackoverflow.com/questions/22713500/iterating-a-range-of-integers-in-go-templates
//If it works, it works.
//N 100 => Iters from 0 to 99(!)
func N(stop uint) (stream chan uint) {
	start := uint(0)
	end := stop - 1
	stream = make(chan uint)
	go func() {
		for i := start; i <= end; i++ {
			stream <- i
		}
		close(stream)
	}()
	return
}

//TODO: Add funcmap support
func CreateFile(FolderPath string, FileName string, TemplateFile string, Data interface{}) {
	TemplatePath := "template/" + TemplateFile

	path := FolderPath + "/" + FileName

	t, err := template.New(TemplateFile).Funcs(template.FuncMap{"N": N}).ParseFiles(TemplatePath)
	if err != nil {
		log.Panic(err)
		return
	}

	f, err := os.Create(path)

	if err != nil {
		log.Panicln("Create File: ", err)
		return
	}

	err = t.ExecuteTemplate(f, TemplateFile, Data)
	if err != nil {
		log.Panicln("Execute: ", err)
		return
	}
}

func FileNameGen(EntityName string) (VHDLFile string, TestFile string) {
	VHDLFile = EntityName + ".vhd"
	TestFile = "test_" + EntityName + ".txt"
	return VHDLFile, TestFile
}

func OverflowCheck8bit(input uint) (output uint, overflow bool) {
	// OverflowMask := byte(0b11111111)
	// byte_input := byte(input)
	// byte_output := byte_input & OverflowMask
	// output = uint(byte_output)
	// booloutput := !(output == input)
	// return output, booloutput
	return OverflowCheckGeneric(input, 8)
}

func OverflowCheck16bit(input uint) (output uint, overflow bool) {
	// OverflowMask := uint16(65535)
	// uint16_input := uint16(input)
	// uint16_output := uint16_input & OverflowMask
	// output = uint(uint16_output)
	// booloutput := !(output == input)
	// return output, booloutput
	return OverflowCheckGeneric(input, 16)
}

func OverflowCheckGeneric(input uint, n uint) (output uint, overflow bool) {
	if n > 64 {
		log.Fatalln("OverflowCheckGeneric not specified for numbers above 64-bit")
	}

	var expo big.Int
	expo.Exp(big.NewInt(2), big.NewInt(int64(n)), big.NewInt(0))
	Maxnumber := expo.Uint64()

	OverflowMask := uint64(Maxnumber - 1)
	uint64_input := uint64(input)
	uint64_output := uint64_input & OverflowMask
	output = uint(uint64_output)
	booloutput := !(output == input)
	return output, booloutput
}

func RemoveDuplicate(Array []VHDLEntity) []VHDLEntity {
	VHDLEntityMap := make(map[string]VHDLEntity)

	for _, Entity := range Array {
		VHDLEntityMap[Entity.ReturnData().EntityName] = Entity
	}

	var v []VHDLEntity

	for _, Entity := range VHDLEntityMap {
		v = append(v, Entity)
	}

	return v
}

//For VHDLEntity recreate N random TestData in form A*B=C for Folderpath
func UniformTestData(Mult VHDLEntityMultiplier, FolderPath string, N uint) {
	BitSize := Mult.ReturnData().BitSize
	OutputSize := Mult.ReturnData().OutputSize
	TestFile := Mult.ReturnData().TestFile

	fmtstr := "%0" + strconv.Itoa(int(BitSize)) + "b %0" + strconv.Itoa(int(BitSize)) + "b %0" + strconv.Itoa(int(OutputSize)) + "b\n"
	path := FolderPath + "/" + TestFile

	file, err := os.Create(path)
	if err != nil {
		log.Println(err)
	}

	writer := bufio.NewWriter(file)

	// maxval := int(math.Exp2(4))

	for i := 0; i < int(N); i++ {
		if i == int((N - 1)) {
			fmtstr = strings.TrimSuffix(fmtstr, "\n")
		}

		A := rand.Intn(int(math.Exp2(float64(BitSize))) - 1)
		B := rand.Intn(int(math.Exp2(float64(BitSize))) - 1)
		C := Mult.ReturnVal(uint(A), uint(B))

		_, err = fmt.Fprintf(writer, fmtstr, A, B, C)
		if err != nil {
			log.Println(err)
		}

	}

	writer.Flush()

}

func NormalTestData(Mult VHDLEntityMultiplier, FolderPath string, N uint) {
	BitSize := Mult.ReturnData().BitSize
	OutputSize := Mult.ReturnData().OutputSize
	TestFile := Mult.ReturnData().TestFile

	fmtstr := "%0" + strconv.Itoa(int(BitSize)) + "b %0" + strconv.Itoa(int(BitSize)) + "b %0" + strconv.Itoa(int(OutputSize)) + "b\n"
	path := FolderPath + "/" + TestFile

	file, err := os.Create(path)
	if err != nil {
		log.Println(err)
	}

	writer := bufio.NewWriter(file)

	// maxval := int(math.Exp2(4))

	for i := 0; i < int(N); i++ {
		if i == int((N - 1)) {
			fmtstr = strings.TrimSuffix(fmtstr, "\n")
		}

		A := RandomNormalInput(int(BitSize))
		B := RandomNormalInput(int(BitSize))
		C := Mult.ReturnVal(uint(A), uint(B))

		_, err = fmt.Fprintf(writer, fmtstr, A, B, C)
		if err != nil {
			log.Println(err)
		}

	}

	writer.Flush()

}

func RandomNormalInput(size int) int {
	var sample float64
	//<div class="csl-entry">Gillani, G. A., Hanif, M. A., Verstoep, B., Gerez, S. H., Shafique, M., &#38; Kokkeler, A. B. J. (2019). MACISH: Designing Approximate MAC Accelerators With Internal-Self-Healing. <i>IEEE Access</i>, <i>7</i>, 77142–77160. https://doi.org/10.1109/ACCESS.2019.2920335</div>
	switch size {
	case 4:
		sample = rand.NormFloat64()*1.5 + 8
		if sample > (math.Exp2(4) - 1) {
			sample = (math.Exp2(4) - 1)
		}
	case 8:
		sample = rand.NormFloat64()*22.5 + 128
		if sample > (math.Exp2(8) - 1) {
			sample = (math.Exp2(8) - 1)
		}
	case 16:
		sample = rand.NormFloat64()*6553 + 32768
		if sample > (math.Exp2(16) - 1) {
			sample = (math.Exp2(16) - 1)
		}
	default:
		sample = 0
	}

	if sample < 0 {
		sample = 0
	}

	return int(sample)
}

//For VHDLEntity try to create MaxSwitching N TestData in form A*B=C for Folderpath
func MaxSwitchingTestData(Mult VHDLEntityMultiplier, FolderPath string, N uint) {
	BitSize := Mult.ReturnData().BitSize
	OutputSize := Mult.ReturnData().OutputSize
	TestFile := Mult.ReturnData().TestFile

	fmtstr := "%0" + strconv.Itoa(int(BitSize)) + "b %0" + strconv.Itoa(int(BitSize)) + "b %0" + strconv.Itoa(int(OutputSize)) + "b\n"
	path := FolderPath + "/" + TestFile

	file, err := os.Create(path)
	if err != nil {
		log.Println(err)
	}

	writer := bufio.NewWriter(file)

	switchbool := false
	for i := 0; i < int(N); i++ {
		if i == int((N - 1)) {
			fmtstr = strings.TrimSuffix(fmtstr, "\n")
		}

		var A uint
		var B uint
		if switchbool {
			switchbool = !switchbool
			A = uint(math.Exp2(float64(BitSize))) - 1
			B = uint(math.Exp2(float64(BitSize))) - 1
		} else {
			switchbool = !switchbool
			A = 0
			B = 0
		}

		C := Mult.ReturnVal(uint(A), uint(B))

		_, err = fmt.Fprintf(writer, fmtstr, A, B, C)
		if err != nil {
			log.Println(err)
		}

	}

	writer.Flush()

}

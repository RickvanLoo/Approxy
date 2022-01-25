package VHDL

import (
	"log"
	"math/big"
	"os"
	"text/template"
)

// Generic functions for badmath/VHDL

type VHDLEntityMultiplier interface {
	VHDLEntity
	Multiplier
}

//VHDLEntity describes an interface for a testable and synthesizable VHDL structure
type VHDLEntity interface {
	ReturnData() *EntityData
	GenerateVHDL(string)
	GenerateTestData(string)
	String() string //MSB -> LSB
}

type Multiplier interface {
	ReturnVal(uint, uint) uint
	Overflow() bool
	MeanAbsoluteError() float64
}

//EntityData encapsulates basic data for a VHDL structure
type EntityData struct {
	EntityName string
	BitSize    uint
	VHDLFile   string
	TestFile   string
}

//TODO: Add funcmap support
func CreateFile(FolderPath string, FileName string, TemplateFile string, Data interface{}) {
	TemplatePath := "template/" + TemplateFile

	path := FolderPath + "/" + FileName

	t, err := template.New(TemplateFile).ParseFiles(TemplatePath)
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

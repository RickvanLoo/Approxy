package VHDL

import (
	"log"
	"os"
	"text/template"
)

// Generic functions for badmath/VHDL

//VHDLEntity describes an interface for a testable and synthesizable VHDL structure
type VHDLEntity interface {
	ReturnData() *EntityData
	GenerateVHDL(string)
	GenerateTestData(string)
	String() string //MSB -> LSB
}

type UnsignedMultiplyer interface {
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
	OverflowMask := byte(0b11111111)
	byte_input := byte(input)
	byte_output := byte_input & OverflowMask
	output = uint(byte_output)
	booloutput := !(output == input)
	return output, booloutput
}

func OverflowCheck16bit(input uint) (output uint, overflow bool) {
	OverflowMask := uint16(65535)
	uint16_input := uint16(input)
	uint16_output := uint16_input & OverflowMask
	output = uint(uint16_output)
	booloutput := !(output == input)
	return output, booloutput
}

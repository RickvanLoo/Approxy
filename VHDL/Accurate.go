package VHDL

import (
	"log"
	"os"
	"text/template"
)

//UnsingedNumericAccurateMultiplyer defines a multiplyer based upon 'accuratebehav.vhd', this is a BitSize A*B=prod multiplyer, using the IEEE Numeric lib
type UnsignedNumericAccurateMultiplyer struct {
	EntityName string
	BitSize    uint
}

//Creates a VHDL file of an UnsignedNumericAccurateMultiplyer.
func UNAM_VHDL(EntityName string, BitSize uint, FolderPath string, FileName string) *UnsignedNumericAccurateMultiplyer {
	m := new(UnsignedNumericAccurateMultiplyer)
	m.EntityName = EntityName
	m.BitSize = BitSize

	templatepath := "template/accuratebehav.vhd"
	templatename := "accuratebehav.vhd"

	path := FolderPath + "/" + FileName

	t, err := template.New(templatename).ParseFiles(templatepath)
	if err != nil {
		log.Print(err)
		return m
	}

	f, err := os.Create(path)

	if err != nil {
		log.Println("create file: ", err)
		return m
	}

	err = t.ExecuteTemplate(f, templatename, m)
	if err != nil {
		log.Print("execute: ", err)
		return m
	}

	return m
}

//GenerateTestData uses the function from New2DUnsignedAcc, since their behaviour is identical.
func (m *UnsignedNumericAccurateMultiplyer) GenerateTestData(FolderPath string, FileName string) {
	Accurate2D := New2DUnsignedAcc(m.BitSize)
	Accurate2D.EntityName = m.EntityName
	Accurate2D.GenerateTestData(FolderPath, FileName)
}

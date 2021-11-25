package VHDL

import (
	"log"
	"os"
	"text/template"
)

type UnsignedAccurateMultiplyer struct {
	EntityName string
	BitSize    uint
}

func UAM_To_VHDL(EntityName string, BitSize uint, FolderPath string, FileName string) *UnsignedAccurateMultiplyer {
	m := new(UnsignedAccurateMultiplyer)
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

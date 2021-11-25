package VHDL

import (
	"log"
	"os"
	"text/template"
)

type Scaler struct {
	LUT        *LUT2D
	LUTName    string
	EntityName string
	BitSize    uint
	ScaleN     uint
}

func CreateScaler(m *LUT2D, N uint, FolderPath string) *Scaler {
	scl := new(Scaler)
	scl.LUT = m
	scl.LUTName = scl.LUT.EntityName
	scl.BitSize = scl.LUT.BitSize
	scl.EntityName = scl.LUTName + "_scaler"
	scl.ScaleN = N

	templatepath := "template/scaler.vhd"
	templatename := "scaler.vhd"
	name := scl.LUTName + "_scaler.vhd"

	path := FolderPath + "/" + name

	t, err := template.New(templatename).ParseFiles(templatepath)
	if err != nil {
		log.Print(err)
		return scl
	}

	f, err := os.Create(path)

	if err != nil {
		log.Println("create file: ", err)
		return scl
	}

	err = t.ExecuteTemplate(f, templatename, scl)
	if err != nil {
		log.Print("execute: ", err)
		return scl
	}

	return scl
}

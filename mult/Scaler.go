package mult

import (
	"log"
	"os"
	"text/template"
)

type Scaler struct {
	Item       *UnsignedAccurateMultiplyer
	ItemName   string
	EntityName string
	BitSize    uint
	ScaleN     uint
}

func CreateScaler(m *UnsignedAccurateMultiplyer, N uint, folder string) *Scaler {
	scl := new(Scaler)
	scl.Item = m
	scl.ItemName = scl.Item.Name
	scl.BitSize = scl.Item.Bitsize
	scl.EntityName = scl.ItemName + "_scaler"
	scl.ScaleN = N

	templatepath := "template/scaler.vhd"
	templatename := "scaler.vhd"
	name := scl.ItemName + "_scaler.vhd"

	path := folder + "/" + name

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

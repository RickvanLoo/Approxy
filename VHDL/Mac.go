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

type MAC struct {
	EntityName   string
	BitSize      uint
	OutputSize   uint
	Multiplier   VHDLEntityMultiplier
	VHDLFile     string
	TestFile     string
	CurrentValue uint
	T            uint
	MaximumValue uint
}

//TODO : Make it the same as Rec4, where GenerateVHDL generates the preceding VHDL files as well.

func NewMAC(Multiplier VHDLEntityMultiplier, T uint) *MAC {
	mac := new(MAC)
	mac.Multiplier = Multiplier
	MultiplierData := mac.Multiplier.ReturnData()
	mac.EntityName = "MAC_" + MultiplierData.EntityName
	mac.BitSize = MultiplierData.BitSize
	mac.VHDLFile, mac.TestFile = FileNameGen(mac.EntityName)
	mac.CurrentValue = 0
	mac.T = T

	maxinput := int(math.Exp2(float64(mac.BitSize)))
	maxval := int(Multiplier.ReturnVal(uint(maxinput-1), uint(maxinput-1)))

	bitvalue := math.Log2(float64(maxval * int(T)))
	mac.OutputSize = uint(math.Ceil(bitvalue))
	mac.MaximumValue = uint(maxval)

	return mac
}

func (mac *MAC) ReturnData() *EntityData {
	d := new(EntityData)
	d.BitSize = mac.BitSize
	d.EntityName = mac.EntityName
	d.OutputSize = mac.OutputSize
	d.TestFile = mac.TestFile
	d.VHDLFile = mac.VHDLFile

	return d
}

func (mac *MAC) GenerateVHDL(FolderPath string) {
	mac.Multiplier.GenerateVHDL(FolderPath)
	CreateFile(FolderPath, mac.VHDLFile, "macbehav.vhd", mac)
}

func (mac *MAC) GenerateTestData(FolderPath string) {
	fmtstr := "%0" + strconv.Itoa(int(mac.BitSize)) + "b %0" + strconv.Itoa(int(mac.BitSize)) + "b %0" + strconv.Itoa(int(mac.OutputSize)) + "b\n"
	path := FolderPath + "/" + mac.TestFile

	file, err := os.Create(path)
	if err != nil {
		log.Println(err)
	}

	writer := bufio.NewWriter(file)

	maxinput := int(math.Exp2(float64(mac.BitSize)))
	a := maxinput - 1
	b := maxinput - 1

	mac.ResetVal()

	for i := 0; i < int(mac.T*2+1); i++ {

		if i == (int(mac.T * 2)) {
			fmtstr = strings.TrimSuffix(fmtstr, "\n")
		}

		out := mac.ReturnVal(uint(a), uint(b))
		_, err = fmt.Fprintf(writer, fmtstr, a, b, out)
		if err != nil {
			log.Println(err)
		}
	}

	// for a := 0; a < maxval; a++ {
	// 	for b := 0; b < maxval; b++ {

	// 		if (a == 255) && (b == 255) {
	// 			fmtstr = strings.TrimSuffix(fmtstr, "\n")
	// 		}

	// 		out := mac.ReturnVal(uint(a), uint(b))

	// 		_, err = fmt.Fprintf(writer, fmtstr, a, b, out)
	// 		if err != nil {
	// 			log.Println(err)
	// 		}

	// 	}
	// }

	writer.Flush()

}

func (mac *MAC) String() string {
	return "MAC -> " + mac.Multiplier.ReturnData().EntityName
}

func (mac *MAC) ReturnVal(a uint, b uint) uint {
	c := mac.Multiplier.ReturnVal(a, b)
	mac.CurrentValue = mac.CurrentValue + c
	mac.CurrentValue, _ = OverflowCheckGeneric(mac.CurrentValue, mac.OutputSize)
	return mac.CurrentValue
}

func (mac *MAC) ResetVal() {
	mac.CurrentValue = 0
}

func (mac *MAC) GenerateVHDLEntityArray() []VHDLEntity {
	//fix me
	var out []VHDLEntity

	out = mac.Multiplier.GenerateVHDLEntityArray()

	out = append([]VHDLEntity{mac}, out...)

	return out
}

// ReturnData() *EntityData
// GenerateVHDL(string)
// GenerateTestData(string)
// String() string //MSB -> LSB

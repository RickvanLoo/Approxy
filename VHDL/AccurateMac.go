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

type UnsignedNumericAccurateMAC struct {
	EntityName   string
	BitSize      uint
	OutputSize   uint
	VHDLFile     string
	TestFile     string
	CurrentValue uint
	T            uint
	MaximumValue uint
}

func NewUnsignedAccurateMAC(BitSize uint, T uint) *UnsignedNumericAccurateMAC {
	mac := new(UnsignedNumericAccurateMAC)
	mac.EntityName = "UAMAC_" + strconv.Itoa(int(BitSize)) + "b" + strconv.Itoa(int(T)) + "T"
	mac.BitSize = BitSize
	mac.VHDLFile, mac.TestFile = FileNameGen(mac.EntityName)
	mac.CurrentValue = 0
	mac.T = T

	maxinput := int(math.Exp2(float64(mac.BitSize)))
	maxval := (maxinput - 1) * (maxinput - 1)

	bitvalue := math.Log2(float64(maxval * int(T)))
	mac.OutputSize = uint(math.Ceil(bitvalue))
	mac.MaximumValue = uint(maxval)

	return mac
}

func (mac *UnsignedNumericAccurateMAC) ReturnData() *EntityData {
	d := new(EntityData)
	d.BitSize = mac.BitSize
	d.EntityName = mac.EntityName
	d.OutputSize = mac.OutputSize
	d.TestFile = mac.TestFile
	d.VHDLFile = mac.VHDLFile

	return d
}

func (mac *UnsignedNumericAccurateMAC) GenerateVHDL(FolderPath string) {
	CreateFile(FolderPath, mac.VHDLFile, "macaccurate.vhd", mac)
}

func (mac *UnsignedNumericAccurateMAC) GenerateTestData(FolderPath string) {
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
	writer.Flush()

}

func (mac *UnsignedNumericAccurateMAC) String() string {
	return mac.EntityName
}

func (mac *UnsignedNumericAccurateMAC) ReturnVal(a uint, b uint) uint {
	c := a * b
	mac.CurrentValue = mac.CurrentValue + c
	mac.CurrentValue, _ = OverflowCheckGeneric(mac.CurrentValue, mac.OutputSize)
	return mac.CurrentValue
}

func (mac *UnsignedNumericAccurateMAC) ResetVal() {
	mac.CurrentValue = 0
}

func (mac *UnsignedNumericAccurateMAC) GenerateVHDLEntityArray() []VHDLEntity {
	//fix me
	var out []VHDLEntity

	out = append([]VHDLEntity{mac}, out...)

	return out
}

func (mac *UnsignedNumericAccurateMAC) MeanAbsoluteError() float64 {
	log.Println("ERROR, MAC MAE NOT IMPLEMENTED")
	return 0
}

func (mac *UnsignedNumericAccurateMAC) Overflow() bool {
	log.Println("ERROR, MAC OVERFLOW CHECK NOT IMPLEMENTED")
	return false
}

// ReturnData() *EntityData
// GenerateVHDL(string)
// GenerateTestData(string)
// String() string //MSB -> LSB

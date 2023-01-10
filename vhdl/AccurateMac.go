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

// UnsignedNumericAccurateMAC is an Accurate Multiply-Accumulator on basis of template "macaccurate.vhd"
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

// NewUnsignedAccurateMAC creates the struct NewUnsignedAccurateMAC
// BitSize: Bit-length of Multiplier inputs
// T: N-amount of Accumulation guaranteed before overflow
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

// ReturnData returns a struct containing metadata of the multiplier
func (mac *UnsignedNumericAccurateMAC) ReturnData() *EntityData {
	d := new(EntityData)
	d.BitSize = mac.BitSize
	d.EntityName = mac.EntityName
	d.OutputSize = mac.OutputSize
	d.TestFile = mac.TestFile
	d.VHDLFile = mac.VHDLFile

	return d
}

// GenerateVHDL creates the VHDL file in FolderPath
func (mac *UnsignedNumericAccurateMAC) GenerateVHDL(FolderPath string) {
	CreateFile(FolderPath, mac.VHDLFile, "macaccurate.vhd", mac)
}

// GenerateTestData creates MAC testing data in plain text. Two inputs and outputs, seperated by \t
// Note that MAC testing data is different, instead of applying an exhaustive test over the full range on inputs, inputs are kept static and overflow behaviour is verified
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

// ReturnVal outputs the current value of the Multiply-Accumulator on basis of input A and B
func (mac *UnsignedNumericAccurateMAC) ReturnVal(a uint, b uint) uint {
	c := a * b
	mac.CurrentValue = mac.CurrentValue + c
	mac.CurrentValue, _ = OverflowCheckGeneric(mac.CurrentValue, mac.OutputSize)
	return mac.CurrentValue
}

// ResetVal resets the accumulator to 0
func (mac *UnsignedNumericAccurateMAC) ResetVal() {
	mac.CurrentValue = 0
}

// GenerateVHDLEntityArray creates an array of potentially multiple VHDLEntities, sorted by priority for synthesizing
// For example: Multiplier A uses a VHDL portmap for the smaller Multiplier B & C. B & C need to be synthesized first, hence A will be last in the array
func (mac *UnsignedNumericAccurateMAC) GenerateVHDLEntityArray() []VHDLEntity {
	//fix me
	var out []VHDLEntity

	out = append([]VHDLEntity{mac}, out...)

	return out
}

// MeanAbsoluteError returns the MeanAbsoluteError of the multiplier in float64
// FIX ME: NOT IMPLEMENTED
func (mac *UnsignedNumericAccurateMAC) MeanAbsoluteError() float64 {
	log.Println("ERROR, MAC MAE NOT IMPLEMENTED")
	return 0
}

// Overflow returns a boolean if any internal overflow has occured
// FIX ME: NOT IMPLEMENTED
func (mac *UnsignedNumericAccurateMAC) Overflow() bool {
	log.Println("ERROR, MAC OVERFLOW CHECK NOT IMPLEMENTED")
	return false
}

// ReturnData() *EntityData
// GenerateVHDL(string)
// GenerateTestData(string)
// String() string //MSB -> LSB

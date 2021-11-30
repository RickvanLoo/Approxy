package VHDL

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"text/template"
)

//LUTArray[0] = AH*BH
//LUTArray[1] = AH*BL
//LUTArray[2] = AL*BH
//LUTArray[3] = AL*BL
type Recursive4 struct {
	EntityName string
	BitSize    uint      //Default to 4
	OutputSize uint      //Default to 8
	LUTArray   [4]*LUT2D //Size of 4
}

func NewRecursive4(LUTArray [4]*LUT2D) *Recursive4 {
	r4 := new(Recursive4)
	r4.BitSize = 4
	r4.OutputSize = 8
	r4.EntityName = "Recursive4"
	r4.LUTArray = LUTArray
	r4.LUTArray[0].EntityName = "AH_BH"
	r4.LUTArray[1].EntityName = "AH_BL"
	r4.LUTArray[2].EntityName = "AL_BH"
	r4.LUTArray[3].EntityName = "AL_BL"

	return r4
}

func (r4 *Recursive4) ReturnVal(a uint, b uint) uint {
	AHBH_LUT := r4.LUTArray[0]
	AHBL_LUT := r4.LUTArray[1]
	ALBH_LUT := r4.LUTArray[2]
	ALBL_LUT := r4.LUTArray[3]

	bin_input := make([]byte, 2)
	bin_input[0] = byte(a)
	bin_input[1] = byte(b)

	maskH := byte(0b00001100)
	maskL := byte(0b00000011)

	AHALBHBL := make([]byte, 4)
	AHALBHBL[0] = (bin_input[0] & maskH) >> 2 //AH
	AHALBHBL[1] = bin_input[0] & maskL        //AL
	AHALBHBL[2] = (bin_input[1] & maskH) >> 2 //BH
	AHALBHBL[3] = bin_input[1] & maskL        //BL

	AH := uint(AHALBHBL[0])
	AL := uint(AHALBHBL[1])
	BH := uint(AHALBHBL[2])
	BL := uint(AHALBHBL[3])

	AHBH := AHBH_LUT.ReturnVal(AH, BH)
	AHBL := AHBL_LUT.ReturnVal(AH, BL)
	ALBH := ALBH_LUT.ReturnVal(AL, BH)
	ALBL := ALBL_LUT.ReturnVal(AL, BL)

	output := ALBL + (ALBH << 2) + (AHBL << 2) + (AHBH << 4)

	return output
}

func (r4 *Recursive4) GenerateTestData(FolderPath string, FileName string) {
	fmtstr := "%0" + strconv.Itoa(int(r4.BitSize)) + "b %0" + strconv.Itoa(int(r4.BitSize)) + "b %0" + strconv.Itoa(int(r4.OutputSize)) + "b\n"
	path := FolderPath + "/" + FileName

	file, err := os.Create(path)
	if err != nil {
		log.Println(err)
	}

	writer := bufio.NewWriter(file)

	maxval := int(math.Exp2(4))

	for a := 0; a < maxval; a++ {
		for b := 0; b < maxval; b++ {

			if (a == 15) && (b == 15) {
				fmtstr = strings.TrimSuffix(fmtstr, "\n")
			}

			out := r4.ReturnVal(uint(a), uint(b))

			_, err = fmt.Fprintf(writer, fmtstr, a, b, out)
			if err != nil {
				log.Println(err)
			}

		}
	}

	writer.Flush()

}

func (r4 *Recursive4) VHDLtoFile(FolderPath string, FileName string) {
	for _, mult := range r4.LUTArray {
		mult.VHDLtoFile(FolderPath, mult.EntityName+".vhd")
	}

	templatepath := "template/rec4behav.vhd"
	templatename := "rec4behav.vhd"

	path := FolderPath + "/" + FileName

	t, err := template.New(templatename).ParseFiles(templatepath)
	if err != nil {
		log.Print(err)
		return
	}

	f, err := os.Create(path)

	if err != nil {
		log.Println("create file: ", err)
		return
	}

	err = t.ExecuteTemplate(f, templatename, r4)
	if err != nil {
		log.Print("execute: ", err)
		return
	}
}

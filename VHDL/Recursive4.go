package VHDL

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
	r4.BitSize = 8
	r4.EntityName = "Recursive4"
	r4.LUTArray = LUTArray

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

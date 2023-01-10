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

// ExternalMult is a struct that defines a Multiplier based on existing synthesizable VHDL files instead of VHDL templates.
// Templates for this model are stored in the folder template/external
// The Entity within the external synthesizable VHDL needs to be modified, info found in the wiki
type ExternalMult struct {
	EntityName      string
	BitSize         uint
	OutputSize      uint
	VHDLFile        string
	TestFile        string
	TemplateFile    string
	Behaviour       [][]uint
	ExtraVHDLEntity []VHDLEntity
}

// NewExternalMult creates a new ExternalMult struct
// EntityName: Name of Entity
// BitSize: Bit-length of input
// Template: Filename of external template, ie. "mult_accurate.vhd"
func NewExternalMult(EntityName string, BitSize uint, Template string) *ExternalMult {
	ext := new(ExternalMult)
	ext.EntityName = EntityName
	ext.BitSize = BitSize
	ext.OutputSize = 2 * ext.BitSize
	ext.TemplateFile = Template
	ext.VHDLFile, ext.TestFile = FileNameGen(ext.EntityName)
	maxval := int(math.Exp2(float64(ext.BitSize)))

	ext.Behaviour = make([][]uint, maxval)
	for i := range ext.Behaviour {
		ext.Behaviour[i] = make([]uint, maxval)
	}
	return ext
}

// ReturnData returns a struct containing metadata of the multiplier
func (ext *ExternalMult) ReturnData() *EntityData {
	d := new(EntityData)
	d.BitSize = ext.BitSize
	d.EntityName = ext.EntityName
	d.OutputSize = ext.OutputSize
	d.TestFile = ext.TestFile
	d.VHDLFile = ext.VHDLFile
	return d
}

// GenerateVHDL creates the VHDL file in FolderPath
func (ext *ExternalMult) GenerateVHDL(FolderPath string) {
	CreateFile(FolderPath, ext.VHDLFile, "external/"+ext.TemplateFile, ext)
}

// AddVHDLEntity is used when external templates require portmap and thus multiple VHDL files
// Create for each portmapped external VHDL file a new ExternalMult struct, and add according to synthesis priority
func (ext *ExternalMult) AddVHDLEntity(Entity VHDLEntity) {
	ext.ExtraVHDLEntity = append(ext.ExtraVHDLEntity, Entity)
}

// GenerateTestData creates a plaintext testdata file containing both inputs and the output in binary seperated by \t
// WARNING: the testdata will be invalid if ParseXSIMOutput has not run yet
func (ext *ExternalMult) GenerateTestData(FolderPath string) {
	fmtstr := "%0" + strconv.Itoa(int(ext.BitSize)) + "b %0" + strconv.Itoa(int(ext.BitSize)) + "b\n"
	path := FolderPath + "/" + ext.TestFile

	file, err := os.Create(path)
	if err != nil {
		log.Println(err)
	}

	writer := bufio.NewWriter(file)

	maxval := int(math.Exp2(float64(ext.BitSize)))

	for a := 0; a < maxval; a++ {
		for b := 0; b < maxval; b++ {

			if (a == (maxval - 1)) && (b == (maxval - 1)) {
				fmtstr = strings.TrimSuffix(fmtstr, "\n")
			}

			_, err = fmt.Fprintf(writer, fmtstr, a, b)
			if err != nil {
				log.Println(err)
			}

		}
	}

	writer.Flush()
}

// GenerateVHDLEntityArray creates an array of potentially multiple VHDLEntities, sorted by priority for synthesizing
// For example: Multiplier A uses a VHDL portmap for the smaller Multiplier B & C. B & C need to be synthesized first, hence A will be last in the array
func (ext *ExternalMult) GenerateVHDLEntityArray() []VHDLEntity {
	var out []VHDLEntity
	out = append(out, ext)
	out = append(out, ext.ExtraVHDLEntity...)
	return out
}

func (ext *ExternalMult) String() string {
	return ext.EntityName
}

// ParseXSIMOutput parses the output file of an XSIM simulation in FolderPath
// The XSIM simulation (done via the Vivado package) applies the whole range of inputs to the multiplier and exports output data
// The data is parsed via this function to build the model behaviour within the Behaviour [][]uint field
func (ext *ExternalMult) ParseXSIMOutput(FolderPath string) {
	file, err := os.Open(FolderPath + "/" + "out_" + ext.TestFile)
	if err != nil {
		log.Println(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text []string

	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	file.Close()

	maxval := int(math.Exp2(float64(ext.BitSize)))

	textindex := 0

	for a := 0; a < maxval; a++ {
		for b := 0; b < maxval; b++ {

			uintout, _ := strconv.ParseUint(text[textindex], 10, 64)
			ext.Behaviour[a][b] = uint(uintout)
			textindex = textindex + 1
		}
	}
}

// ReturnVal returns the output of the multiplier
func (ext *ExternalMult) ReturnVal(a, b uint) uint {
	return ext.Behaviour[a][b]
}

// MeanAbsoluteError returns the MeanAbsoluteError of the multiplier in float64, for a uniform distribution
func (ext *ExternalMult) MeanAbsoluteError() float64 {
	maxval := int(math.Exp2(float64(ext.BitSize)))
	accum := float64(0)
	for a := 1; a < maxval; a++ {
		for b := 1; b < maxval; b++ {
			accResult := float64(a * b)
			r4Result := ext.ReturnVal(uint(a), uint(b))
			accum += math.Abs(float64(r4Result) - accResult)
		}
	}

	return float64(1.0/math.Exp2(float64(ext.OutputSize))) * accum
}

// MeanAbsoluteErrorNormalDist returns the MeanAbsoluteError of the multiplier in float64, for a normal distribution
// N = the amount of random input/output combinations
func (ext *ExternalMult) MeanAbsoluteErrorNormalDist(N int) float64 {
	accum := float64(0)

	for i := 0; i < N; i++ {
		a := RandomNormalInput(int(ext.BitSize))
		b := RandomNormalInput(int(ext.BitSize))
		accResult := float64(a * b)
		r4Result := ext.ReturnVal(uint(a), uint(b))
		accum += math.Abs(float64(r4Result) - accResult)
	}

	return float64(1.0/float64(N)) * accum
}

// AverageRelativeError returns the Average Relative Error of the multiplier in float64
func (ext *ExternalMult) AverageRelativeError() float64 {
	maxval := int(math.Exp2(float64(ext.BitSize)))
	accum := float64(0)
	for a := 1; a < maxval; a++ {
		for b := 1; b < maxval; b++ {
			accResult := float64(a * b)
			r4Result := ext.ReturnVal(uint(a), uint(b))
			accum += math.Abs((float64(r4Result) - accResult) / accResult)
		}
	}

	return float64(1.0/math.Exp2(float64(ext.OutputSize))) * accum

}

// Overflow returns a boolean for internal overflow, not supported for ExternalMult
func (ext *ExternalMult) Overflow() bool {
	log.Println("Overflow not supported for External Multipliers")
	return false
}

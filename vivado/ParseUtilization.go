package vivado

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"

	VHDL "github.com/RickvanLoo/Approxy/vhdl"
)

// Utilization is struct describing the physicial element utilization for a design
type Utilization struct {
	EntityName string `json:"-"` //For identification, ignored when unmarshaling to JSON file
	TotalLUT   int
	LogicLUT   int
	LUTRAMs    int
	SRLs       int
	FFs        int
	RAMB36     int
	RAMB18     int
	DSP        int
	CARRY      int
}

// ParseUtilizationReport parses the exported Vivado utilization report using hacky string manipulation
// This is needed because Vivado can only report a spreadsheet report in GUI mode, not in batch mode (?)
// Heavy usage of horrible string manipulation in the next function, try to touch as little as possible
func ParseUtilizationReport(FolderPath string, Entity VHDL.VHDLEntity) *Utilization {
	util := new(Utilization)
	filextension := "_post_place_ult.rpt"
	filename := Entity.ReturnData().EntityName + filextension

	file, err := os.Open(FolderPath + "/" + filename)

	if err != nil {
		//log.Printf("Warning, Returning ZERO: failed opening file: %s", err)
		return util
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	file.Close()

	var entitylines []string

	for _, line := range txtlines {
		if strings.Contains(line, Entity.ReturnData().EntityName) {
			entitylines = append(entitylines, line)
		}
	}
	UtilizationStr := entitylines[2]                             //wow
	UtilizationStr = strings.ReplaceAll(UtilizationStr, " ", "") //Remove spaces
	Results := strings.Split(UtilizationStr, "|")                //Hmm

	//[0] & [11] are empty strings due to the | at the beginning and end
	util.EntityName = Results[1] //Skip [2], not needed info
	util.TotalLUT, _ = strconv.Atoi(Results[3])
	util.LogicLUT, _ = strconv.Atoi(Results[4])
	util.LUTRAMs, _ = strconv.Atoi(Results[5])
	util.SRLs, _ = strconv.Atoi(Results[6])
	util.FFs, _ = strconv.Atoi(Results[7])
	util.RAMB36, _ = strconv.Atoi(Results[8])
	util.RAMB18, _ = strconv.Atoi(Results[9])
	util.DSP, _ = strconv.Atoi(Results[10])

	///Carry

	filextension = "_primitive.rpt"

	filename = Entity.ReturnData().EntityName + filextension

	file, err = os.Open(FolderPath + "/" + filename)

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner = bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var primlines []string

	for scanner.Scan() {
		primlines = append(primlines, scanner.Text())
	}
	util.CARRY, _ = strconv.Atoi(primlines[0])

	file.Close()

	return util
}

package vivado

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	VHDL "github.com/RickvanLoo/Approxy/vhdl"
)

// Timing is a struct that contains timing information for a single design, only the worst or critical path is reported.
type Timing struct {
	EntityName string `json:"-"` //For identification, ignored when unmarshaling to JSON file
	EndPoint   string
	WorstPath  float64 //ns
	MaxFreq    float64 //MHz
}

// ParseTimingReport parses the timing report exported by Vivado, it uses a hacky approach of string manipulation
// The timing report is named EntityName + "_post_place_time.rpt"
func ParseTimingReport(FolderPath string, Entity VHDL.VHDLEntity) *Timing {
	time := new(Timing)
	filextension := "_post_place_time.rpt"
	filename := Entity.ReturnData().EntityName + filextension

	file, err := os.Open(FolderPath + "/" + filename)

	if err != nil {
		//log.Printf("Warning, Returning ZERO: failed opening file: %s", err)
		return time
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	file.Close()

	output := strings.Fields(txtlines[15]) //hacky

	time.EndPoint = output[0]
	time.WorstPath, _ = strconv.ParseFloat(output[1], 64)
	time.MaxFreq = (1 / (time.WorstPath * 1e-9)) * 1e-6
	time.EntityName = Entity.ReturnData().EntityName

	return time

}

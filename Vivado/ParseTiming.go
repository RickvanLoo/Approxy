package Vivado

import (
	"approxy/VHDL"
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Timing struct {
	EntityName string `json:"-"`
	EndPoint   string
	WorstPath  float64 //ns
	MaxFreq    float64 //MHz
}

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

package vivado

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	VHDL "github.com/RickvanLoo/Approxy/vhdl"
)

// Report is a struct containing all parsed Vivado report data of a single design, ie. a single multiplier, a scaled N=100 multiplier or a MAC
type Report struct {
	EntityName string
	Util       *Utilization
	Power      *PowerReport
	Timing     *Timing
	Other      []Data
}

// Run is a collection of Reports, used for comparing different multiplier designs. Reports can be added to the Run, the Run can be exported and modified to a local JSON file at any time.
// A Run is basically a very simple database for Vivado Reports of different multiplier designs.
type Run struct {
	Name       string
	Reports    []Report
	ResultPath string `json:"-"`
	ReportPath string `json:"-"`
	Other      []Data
}

// Data is a struct containing Key/Value information. Used to add non-Vivado results to a Report, for example MeanAbsoluteError
// Can also be used to add misc. data to a Run, ie. parameters
// TODO : Replace by Go maps to improve efficiency
type Data struct {
	Key   string
	Value string
}

// CreateReport creates a single report on basis of a VHDLEntity, requires Vivado reports to be available.
func CreateReport(FolderPath string, Entity VHDL.VHDLEntity) *Report {
	Report := new(Report)

	Report.EntityName = Entity.ReturnData().EntityName
	Report.Power = ParsePowerReport(FolderPath, Entity)
	Report.Timing = ParseTimingReport(FolderPath, Entity)
	Report.Util = ParseUtilizationReport(FolderPath, Entity)

	return Report
}

// AddData adds a Key/Value data element to a Report
func (r *Report) AddData(key string, value string) {
	var data Data
	data.Key = key
	data.Value = value
	r.Other = append(r.Other, data)
}

// StartRun creates a new run if non-existant. Opens and modifies old run if does exist based upon name and reportpath
func StartRun(ResultPath string, ReportPath string, Name string) *Run {
	Run := new(Run)
	Run.ResultPath = ResultPath
	Run.ReportPath = ReportPath
	Run.Name = Name

	Filepath := Run.ResultPath + "/" + Run.Name + ".json"

	if _, err := os.Stat(Filepath); err == nil {
		//Run exists
		Reset := "\033[0m"
		Yellow := "\033[33m"
		log.Printf(Yellow + "Warning, run: " + Name + " exists!\n" + Reset)

		jsonFile, err := os.Open(Filepath)
		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &Run)

	} else if errors.Is(err, os.ErrNotExist) {
		//Run does not exist
		file, _ := json.Marshal(Run)
		_ = os.WriteFile(Filepath, file, 0644)

	} else {
		log.Println(err)
	}

	return Run
}

// AddReport adds an Report to a Run, automatically removes duplicates on basis of EntityName, and updates JSON file
func (r *Run) AddReport(Report Report) {
	r.Reports = append(r.Reports, Report)
	r.removeDuplicates()
	r.updateJSON()
}

// Exists checks if a Report exists within the Run on basis of their EntityName
func (r *Run) Exists(EntityName string) bool {
	for _, report := range r.Reports {
		if report.EntityName == EntityName {
			return true
		}
	}
	return false
}

func (r *Run) updateJSON() {
	Filepath := r.ResultPath + "/" + r.Name + ".json"
	file, _ := json.Marshal(r)
	_ = os.WriteFile(Filepath, file, 0644)
}

func (r *Run) removeDuplicates() {
	CurrentReports := make(map[string]Report)

	for _, report := range r.Reports {
		CurrentReports[report.EntityName] = report
	}

	var NewReports []Report

	for _, report := range CurrentReports {
		NewReports = append(NewReports, report)
	}

	r.Reports = NewReports
}

// AddData adds key/value data to a single Run
func (r *Run) AddData(key string, value string) {
	var data Data
	data.Key = key
	data.Value = value
	r.Other = append(r.Other, data)
	r.updateJSON()
}

// ClearData clears ALL key/value data of a single Run
func (r *Run) ClearData() {
	r.Other = nil
}

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

type Report struct {
	EntityName string
	Util       *Utilization
	Power      *PowerReport
	Timing     *Timing
	Other      []Data
}

type Run struct {
	Name       string
	Reports    []Report
	ResultPath string `json:"-"`
	ReportPath string `json:"-"`
	Other      []Data
}

type Data struct {
	Key   string
	Value string
}

func CreateReport(FolderPath string, Entity VHDL.VHDLEntity) *Report {
	Report := new(Report)

	Report.EntityName = Entity.ReturnData().EntityName
	Report.Power = ParsePowerReport(FolderPath, Entity)
	Report.Timing = ParseTimingReport(FolderPath, Entity)
	Report.Util = ParseUtilizationReport(FolderPath, Entity)

	return Report
}

func (r *Report) AddData(key string, value string) {
	var data Data
	data.Key = key
	data.Value = value
	r.Other = append(r.Other, data)
}

// Creates new run if non-existant. Opens old run if does exist based upon Name
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

func (r *Run) AddReport(Report Report) {
	r.Reports = append(r.Reports, Report)
	r.RemoveDuplicates()
	r.updateJSON()
}

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

func (r *Run) RemoveDuplicates() {
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

func (r *Run) AddData(key string, value string) {
	var data Data
	data.Key = key
	data.Value = value
	r.Other = append(r.Other, data)
	r.updateJSON()
}

func (r *Run) ClearData() {
	r.Other = nil
}

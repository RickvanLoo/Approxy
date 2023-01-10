package Vivado

import (
	"approxy/vhdl"
	"encoding/xml"
	"io/ioutil"
	"os"
	"strconv"
)

// RptDoc was generated 2022-05-03 14:57:01 by rick on rivalo.
// Using 'zek' to convert XML to Golang
// This needs refactoring into a readable XML struct
type RptDoc struct {
	XMLName     xml.Name `xml:"RptDoc"`
	Text        string   `xml:",chardata"`
	Cmdline     string   `xml:"cmdline,attr"`
	Designname  string   `xml:"designname,attr"`
	Designstate string   `xml:"designstate,attr"`
	Partname    string   `xml:"partname,attr"`
	Speedfile   string   `xml:"speedfile,attr"`
	Title       string   `xml:"title,attr"`
	Version     string   `xml:"version,attr"`
	Section     []struct {
		Text  string `xml:",chardata"`
		Class string `xml:"class,attr"`
		Title string `xml:"title,attr"`
		Table struct {
			Text               string `xml:",chardata"`
			Class              string `xml:"class,attr"`
			Style              string `xml:"style,attr"`
			Title              string `xml:"title,attr"`
			UseFootnoteNumbers string `xml:"useFootnoteNumbers,attr"`
			Footnote           struct {
				Text     string `xml:",chardata"`
				AttrText string `xml:"text,attr"`
			} `xml:"footnote"`
			Tablerow []struct {
				Text           string `xml:",chardata"`
				Class          string `xml:"class,attr"`
				Suppressoutput string `xml:"suppressoutput,attr"`
				Wordwrap       string `xml:"wordwrap,attr"`
				Tablecell      []struct {
					Text     string `xml:",chardata"`
					Class    string `xml:"class,attr"`
					Contents string `xml:"contents,attr"`
					Halign   string `xml:"halign,attr"`
					Type     string `xml:"type,attr"`
				} `xml:"tablecell"`
			} `xml:"tablerow"`
		} `xml:"table"`
		Section []struct {
			Text  string `xml:",chardata"`
			Class string `xml:"class,attr"`
			Title string `xml:"title,attr"`
			Table struct {
				Text               string `xml:",chardata"`
				Class              string `xml:"class,attr"`
				Style              string `xml:"style,attr"`
				Title              string `xml:"title,attr"`
				UseFootnoteNumbers string `xml:"useFootnoteNumbers,attr"`
				Tablerow           []struct {
					Text           string `xml:",chardata"`
					Class          string `xml:"class,attr"`
					Suppressoutput string `xml:"suppressoutput,attr"`
					Wordwrap       string `xml:"wordwrap,attr"`
					Tableheader    []struct {
						Text     string `xml:",chardata"`
						Class    string `xml:"class,attr"`
						Contents string `xml:"contents,attr"`
						Halign   string `xml:"halign,attr"`
						Width    string `xml:"width,attr"`
					} `xml:"tableheader"`
					Tablecell []struct {
						Text          string `xml:",chardata"`
						Class         string `xml:"class,attr"`
						Contents      string `xml:"contents,attr"`
						Halign        string `xml:"halign,attr"`
						Type          string `xml:"type,attr"`
						Decimalplaces string `xml:"decimalplaces,attr"`
					} `xml:"tablecell"`
				} `xml:"tablerow"`
			} `xml:"table"`
		} `xml:"section"`
	} `xml:"section"`
}

type PowerReport struct {
	EntityName       string `json:"-"`
	Total_Power      float64
	Dynamic_Power    float64
	Static_Power     float64
	Confidence_Level string
}

func ParsePowerReport(FolderPath string, Entity vhdl.VHDLEntity) *PowerReport {

	filextension := "_post_place_power.rpt"
	filename := Entity.ReturnData().EntityName + filextension

	file, err := os.Open(FolderPath + "/" + filename)

	if err != nil {
		//log.Printf("Warning, Returning ZERO: failed opening file: %s", err)
		Report := new(PowerReport)
		return Report
	}

	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)
	var PowerXML RptDoc
	xml.Unmarshal(byteValue, &PowerXML)

	Report := new(PowerReport)
	//https://jsonformatter.org/xml-viewer
	TotalOnChipW_contents := PowerXML.Section[0].Table.Tablerow[0].Tablecell[1].Contents
	DynamicW_contents := PowerXML.Section[0].Table.Tablerow[3].Tablecell[1].Contents
	StaticW_contents := PowerXML.Section[0].Table.Tablerow[4].Tablecell[1].Contents
	ConfidenceLevel_contents := PowerXML.Section[0].Table.Tablerow[8].Tablecell[1].Contents

	TotalOnChipW_float, _ := strconv.ParseFloat(TotalOnChipW_contents, 64)
	DynamicW_float, _ := strconv.ParseFloat(DynamicW_contents, 64)
	StaticW_float, _ := strconv.ParseFloat(StaticW_contents, 64)

	Report.Confidence_Level = ConfidenceLevel_contents
	Report.Dynamic_Power = DynamicW_float
	Report.Total_Power = TotalOnChipW_float
	Report.Static_Power = StaticW_float
	Report.EntityName = Entity.ReturnData().EntityName

	return Report
}

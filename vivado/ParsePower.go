package vivado

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/RickvanLoo/Approxy/vhdl"
)

// RptDoc encapsulates the power report XML exported by Vivado
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

// PowerReport is a struct containing power results for a single design
type PowerReport struct {
	EntityName      string  `json:"-"` //For identification, ignored when unmarshaling to JSON file
	TotalPower      float64 //W
	DynamicPower    float64 //W
	StaticPower     float64 //W
	ConfidenceLevel string  //High required
}

// ParsePowerReport parses an Vivado XML power report, ignores irrelevant info and export data in the PowerReport struct
// The power report is named EntityName + "_post_place_power.rpt"
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
	TotalOnChipWcontents := PowerXML.Section[0].Table.Tablerow[0].Tablecell[1].Contents
	DynamicWcontents := PowerXML.Section[0].Table.Tablerow[3].Tablecell[1].Contents
	StaticWcontents := PowerXML.Section[0].Table.Tablerow[4].Tablecell[1].Contents
	ConfidenceLevelcontents := PowerXML.Section[0].Table.Tablerow[8].Tablecell[1].Contents

	TotalOnChipWfloat, _ := strconv.ParseFloat(TotalOnChipWcontents, 64)
	DynamicWfloat, _ := strconv.ParseFloat(DynamicWcontents, 64)
	StaticWfloat, _ := strconv.ParseFloat(StaticWcontents, 64)

	Report.ConfidenceLevel = ConfidenceLevelcontents
	Report.DynamicPower = DynamicWfloat
	Report.TotalPower = TotalOnChipWfloat
	Report.StaticPower = StaticWfloat
	Report.EntityName = Entity.ReturnData().EntityName

	return Report
}

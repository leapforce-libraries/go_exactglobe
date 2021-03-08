package exactglobe

import (
	"encoding/xml"
	"fmt"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type Response struct {
	XMLName xml.Name `xml:"feed"`
	Entries []Entry  `xml:"entry"`
	Link    []Link   `xml:"link"`
}

type Entry struct {
	XMLName xml.Name `xml:"entry"`
	ID      string   `xml:"id"`
	Content Content  `xml:"content"`
}

type Link struct {
	XMLName xml.Name `xml:"link"`
	Rel     string   `xml:"rel,attr"`
	Title   string   `xml:"title,attr"`
	Href    string   `xml:"href,attr"`
}

type Content struct {
	XMLName    xml.Name   `xml:"content"`
	Properties Properties `xml:"properties"`
}

type Properties struct {
	XMLName           xml.Name `xml:"properties"`
	AllocationLevel   int32
	Class1            string
	Class2            string
	Class3            string
	Class4            string
	ClassDescription1 string
	ClassDescription2 string
	ClassDescription3 string
	ClassDescription4 string
	Code              string
	CompanyCode       string
	CompanyName       string
	CreatedDate       string
	Creator           int32
	CreatorName       string
	Description       string
	Description1      string
	Description2      string
	Description3      string
	Description4      string
	DirectManager     int32
	DirectManagerName string
	Enabled           bool
	GLAccount         string
	GLOffsetAccount   string
	ID                int32
	ModifiedDate      string
	Modifier          int32
	ModifierName      string
	StandardRate      float64
	TextFreeField1    string
	TextFreeField2    string
	TextFreeField3    string
	TextFreeField4    string
	TextFreeField5    string
	NumberFreeField1  float64
	NumberFreeField2  float64
	NumberFreeField3  float64
	NumberFreeField4  float64
	NumberFreeField5  float64
}

type GetCostCentersConfig struct {
	Top *uint
}

// GetCostCenters returns all entries
//
func (service *Service) GetCostCenters(getCostCentersConfig *GetCostCentersConfig) (*[]Entry, *errortools.Error) {
	values := url.Values{}

	var top uint = 1
	var skiptoken *string = nil

	if getCostCentersConfig != nil {
		if getCostCentersConfig.Top != nil {
			top = *getCostCentersConfig.Top
		}
	}
	values.Set("$top", fmt.Sprintf("%v", top))

	entries := []Entry{}

	for true {
		if skiptoken != nil {
			values.Set("$skiptoken", *skiptoken)
		}

		response := Response{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf("CostCenter/?%s", values.Encode())),
			ResponseModel: &response,
		}

		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		entries = append(entries, response.Entries...)

		skiptoken = service.extractSkiptoken(&response.Link)

		if skiptoken == nil {
			break
		}
	}

	return &entries, nil
}

package exactglobe

import (
	"encoding/xml"
	"fmt"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type Response struct {
	XMLName xml.Name     `xml:"feed"`
	Entries []CostCenter `xml:"entry"`
}

type CostCenter struct {
	XMLName xml.Name `xml:"entry"`
	ID      string   `xml:"id"`
	Content Content  `xml:"content"`
}

type Content struct {
	XMLName    xml.Name   `xml:"content"`
	Properties Properties `xml:"properties"`
}

type Properties struct {
	XMLName     xml.Name `xml:"properties"`
	CreatedDate string   `xml:"CreatedDate"`
	ID          string   `xml:"ID"`
}

type GetCostCentersConfig struct {
	Top  *uint
	Skip *uint
}

// GetCostCenters returns all costCenters
//
func (service *Service) GetCostCenters(getCostCentersConfig *GetCostCentersConfig) (*[]CostCenter, *errortools.Error) {
	values := url.Values{}

	skip := uint(0)
	top := uint(100)

	if getCostCentersConfig != nil {
		if getCostCentersConfig.Top != nil {
			top = *getCostCentersConfig.Top
		}
		if getCostCentersConfig.Skip != nil {
			skip = *getCostCentersConfig.Skip
		}
	}
	values.Set("$top", fmt.Sprintf("%v", top))

	costCenters := []CostCenter{}

	for true {
		values.Set("$skip", fmt.Sprintf("%v", skip))

		response := Response{}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf("CostCenter/?%s", values.Encode())),
			ResponseModel: &response,
		}
		fmt.Println(service.url(fmt.Sprintf("CostCenter/?%s", values.Encode())))

		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		costCenters = append(costCenters, response.Entries...)

		if len(response.Entries) < int(top) {
			break
		}
		skip += top
	}

	return &costCenters, nil
}

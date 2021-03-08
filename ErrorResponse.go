package exactglobe

import "encoding/xml"

// ErrorResponse stores general error response
//

type ErrorResponse struct {
	XMLName xml.Name `xml:"error"`
	Code    string   `xml:"code"`
	Message string   `xml:"message"`
}

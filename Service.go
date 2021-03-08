package exactglobe

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/Azure/go-ntlmssp"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

const (
	apiPath string = "services/Exact.Entity.REST.EG"
)

// type
//
type Service struct {
	host         string
	serverName   string
	databaseName string
	username     string
	password     string
	httpService  *go_http.Service
}

type ServiceConfig struct {
	Host         string
	ServerName   string
	DatabaseName string
	Username     string
	Password     string
}

func NewService(serviceConfig *ServiceConfig) (*Service, *errortools.Error) {
	if serviceConfig == nil {
		return nil, errortools.ErrorMessage("ServiceConfig must not be a nil pointer")
	}

	if serviceConfig.Host == "" {
		return nil, errortools.ErrorMessage("Host not provided")
	}

	if serviceConfig.ServerName == "" {
		return nil, errortools.ErrorMessage("ServerName not provided")
	}

	if serviceConfig.DatabaseName == "" {
		return nil, errortools.ErrorMessage("DatabaseName not provided")
	}

	if serviceConfig.Username == "" {
		return nil, errortools.ErrorMessage("Username not provided")
	}

	if serviceConfig.Password == "" {
		return nil, errortools.ErrorMessage("Password not provided")
	}

	httpServiceConfig := go_http.ServiceConfig{
		HTTPClient: &http.Client{
			Transport: ntlmssp.Negotiator{
				RoundTripper: &http.Transport{},
			},
		},
	}

	return &Service{
		host:         serviceConfig.Host,
		serverName:   serviceConfig.ServerName,
		databaseName: serviceConfig.DatabaseName,
		username:     serviceConfig.Username,
		password:     serviceConfig.Password,
		httpService:  go_http.NewService(&httpServiceConfig),
	}, nil
}

func (service *Service) httpRequest(httpMethod string, requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	// add authentication header
	header := http.Header{}
	header.Set("ServerName", service.serverName)
	header.Set("DatabaseName", service.databaseName)
	header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(service.username+":"+service.password)))
	(*requestConfig).NonDefaultHeaders = &header

	// add error model
	errorResponse := ErrorResponse{}
	(*requestConfig).ErrorModel = &errorResponse

	request, response, e := service.httpService.HTTPRequest(httpMethod, requestConfig)
	if errorResponse.Message != "" {
		e.SetMessage(errorResponse.Message)
	}

	return request, response, e
}

func (service *Service) url(path string) string {
	return fmt.Sprintf("%s/%s/%s", service.host, apiPath, path)
}

func (service *Service) get(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	return service.httpRequest(http.MethodGet, requestConfig)
}

func (service *Service) post(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	return service.httpRequest(http.MethodPost, requestConfig)
}

func (service *Service) put(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	return service.httpRequest(http.MethodPut, requestConfig)
}

func (service *Service) delete(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	return service.httpRequest(http.MethodDelete, requestConfig)
}

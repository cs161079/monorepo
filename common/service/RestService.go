package service

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type RestService interface {
	OasaRequestApi00(action string, extraParams map[string]interface{}) *OasaResponse
	OasaRequestApi02(action string) *OasaResponse
}

type restService struct {
}

const (
	oasaApplicationHost = "http://telematics.oasa.gr"
	testApplicationHost = "http://localhost:8080"
	geoapifyApplication = "https://api.geoapify.com"
)

type OpswHttpRequest struct {
	Method   string
	Headers  map[string]string
	Body     io.Reader
	Endpoint string
}

type OasaResponse struct {
	Error error
	Data  any
}

func checkFields(request *OpswHttpRequest) error {
	if request.Endpoint == "" {
		return fmt.Errorf("REQUEST ENDPOINT IS NOT SET")
	}
	if request.Method == "" {
		return fmt.Errorf("REQUEST HTTP METHOD IS NOT SET")
	}
	return nil
}

func getProperty(v interface{}, property string) any {
	if v != nil {
		if reflect.TypeOf(v).Kind() == reflect.Slice {
			return nil
		} else {
			result := v.(map[string]any)[property]
			return result
		}
	}
	return nil
}

// *********************************************************************************************************************
// ***************** Είναι διαδικασία η οποία δεν κάνει απολύτως τίποτα από θέμα λογικής και ελέγχων *******************
// ***************** Υλοποιεί μόνο ένα HttpRequest  στα πλαίσια της GoLang και επιστρέφει            *******************
// *****************                              *http.Response και error                           *******************
// *********************************************************************************************************************
func httpRequest(request *OpswHttpRequest) (*http.Response, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	if request == nil {
		return nil, fmt.Errorf("REQUEST OBJECT-STRUCT IS NIL OR IS NOT SET CORRECTLY")
	}

	var err = checkFields(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(request.Method, request.Endpoint, request.Body)
	if err != nil {
		return nil, err
	}

	if request.Headers != nil && len(request.Headers) > 0 {
		for key, value := range request.Headers {
			req.Header.Set(key, value)
		}
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r restService) OasaRequestApi02(action string) *OasaResponse {
	var req OpswHttpRequest = OpswHttpRequest{
		Method:   http.MethodGet,
		Endpoint: fmt.Sprintf("%s/api/?act=%s", oasaApplicationHost, action),
		Headers:  map[string]string{"Accept-Encoding": "gzip, deflate"},
	}
	responseByte, err := internalHttpRequest(req)

	if err != nil {
		return &OasaResponse{Error: err, Data: nil}
	}

	// if response.StatusCode == http.StatusInternalServerError {
	// 	return &OasaResponse{Error: fmt.Errorf(models.INTERNALL_SERVER_ERROR), Data: nil}
	// 	// return nil, fmt.Errorf(models.INTERNALL_SERVER_ERROR)
	// }

	reader, err := gzip.NewReader(bytes.NewReader(responseByte))

	defer reader.Close()
	if err != nil {
		return &OasaResponse{Error: err, Data: nil}
		// return nil, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	responseStr := buf.String()
	if responseStr == "" {
		return &OasaResponse{Error: fmt.Errorf("Server response empty."), Data: nil}
	}

	responseStr = responseStr[1 : len(responseStr)-2]
	syncDataLines := strings.Split(responseStr, "),(")

	return &OasaResponse{Data: syncDataLines, Error: err}
}

func (r restService) OasaRequestApi00(action string, extraParams map[string]interface{}) *OasaResponse {
	var oasaResult OasaResponse = OasaResponse{}
	var extraparamUrl string = ""
	for k := range extraParams {
		extraparamUrl = extraparamUrl + fmt.Sprintf("&%s=%v", k, extraParams[k])
	}
	var req OpswHttpRequest = OpswHttpRequest{
		Method:   http.MethodGet,
		Endpoint: fmt.Sprintf("%s/api/?act=%s%s", oasaApplicationHost, action, extraparamUrl),
	}
	//Error Code for error occured in Request Creation
	//var tries = 1
	var resp []byte = nil
	var err error = nil

	resp, err = internalHttpRequest(req)

	if err != nil {
		oasaResult.Error = err
		return &oasaResult
	}

	var tmpResult interface{}
	err = json.Unmarshal(resp, &tmpResult)
	if err != nil {
		oasaResult.Error = fmt.Errorf("AN ERROR OCCURED WHEN CONVERT JSON STRING TO INTERFACE. %s \n %+v", err.Error(), resp)
		return &oasaResult
	}
	hasError := getProperty(tmpResult, "error")
	if hasError != nil {
		oasaResult.Error = fmt.Errorf("SERVER RESPONSES ERROR. %s", hasError)
		return &oasaResult
	}

	oasaResult.Data = tmpResult
	return &oasaResult
}

func internalHttpRequest(req OpswHttpRequest) ([]byte, error) {
	response, err := httpRequest(&req)
	if err != nil {
		return nil, err
	}
	// if transferEncoding := response.Header.Get("Transfer-Encoding"); transferEncoding == "" {
	// 	// Create a map to hold the response body
	// 	var result []map[string]any

	// 	// Decode the JSON response body into the map
	// 	err = json.NewDecoder(response.Body).Decode(&result)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	// if _, ok := result["error"]; ok {
	// 	// 	return nil, fmt.Errorf("Http Error Response [%s]", result["error"])
	// 	// }

	// }
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		var returnedError = fmt.Errorf("AN ERROR OCCURED ANALYZE RESPONSE BODY. %s", err.Error())
		return nil, returnedError
	}
	if response.StatusCode >= http.StatusBadRequest && response.StatusCode <= http.StatusUnavailableForLegalReasons {
		//fmt.Println("Client Error Response from Server")
		var returnedError = fmt.Errorf("%s %s", response.Status, responseBody)
		return nil, returnedError
	}
	if response.StatusCode >= http.StatusInternalServerError && response.StatusCode <= http.StatusNetworkAuthenticationRequired {
		var returnedError = fmt.Errorf("%s %s", response.Status, responseBody)
		//logger.ERROR(string(responseBody))
		return nil, returnedError
	}
	return responseBody, nil
}

func NewRestService() RestService {
	return restService{}
}

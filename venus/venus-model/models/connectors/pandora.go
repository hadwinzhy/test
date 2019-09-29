package connectors

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/parnurzeal/gorequest"
)

var httpClient = &http.Client{}
var httpBaseURL = ""

// ConnectorInit Init connectors
func ConnectorInit(baseURL string) {
	httpBaseURL = baseURL
}

// HTTPRequest is for http request
func HTTPRequest(method, url string, body io.Reader) (response http.Response, err error) {
	request, err := http.NewRequest(method, httpBaseURL+url, body)

	if strings.ToUpper(method) == "POST" {
		request.Header.Set("Content-type", "application/json")
	}

	if err != nil {
		return
	}

	responsePtr, err := httpClient.Do(request)

	if responsePtr != nil {
		response = *responsePtr
	}

	return
}

func CreateDevicePack(uuid string, delta uint, hostName string) error {
	request := gorequest.New()
	sendString := `{"device_pack_uuid": "%s","delta": %d,"host_name": "%s"}`
	response, body, _ := request.Post(httpBaseURL + "/v1/api/device_packs").
		Send(fmt.Sprintf(sendString, uuid, 24, hostName)).
		End()

	if response.StatusCode == http.StatusOK {
		return nil
	}
	return errors.New("request pandora error. " + strconv.Itoa(response.StatusCode) + body)
}

package apiclient

import (
	"appstax-cli/appstax/config"
	"appstax-cli/appstax/fail"
	"appstax-cli/appstax/log"
	"appstax-cli/appstax/session"
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"os"
)

var defaultBaseUrl = "https://appstax.com/api/latest"
var configBaseUrl = ""
var optionBaseUrl = ""

func SetBaseUrl(url string) {
	optionBaseUrl = url
}

func selectBaseUrl() string {
	url := defaultBaseUrl
	configBaseUrl = config.Read().ApiBaseUrl
	if configBaseUrl != "" {
		url = configBaseUrl
	}
	if optionBaseUrl != "" {
		url = optionBaseUrl
	}
	log.Infof("Using base url: %s", url)
	return url
}

func Url(url string, params ...string) string {
	encodedParams := make([]interface{}, len(params))
	for i, p := range params {
		encodedParams[i] = neturl.QueryEscape(p)
	}
	url = fmt.Sprintf(url, encodedParams...)
	return selectBaseUrl() + url
}

func PostFile(url string, path string, progressWriter io.Writer) ([]byte, *http.Response, error) {
	log.Debugf("HTTP POST FILE: %s", url)
	file, err := os.Open(path)
	fail.Handle(err)
	defer file.Close()

	fileReader := bufio.NewReader(file)
	teeReader := io.TeeReader(fileReader, progressWriter)

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, teeReader)
	addHeaders(req)

	resp, err := client.Do(req)
	fail.Handle(err)
	return handleResult(resp, err)
}

func Post(url string, data interface{}) ([]byte, *http.Response, error) {
	log.Debugf("HTTP POST: %s", url)
	json, err := json.Marshal(data)
	fail.Handle(err)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json))
	addHeaders(req)
	resp, err := client.Do(req)
	fail.Handle(err)
	return handleResult(resp, err)
}

func Put(url string, data interface{}) ([]byte, *http.Response, error) {
	log.Debugf("HTTP PUT: %s", url)
	json, err := json.Marshal(data)
	fail.Handle(err)
	log.Debugf("HTTP PUT JSON: %s", json)
	client := &http.Client{}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(json))
	addHeaders(req)
	resp, err := client.Do(req)
	fail.Handle(err)
	return handleResult(resp, err)
}

func Get(url string) ([]byte, *http.Response, error) {
	log.Debugf("HTTP GET: %s", url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	addHeaders(req)
	resp, err := client.Do(req)
	fail.Handle(err)
	return handleResult(resp, err)
}

func ParseStringMap(data []byte) map[string]string {
	var result map[string]string
	json.Unmarshal(data, &result)
	return result
}

func ParseMap(data []byte) map[string]interface{} {
	var result map[string]interface{}
	json.Unmarshal(data, &result)
	return result
}

func addHeaders(req *http.Request) {
	sessionID := session.ReadSessionID()
	appKey := config.Read().AppKey
	
	if sessionID != "" {
		req.Header.Add("x-appstax-sessionid", sessionID)
	}
	if appKey != "" {
		req.Header.Add("x-appstax-appkey", appKey)
	}
}

func handleResult(resp *http.Response, err error) ([]byte, *http.Response, error) {
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("HTTP response (%d): %s", resp.StatusCode, body)
	if err == nil && resp.StatusCode/100 == 2 {
		return body, resp, nil
	} else {
		return nil, nil, errors.New(getErrorMessage(resp, body, err))
	}
}

func getErrorMessage(resp *http.Response, body []byte, err error) string {
	message := ParseStringMap(body)["errorMessage"]

	if err != nil {
		message = err.Error()
	} else if message == "" && resp.StatusCode == 401 {
		switch resp.StatusCode {
		case 401:
			message = "Not authorized"
		}
	}

	if message == "" {
		message = "Ooops! Error communicating with appstax server."
	}

	return message
}

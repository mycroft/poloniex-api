package poloniexapi

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (api *PoloniexApi) query(url string, params url.Values, with_signature bool) ([]byte, error) {
	headers := map[string]string{}
	method := "GET"

	if with_signature {
		params.Set("nonce", fmt.Sprintf("%d", time.Now().UnixNano()/1000))

		signature := createPoloniexSignature(params, api.secret)

		headers["Key"] = api.Key
		headers["Sign"] = signature

		method = "POST"
	}

	if api.UserAgent != "" {
		headers["User-Agent"] = api.UserAgent
	}

	headers["Content-Type"] = "application/x-www-form-urlencoded"

	return executeHttpQuery(method, url, headers, params)
}

func (api *PoloniexApi) parse(resp []byte, out interface{}) (interface{}, error) {
	var response interface{}

	if nil != out {
		response = out
	}

	err := json.Unmarshal(resp, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (api *PoloniexApi) queryparse(url string, params url.Values, with_signature bool, out interface{}) (interface{}, error) {
	var response interface{}

	resp, err := api.query(url, params, with_signature)
	if err != nil {
		return nil, err
	}

	if out != nil {
		response = out
	}

	fmt.Println("Response:")
	fmt.Println(string(resp))

	_, err = api.parse(resp, &response)
	if err != nil {
		return nil, err
	}

	return response, err
}

func executeHttpQuery(method string, url string, headers map[string]string, values url.Values) ([]byte, error) {
	var bodyReader io.Reader

	client := &http.Client{}

	if method == "GET" {
		bodyReader = nil
		url = url + "?" + values.Encode()
	} else {
		bodyReader = strings.NewReader(values.Encode())
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("Could not execute request! (%s)", err.Error())
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Could not execute request! (%s)", err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Could not execute request! (%s)", err.Error())
	}

	return body, nil
}

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	gateway string
	http    *http.Client
}

func NewClient(gateway string) *Client {
	return &Client{
		gateway: gateway,
		http:    &http.Client{},
	}
}

func (this *Client) Post(endPoint string, req interface{}) ([]byte, error) {
	buf, err := RequestBuffer(req)
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("%s%s", this.gateway, endPoint)
	request, err := http.NewRequest("POST", uri, buf)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-type", "application/json")
	response, err := this.http.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	respBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	err = checkErr(respBytes)
	if err != nil {
		return nil, err
	}
	return respBytes, nil
}

func (this *Client) Get(endPoint string) ([]byte, error) {
	uri := fmt.Sprintf("%s%s", this.gateway, endPoint)
	request, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Connection", "keep-alive")
	response, err := this.http.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	respBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	err = checkErr(respBytes)
	if err != nil {
		return nil, err
	}
	return respBytes, nil
}

func checkErr(js []byte) error {
	var e Error
	err := json.Unmarshal(js, &e)
	if err != nil {
		return err
	}
	if e.Code != 0 {
		return e
	}
	return nil
}

func RequestBuffer(req interface{}) (*bytes.Reader, error) {
	js, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(js)
	return reader, nil
}

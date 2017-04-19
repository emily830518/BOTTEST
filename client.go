package main

import (
	"io/ioutil"
	"net/http"
)

type client struct {
	url string
}

func NewClient(url string) *client {
	c := new(client)
	c.url = url
	return c
}

func (c *client) GetHttpRes() ([]byte, error) {
	return getHttpResponse(c.url)
}

func getHttpResponse(url string) ([]byte, error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}
	return body, nil
}
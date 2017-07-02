package util

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/context/ctxhttp"
)

func Request(method string, url string, bodyType string, body io.Reader, header http.Header, timeout int) ([]byte, error) {
	var resp *http.Response
	var err error
	Umethod := strings.ToUpper(method)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	switch Umethod {
	case "GET":
		resp, err = ctxhttp.Get(ctx, nil, url)
		if err != nil {
			return nil, err
		}
	case "POST":
		resp, err = ctxhttp.Post(ctx, nil, url, bodyType, body)
		if err != nil {
			return nil, err
		}
	case "DELETE":
		client := &http.Client{}
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			return nil, err
		}
		resp, err = client.Do(req)
		if err != nil {
			return nil, err
		}
	case "HEAD": //用于设置head
		client := &http.Client{}
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header = header
		resp, err = client.Do(req)
		if err != nil {
			return nil, err
		}
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

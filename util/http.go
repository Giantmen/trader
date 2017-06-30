package proto

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/context/ctxhttp"
)

func Request(method string, url string, bodyType string, body io.Reader, timeout int) ([]byte, error) {
	var resp *http.Response
	var err error
	Umethod := strings.ToUpper(method)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	if Umethod == "GET" {
		resp, err = ctxhttp.Get(ctx, nil, url)
		if err != nil {
			return nil, err
		}
	} else if Umethod == "POST" {
		resp, err = ctxhttp.Post(ctx, nil, url, bodyType, body)
		if err != nil {
			return nil, err
		}
	} else if Umethod == "DELETE" {
		client := &http.Client{}
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			return nil, err
		}
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

package util

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/context/ctxhttp"
	"fmt"
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
		if header != nil { //用于设置head
			client := &http.Client{}
			req, err := http.NewRequest("POST", url, body)
			if err != nil {
				return nil, err
			}
			req.Header = header
			req.Header.Set("Content-Type", bodyType)
			resp, err = client.Do(req)
			if err != nil {
				return nil, err
			}
		} else {
			resp, err = ctxhttp.Post(ctx, nil, url, bodyType, body)
			if err != nil {
				return nil, err
			}
		}

	case "DELETE":
		client := &http.Client{}
		req, err := http.NewRequest("DELETE", url, body)
		if err != nil {
			return nil, err
		}
		req.Header=header
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

func HttpGet(url string,header *http.Header) ([]byte,error) {
	client := http.Client{}
	req,err := http.NewRequest("GET",url,nil)
	if err!= nil {
		fmt.Println(err)
	}
	if header != nil {
		req.Header = *header
	}
	//req.Header.Add("X-MBX-APIKEY","XfZyiOPMQRtVAYl8rgQPRIrpB5RAxsHPkgZvdR1lURP3qbZ7zN0claHSGW1yQdt0")
	resp,err := client.Do(req)
	if err !=nil {
		return nil,err
	}
	b,err := ioutil.ReadAll(resp.Body)
	if err !=nil {
		return nil,err
	}
	//fmt.Println(string(b))
	return b,nil
}

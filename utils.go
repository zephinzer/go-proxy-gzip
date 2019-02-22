package main

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"net/http"
)

func getRequestBody(request *http.Request) []byte {
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		panic(err)
	}
	return requestBody
}

func getResponseBody(response *http.Response) []byte {
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	return responseBody
}

func gzipEncode(content []byte) []byte {
	var gzippedContent bytes.Buffer
	gzipper := gzip.NewWriter(&gzippedContent)
	gzipper.Write(content)
	gzipper.Close()
	return gzippedContent.Bytes()
}

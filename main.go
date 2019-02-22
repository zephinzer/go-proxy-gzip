package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/sirupsen/logrus"
	config "github.com/spf13/viper"
)

func main() {
	initConfiguration()

	logger := createLogger("main", config.GetString("log_format"))
	logger.Debugf("ADDR         : %s", config.GetString("addr"))
	logger.Debugf("PORT         : %s", config.GetString("port"))
	logger.Debugf("FORWARD_TO   : %s", config.GetString("forward_to"))
	logger.Debugf("CONTENT_TYPE : %s", config.GetString("content_type"))
	logger.Debugf("LOG_FORMAT   : %s", config.GetString("log_format"))

	logger.Debug("initialising proxy forwarder...")
	proxyGzip := Init(&InitConfiguration{
		address:      config.GetString("addr"),
		contentType:  config.GetString("content_type"),
		forwardTo:    config.GetString("forward_to"),
		logFormatter: config.GetString("log_format"),
		port:         config.GetString("port"),
	})

	logger.Info("starting proxy forwarder server...")
	proxyGzip.Listen()
}

// InitConfiguration exists to shorten Init() parameters
type InitConfiguration struct {
	// address to listen on (eg. 0.0.0.0)
	address string
	// enforces a content type if http.detectContentType cannot figure it out
	contentType string
	// endpoint for the proxy
	forwardTo string
	// text/json for development/production
	logFormatter string
	// port to listen on (eg. 8888)
	port string
}

// Init initialises the ProxyGzip server
func Init(config *InitConfiguration) *ProxyGzip {
	logger := createLogger("ProxyGzip", config.logFormatter)
	proxyGzip := &ProxyGzip{
		config: config,
		logger: logger,
	}
	proxyGzip.setupMux()
	return proxyGzip
}

// ProxyGzip server
type ProxyGzip struct {
	config *InitConfiguration
	mux    *http.ServeMux
	logger *logrus.Entry
}

// Listen exposes the server starting method
func (pgz *ProxyGzip) Listen() {
	listenAddress := fmt.Sprintf("%s:%s", pgz.config.address, pgz.config.port)
	pgz.logger.Infof("listening on %s", listenAddress)
	err := http.ListenAndServe(listenAddress, pgz.mux)
	if err != nil {
		pgz.logger.Panicln(err)
	}
}

// createProxiedRequestFromIncomingRequest
func (pgz *ProxyGzip) createProxiedRequestFromIncomingRequest(incomingRequest *http.Request, body []byte) *http.Request {
	gzippedBody := gzipEncode(body)
	proxiedRequest, err := http.NewRequest(
		incomingRequest.Method,
		fmt.Sprintf("%s%s", pgz.config.forwardTo, incomingRequest.URL.Path),
		bytes.NewReader(gzippedBody),
	)
	if err != nil {
		pgz.logger.Panic(err)
	}
	contentType := pgz.config.contentType
	if len(contentType) == 0 {
		contentType = http.DetectContentType(body)
	}
	for headerKey, headerValues := range incomingRequest.Header {
		for _, headerValue := range headerValues {
			proxiedRequest.Header.Add(headerKey, headerValue)
		}
	}
	contentLength := strconv.FormatUint(uint64(len(gzippedBody)), 10)
	proxiedRequest.Header.Add("Content-Type", contentType)
	proxiedRequest.Header.Add("Content-Length", contentLength)
	proxiedRequest.Header.Add("Content-Encoding", "gzip")
	return proxiedRequest
}

// forwardRequest is here for calling to forward the incoming request to the next-hop server
func (pgz *ProxyGzip) forwardRequest(r *http.Request, body []byte) *http.Response {
	httpClient := &http.Client{}
	request := pgz.createProxiedRequestFromIncomingRequest(r, body)
	response, err := httpClient.Do(request)
	if err != nil {
		panic(err)
	}
	return response
}

// isRequestForwardingEnabled encapsulates logic for checking whether we should forward the request
// or just be a normal ol' echoserver (nothing wrong with that!)
func (pgz *ProxyGzip) isRequestForwardingEnabled() bool {
	return len(pgz.config.forwardTo) > 0
}

// logIncomingRequest is here to make your life easier by displaying the incoming request
// in the logs
func (pgz *ProxyGzip) logIncomingRequest(request *http.Request, requestID string, requestBody []byte) {
	requestBodyHeader := string(requestBody)
	if len(requestBody) > 100 {
		requestBodyHeader = requestBodyHeader[100:] + "... [truncated]"
	}
	fields := map[string]interface{}{
		"requestID":     requestID,
		"requestBody":   string(requestBody),
		"requestHost":   request.Host,
		"requestMethod": request.Method,
		"requestPath":   request.URL.EscapedPath(),
	}
	for headerKey, headerValues := range request.Header {
		for _, headerValue := range headerValues {
			fields[headerKey] = headerValue
		}
	}
	pgz.logger.WithFields(fields).Info("incoming request")
}

// setupMux initialises the server routing component
func (pgz *ProxyGzip) setupMux() {
	pgz.mux = http.NewServeMux()
	pgz.mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				err := fmt.Sprintf("%v", r)
				pgz.logger.Error(fmt.Sprintf("err:%s", err))
				writer.WriteHeader(500)
				writer.Write([]byte(err))
			}
		}()

		requestID := uuid.New().String()
		requestBody := getRequestBody(request)
		pgz.logIncomingRequest(request, requestID, requestBody)
		pgz.setConfigResponseHeaders(writer)
		pgz.setIncomingRequestResponseHeaders(writer, request, requestID, requestBody)
		if pgz.isRequestForwardingEnabled() {
			response := pgz.forwardRequest(request, requestBody)
			defer response.Body.Close()
			responseBody := getResponseBody(response)
			pgz.logger.Infof("finished handling request for request with id '%s'", requestID)
			pgz.setOutgoingResponseHeaders(writer, response, requestID)
			writer.Write(responseBody)
		} else {
			pgz.logger.Infof("finished handling request for request with id '%s'", requestID)
			writer.Write(requestBody)
		}
	})
}

func (pgz *ProxyGzip) setConfigResponseHeaders(w http.ResponseWriter) {
	w.Header().Set("PGZ-Config-Forward-To", pgz.config.forwardTo)
	w.Header().Set("PGZ-Config-Content-Type", pgz.config.contentType)
	w.Header().Set("PGZ-Config-Address", pgz.config.address)
	w.Header().Set("PGZ-Config-Port", pgz.config.port)
}

// setIncomingRequestResponseHeaders is here to put meta-data from the incoming request into
// the response headers so that we can view what was being sent to it - useful for debugging
func (pgz *ProxyGzip) setIncomingRequestResponseHeaders(w http.ResponseWriter, request *http.Request, requestID string, requestBody []byte) {
	requestBodyHeader := string(requestBody)
	if len(requestBody) > 100 {
		requestBodyHeader = requestBodyHeader[0:256] + "... [truncated]"
	}
	w.Header().Set("PGZ-Request-ID", requestID)
	w.Header().Set("PGZ-Request-Body", requestBodyHeader)
	w.Header().Set("PGZ-Request-Host", request.Host)
	w.Header().Set("PGZ-Request-Method", request.Method)
	w.Header().Set("PGZ-Request-Path", request.URL.EscapedPath())
	for headerKey, headerValues := range request.Header {
		for index, headerValue := range headerValues {
			w.Header().Set(fmt.Sprintf("PGZ-Request-Header-%v-%s", index, headerKey), headerValue)
		}
	}
}

// setOutgoingResponseHeaders is here to format the response back to the requesting client
// given the response from the next-hop server
func (pgz *ProxyGzip) setOutgoingResponseHeaders(writer http.ResponseWriter, responseFromNextHop *http.Response, requestID string) {
	for headerKey, headerValues := range responseFromNextHop.Header {
		for _, headerValue := range headerValues {
			writer.Header().Set(fmt.Sprintf("%s", headerKey), headerValue)
		}
	}
	writer.WriteHeader(responseFromNextHop.StatusCode)

}

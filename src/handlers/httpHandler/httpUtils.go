package httpHandler

import (
	"io"
	"net/http"
	"net/http/httputil"
	"proxy/src/database"
	"proxy/src/handlers"
	"strings"
)

func (handler *HttpHandler) doRequest() error {
	// Creating HTTP client to resend original request
	httpClient := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// Copying original request
	proxyRequest, err := http.NewRequest(handler.clientRequest.Method, handler.clientRequest.RequestURI, handler.clientRequest.Body)
	if err != nil {
		return err
	}

	for header, values := range handler.clientRequest.Header {
		for _, value := range values {
			var requestHeaderKey = strings.ToLower(header)
			if (requestHeaderKey != "proxy-connection") {
				proxyRequest.Header.Add(requestHeaderKey, value)
			}
		}
	}

	// Dumping request to database
	requestDump, err := httputil.DumpRequest(proxyRequest, true)
	if err != nil {
		return err
	}

	requestModel := database.Request{
		Host: handler.clientRequest.RequestURI,
		Request: string(requestDump),
	}
	handler.dbConn.AddRequest(requestModel)

	// Send request
	handler.proxyResponse, err = httpClient.Do(proxyRequest)
	if err != nil {
		return err
	}

	return nil
}

func (handler *HttpHandler) returnResponse() error {
	handlers.CopyHttpHeaders(handler.proxyResponse.Header, handler.responseWriter.Header())
	handler.responseWriter.WriteHeader(handler.proxyResponse.StatusCode)
	handler.responseWriter.Header().Add("Connection", handler.proxyResponse.Header.Get("Connection"))

	_, err := io.Copy(handler.responseWriter, handler.proxyResponse.Body)
	if err != nil {
		return err
	}

	return nil
}

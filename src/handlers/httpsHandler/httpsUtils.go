package httpsHandler

import (
	"bufio"
	"net/http"
	"net/http/httputil"
	"proxy/src/database"
)

func (handler *HttpsHandler) getRequest() error {
	reader := bufio.NewReader(handler.clientConnection)
	request, err := http.ReadRequest(reader)
	if err != nil {
		return err
	}

	handler.clientRequest = request

	return nil
}

func (handler *HttpsHandler) doRequest() error {
	// Saving request to database
	requestDump, err := httputil.DumpRequest(handler.clientRequest, true)
	if err != nil {
		return err
	}

	requestModel := database.Request{
		Host: handler.clientRequest.RequestURI,
		Request: string(requestDump),
	}
	handler.dbConn.AddRequest(requestModel)

	// sending request to server
	_, err = handler.serverConnection.Write(requestDump)
	if err != nil {
		return err
	}

	// Getting response from server
	writer := bufio.NewReader(handler.serverConnection)
	response, err := http.ReadResponse(writer, handler.clientRequest)
	if err != nil {
		return err
	}

	handler.proxyResponse = response

	return nil
}

func (handler *HttpsHandler) returnResponse() error {
	rawResp, err := httputil.DumpResponse(handler.proxyResponse, true)
	_, err = handler.clientConnection.Write(rawResp)
	if err != nil {
		return err
	}

	return nil
}
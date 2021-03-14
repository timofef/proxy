package httpHandler

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"proxy/src/database"
)

type HttpHandler struct {
	responseWriter http.ResponseWriter
	clientRequest  *http.Request
	proxyResponse  *http.Response
	dbConn         *database.Database
}

func NewHttpHandler(responseWriter http.ResponseWriter, clientRequest *http.Request, db *database.Database) *HttpHandler {
	return &HttpHandler{
		responseWriter: responseWriter,
		clientRequest:  clientRequest,
		dbConn:         db,
	}
}

func (handler *HttpHandler) ProxyRequest() error {
	err:= handler.doRequest()
	if err != nil {
		logrus.Error(err)
	}

	err = handler.returnResponse()
	if err != nil {
		logrus.Error(err)
	}

	return nil
}

func (handler *HttpHandler) Defer() {
	handler.dbConn.CloseConnection()
	handler.proxyResponse.Body.Close()
}

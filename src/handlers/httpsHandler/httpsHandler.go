package httpsHandler

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"proxy/src/database"
)

type HttpsHandler struct {
	responseWriter http.ResponseWriter
	clientRequest  *http.Request
	proxyResponse  *http.Response
	dbConn         *database.Database

	tlsConfig        *tls.Config
	url              *url.URL
	connectRequest   *http.Request
	clientConnection net.Conn
	serverConnection *tls.Conn
}

func NewHttpsHandler(responseWriter http.ResponseWriter, connectRequest *http.Request, db *database.Database) (*HttpsHandler, error) {
	handler := &HttpsHandler{}

	handler.responseWriter = responseWriter
	handler.connectRequest = connectRequest
	handler.dbConn = db

	var err error

	handler.url, err = url.Parse(connectRequest.RequestURI)
	if err != nil {
		return nil, err
	}

	err = handler.setupHttps()
	if err != nil {
		return nil, err
	}

	err = handler.setupClientConnection()
	if err != nil {
		return nil, err
	}

	err = handler.setupServerConnection()
	if err != nil {
		return nil, err
	}

	return handler, nil
}

func (handler *HttpsHandler) ProxyRequest() error {
	err := handler.getRequest()
	if err != nil {
		return err
	}

	err = handler.doRequest()
	if err != nil {
		return err
	}

	err = handler.returnResponse()
	if err != nil {
		return err
	}

	return nil
}

func (handler *HttpsHandler) Defer() {
	handler.serverConnection.Close()
	handler.clientConnection.Close()
	handler.dbConn.CloseConnection()
	handler.proxyResponse.Body.Close()
}



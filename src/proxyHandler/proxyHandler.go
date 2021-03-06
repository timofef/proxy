package proxyHandler

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"proxy/src/database"
	"proxy/src/handlers"
	"proxy/src/handlers/httpHandler"
)

func Serve(responseWriter http.ResponseWriter, request *http.Request) {
	logrus.Info("Request: " + request.RequestURI)

	db, err := database.InitConnection()
	if err != nil {
		logrus.Warn("Can't connect to database")
		logrus.Fatal(err)
	}

	var handler handlers.HandlerInterface

	/*if request.Method == http.MethodConnect {
		handler, err = httpsHandler.NewHttpsHandler(responseWriter, request, db)
		logrus.Error(err)
	} else {
		handler = httpHandler.NewHttpHandler(responseWriter, request, db)
	}*/

	handler = httpHandler.NewHttpHandler(responseWriter, request, db)

	err = handler.ProxyRequest()
	if err != nil {
		logrus.Error(err)
	}

	defer handler.Defer()
}

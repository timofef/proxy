package repeatRequestHandlers

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"proxy/src/database"
)

func ShowRequest(responseWriter http.ResponseWriter, request *http.Request) {
	dbRequest, err := getRequestFromDatabase(responseWriter, request)

	if err != nil {
		logrus.Error(err)
		return
	}

	if (dbRequest == database.Request{}) {
		return
	}

	fmt.Fprintf(responseWriter, "\n----------- Request %d -----------\n\n\n", dbRequest.Id)
	_, _ = fmt.Fprintf(responseWriter,
		"Host: %s\n\n%s"+"----------- Repeat request: http://127.0.0.1/repeat?id=%d\n\n\n",
		dbRequest.Host,
		dbRequest.Request,
		dbRequest.Id)
}

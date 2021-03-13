package repeatRequestHandlers

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"proxy/src/database"
	"strconv"
)

func ShowRequest(responseWriter http.ResponseWriter, request *http.Request) {
	db, err := database.InitConnection()
	if err != nil {
		logrus.Warn("Can't connect to database")
		logrus.Error(err)
	}
	defer db.CloseConnection()

	requestId := request.URL.Query()["id"]
	if len(requestId) < 1 {
		_, _ = fmt.Fprintf(responseWriter, "To repeat request you need to enter id parameter\n"+
			"To see all requests visit:        http://127.0.0.1:8082\n")
		return
	}

	id, err := strconv.Atoi(requestId[0])
	if err != nil {
		logrus.Error(err)
	}

	dbRequest, err := db.GetRequest(id)
	if err != nil {
		logrus.Error(err)
	}

	if (dbRequest == database.Request{}) {
		_, _ = fmt.Fprintf(responseWriter, "No request with id = %d\n", id)
		return
	}

	fmt.Fprintf(responseWriter, "\n----------- Request %d -----------\n\n\n", dbRequest.Id)
	_, _ = fmt.Fprintf(responseWriter,
		"Host: %s\n\n%s"+"----------- Repeat request: http://127.0.0.1/repeat?id=%d\n\n\n",
		dbRequest.Host,
		dbRequest.Request,
		dbRequest.Id)
}

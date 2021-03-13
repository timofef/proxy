package repeatRequestHandlers

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"proxy/src/database"
)

func ListAllRequests(responseWriter http.ResponseWriter, _ *http.Request) {
	db, err := database.InitConnection()
	if err != nil {
		logrus.Warn("Can't connect to database")
		logrus.Error(err)
	}
	defer db.CloseConnection()

	list, err := db.GetRequestList()
	if err != nil {
		logrus.Warn("Can't get data from DB")
		_, _ = fmt.Fprintf(responseWriter, "Can't get info\n")
		return
	}

	if len(list) == 0 {
		_, _ = fmt.Fprintf(responseWriter, "No requests yet\n")
		return
	}

	fmt.Fprintf(responseWriter, "\n----------- Saved requests -----------\n\n\n")
	for i, request := range list {
		_, _ = fmt.Fprintf(responseWriter,
			"%d) Host: %s\n\n%s"+"----------- Repeat request: http://127.0.0.1/repeat?id=%d\n\n\n",
			i+1,
			request.Host,
			request.Request,
			request.Id)
	}
}

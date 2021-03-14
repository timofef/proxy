package repeatRequestHandlers

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"proxy/src/database"
	"proxy/src/handlers"
	"strconv"
	"strings"
)

func getRequestFromDatabase(responseWriter http.ResponseWriter, request *http.Request) (database.Request, error){
	db, err := database.InitConnection()
	if err != nil {
		logrus.Warn("Can't connect to database")
		return database.Request{}, err
	}
	defer db.CloseConnection()

	requestId := request.URL.Query()["id"]
	if len(requestId) < 1 {
		_, _ = fmt.Fprintf(responseWriter, "To repeat request you need to enter id parameter\n"+
			"To see all requests visit:        http://127.0.0.1:8082\n")
		return database.Request{}, err
	}

	id, err := strconv.Atoi(requestId[0])
	if err != nil {
		return database.Request{}, err
	}

	dbRequest, err := db.GetRequest(id)

	if (dbRequest == database.Request{}) {
		_, _ = fmt.Fprintf(responseWriter, "No request with id = %d\n", id)
		return database.Request{}, err
	}

	return dbRequest, nil
}

func formRequest(dbRequest database.Request) (*http.Request, error) {
	reqReader := bufio.NewReader(strings.NewReader(dbRequest.Request))
	buffer, err := http.ReadRequest(reqReader)
	if err != nil {
		return &http.Request{}, err
	}

	httpReq, err := http.NewRequest(buffer.Method, dbRequest.Host, buffer.Body)
	if err != nil {
		return &http.Request{}, err
	}

	handlers.CopyHttpHeaders(buffer.Header, httpReq.Header)

	return httpReq, nil
}
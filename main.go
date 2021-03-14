package main

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"proxy/src/proxyHandler"
	"proxy/src/repeatRequestHandlers"
)

func init() {
	logrus.SetLevel(logrus.InfoLevel)

	// Generate root certificate
	/*pwd, err := os.Getwd()
	if err != nil {
		logrus.Error(err)
	}

	genCmd := exec.Command(pwd+"/cert")
	_, err = genCmd.CombinedOutput()
	if err != nil {
		logrus.Error(err)
	}*/
}

func main() {

	// Proxy server config
	server := &http.Server{
		Addr:    ":8081",
		Handler: http.HandlerFunc(proxyHandler.Serve),
	}

	// Repeat server config
	http.HandleFunc("/", repeatRequestHandlers.ListAllRequests)
	http.HandleFunc("/request", repeatRequestHandlers.ShowRequest)
	http.HandleFunc("/repeat", repeatRequestHandlers.RepeatRequest)
	http.HandleFunc("/scan", repeatRequestHandlers.ScanRequestForXXE)

	// Start repeat server
	go http.ListenAndServe(":8082", nil)

	// Start proxy server
	logrus.Fatal(server.ListenAndServe())
}

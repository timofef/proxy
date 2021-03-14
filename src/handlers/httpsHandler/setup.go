package httpsHandler

import (
	"crypto/tls"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

func (handler *HttpsHandler) setupHttps() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	scriptDir := pwd + "/cert"
	certsDir := scriptDir + "/certs/"

	certFilename := certsDir + handler.url.Scheme + ".crt"

	// Generating client certificate if it does not exist yet
	_, err = os.Stat(certFilename)
	if os.IsNotExist(err) {
		err = genCert(scriptDir, handler.url.Scheme, certsDir)
		if err != nil {
			return err
		}
	}

	cert, err := tls.LoadX509KeyPair(certFilename, scriptDir+"/cert.key")
	if err != nil {
		return err
	}

	config := new(tls.Config)
	config.Certificates = []tls.Certificate{cert}
	config.ServerName = handler.url.Scheme

	handler.tlsConfig = config

	return nil
}

func genCert(scriptPath, host, savePath string) error {
	genCmd := exec.Command(scriptPath+"/gen_cert.sh", host, scriptPath, strconv.Itoa(rand.Intn(10000)), savePath)
	_, err := genCmd.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}

func (handler *HttpsHandler) setupClientConnection() error {
	// Hijacking client request
	raw, _, err := handler.responseWriter.(http.Hijacker).Hijack()
	if err != nil {
		return err
	}

	// Immediately returning 200
	_, err = raw.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
	if err != nil {
		raw.Close()
		return err
	}

	clientConnection := tls.Server(raw, handler.tlsConfig)
	err = clientConnection.Handshake()
	if err != nil {
		clientConnection.Close()
		raw.Close()
		return err
	}

	handler.clientConnection = clientConnection

	return nil
}

func (handler *HttpsHandler) setupServerConnection() error {
	serverConnection, err:= tls.Dial("tcp", handler.connectRequest.Host, handler.tlsConfig)
	if err != nil {
		return err
	}

	handler.serverConnection = serverConnection

	return nil
}

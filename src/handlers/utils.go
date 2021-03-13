package handlers

import "net/http"

func CopyHttpHeaders(from http.Header, to http.Header) {
	for header, values := range from {
		for _, value := range values {
			to.Add(header, value)
		}
	}
}
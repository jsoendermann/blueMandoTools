package main

import (
	"fmt"
	"net/http"
	"strings"
)

func getLastPathComponent(request *http.Request) string {
	// get the word or sentence from the path
	pathElements := strings.Split(request.URL.Path, "/")
	requestData := pathElements[len(pathElements)-1]

	return requestData
}

func getColors(request *http.Request) []string {
	colors := make([]string, 5)
	for i := 0; i <= 4; i++ {
		colors[i] = request.FormValue(fmt.Sprintf("tone%d", i))
	}

	return colors
}

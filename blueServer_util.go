package main

import (
	"fmt"
	"net/http"
	"strings"
)

// getRequestDataAndColors gets the word or sentence from the url
// and the color array from the post data in the request
func getRequestDataAndColors(r *http.Request) (string, []string) {
	// get the word or sentence from the path
	pathElements := strings.Split(r.URL.Path, "/")
	requestData := pathElements[len(pathElements)-1]

	// get colors
	colors := make([]string, 5)
	for i := 0; i <= 4; i++ {
		colors[i] = r.FormValue(fmt.Sprintf("tone%d", i))
	}

	return requestData, colors
}

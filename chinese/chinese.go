// Package chinese defines constants and other useful stuff that other
// chinese-related packages rely on.
package chinese

import (
    "fmt"
    "net/http"
    "strings"
    "io/ioutil"
    "encoding/json"
)

type CharSet int

const (
    Trad = iota
    Simp
)

func (cs CharSet) String() string {
    if cs == Trad {
        return "trad"
    } else if cs == Simp {
        return "simp"
    }
    return "UNKNOWN CHARACTER SET"
}

// Different types of dictionary records can implement this interface
// to provide a way to convert the data they hold to html for use on
// a card. The colors must be in the form "#FFFFFF"
type ToHTMLer interface {
  ToHTML(toneColors []string) string
}

// TODO Maybe add a ToCarder interface that converts records to 
// csv rows.


func SplitChineseTextIntoWords(text string, charSet CharSet) ([]string, error) {
    charSetString := charSet.String()

    requestData := `{"charset": "` + charSetString + `", "text": "` + strings.Replace(text, `"`, `\"`, -1) + `"}`

    post_data := strings.NewReader(requestData)
    response, err := http.Post("http://split.zaoyin.eu/split", "application/json", post_data)
    if err != nil {
        return nil, err
    }
    defer response.Body.Close()

    responseString, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return nil, err
    }

    type SplitResponse struct {
        Result string
        Split_text []string
    }

    var sr SplitResponse
    err = json.Unmarshal(responseString, &sr)
    if err != nil {
        return nil, err
    }

    if sr.Result != "ok" {
        return nil, fmt.Errorf("Result: " + sr.Result)
    }

    return sr.Split_text, nil
}

/*
TODO package description
*/
package moedict

import (
	"fmt"
	// "github.com/yangchuanzhang/chinese"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// These three structs reflect the json of the api at https://www.moedict.tw/uni/
// which is documented (in Chinese) here: https://hackpad.com/3du.tw-API-Design-95jKjray8uR
// Fields that are commented out are not used at the moment.
type Entry struct {
	Title string
	//	Radical                  string
	//	Stroke_count             int
	//	Non_radical_stroke_count int

	Heteronyms []Heteronym
}

type Heteronym struct {
	Pinyin   string
	Bopomofo string
	//	Bopomofo2 string

	Definitions []Definition
}

type Definition struct {
	Def     string
	Quote   []string
	Example []string
	//	DefType  string `json:"type"` // this field is called "type" in the output of the server
	//	Link     string
	Synonyms string
}


// FindEntry queries http://www.moedict.tw/uni/ and loads the data into a variable of type Entry
// a pointer to which it returns. This method blocks during the http request.
func FindEntry(word string) (*Entry, error) {
	// make http request, check for errors and defer close
	resp, err := http.Get("http://www.moedict.tw/uni/" + word)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// read data into variable
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// unmarshal json into e
	var e Entry

	jsonErr := json.Unmarshal(body, &e)
	if jsonErr != nil {
		return nil, jsonErr
	}

  fmt.Println(e.Title)
  if e.Title == "" { // Word could not be found in dict
    return nil, fmt.Errorf("Word doesn't exist in dictionary")
  }

	return &e, nil
}

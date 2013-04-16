package main

//package moedict

import (
	"fmt"
	// "github.com/yangchuanzhang/chinese"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// these three structs reflect the json of the api at https://www.moedict.tw/uni/
// which is documented (in Chinese) here: https://hackpad.com/3du.tw-API-Design-95jKjray8uR
type Record struct {
	title                    string
	radical                  string
	stroke_count             int
	non_radical_stroke_count int

	heteronyms []heteronym
}

type heteronym struct {
	pinyin    string
	bopomofo  string
	bopomofo2 string

	definitions []definition
}

type definition struct {
	def      string
	quote    []string
	example  []string
	defType  string // this field is called "type" in the output of the server
	link     string
	synonyms string
}

func FindRecord(word string) (*Record, error) {
	resp, err := http.Get("http://www.moedict.tw/uni/" + word)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data interface{}
	jsonErr := json.Unmarshal(body, &data)
	if jsonErr != nil {
		return nil, jsonErr
	}

	r := new(Record)

	dataDict := data.(map[string]interface{})

	r.title = dataDict["title"].(string)
	r.stroke_count = dataDict["stroke_count"].(int)
	// r.radical = data["radical"].(string)
	// r.non_radical_stroke_count = data["non_radical_stroke_count"].(int)

	fmt.Println(r)
	return nil, nil
}

func main() {
	FindRecord("如")
}

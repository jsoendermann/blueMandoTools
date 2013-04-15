package moedict

import (
  "fmt"
  "github.com/yangchuanzhang/chinese"
  "ioutil"
)


// these three structs reflect the json of the api at https://www.moedict.tw/uni/
// which is documented (in Chinese) here: https://hackpad.com/3du.tw-API-Design-95jKjray8uR
type Record struct {
  word string
  radical string
  strokes int
  nonRadicalStrokes int

  heteronyms []heteronym
}

type heteronym struct {
  pinyin string
  bopomofo string
  bopomofo2 string

  definitions []definition
}

type definition struct {
  def string
  quote []string
  example []string
  defType string  // this field is called "type" in the output of the server
  link string
  synonyms string
}

func FindRecord(word string) (Record, error) {
  resp, err := http.Get("http://www.moedict.tw/uni/"+word)
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()

  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return nil, err
  }



}

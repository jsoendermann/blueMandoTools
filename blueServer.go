/*
TODO add package description
*/
package main

import (
  "fmt"
  "io/ioutil"
  "net/http"
//  "os"
//  "html/template"
  "github.com/yangchuanzhang/cedict"
  "github.com/yangchuanzhang/pinyin"
  "github.com/yangchuanzhang/chinese"
)

const (
  vocabPath = "/vocab/"
  sentencesPath = "/sentences/"

  vocabLookupPath = "/vocab/lookup/"
  sentencesLookupPath = "/sentences/lookup/"

  assetsPath = "/assets/"
)

const (
  vocabLookupPathLength = len(vocabLookupPath)
  sentencesLookupPathLength = len(sentencesLookupPath)
)


var vocabHtml, sentencesHtml string


func indexHandler(w http.ResponseWriter, r *http.Request) {
  http.Redirect(w, r, sentencesPath, http.StatusFound)
}

func vocabHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, vocabHtml)
}

func vocabLookupHandler(w http.ResponseWriter, r *http.Request) {
  word := r.URL.Path[vocabLookupPathLength:]
  fmt.Println(word)

  //FIXME handle error
  records, _ := cedict.FindRecords(word, chinese.Simp)
  if records == nil {
    records, _ = cedict.FindRecords(word, chinese.Trad)
  }

  if len(records) == 0 {
    fmt.Fprintf(w, `{"error": "No matches found", "word": "`+word+`"}`)
    return
  }

  var output string

  //FIXME explain \\t how it's not valid json
  output += records[0].Simp
  output += "\\t"
  output += records[0].Trad
  output += "\\t"

  for _,record := range records {
    //FIXME add colors
    output += pinyin.Num2Dia(record.Pinyin, "&nbsp;")
    //FIXME explain this
    output += "&nbsp;&nbsp;&nbsp; "
    //FIXME mask "s
    output += record.English
    output += "<br />"
  }



  fmt.Fprintf(w, `{"error":"nil", "response":"`+output+`"}`)
}

func sentenceHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
  fmt.Println("Welcome to the Blue Mandarin Lab Card Generator Server.")


  vocabHtmlData, err := ioutil.ReadFile("vocab.html")
  if err != nil {
    panic("vocab.html could not be opened")
  }
  vocabHtml = string(vocabHtmlData)

//  sentencesHtmlData, err := ioutil.ReadFile("sentences.html")
//  if err != nil {
//    panic("sentences.html could not be opened")
//  }
//  sentencesHtml = string(sentencesHtmlData)


    http.HandleFunc("/", indexHandler)
    http.HandleFunc(vocabPath, vocabHandler)
    http.HandleFunc(sentencesPath, sentenceHandler)

    http.HandleFunc(vocabLookupPath, vocabLookupHandler)

    http.Handle(assetsPath, http.FileServer(http.Dir(".")))

    http.ListenAndServe(":8080", nil)
}

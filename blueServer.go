package main

import (
	"fmt"
	"github.com/yangchuanzhang/cedict"
	"github.com/yangchuanzhang/chinese"
	"github.com/yangchuanzhang/pinyin"
	"io/ioutil"
	"net/http"
  "encoding/json"
)

// Paths on the server
const (
	vocabPath     = "/vocab/"
	sentencesPath = "/sentences/"

	vocabLookupPath     = "/vocab/lookup/"
	sentencesLookupPath = "/sentences/lookup/"

	assetsPath = "/assets/"
)

// The length of the lookup paths, used to separate
// the words from the url
const (
	vocabLookupPathLength     = len(vocabLookupPath)
	sentencesLookupPathLength = len(sentencesLookupPath)
)

// The HTML for the two pages is loaded into these variables
// to avoid having to load them on each request
var vocabHtml, sentencesHtml string

// The indexHandler redirects the user to the sentences page
func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, sentencesPath, http.StatusFound)
}

// The vocabHandler is a static page that gets written to the
// ResponseWriter
func vocabHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, vocabHtml)
}

// FIXME always use json.Marshal instead of manually constructing json
func vocabLookupHandler(w http.ResponseWriter, r *http.Request) {
	// get the word from the path
	word := r.URL.Path[vocabLookupPathLength:]

  // get colors
  colors := make([]string, 5)
  for i := 0; i <= 4; i ++ {
    colors[i] = r.FormValue(fmt.Sprintf("tone%d", i))
  }

	// search the db for records (simp first, if unsuccessful, try trad)
	// and send errors back to client if any arise
	records, err := cedict.FindRecords(word, chinese.Simp)
	if err != nil {
		fmt.Fprintf(w, `{"error": "`+err.Error()+`", "word": "`+word+`"}`)
		return
	}
	if records == nil {
		records, err = cedict.FindRecords(word, chinese.Trad)
		if err != nil {
			fmt.Fprintf(w, `{"error": "`+err.Error()+`", "word": "`+word+`"}`)
			return
		}

	}

	// check if there were no matches in the db
	if len(records) == 0 {
		fmt.Fprintf(w, `{"error": "No matches found", "word": "`+word+`"}`)
		return
	}

	// construct csv row
	var output string

	output += records[0].Simp
	output += "\t"
	output += records[0].Trad
	output += "\t"

	for _, record := range records {
		output += pinyin.Num2DiaCol(record.Pinyin, colors, "&nbsp;")
		// \t can't be used for this because it separates the columns
		// in the csv. Add another real space character at the end
		// to make the line break between pinyin and definition on small devices
		output += "&nbsp;&nbsp;&nbsp; "
		output += record.English
		output += "<br />"
	}

  // use json.Marshal with an anonymous variable to escape the \t and " characters
  // in the response
  j, err := json.Marshal(map[string]interface{}{
    "error": "nil",
    "csv": output,
  })
  if err != nil {
    fmt.Fprintf(w, `{"error": "`+err.Error()+`", "word": "`+word+`"}`)
    return
  }

	fmt.Fprintf(w, string(j))
}

// The sentenceHandler is a static page that gets written to the
// ResponseWriter
func sentenceHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, sentencesHtml)
}

func sentencesLookupHandler(w http.ResponseWriter, r *http.Request) {
	//TODO implement this method
}

func main() {
	fmt.Println("Welcome to the Blue Mandarin Lab Card Generator Server.")

	// Load Db, panic on error and defer close
	err := cedict.LoadDb()
	if err != nil {
		panic(err)
	}
	defer cedict.CloseDb()

	// Load the html for the two pages into memory, panic on error
	vocabHtmlData, err := ioutil.ReadFile("vocab.html")
	if err != nil {
		panic("vocab.html could not be opened")
	}
	vocabHtml = string(vocabHtmlData)

	sentencesHtmlData, err := ioutil.ReadFile("sentences.html")
	if err != nil {
		panic("sentences.html could not be opened")
	}
	sentencesHtml = string(sentencesHtmlData)

	// root
	http.HandleFunc("/", indexHandler)

	// static pages
	http.HandleFunc(vocabPath, vocabHandler)
	http.HandleFunc(sentencesPath, sentenceHandler)

	// json api
	http.HandleFunc(vocabLookupPath, vocabLookupHandler)
	http.HandleFunc(sentencesLookupPath, sentencesLookupHandler)

	// assets file server
	http.Handle(assetsPath, http.FileServer(http.Dir(".")))

	// start server
	http.ListenAndServe(":8080", nil)
}

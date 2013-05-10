package main

import (
	"encoding/json"
	"fmt"
	"github.com/hoisie/mustache"
	"github.com/yangchuanzhang/cedict"
	"github.com/yangchuanzhang/chinese"
	"github.com/yangchuanzhang/moedict"
	"github.com/yangchuanzhang/pinyin"
	"net/http"
	"regexp"
	"strings"
)

// This regexp is used to find words in sentences
var sentencesRegexp *regexp.Regexp

func main() {
	fmt.Println("Welcome to the Blue Mandarin Lab Flash Card Server.")

	err := cedict.LoadDb()
	if err != nil {
		panic(err)
	}
	defer cedict.CloseDb()

	sentencesRegexp, err = regexp.Compile("\\[(.*?)\\]")
	if err != nil {
		panic(err)
	}

	vocabHtml := mustache.RenderFileInLayout("vocab.html", "layout.html")
	sentencesHtml := mustache.RenderFileInLayout("sentences.html", "layout.html")

	// FIXME reimplement this
	// set active class in navbar
	// FIXME find a better way to do this
	// vocabHtml = strings.Replace(vocabHtml, "<li id='vocab-link'>", "<li id='vocab-link' class='active'>", 1)
	// sentencesHtml = strings.Replace(sentencesHtml, "<li id='sentences-link'>", "<li id='sentences-link' class='active'>", 1)

	// Set up the http server

	// the root is handled by an anonymous function that redirects to "/sentences/"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, sentencesPath, http.StatusFound)
	})

	// /vocab/ and /sentences/ are both handled by simple, anonymous functions that
	// write static html to the ResponseWriter
	http.HandleFunc(vocabPath, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, vocabHtml)
	})
	http.HandleFunc(sentencesPath, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, sentencesHtml)
	})

	http.HandleFunc(vocabLookupPath, vocabLookupHandler)
	http.HandleFunc(sentencesLookupPath, sentencesLookupHandler)

	// assets file server
	http.Handle(assetsPath, http.FileServer(http.Dir(".")))

	http.ListenAndServe(":8080", nil)
}

// This function is responsible for paths of the form
// /vocab/lookup/<word> and returns a json dictionary
// containing the csv to be added to the output text area
// or an "error" field of something other than "nil" if an
// error occured during execution.
func vocabLookupHandler(w http.ResponseWriter, r *http.Request) {
	word, colors := getRequestDataAndColors(r)

	// search the db for records (simp first, if unsuccessful, try trad)
	// and send errors back to client if any occur
	records, err := cedict.FindRecords(word, chinese.Simp)
	if err != nil {
		fmt.Fprintf(w, `{"error": "`+err.Error()+`", "word": "`+word+`"}`)
		return
	}
	if len(records) == 0 {
		records, err = cedict.FindRecords(word, chinese.Trad)
		if err != nil {
			fmt.Fprintf(w, `{"error": "`+err.Error()+`", "word": "`+word+`"}`)
			return
		}

	}

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
		// Add another real space character at the end
		// to make the line break between pinyin and definition on small screens
		output += "&nbsp;&nbsp;&nbsp; "
		output += record.English
		output += "<br />"
	}

	// use json.Marshal with an anonymous variable to escape the \t and " characters
	// in the response
	j, err := json.Marshal(map[string]interface{}{
		"error": "nil",
		"csv":   output,
	})
	if err != nil {
		fmt.Fprintf(w, `{"error": "`+err.Error()+`", "word": "`+word+`"}`)
		return
	}

	fmt.Fprintf(w, string(j))
}

func sentencesLookupHandler(w http.ResponseWriter, r *http.Request) {
	sentence, colors := getRequestDataAndColors(r)

	// get words in sentence
	wordsRaw := sentencesRegexp.FindAllStringSubmatch(sentence, -1)
	words := make([]string, len(wordsRaw))
	for i, w := range wordsRaw {
		words[i] = w[1]
	}

	charSet := cedict.DetermineCharSet(sentence)

	// get moe entries
	// this array might end up being larger than len(words) if
	// one simplified word maps to more than one traditional word
	moeEntries := make([]moedict.Entry, len(words))

	// TODO explain diff. between trad and simp
	if charSet == chinese.Trad {
		for i, word := range words {
			entry, err := moedict.FindEntry(word)
			if err != nil {
				fmt.Fprintf(w, `{"error": "`+err.Error()+`", "sentence": "`+sentence+`"}`)
				return
			}

			moeEntries[i] = *entry
		}
	}
	// if the sentences is in simplified characters, find cedict records first
	// (there may be multiple records per word) and find moedict records using
	// the traditional word of the cedict records.
	if charSet == chinese.Simp {
		for _, word := range words {

			records, err := cedict.FindRecords(word, chinese.Simp)
			if err != nil {
				fmt.Fprintf(w, `{"error": "`+err.Error()+`", "sentence": "`+sentence+`"}`)
				return
			}

			for _, record := range records {
				entry, err := moedict.FindEntry(record.Trad)
				if err != nil {
					fmt.Fprintf(w, `{"error": "`+err.Error()+`", "sentence": "`+sentence+`"}`)
					return
				}

				moeEntries = append(moeEntries, *entry)
			}
		}
	}

	// construct csv row
	var output string

	var outputSentence string
	outputSentence = strings.Replace(sentence, "[", "", -1)
	outputSentence = strings.Replace(outputSentence, "]", "", -1)

	output += outputSentence
	output += "\t"

	for _, moeEntry := range moeEntries {
		output += moeEntry.ToHTML(colors)
		output += "<br>"
	}

	// TODO turn this into a function

	// use json.Marshal with an anonymous variable to escape the \t and " characters
	// in the response
	j, err := json.Marshal(map[string]interface{}{
		"error": "nil",
		"csv":   output,
	})
	if err != nil {
		fmt.Fprintf(w, `{"error": "`+err.Error()+`", "sentence": "`+sentence+`"}`)
		return
	}

	fmt.Fprintf(w, string(j))
}

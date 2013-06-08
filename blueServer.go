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
var findWordsInSentencesRegexp *regexp.Regexp
var findCharInMcdsRegexp *regexp.Regexp

func main() {
	fmt.Println("Welcome to the Blue Mandarin Lab Flash Card Server.")

	err := cedict.LoadDb()
	if err != nil {
		panic(err)
	}
	defer cedict.CloseDb()

	findWordsInSentencesRegexp, err = regexp.Compile("\\[(.*?)\\]")
	if err != nil {
		panic(err)
	}
	findCharInMcdsRegexp, err = regexp.Compile(clozeBegin + "(.*?)" + clozeEnd)
	if err != nil {
		panic(err)
	}

	vocabHtml := mustache.RenderFileInLayout("vocab.html", "layout.html")
	sentencesHtml := mustache.RenderFileInLayout("sentences.html", "layout.html")
	mcdsHtml := mustache.RenderFileInLayout("mcds.html", "layout.html")

	// FIXME reimplement this
	// set active class in navbar
	// FIXME find a better way to do this
	// vocabHtml = strings.Replace(vocabHtml, "<li id='vocab-link'>", "<li id='vocab-link' class='active'>", 1)
	// sentencesHtml = strings.Replace(sentencesHtml, "<li id='sentences-link'>", "<li id='sentences-link' class='active'>", 1)

	// Set up the http server

	// the root is handled by an anonymous function that redirects to "/sentences/"
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer, request, sentencesPath, http.StatusFound)
	})

	// /vocab/ and /sentences/ are both handled by simple, anonymous functions that
	// write static html to the ResponseWriter
	http.HandleFunc(vocabPath, func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, vocabHtml)
	})
	http.HandleFunc(sentencesPath, func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, sentencesHtml)
	})
	http.HandleFunc(mcdsPath, func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, mcdsHtml)
	})

	http.HandleFunc(vocabLookupPath, vocabLookupHandler)
	http.HandleFunc(sentencesLookupPath, sentencesLookupHandler)
	http.HandleFunc(mcdsLookupPath, mcdsLookupHandler)

	// assets file server
	http.Handle(assetsPath, http.FileServer(http.Dir(".")))

	http.ListenAndServe(":8080", nil)
}

// This function is responsible for paths of the form
// /vocab/lookup/<word> and returns a json dictionary
// containing the csv to be added to the output text area
// or an "error" field of something other than "nil" if an
// error occured during execution.
func vocabLookupHandler(writer http.ResponseWriter, request *http.Request) {
	word := getLastPathComponent(request)
	colors := getColors(request)

	// search the db for records (simp first, if unsuccessful, try trad)
	// and send errors back to client if any occur
	records, err := cedict.FindRecords(word, chinese.Simp)
	if err != nil {
		fmt.Fprintf(writer, `{"error": "`+err.Error()+`", "word": "`+word+`"}`)
		return
	}
	if len(records) == 0 {
		records, err = cedict.FindRecords(word, chinese.Trad)
		if err != nil {
			fmt.Fprintf(writer, `{"error": "`+err.Error()+`", "word": "`+word+`"}`)
			return
		}

	}

	if len(records) == 0 {
		fmt.Fprintf(writer, `{"error": "No matches found", "word": "`+word+`"}`)
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
		fmt.Fprintf(writer, `{"error": "`+err.Error()+`", "word": "`+word+`"}`)
		return
	}

	fmt.Fprintf(writer, string(j))
}

func sentencesLookupHandler(writer http.ResponseWriter, request *http.Request) {
	sentence := getLastPathComponent(request)
	colors := getColors(request)

	// get words in sentence
	wordsRaw := findWordsInSentencesRegexp.FindAllStringSubmatch(sentence, -1)
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
				fmt.Fprintf(writer, `{"error": "`+err.Error()+`", "sentence": "`+sentence+`"}`)
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
				fmt.Fprintf(writer, `{"error": "`+err.Error()+`", "sentence": "`+sentence+`"}`)
				return
			}

			for _, record := range records {
				entry, err := moedict.FindEntry(record.Trad)
				if err != nil {
					fmt.Fprintf(writer, `{"error": "`+err.Error()+`", "sentence": "`+sentence+`"}`)
					return
				}

				if entry != nil {
					moeEntries = append(moeEntries, *entry)
				}
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

	// TODO add cedict records

	// TODO turn this into a function

	// use json.Marshal with an anonymous variable to escape the \t and " characters
	// in the response
	j, err := json.Marshal(map[string]interface{}{
		"error": "nil",
		"csv":   output,
	})
	if err != nil {
		fmt.Fprintf(writer, `{"error": "`+err.Error()+`", "sentence": "`+sentence+`"}`)
		return
	}

	fmt.Fprintf(writer, string(j))
}

func mcdsLookupHandler(writer http.ResponseWriter, request *http.Request) {
	mcd := getLastPathComponent(request)
	colors := getColors(request)

	// get clozed char and original text
	mcd = strings.Replace(mcd, "@SLASH@", "/", -1)
	back := strings.Split(mcd, "\t")[1]
	chars := findCharInMcdsRegexp.FindAllStringSubmatch(back, -1)
	clozeChar := chars[0][1]
	originalText := strings.Replace(back, clozeBegin+clozeChar+clozeEnd, clozeChar, -1)

	splitText, err := cedict.SplitChineseTextIntoWords(originalText)
	if err != nil {
		fmt.Fprintf(writer, `{"error": "`+err.Error()+`", "mcd": "`+mcd+`"}`)
		return
	}

	charSet := cedict.DetermineCharSet(originalText)

	words := make([]string, 0)
	cedictRecords := make([]cedict.Record, 0)

	for _, ctw := range splitText {
		if ctw.T == cedict.WordTypeRecords {
			if cw := ctw.R[0].WordByCharSet(charSet); strings.Index(cw, clozeChar) != -1 {
				//fmt.Println(cw)
				words = append(words, cw)
				for _, r := range ctw.R {
					cedictRecords = append(cedictRecords, r)
				}
			}
		}
	}

	// get moe entries
	// this array might end up being larger than len(words) if
	// one simplified word maps to more than one traditional word
	moeEntries := make([]moedict.Entry, len(words))

	// TODO explain diff. between trad and simp
	if charSet == chinese.Trad {
		for i, word := range words {
			entry, err := moedict.FindEntry(word)
			if err != nil {
				fmt.Fprintf(writer, `{"error": "`+err.Error()+`", "mcd": "`+mcd+`"}`)
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
				fmt.Println("2")
				fmt.Fprintf(writer, `{"error": "`+err.Error()+`", "mcd": "`+mcd+`"}`)
				return
			}

			for _, record := range records {
				entry, err := moedict.FindEntry(record.Trad)
				if err != nil {
					fmt.Fprintf(writer, `{"error": "`+err.Error()+`", "mcd": "`+mcd+`"}`)
					return
				}

				if entry != nil {
					moeEntries = append(moeEntries, *entry)
				}
			}
		}
	}

	// construct csv row
	var output string

	output += mcd
	output += "\t"

	for _, moeEntry := range moeEntries {
		output += moeEntry.ToHTML(colors)
		output += "<br>"
	}

	// TODO add cedict records

	// TODO turn this into a function

	// use json.Marshal with an anonymous variable to escape the \t and " characters
	// in the response
	j, err := json.Marshal(map[string]interface{}{
		"error": "nil",
		"csv":   output,
	})
	if err != nil {
		fmt.Fprintf(writer, `{"error": "`+err.Error()+`", "mcd": "`+mcd+`"}`)
		return
	}

	fmt.Fprintf(writer, string(j))

}

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
	fmt.Println("### Welcome to the Blue Mandarin Lab Flash Card Server ###")

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

  // Load html files, the array at the end contains the js files to be loaded
	vocabHtml := mustache.RenderFileInLayout(    "vocab.html",      layoutFile, map[string]interface{}{"jsfiles": []string{"vocab"}})
  moeVocabHtml := mustache.RenderFileInLayout( "moe-vocab.html",  layoutFile, map[string]interface{}{"jsfiles": []string{"moe-vocab"}})
	sentencesHtml := mustache.RenderFileInLayout("sentences.html",  layoutFile, map[string]interface{}{"jsfiles": []string{"sentences"}})
	mcdsHtml := mustache.RenderFileInLayout(     "mcds.html",       layoutFile, map[string]interface{}{"jsfiles": []string{"mcds", "mcds-dict"}})
  settingsHtml := mustache.RenderFileInLayout( "settings.html",   layoutFile, map[string]interface{}{"jsfiles": []string{"settings"}})
  helpAboutHtml := mustache.RenderFileInLayout("help-about.html", layoutFile)

	// FIXME reimplement this
	// set active class in navbar
	// FIXME find a better way to do this
	// vocabHtml = strings.Replace(vocabHtml, "<li id='vocab-link'>", "<li id='vocab-link' class='active'>", 1)
	// sentencesHtml = strings.Replace(sentencesHtml, "<li id='sentences-link'>", "<li id='sentences-link' class='active'>", 1)

	// Set up the http server

	// the root is handled by an anonymous function that redirects to "/sentences/"
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer, request, helpAboutPath, http.StatusFound)
	})

	// /vocab/, /sentences/ and /mcds/ are both handled by simple, anonymous functions that
	// write static html to the ResponseWriter
	http.HandleFunc(vocabPath, func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, vocabHtml)
	})
  http.HandleFunc(moeVocabPath, func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, moeVocabHtml)
	})
	http.HandleFunc(sentencesPath, func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, sentencesHtml)
	})
	http.HandleFunc(mcdsPath, func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, mcdsHtml)
	})

	http.HandleFunc(vocabLookupPath, vocabLookupHandler)
  http.HandleFunc(moeVocabLookupPath, moeVocabLookupHandler)
	http.HandleFunc(sentencesLookupPath, sentencesLookupHandler)
	http.HandleFunc(mcdsLookupPath, mcdsLookupHandler)

  // /settings/ and /help-about/ handlers
  http.HandleFunc(settingsPath, func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, settingsHtml)
	})
  http.HandleFunc(helpAboutPath, func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, helpAboutHtml)
	})

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

func moeVocabLookupHandler(writer http.ResponseWriter, request *http.Request) {
  word := getLastPathComponent(request)
	colors := getColors(request)

  // convert to trad
  tradWord, err := cedict.Simp2Trad(word)
  if err != nil {
    fmt.Fprintf(writer, `{"error": "`+err.Error()+`", "word": "`+word+`"}`)
		return
	}

  moeEntry, err := moedict.FindEntry(tradWord)
  if moeEntry == nil {
		fmt.Fprintf(writer, `{"error": "No matches found", "word": "`+word+`"}`)
		return
	}
  if err != nil {
    fmt.Fprintf(writer, `{"error": "`+err.Error()+`", "word": "`+word+`"}`)
		return
	}

	// construct csv row
	var output string

	output += tradWord
	output += "\t"
  output += moeEntry.ToHTML(colors)

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

	moeEntries, err := findMoeEntriesForWords(words, charSet)
	if err != nil {
		fmt.Fprintf(writer, `{"error": "`+err.Error()+`", "sentence": "`+sentence+`"}`)
		return
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
				wordInArray := false
				for _, ew := range words {
					if ew == cw {
						wordInArray = true
					}
				}
				if !wordInArray {
					words = append(words, cw)
					for _, r := range ctw.R {
						cedictRecords = append(cedictRecords, r)
					}
				}
			}
		}
	}

	moeEntries, err := findMoeEntriesForWords(words, charSet)
	if err != nil {
		fmt.Fprintf(writer, `{"error": "`+err.Error()+`", "mcd": "`+mcd+`"}`)
		return
	}

	// construct csv row
	var output string

	output += mcd
	output += "\t"

	for _, moeEntry := range moeEntries {
		output += moeEntry.ToHTML(colors)
		output += "<br>"
	}

	output += "\t"

	for _, cr := range cedictRecords {
		output += cr.ToHTML(colors)
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
		fmt.Fprintf(writer, `{"error": "`+err.Error()+`", "mcd": "`+mcd+`"}`)
		return
	}

	fmt.Fprintf(writer, string(j))

}

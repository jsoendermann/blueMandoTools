package main

import (
	"encoding/json"
	"fmt"
	"github.com/hoisie/mustache"
	"github.com/yangchuanzhang/cedict"
	"github.com/yangchuanzhang/chinese"
	"github.com/yangchuanzhang/blueMandoTools/moedict"
	"github.com/yangchuanzhang/pinyin"
	"github.com/yangchuanzhang/zhDicts"
	"net/http"
	"regexp"
	"strings"
)

// This regexp is used to find words in sentences
var findWordsInSentencesRegexp *regexp.Regexp
var findCharInMcdsRegexp *regexp.Regexp

func main() {
	fmt.Println("### Welcome to the Blue Mandarin Lab Flash Card Server ###")

	err := zhDicts.LoadDb()
	if err != nil {
		panic(err)
	}
	defer zhDicts.CloseDb()

	err = cedict.Initialize()
	if err != nil {
		panic(err)
	}
	err = moedict.Initialize()
	if err != nil {
		panic(err)
	}

	findWordsInSentencesRegexp, err = regexp.Compile("\\[(.*?)\\]")
	if err != nil {
		panic(err)
	}
	findCharInMcdsRegexp, err = regexp.Compile(clozeBegin + "(.*?)" + clozeEnd)
	if err != nil {
		panic(err)
	}

	// Load html files, the array at the end contains the js files to be loaded
	vocabHtml := mustache.RenderFileInLayout("vocab.html", layoutFile, map[string]interface{}{"jsfiles": []string{"vocab"}})
	moeVocabHtml := mustache.RenderFileInLayout("moe-vocab.html", layoutFile, map[string]interface{}{"jsfiles": []string{"moe-vocab"}})
	htmlLookupHtml := mustache.RenderFileInLayout("html-lookup.html", layoutFile, map[string]interface{}{"jsfiles": []string{"html-lookup"}})
	convertHtml := mustache.RenderFileInLayout("convert.html", layoutFile, map[string]interface{}{"jsfiles": []string{"convert"}})
	sentencesHtml := mustache.RenderFileInLayout("sentences.html", layoutFile, map[string]interface{}{"jsfiles": []string{"sentences"}})
	mcdsHtml := mustache.RenderFileInLayout("mcds.html", layoutFile, map[string]interface{}{"jsfiles": []string{"mcds", "mcds-dict"}})
	settingsHtml := mustache.RenderFileInLayout("settings.html", layoutFile, map[string]interface{}{"jsfiles": []string{"settings"}})
	helpAboutHtml := mustache.RenderFileInLayout("help-about.html", layoutFile)

	// FIXME set active class in navbar

	// Set up the http server

	// the root is handled by an anonymous function that redirects to "/sentences/"
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer, request, mcdsPath, http.StatusFound)
	})

	// small helper function (used below)
	addStaticHtmlHandler := func(path string, html string) {
		http.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
			fmt.Fprintf(writer, html)
		})
	}

	addStaticHtmlHandler(vocabPath, vocabHtml)
	addStaticHtmlHandler(moeVocabPath, moeVocabHtml)
	addStaticHtmlHandler(htmlLookupPath, htmlLookupHtml)
	addStaticHtmlHandler(convertPath, convertHtml)
	addStaticHtmlHandler(sentencesPath, sentencesHtml)
	addStaticHtmlHandler(mcdsPath, mcdsHtml)
	addStaticHtmlHandler(settingsPath, settingsHtml)
	addStaticHtmlHandler(helpAboutPath, helpAboutHtml)

	http.HandleFunc(vocabLookupPath, vocabLookupHandler)
	http.HandleFunc(moeVocabLookupPath, moeVocabLookupHandler)
	http.HandleFunc(convertLookupPath, convertLookupHandler)
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
		fmt.Fprint(writer, `{"error": "`+err.Error()+`", "word": "`+word+`"}`)
		return
	}
	if len(records) == 0 {
		records, err = cedict.FindRecords(word, chinese.Trad)
		if err != nil {
			fmt.Fprint(writer, `{"error": "`+err.Error()+`", "word": "`+word+`"}`)
			return
		}

	}

	if len(records) == 0 {
		fmt.Fprint(writer, `{"error": "No matches found", "word": "`+word+`"}`)
		return
	}

	// this string has the same number of characters as the word.
	// in places where the simplified and the traditional characters
	// are the same, it has a space, otherwise it has the simplified character
	simpChars := ""
	tradChars := ""
	for is, cs := range records[0].Simp {
		for it, ct := range records[0].Trad {
			if is == it {
				if cs == ct {
					simpChars += string('\u3000')
					tradChars += string('\u3000')
				} else {
					simpChars += string(cs)
					tradChars += string(ct)
				}
			}
		}
	}

	// construct csv row
	var output string

	output += records[0].Trad
	output += "\t"
	output += records[0].Simp
	output += "\t"
	// This dot is necessary because Anki trims whitespace when importing.
	// For more details, see the card layout of the shengci deck
	output += "." + tradChars
	output += "\t"
	output += "." + simpChars
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
		fmt.Fprint(writer, `{"error": "`+err.Error()+`", "word": "`+word+`"}`)
		return
	}

	fmt.Fprint(writer, string(j))
}

func moeVocabLookupHandler(writer http.ResponseWriter, request *http.Request) {
	word := getLastPathComponent(request)
	colors := getColors(request)

	// convert to trad
	tradWord, err := cedict.Simp2Trad(word)
	if err != nil {
		fmt.Fprint(writer, `{"error": "`+err.Error()+`", "word": "`+word+`"}`)
		return
	}

	moeEntry, err := moedict.FindEntry(tradWord)
	if moeEntry == nil {
		fmt.Fprint(writer, `{"error": "No matches found", "word": "`+word+`"}`)
		return
	}
	if err != nil {
		fmt.Fprint(writer, `{"error": "`+err.Error()+`", "word": "`+word+`"}`)
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
		fmt.Fprint(writer, `{"error": "`+err.Error()+`", "word": "`+word+`"}`)
		return
	}

	fmt.Fprint(writer, string(j))
}

func convertLookupHandler(writer http.ResponseWriter, request *http.Request) {} //TODO

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
		fmt.Fprint(writer, `{"error": "`+err.Error()+`", "sentence": "`+sentence+`"}`)
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
		fmt.Fprint(writer, `{"error": "`+err.Error()+`", "sentence": "`+sentence+`"}`)
		return
	}

	fmt.Fprint(writer, string(j))
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
		fmt.Fprint(writer, `{"error": "`+err.Error()+`", "mcd": "`+mcd+`"}`)
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
		fmt.Fprint(writer, `{"error": "`+err.Error()+`", "mcd": "`+mcd+`"}`)
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
		fmt.Fprint(writer, `{"error": "`+err.Error()+`", "mcd": "`+mcd+`"}`)
		return
	}

	fmt.Fprint(writer, string(j))

}

package main

import (
	"encoding/json"
	"fmt"
	"github.com/hoisie/mustache"
	"github.com/jsoendermann/blueMandoTools/cedict"
	"github.com/jsoendermann/blueMandoTools/chinese"
	"github.com/jsoendermann/blueMandoTools/moedict"
	"github.com/jsoendermann/blueMandoTools/pinyin"
	"github.com/jsoendermann/blueMandoTools/zhDicts"
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
    homeHtml := mustache.RenderFileInLayout("home.html", layoutFile, map[string]interface{}{})
	vocabHtml := mustache.RenderFileInLayout("vocab.html", layoutFile, map[string]interface{}{"jsfiles": []string{"vocab"}})
	moeVocabHtml := mustache.RenderFileInLayout("moe-vocab.html", layoutFile, map[string]interface{}{"jsfiles": []string{"moe-vocab"}})
	sentencesHtml := mustache.RenderFileInLayout("sentences.html", layoutFile, map[string]interface{}{"jsfiles": []string{"sentences"}})
	mcdsHtml := mustache.RenderFileInLayout("mcds.html", layoutFile, map[string]interface{}{"jsfiles": []string{"mcds", "mcds-dict"}})
	settingsHtml := mustache.RenderFileInLayout("settings.html", layoutFile, map[string]interface{}{"jsfiles": []string{"settings"}})

	// FIXME set active class in navbar


	// Set up the http server

	// small helper function (used below)
	addStaticHtmlHandler := func(path string, html string) {
		http.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
			fmt.Fprintf(writer, html)
		})
	}

    addStaticHtmlHandler(homePath, homeHtml)
	addStaticHtmlHandler(vocabPath, vocabHtml)
	addStaticHtmlHandler(moeVocabPath, moeVocabHtml)
	addStaticHtmlHandler(sentencesPath, sentencesHtml)
	addStaticHtmlHandler(mcdsPath, mcdsHtml)
	addStaticHtmlHandler(settingsPath, settingsHtml)

	http.HandleFunc(vocabLookupPath, vocabLookupHandler)
	http.HandleFunc(moeVocabLookupPath, moeVocabLookupHandler)
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

    fmt.Println("vocabLookupHandler: " + word)

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

    fmt.Println("moeVocabLookupHandler: " + word)

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

    fmt.Println("sentencesLookupHandler: " + string([]rune(sentence)[0:10]) + "。。。")

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

    cedictRecords := make([]cedict.Record, 0)
    for _, word := range words {
        records, err := cedict.FindRecords(word, charSet)
        if err != nil {
            fmt.Fprint(writer, `{"error": "`+err.Error()+`", "sentence": "`+sentence+`"}`)
            return
        }
        for _, record := range records {
            cedictRecords = append(cedictRecords, record)
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

	output += "\t"

    for _, cedictRecord := range cedictRecords {
        output += cedictRecord.ToHTML(colors)
        output += "<br>"
    }

    // This adds an additional field for the audio file
    output += "\t"
    // This is necessary for Anki to recognise all fields
    output += "&nbsp;"

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
	text := getLastPathComponent(request)
    text = strings.Replace(text, "@SLASH@", "/", -1)
    text = strings.Replace(text, "\n", "<br />", -1)
    text = strings.Replace(text, "\t", "&nbsp;&nbsp;&nbsp;&nbsp;", -1)

    notes := request.FormValue("notes")
    notes = strings.Replace(notes, "@SLASH@", "/", -1)
    notes = strings.Replace(notes, "\n", "<br />", -1)
    notes = strings.Replace(notes, "\t", "&nbsp;&nbsp;&nbsp;&nbsp;", -1)

    charsRaw := request.FormValue("chars")
    charsArrayRaw := strings.Split(charsRaw, " ")
    // Construct new slice without empty strings
    chars := make([]string, 0)
    for _, char := range charsArrayRaw {
        if char != "" {
            chars = append(chars, char)
        }
    }

	colors := getColors(request)

    charSet := cedict.DetermineCharSet(text)
    splitText, err := chinese.SplitChineseTextIntoWords(text, charSet)
    if err != nil {
		fmt.Fprint(writer, `{"error": "`+err.Error()+`}`)
		return
	}

    var output string

    // Go through all chars to be clozed. Every iteration adds one line to the output var (i.e. one card)
    for _, char := range chars {
        var front, back string
        wordsToBeLookedUp := make([]string, 0)

        // This loop goes through all words in the text. This could be optimised by only going
        // through the text once and constructing all cards simultaneously
        for _, wordInText := range splitText {
            indexOfChar := strings.Index(wordInText, char)
            if indexOfChar == -1 {
                front += wordInText
                back += wordInText
            } else {
                front += strings.Replace(wordInText, char, clozeBegin + clozeChar + clozeEnd, -1)
                back += strings.Replace(wordInText, char, clozeBegin + char + clozeEnd, -1)
            
                wordsToBeLookedUp = append(wordsToBeLookedUp, wordInText)
            }
        }

        // Fine unique words to be looked up in moedict
        wordsToBeLookedUnique := make([]string, 0)
        for _, word := range wordsToBeLookedUp {
            alreadyInWordsToBeLookedUpUnique := false
            for _, wordUnique := range wordsToBeLookedUnique {
                if word == wordUnique {
                    alreadyInWordsToBeLookedUpUnique = true
                }
            }
            if !alreadyInWordsToBeLookedUpUnique {
                wordsToBeLookedUnique = append(wordsToBeLookedUnique, word)
            }
        }

        cedictRecords := make([]cedict.Record,0)
        for _, wordToBeLookedUpInCedict := range wordsToBeLookedUnique {
            records, _ := cedict.FindRecords(wordToBeLookedUpInCedict, charSet)
            for _, record := range records {
                cedictRecords = append(cedictRecords, record)
            }
        }

        moeEntries, err := findMoeEntriesForWords(wordsToBeLookedUnique, charSet)
	    if err != nil {
		    fmt.Fprint(writer, `{"error": "`+err.Error()+`}`)
		    return
	    }

        output += front
        output += "\t"
        output += back
        output += "\t"
        output += notes
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

        output += "\n"
    }

    // TODO turn this into a function

	// use json.Marshal with an anonymous variable to escape the \t and " characters
	// in the response
	j, err := json.Marshal(map[string]interface{}{
		"error": "nil",
		"result":   output,
	})
	if err != nil {
		fmt.Fprint(writer, `{"error": "`+err.Error()+`}`)
		return
	}

	fmt.Fprint(writer, string(j))
    


}

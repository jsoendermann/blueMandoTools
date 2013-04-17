package main

import (
  "fmt"
  "github.com/yangchuanzhang/cedict"
  "github.com/yangchuanzhang/chinese"
  "github.com/yangchuanzhang/pinyin"
  "io/ioutil"
  "net/http"
  "encoding/json"
  "strings"
  "regexp"
  "github.com/yangchuanzhang/moedict"
)

//TODO add printlns

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

// This regexp is used to find words in sentences
var sentencesRegexp *regexp.Regexp





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
  // get the sentence from the path
  sentence := r.URL.Path[sentencesLookupPathLength:]

  fmt.Println("Received lookup request for ", sentence)

  // get colors
  colors := make([]string, 5)
  for i := 0; i <= 4; i ++ {
    colors[i] = r.FormValue(fmt.Sprintf("tone%d", i))
  }

  // get words in sentence
  wordsRaw := sentencesRegexp.FindAllStringSubmatch(sentence, -1)
  words := make([]string, len(wordsRaw))
  for i,w := range wordsRaw {
    words[i] = w[1]
  }

  fmt.Println(" Word list: ", words)

  // determine char set
  charSet := cedict.DetermineCharSet(sentence)

  // get moe entries
  // this array might end up being larger than len(words) if
  // one simplified word maps to more than one traditional word
  moeEntries := make([]moedict.Entry, len(words))

  // TODO explain diff. between trad and simp
  if charSet == chinese.Trad {
    for i,word := range words {
      // find entry, return on error
      entry, err := moedict.FindEntry(word)
      if err != nil {
        fmt.Fprintf(w, `{"error": "`+err.Error()+`", "word": "`+word+`"}`)
        return
      }

      moeEntries[i] = *entry
    }
  }
  // if the sentences is in simplified characters, find cedict records first
  // (there may be multiple records per word) and find moedict records using
  // the traditional word of the cedict records.
  if charSet == chinese.Simp {
    // go through words
    for _,word := range words{

      // find cedict records
      records, err := cedict.FindRecords(word, chinese.Simp)
      if err != nil {
        fmt.Fprintf(w, `{"error": "`+err.Error()+`", "word": "`+word+`"}`)
        return
      }

      // find moedict records for record
      for _,record := range records {
        entry, err := moedict.FindEntry(record.Trad)
        if err != nil {
          fmt.Fprintf(w, `{"error": "`+err.Error()+`", "word": "`+word+`"}`)
          return
        }

        moeEntries = append(moeEntries, *entry)
      }
    }
  }

  // construct csv row
  var output string

  // remove brackets from string
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
    "csv": output,
  })
  if err != nil {
    fmt.Fprintf(w, `{"error": "`+err.Error()+`", "word": "`+sentence+`"}`)
    return
  }

  fmt.Fprintf(w, string(j))
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
  applicationHtml := loadFilePanicOnError("application.html")
  vocabHtmlView := loadFilePanicOnError("vocab.html")
  sentencesHtmlView := loadFilePanicOnError("sentences.html")

  // combine the layout and the two views into complete html
  vocabHtml = strings.Replace(applicationHtml, "@yield@", vocabHtmlView, 1)
  sentencesHtml = strings.Replace(applicationHtml, "@yield@", sentencesHtmlView, 1)

  // set active class in navbar
  // FIXME move this into compile stage
  vocabHtml = strings.Replace(vocabHtml, "<li id='vocab-link'>", "<li id='vocab-link' class='active'>", 1)
  sentencesHtml = strings.Replace(sentencesHtml, "<li id='sentences-link'>", "<li id='sentences-link' class='active'>", 1)

  // compile regexp
  sentencesRegexp, err = regexp.Compile("\\[(.*?)\\]")
  if err != nil {
    panic(err)
  }


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

func loadFilePanicOnError(filename string) string {
  data, err := ioutil.ReadFile(filename)
  if err != nil {
    panic(err)
  }
  return string(data)
}

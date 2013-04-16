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
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
  http.Redirect(w, r, "/sentence-cards/", http.StatusFound)
}


func vocabHandler(w http.ResponseWriter, r *http.Request) {
 // t, _ := template.ParseFiles("vocab.html")
 // t.Execute(w)
  vocabData, _ := ioutil.ReadFile("vocab.html")
  fmt.Fprintf(w, string(vocabData))
}

func sentenceHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
  fmt.Println("hallo")
    http.HandleFunc("/", indexHandler)
    http.HandleFunc("/vocab-cards/", vocabHandler)
    http.HandleFunc("/sentence-cards/", sentenceHandler)

    http.Handle("/assets/", http.FileServer(http.Dir(".")))

    http.ListenAndServe(":8080", nil)

}

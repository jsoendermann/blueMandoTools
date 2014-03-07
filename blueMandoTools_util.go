package main

import (
	"fmt"
	"github.com/jsoendermann/blueMandoTools/cedict"
	"github.com/jsoendermann/blueMandoTools/chinese"
	"github.com/jsoendermann/blueMandoTools/moedict"
	"net/http"
	"strings"
)

func getLastPathComponent(request *http.Request) string {
	// get the word or sentence from the path
	pathElements := strings.Split(request.URL.Path, "/")
	requestData := pathElements[len(pathElements)-1]

	return requestData
}

func getColors(request *http.Request) []string {
	colors := make([]string, 5)
	for i := 0; i <= 4; i++ {
		colors[i] = request.FormValue(fmt.Sprintf("tone%d", i))
	}

	return colors
}

func findMoeEntriesForWords(words []string, charSet chinese.CharSet) ([]moedict.Entry, error) {
	moeEntries := make([]moedict.Entry, 0)

	// TODO explain diff. between trad and simp
	if charSet == chinese.Trad {
		for _, word := range words {
			entry, err := moedict.FindEntry(word)
			if err != nil {
				return nil, err
			}
      if entry != nil {
        moeEntries = append(moeEntries, *entry)
      }
		}
	}

	// if the sentences is in simplified characters, find cedict records first
	// (there may be multiple records per word) and find moedict records using
	// the traditional word of the cedict records.
	if charSet == chinese.Simp {
		for _, word := range words {

			records, err := cedict.FindRecords(word, chinese.Simp)
			if err != nil {
				return nil, err
			}

			individualRecords := make([]cedict.Record, 0)
			for _, r := range records {
				recordInSlice := false
				for _, ir := range individualRecords {
					if r.Trad == ir.Trad {
						recordInSlice = true
					}
				}
				if !recordInSlice {
					individualRecords = append(individualRecords, r)
				}
			}

			for _, record := range individualRecords {
				entry, err := moedict.FindEntry(record.Trad)
				if err != nil {
					return nil, err
				}

				if entry != nil {
					moeEntries = append(moeEntries, *entry)
				}
			}
		}
	}

	return moeEntries, nil
}

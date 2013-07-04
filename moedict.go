/*
TODO package description
*/
package moedict

import (
	"fmt"
        "strings"
	"github.com/yangchuanzhang/bopomofo"
	"github.com/yangchuanzhang/zhDicts"
)

// These three structs reflect the json of the api at https://www.moedict.tw/uni/
// which is documented (in Chinese) here: https://hackpad.com/3du.tw-API-Design-95jKjray8uR
// Fields that are commented out are not used at the moment.
type Entry struct {
	Title                    string
	Radical                  string
	Stroke_count             int
	Non_radical_stroke_count int

	Heteronyms []Heteronym
}

type Heteronym struct {
	Pinyin    string
	Bopomofo  string
	Bopomofo2 string

	Definitions []Definition
}

type Definition struct {
	Def      string
	Quote    []string
	Example  []string
	DefType  string `json:"type"` // this field is called "type" in the output of the server
	Link     []string
	Synonyms string
	Antonyms string
}

// TODO implement Stringer interface for Entry

func FindEntry(word string) (*Entry, error) {
	var e Entry
	db := zhDicts.Db()

	// Find Entry
	eRow, err := db.Query("SELECT * FROM md_entries WHERE title = '" + word + "'")
	if err != nil {
		return nil, err
	}
	defer eRow.Close()

	var entry_id int
	if eRow.Next() {
		eRow.Scan(&entry_id, &e.Title, &e.Radical, &e.Stroke_count, &e.Non_radical_stroke_count)
	} else {
		return nil, nil
	}

	e.Heteronyms, err = findHeteronyms(entry_id)
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func findHeteronyms(entry_id int) ([]Heteronym, error) {
        db := zhDicts.Db()
	hs := make([]Heteronym, 0)

	// Find heteronyms
	hRows, err := db.Query(fmt.Sprintf("SELECT id, pinyin, bopomofo, bopomofo2 FROM md_heteronyms WHERE entry_id = %d ORDER BY idx", entry_id))
	if err != nil {
		return nil, err
	}
	defer hRows.Close()

	for hRows.Next() {
		var h Heteronym
		var heteronym_id int

		hRows.Scan(&heteronym_id, &h.Pinyin, &h.Bopomofo, &h.Bopomofo2)

		h.Definitions, err = findDefinitions(heteronym_id)
		if err != nil {
			return nil, err
		}

		hs = append(hs, h)
	}
	return hs, nil
}

func findDefinitions(heteronym_id int) ([]Definition, error) {
        db := zhDicts.Db()
	ds := make([]Definition, 0)

	dRows, err := db.Query(fmt.Sprintf("SELECT def, quotes, examples, type, link, synonyms, antonyms FROM md_definitions WHERE heteronym_id = %d ORDER BY idx", heteronym_id))
	if err != nil {
		return nil, err
	}
	defer dRows.Close()

	for dRows.Next() {
		var d Definition

		quotes := ""
		examples := ""
		links := ""

		dRows.Scan(&d.Def, &quotes, &examples, &d.DefType, &links, &d.Synonyms, &d.Antonyms)

		d.Quote = strings.Split(quotes, "|||")
		d.Example = strings.Split(examples, "|||")
		d.Link = strings.Split(links, "|||")

		ds = append(ds, d)
	}

	return ds, nil
}

func Initialize() error {
	return nil
}

// Implement chinese.ToHTMLer
func (e Entry) ToHTML(toneColors []string) string {
	var html string

	for _, heteronym := range e.Heteronyms {
		// title nice and large
		html += `<span style="font-family: Arial; font-size: 32px; color: #000000; white-space: pre-wrap;">` + e.Title + `</span>`

		html += "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;"

		// bopomofo
		html += bopomofo.Bop2Col(heteronym.Bopomofo, toneColors, "&nbsp;")

		html += "<br>"

		// definitions
		for _, definition := range heteronym.Definitions {
			nonEmptyDefinition := false
			if definition.Def != "" {
				nonEmptyDefinition = true

				html += "•"
				html += definition.Def
				html += "<br>"
			}

			// examples
			for _, example := range definition.Example {
				nonEmptyDefinition = true

				html += `<span style="color:#970000;">例</span>: `
				html += example
				html += "<br>"
			}

			// quotes
			for _, quote := range definition.Quote {
				nonEmptyDefinition = true

				html += `<span style="color:#BBBBBB;">`
				html += quote
				html += "</span><br>"
			}
			if nonEmptyDefinition {
				html += "<br>"
			}
		}

		// TODO Add more fields to html output
	}

	return html
}

package cedict

import (
	"github.com/jsoendermann/blueMandoTools/chinese"
	"github.com/jsoendermann/blueMandoTools/zhDicts"
)

type WordType int

const (
	WordTypeString = iota
	WordTypeRecords
)

type ChineseTextWord struct {
	T WordType
	S string
	R []Record
}

// Once implemented, this method will split a string of chinese text
// into a slice of words of type WordType.
func SplitChineseTextIntoWords(text string) ([]ChineseTextWord, error) {
	output := make([]ChineseTextWord, 0)

	charSet := DetermineCharSet(text)
	charSetString := ""
	if charSet == chinese.Trad {
		charSetString = "trad"
	} else {
		charSetString = "simp"
	}

	index := 0

	for index < len([]rune(text)) {
		// get next string of length maxRunecount
		var substring string
		if index+maxRunecount > len([]rune(text))-1 {
			substring = string([]rune(text)[index:])
		} else {
			substring = string([]rune(text)[index : index+maxRunecount])
		}

		sql := "SELECT * FROM cedict WHERE position(" + charSetString + " in '" + substring + "') = 1 AND length(" + charSetString + ") = (SELECT length(" + charSetString + ") FROM cedict WHERE position(" + charSetString + " in '" + substring + "') = 1 ORDER BY length(" + charSetString + ") DESC LIMIT 1)"

		rows, err := zhDicts.Db().Query(sql)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		// create slice to hold records
		records := make([]Record, 0)

		var runecount int
		// populate records with the data from the db query
		for rows.Next() {
			var id int
			var trad, simp, pinyin, english string
			rows.Scan(&id, &trad, &simp, &pinyin, &english, &runecount)
			records = append(records, Record{Trad: trad, Simp: simp, Pinyin: pinyin, English: english})
		}

		if len(records) == 0 {
			output = append(output, ChineseTextWord{T: WordTypeString, S: string([]rune(text)[index]), R: nil})
			index += 1
		} else {
			output = append(output, ChineseTextWord{T: WordTypeRecords, S: "", R: records})
			index += runecount
		}
	}

	return output, nil
}

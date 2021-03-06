package cedict

import (
	"github.com/jsoendermann/blueMandoTools/chinese"
)

// DetermineCharSet takes a string and returns a variable
// of type chinese.CharSet indicating whether the text is
// in simplified or traditional characters. In ambiguous cases,
// where all characters can be interpreted as both simplified and
// traditional (such as "你好"), this method returns chinese.Trad.
// It also returns chinese.Trad for characters where one simplified
// character maps to multiple traditional ones like "后".
// For texts that are more than 1-2 sentences in length, this method
// is usually very accurate.
func DetermineCharSet(text string) chinese.CharSet {
	// go through all the runes in the string and check for each
	// whether there's a simplified but no traditional match in
	// the db. If there is, the text is in simplified characters.
	for _, c := range text {

		// search for a traditional record first, if there is,
		// skip to the next rune
		hasTradRecord := false
		tradRecords, _ := FindRecords(string(c), chinese.Trad)
		if len(tradRecords) > 0 {
			hasTradRecord = true
		}

		// if there's no traditional record, search for simplified records
		if !hasTradRecord {
			simpRecords, _ := FindRecords(string(c), chinese.Simp)
			if len(simpRecords) > 0 {
				return chinese.Simp
			}
		}
	}

	return chinese.Trad
}

// TODO improve this function by always choosing the simpler
//      character when there are multiple records
func convertBetweenCharSets(text string, conversionTarget chinese.CharSet) (string, error) {
	if DetermineCharSet(text) == conversionTarget {
		return text, nil
	}

	var conversionOrigin chinese.CharSet
	if conversionTarget == chinese.Simp {
		conversionOrigin = chinese.Trad
	} else {
		conversionOrigin = chinese.Simp
	}

	t, err := chinese.SplitChineseTextIntoWords(text, conversionOrigin)
	if err != nil {
		return "", err
	}

	// turn t into a string
	output := ""
	for _, w := range t {
		records, _ := FindRecords(w, conversionOrigin)
		if len(records) == 0 {
			output += w
		} else {
			output += records[0].WordByCharSet(conversionTarget)
		}
	}

	return output, nil

}

func Simp2Trad(simp string) (string, error) {
	return convertBetweenCharSets(simp, chinese.Trad)
}

func Trad2Simp(trad string) (string, error) {
	return convertBetweenCharSets(trad, chinese.Simp)
}

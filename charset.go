package cedict

import (
  "github.com/yangchuanzhang/chinese"
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
  for _,c := range text {

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

// TODO comment this function
func Simp2Trad(simp string) (string, error) {
  t := SplitChineseTextIntoWords(simp)
  output := ""

  for _,w := range t {
    if w.T == WordTypeString {
      output += w.S
    } else {
      output += w.R[0].Trad
    }
  }

  return output, nil
}


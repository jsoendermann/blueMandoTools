package pinyin

import (
  "testing"
  "unicode/utf8"
)

func TestIndexOfRune(t *testing.T) {
  var tests = []struct {
    letter string
    index int
  }{
    {"a",0},
    {"ü",5},
    {"U",10},
    {"Ü",11},
    {"z",-1},
    {"",-1},
  }

  for _,u := range tests {
    letterRune, _ := utf8.DecodeRuneInString(u.letter)
    i := indexOfRune(letterRune)

    if i != u.index {
      t.Errorf("indexOfRune(%q) == %q, want %q", u.letter, i, u.index)
    }
  }
}

func TestDiacriticRuneForRuneAndTone(t *testing.T) {
  var tests = []struct {
    letter string
    tone int
    diacriticLetter string
  }{
    {"a",1,"ā"},
    {"ü",3,"ǚ"},
    {"Ü",2,"Ǘ"},
  }

  for _,u := range tests {
    letterRune, _ := utf8.DecodeRuneInString(u.letter)
    diacriticLetterRune, _ := utf8.DecodeRuneInString(u.diacriticLetter)

    result := diacriticRuneForRuneAndTone(letterRune, u.tone)

    if  result != diacriticLetterRune {
      t.Errorf("diacriticRuneForRuneAndTone(%q, %d) == %q, want %q", u.letter, u.tone, result, u.diacriticLetter)
    }
  }
}

//TODO write tests for indexOfLastVowel
//func TestIndexOfLastVowel(t *testing.T) {
//
//
//
//}

//TODO write tests for runeByIndex
//TODO write tests for substringRuneIndex
//TODO write tests for runecount

func TestAddDiacritics(t *testing.T) {
  var tests = []struct {
    syllable string
    tone int
    syllableWithDiacritic string
  }{
    {"John", -1, "John"},
    {"ni", 3, "nǐ"},
    {"bu", 4, "bù"},
    {"yao", 4, "yào"},
    {"Bei", 3, "Běi"},
    {"jing", 1, "jīng"},
    {"ma", 0, "ma"},
    {"gou", 3, "gǒu"},
    {"shuang", 1, "shuāng"},
    {"lve", 4, "lvè"},
  }

  for _,u := range tests {
    result := addDiacritic(u.syllable, u.tone)

    if result != u.syllableWithDiacritic {
      t.Errorf("addDiacritic(%q, %d) == %q, want %q", u.syllable, u.tone, result, u.syllableWithDiacritic)
    }
  }

}

func TestNumbersToDiacritics(t *testing.T) {
  var tests = []struct {
    numbers, diacritics string
  }{
    {"ni3", "nǐ"},
    //TODO more tests
  }

  for _,u := range tests {
    result := NumbersToDiacritics(u.numbers)

    if result != u.diacritics {
      t.Errorf("NumbersToDiacritics(%q) == %q, want %q", u.numbers, result, u.diacritics)
    }
  }
}

//TODO write tests for splitNumbersString
//TODO write tests for NumbersToDiacriticsAndHtmlColors

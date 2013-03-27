package pinyin

import (
  "testing"
  "unicode/utf8"
)

func TestIndexOfRuneInTonesArray(t *testing.T) {
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
    i := indexOfRuneInTonesArray(letterRune)

    if i != u.index {
      t.Errorf("indexOfRuneInTonesArray(%q) == %q, want %q", u.letter, i, u.index)
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
    {"e",0,"e"},
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

func TestIndexOfLastVowel(t *testing.T) {
  var tests = []struct {
    s string
    i int
  }{
    {"bba",2},
    {"a",0},
    {"züzüü",6},
    {"",-1},
    {"bb",-1},
  }

  for _,u := range tests {
    result := indexOfLastVowel(u.s)

    if result != u.i {
      t.Errorf("indexOfLastVowel(%q) == %d, want %d", u.s, result, u.i)
    }
  }
}





//TODO write tests for runeByIndex

func TestAddDiacritics(t *testing.T) {
  var tests = []struct {
    w word
    syllableWithDiacritic string
  }{
    {word{syllable: "John",   tone: -1}, "John"},
    {word{syllable: "ni",     tone: 3}, "nǐ"},
    {word{syllable: "bu",     tone: 4}, "bù"},
    {word{syllable: "yao",    tone: 4}, "yào"},
    {word{syllable: "Bei",    tone: 3}, "Běi"},
    {word{syllable: "jing",   tone: 1}, "jīng"},
    {word{syllable: "ma",     tone: 0}, "ma"},
    {word{syllable: "gou",    tone: 3}, "gǒu"},
    {word{syllable: "shuang", tone: 1}, "shuāng"},
    {word{syllable: "lve",    tone: 4}, "lvè"},
    {word{syllable: "lüe",    tone: 4}, "lüè"},
  }

  for _,u := range tests {
    result := addDiacritic(u.w)

    if result != u.syllableWithDiacritic {
      t.Errorf("addDiacritic(%q) == %q, want %q", u.w, result, u.syllableWithDiacritic)
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

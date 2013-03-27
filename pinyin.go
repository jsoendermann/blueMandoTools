package pinyin

import (
  "strconv"
  "strings"
  "unicode"
  "unicode/utf8"
)

// This structure holds a syllable and a tone number between
// -1 (no tone) and 4 with 0 being the neutral tone
type word struct {
  syllable string
  tone     int
}

// This array holds all the necessary combination of diacritics and vowels
var tones = [][]string{
  []string{"a", "e", "i", "o", "u", "ü", "A", "E", "I", "O", "U", "Ü"},
  []string{"ā", "ē", "ī", "ō", "ū", "ǖ", "Ā", "Ē", "Ī", "Ō", "Ū", "Ǖ"},
  []string{"á", "é", "í", "ó", "ú", "ǘ", "Á", "É", "Í", "Ó", "Ú", "Ǘ"},
  []string{"ǎ", "ě", "ǐ", "ǒ", "ǔ", "ǚ", "Ǎ", "Ě", "Ǐ", "Ǒ", "Ǔ", "Ǚ"},
  []string{"à", "è", "ì", "ò", "ù", "ǜ", "À", "È", "Ì", "Ò", "Ù", "Ǜ"},
}

// indexOfRune retuns the index of r in the array tones[0], i.e. the column
// of the vowel, or -1 if tones[0] doesn't contain r.
func indexOfRuneInTonesArray(r rune) int {
  for i, v := range tones[0] {
    if vr, _ := utf8.DecodeRuneInString(v); vr == r {
      return i
    }
  }
  return -1
}

// diacriticRuneForRuneAndTone combines r and t into a diacritic letter.
// t must be between 0 and 4
func diacriticRuneForRuneAndTone(r rune, t int) rune {
  // get column of r in tones
  i := indexOfRuneInTonesArray(r)

  for j, v := range tones[t] {
    if j == i {
      r, _ := utf8.DecodeRuneInString(v)
      return r
    }
  }
  //FIXME return error
  return r
}

// indexOfLastVowel retuns the index of the last vowel (including ü) in s or
// -1 if s contains no vowel. 
func indexOfLastVowel(s string) int {
  vowels := "aeiouüAEIOUÜ"
  index := -1

  for i, v := range s {
    if strings.Index(vowels, string(v)) != -1 {
      index = i
    }
  }

  return index
}

// addDiacritic adds correct diacritic based on the tone t to the correct vowel in
// syllable s and retuns the result
// Pinyin rules for adding diacritics:
// 1. If s contains an 'a' or an 'e', it takes the diacritic
// 2. Otherwise, if s contains 'ou', the 'o' takes the tone mark
// 3. Otherwise, the last vowel takes the tone mark
func addDiacritic(w word) string {
  if w.tone < 1 || w.tone > 4 {
    return w.syllable
  }

  var result string

  // check for 'a' and 'e' by trying to replace them
  for _, v := range "aeAE" {
    result = strings.Replace(w.syllable, string(v), string(diacriticRuneForRuneAndTone(v, w.tone)), -1)
    if result != w.syllable {
      return result
    }
  }

  // check for ou by trying to replace it
  oRune, _ := utf8.DecodeRuneInString("o")
  ORune, _ := utf8.DecodeRuneInString("O")

  result = strings.Replace(w.syllable, "ou", string(diacriticRuneForRuneAndTone(oRune, w.tone))+"u", -1)
  if result != w.syllable {
    return result
  }
  result = strings.Replace(w.syllable, "Ou", string(diacriticRuneForRuneAndTone(ORune, w.tone))+"u", -1)
  if result != w.syllable {
    return result
  }

  // put diacritic on last vowel in string
  lastVowelIndex := indexOfLastVowel(w.syllable)
  //FIXME deal with error
  lastVowelRune, _ := utf8.DecodeRuneInString(w.syllable[lastVowelIndex:]) 
  lastVowelDiacriticRune := diacriticRuneForRuneAndTone(lastVowelRune, w.tone)

  result = w.syllable[0:lastVowelIndex]
  result += string(lastVowelDiacriticRune)
  result += w.syllable[lastVowelIndex+1 : utf8.RuneCountInString(w.syllable)]

  return result
}

// splitNumbersString splits a string of pinyin into []word
// FIXME preserve whitespace
func splitNumbersString(pinyin string) []word {
  words := strings.Fields(pinyin)

  splitPinyin := make([]word, len(words))

  for i, v := range words {
    lastChar := v[len(v)-1:]

    if unicode.IsNumber(rune(lastChar[0])) {
      n, _ := strconv.Atoi(lastChar)

      // convert 5's to 0's for neutral tones
      if n == 5 {
        n = 0
      }
      splitPinyin[i] = word{syllable: v[0 : len(v)-1], tone: n}
    } else {
      splitPinyin[i] = word{syllable: v, tone: -1}
    }
  }
  return splitPinyin
}

func addHTMLColors(s, color string) string {
  return `<span style="color:` + color + `;">` + s + "</span>"
}


// Num2DiaCol takes a pinyin string and a slice of html colors, converts
// the pinyin to diacritics and adds html for the colors
func Num2DiaCol(pinyinWithNumbers string, colors []string) string {
  pinyinWithNumbers = strings.Replace(pinyinWithNumbers, "v", "ü", -1)
  pinyinWithNumbers = strings.Replace(pinyinWithNumbers, "V", "Ü", -1)

  splitString := splitNumbersString(pinyinWithNumbers)

  s := ""

  for _, w := range splitString {
    t := addDiacritic(w)
    if colors != nil {
      if w.tone != -1 {
        t = addHTMLColors(t, colors[w.tone])
      }
    }
    s += t
  }
  return s
}

func Num2Dia(pinyinWithNumbers string) string {
  return Num2DiaCol(pinyinWithNumbers, nil)
}



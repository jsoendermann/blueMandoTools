package pinyin

import (
  //"fmt"
  "strings"
  "strconv"
  "unicode"
  "unicode/utf8"
  "github.com/jsoendermann/util"
)

type word struct {
  syllable string
  tone int
}

var tones = [][]string{
[]string{"a", "e", "i", "o", "u", "ü", "A", "E", "I", "O", "U", "Ü"}, 
[]string{"ā", "ē", "ī", "ō", "ū", "ǖ", "Ā", "Ē", "Ī", "Ō", "Ū", "Ǖ"},
[]string{"á", "é", "í", "ó", "ú", "ǘ", "Á", "É", "Í", "Ó", "Ú", "Ǘ"},
[]string{"ǎ", "ě", "ǐ", "ǒ", "ǔ", "ǚ", "Ǎ", "Ě", "Ǐ", "Ǒ", "Ǔ", "Ǚ"},
[]string{"à", "è", "ì", "ò", "ù", "ǜ", "À", "È", "Ì", "Ò", "Ù", "Ǜ"},
}

// indexOfRune retuns the index of r in the array tones[0].
func indexOfRune(r rune) int {
  for i,v := range tones[0] {
    if vr,_ := utf8.DecodeRuneInString(v); vr == r {
      return i
    }
  }
  return -1
}

// diacriticRuneForRuneAndTone combines r and t into a diacritic letter.
func diacriticRuneForRuneAndTone(r rune, t int) rune {
  // get horizontal index of r in tones
  i := indexOfRune(r)

  for j,v := range tones[t] {
    if j == i {
      r,_ := utf8.DecodeRuneInString(v)
      return r
    }
  }
  //FIXME return error
  return r
}

// indexOfLastVowel retuns the index of the last vowel (including ü) in s or
// -1 if s contains no vowel.
//FIXME this retunes rune index not byte position
func indexOfLastVowel(s string) int {
  vowels := "aeiouüAEIOUÜ"
  index := -1

  for i,v := range s {
    if strings.Index(vowels, string(v)) != -1 {
      index = i
    }
  }

  return index
}

// This function retuns the rune at rune index i
func runeByIndex(s string, i int) rune {
  for j,v := range s {
    if j == i {
      return v
    }
  }
  //TODO maybe throw an error?
  panic("Rune index out of bounds")
}




// TODO dry these two functions
func NumbersToDiacriticsAndHtmlColors(pinyinWithNumbers string, colors []string) {
  pinyinWithNumbers = strings.Replace(pinyinWithNumbers, "v", "ü", -1)
  pinyinWithNumbers = strings.Replace(pinyinWithNumbers, "V", "Ü", -1)

  splitString := splitNumbersString(pinyinWithNumbers)

  s := ""

  for _, v := range splitString {
    t := addDiacritic(v.syllable, v.tone)
    if v.tone != -1 {
      t = colorize(t, colors[v.tone])
    }
    s += t
  }
}

func NumbersToDiacritics(pinyinWithNumbers string) string {
  pinyinWithNumbers = strings.Replace(pinyinWithNumbers, "v", "ü", -1)
  pinyinWithNumbers = strings.Replace(pinyinWithNumbers, "V", "Ü", -1)

  splitString := splitNumbersString(pinyinWithNumbers)

  s := ""

  for _,v := range splitString {
    s += addDiacritic(v.syllable, v.tone)
    s += " "
  }

  return strings.TrimSpace(s)
}

func colorize(s, color string) string {
  return `<span style="color:`+color+`;">`+s+"</span>"
}


// addDiacritic adds correct diacritic based on the tone t to the correct vowel in
// syllable s and retuns the result
// Pinyin rules for adding diacritics:
// 1. If s contains an 'a' or an 'e', it takes the diacritic
// 2. Otherwise, if s contains 'ou', the 'o' takes the tone mark
// 3. Otherwise, the last vowel takes the tone mark
func addDiacritic(s string, t int) string {
  if t < 1 || t > 4 {
    return s
  }

  var result string

  // check for 'a' and 'e' by trying to replace them
  for _,v := range "aeAE" {
    result = strings.Replace(s, string(v), string(diacriticRuneForRuneAndTone(v,t)), -1)
    if result != s {
      return result
    }
  }

  // check for ou by trying to replace it
  oRune, _ := utf8.DecodeRuneInString("o")
  ORune, _ := utf8.DecodeRuneInString("O")

  result = strings.Replace(s, "ou", string(diacriticRuneForRuneAndTone(oRune,t))+"u", -1)
  if result != s {
    return result
  }
  result = strings.Replace(s, "Ou", string(diacriticRuneForRuneAndTone(ORune,t))+"u", -1)
  if result != s {
    return result
  }

  
  lastVowelIndex := indexOfLastVowel(s)
  lastVowelRune := runeByIndex(s, lastVowelIndex) 
  lastVowelDiacriticRune := diacriticRuneForRuneAndTone(lastVowelRune, t)

  result = util.SubstringByRuneIndex(s, 0, lastVowelIndex)
  result += string(lastVowelDiacriticRune)
  result += util.SubstringByRuneIndex(s, lastVowelIndex + 1, utf8.RuneCountInString(s))

  return result
}

//FIXME return tones from -1 to 4
func splitNumbersString(pinyin string) []word {
  words := strings.Fields(pinyin)

  splitPinyin := make([]word, len(words))

  for i,v := range words {
    lastChar := v[len(v)-1:]


    if unicode.IsNumber(rune(lastChar[0])) {
      n,_ := strconv.Atoi(lastChar)
      splitPinyin[i] = word{syllable: v[0:len(v)-1], tone: n}
    } else {
      splitPinyin[i] = word{syllable: v, tone: -1}
    }
  }
  return splitPinyin
}

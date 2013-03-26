package pinyin

import (
  "fmt"
  "strings"
  "strconv"
  "unicode"
  //"unicode/utf8"
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

func indexOfLetter(letter string) int {
  for i,v := range tones[0] {
    if letter == v {
      return i
    }
  }
  return -1
}

func diacriticLetterForLetterAndTone(letter string, tone int) string {
  i := indexOfLetter(letter)
  return tones[tone][i]
}

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



func Numbers2diacritics(pinyinWithNumbers string) string {
  pinyinWithNumbers = strings.Replace(pinyinWithNumbers, "v", "ü", -1)
  pinyinWithNumbers = strings.Replace(pinyinWithNumbers, "V", "Ü", -1)

  splitString := splitNumbersString(pinyinWithNumbers)

  s := ""

  for _,v := range splitString {
    if v.tone == -1 || v.tone == 5 {
      s += v.syllable
    } else {
      s += addDiacritic(v.syllable, v.tone)
    }
    s += " "
  }

  return s
}

// t within [1..4]
func addDiacritic(s string, t int) string {
  var newS string
  for _,v := range []string{"a", "e", "A", "E"} {
    newS = strings.Replace(s, v, diacriticLetterForLetterAndTone(v,t), -1)
    if newS != s {
      return newS
    }
  }
  newS = strings.Replace(s, "ou", diacriticLetterForLetterAndTone("o",t)+"u", -1)
  if newS != s {
    return newS
  }
  newS = strings.Replace(s, "Ou", diacriticLetterForLetterAndTone("O",t)+"u", -1)
  if newS != s {
    return newS
  }

  lastVowelIndex := indexOfLastVowel(s)
  fmt.Println(lastVowelIndex)
  lastVowel := string(s[lastVowelIndex])

  fmt.Println(lastVowelIndex, lastVowel)

  newS = strings.Replace(s, lastVowel, diacriticLetterForLetterAndTone(lastVowel,t), -1)

  
  return newS
}

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

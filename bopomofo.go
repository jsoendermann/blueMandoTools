/*
Package bopomofo implements adding colors to bopomofo strings.
*/
package bopomofo

import (
  "strings"
)

const (
  neutralRune = "˙"
  secondRune = "ˊ"
  thirdRune = "ˇ"
  forthRune = "ˋ"
)

// Bop2Col adds html colors to bopomofo.
func Bop2Col(bopomofo string, colors []string, separator string) string {
  // FIXME preserve whitespace
  words := strings.Fields(bopomofo)

  var output string

  for i,word := range words {
    // FIXME check if word is a valid bopomofo string
    if strings.Index(word, neutralRune) != -1 {
      output += addHTMLColors(word, colors[0])
    } else if strings.Index(word, secondRune) != -1 {
      output += addHTMLColors(word, colors[2])
    } else if strings.Index(word, thirdRune) != -1 {
      output += addHTMLColors(word, colors[3])
    } else if strings.Index(word, forthRune) != -1 {
      output += addHTMLColors(word, colors[4])
    } else {
      output += addHTMLColors(word, colors[1])
    }

    if i < len(words) - 1 {
      output += separator
    }
  }

  return output
}

func addHTMLColors(s, color string) string {
  return `<span style="color:` + color + `;">` + s + "</span>"
}

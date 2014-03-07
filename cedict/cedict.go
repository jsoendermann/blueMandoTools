/*
TODO: add package description
*/
package cedict

import (
  _ "github.com/mattn/go-sqlite3"
  "fmt"
  "github.com/jsoendermann/blueMandoTools/chinese"
  "github.com/jsoendermann/blueMandoTools/pinyin"
  "github.com/jsoendermann/blueMandoTools/zhDicts"
)

type Record struct {
  Simp string
  Trad string
  Pinyin string
  English string
}

func (r Record) WordByCharSet(c chinese.CharSet) string {
  if c == chinese.Trad {
    return r.Trad
  }
  return r.Simp
}

// This method implements the Stringer interface for Record
func (r Record) String() string {
  return fmt.Sprintf("[simp: %q  trad: %q  pinyin: %q  english: %q]", r.Simp, r.Trad, r.Pinyin, r.English)
}

// Implement chinese.ToHTMLer
func (r Record) ToHTML(toneColors []string) string {
  var html string

  html += r.Trad
  html += "&nbsp;&nbsp;&nbsp; "
  html += pinyin.Num2DiaCol(r.Pinyin, toneColors, "&nbsp;")
  html += "&nbsp;&nbsp;&nbsp; "
  html += r.English

  return html
}

var maxRunecount int

func Initialize() error {
    db := zhDicts.Db()
    if db == nil {
        return fmt.Errorf("Database not loaded")
    }

    // get max runecount
    sqlMaxRunecount := "SELECT MAX(runecount) AS maxRunecount FROM cedict"

    rows, err := db.Query(sqlMaxRunecount)
    if err != nil {
      return err
    }
    defer rows.Close()

    rows.Next()
    rows.Scan(&maxRunecount)

    return nil
}



// FindRecords searches the database of cedict records and returns a slice of type
// []Record and an error. It returns an empty slice if no matches could be found.
func FindRecords(word string, charSet chinese.CharSet) ([]Record, error) {
  // construct db query based on charSet
  sql := "SELECT * FROM cedict "

  switch charSet {
    case chinese.Trad: 
    sql += "WHERE trad LIKE '"+word+"'"
    case chinese.Simp: 
    sql += "WHERE simp LIKE '"+word+"'"
  default:
    return nil, fmt.Errorf("cedict: unrecognized character set")
  }

  // execute the query and defer closing it
  rows, err := zhDicts.Db().Query(sql)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  // create slice to hold records
  records := make([]Record, 0)

  // populate records with the data from the db query
  for rows.Next() {
    var id, runecount int
    var trad, simp, pinyin, english string
    rows.Scan(&id, &trad, &simp, &pinyin, &english, &runecount)
    records = append(records, Record{Trad: trad, Simp: simp, Pinyin: pinyin, English: english})
  }

  return records, nil
}



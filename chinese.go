// Package chinese defines constants and other useful stuff that other
// chinese-related packages rely on.
package chinese


type CharSet int

const (
    Trad = iota
    Simp
)

// Different types of dictionary records can implement this interface
// to provide a way to convert the data they hold to html for use on
// a card. The colors must be in the form "#FFFFFF"
// TODO This is commented out because in its current form, this interface
// makes no sense. The differences between dictionaries and between the
// desired html output depending on the use case are too big to unify them
// like this.
//type ToHTMLer interface {
//  ToHTML(toneColors []string) string
//}

// TODO Maybe add a ToCarder interface that converts records to 
// csv rows.

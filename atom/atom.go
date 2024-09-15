package atom

import "encoding/xml"

type Atom struct {
	XMLName   xml.Name  `xml:"feed"`
	Title     string    `xml:"title"`
	Author    string    `xml:"author>name"`
	Published string    `xml:"published"`
	Updated   string    `xml:"updated"`
	Hyperlink Hyperlink `xml:"link"`
	Entries   []Entry   `xml:"entry"`
}

type Entry struct {
	Title     string    `xml:"title"`
	Hyperlink Hyperlink `xml:"link"`
	Published string    `xml:"published"`
	Updated   string    `xml:"updated"`
}

type Hyperlink struct {
	Href string `xml:"href,attr"`
}

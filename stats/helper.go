package stats

import (
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type entryIterator func(*goquery.Selection) (Entry, error)

// EachWithError iterates over a set of HTML nodes with a provided function
func EachWithError(s *goquery.Selection, f entryIterator) ([]Entry, error) {
	var err error

	el := make([]Entry, len(s.Nodes))

	for i, n := range s.Nodes {
		ns := &goquery.Selection{[]*html.Node{n}, s.Document, nil}
		el[i], err = f(ns)
		if err != nil {
			return el, err
		}
	}
	return el, nil
}

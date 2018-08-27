package stats

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type entryIterator func(*goquery.Selection) (Entry, error)

func EachWithError(s *goquery.Selection, f entryIterator) ([]Entry, error) {
	var err error

	el := make([]Entry, len(s.Nodes))

	for i, n := range s.Nodes {
		ns := &goquery.Selection{[]*html.Node{n}, s.document, nil}
		el[i], err = f(ns)
		if err != nil {
			return el, err
		}
	}
	return el, nil
}

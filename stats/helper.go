package stats

import (
	"github.com/PuerkitoBio/goquery"
)

type entryIterator func(*goquery.Selection) (Entry, error)

// EachWithError iterates over a set of HTML nodes with a provided function
func EachWithError(selectionSet *goquery.Selection, f entryIterator) ([]Entry, error) {
	var err error
	el := make([]Entry, len(selectionSet.Nodes))
	selectionSet.EachWithBreak(func(i int, s *goquery.Selection) bool {
		el[i], err = f(s)
		return err != nil
	})
	return el, err
}

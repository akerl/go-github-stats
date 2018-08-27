package stats

import (
	"time"
)

const (
	magicNumber = 3.77972616981
)

// Entry describes a date/score pair
type Entry struct {
	Date  time.Time
	Score int
}

// User describes the GitHub stats for a user
type User struct {
	Name    string
	Entries []Entry
}

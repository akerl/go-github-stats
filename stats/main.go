package stats

// User describes the GitHub stats for a user
type User struct {
	Name string
}

// LookupUser loads a user given their name
func LookupUser(name string) (User, error) {
	return User{}, nil
}

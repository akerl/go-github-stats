package stats

import (
	"fmt"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/tcnksm/go-gitconfig"
)

const (
	defaultGithubHost = "https://github.com"
	githubPathFormat  = "users/%s/contributions"
	dateFormat        = "%Y/%m/%d"
)

// LookupUser loads a user given their name
func LookupUser(name string) (User, error) {
	var err error

	if name == "" {
		name, err = getUserName()
		if err != nil {
			return User{}, err
		}
	}

	return User{}, nil
}

func getUserName() (string, error) {
	user, err := gitconfig.GithubUser()
	if err == nil && user != "" {
		return user, nil
	}
	user := os.GetEnv("USER")
	if user != "" {
		return user, nil
	}
	return "", fmt.Errorf("No user given and lookup failed")
}

func getURL(name string) string {
	host := os.GetEnv("GITHUB_URL")
	if host == "" {
		host = defaultGithubHost
	}
	path := fmt.Sprintf(githubPathFormat, name)
	url := fmt.Sprintf("%s/%s", host, path)
	return url, nil
}

func getResponse(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			panic(err)
		}
	}()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Failed loading url %s (%d)", url, resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func entryFromDiv(s *goquery.Selection) (Entry, error) {
	strDate, ok := s.Attr("data-date")
	if !ok {
		return Entry{}, fmt.Errorf("No date found in attr")
	}
	score, ok := s.Attr("data-count")
	if !ok {
		return Entry{}, fmt.Errorf("No count found in attr")
	}

	date, err := time.Parse(dateFormat, strDate)
	if err != nil {
		return Entry{}, err
	}
	e := Entry{
		Date:  date,
		Score: score,
	}
	return e, nil
}

func getContribData(name string) (User, error) {
	baseURL := getURL(name)

	// TODO: Handle longer-than-1-year streaks
	page, err := getResponse(baseURL)
	if err != nil {
		return User{}, err
	}
	reader := bytes.NewReader(page)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return User{}, err
	}

	divs := doc.Find(".day")
	el, err := EachWithError(divs, entriesFromDivFunc)
	if err != nil {
		return User{}, err
	}

	u := User{
		Name:    name,
		Entries: el,
	}
	return u, nil
}

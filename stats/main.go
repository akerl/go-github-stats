package stats

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	gitconfig "github.com/tcnksm/go-gitconfig"
)

const (
	defaultGithubHost = "https://github.com"
	githubPathFormat  = "users/%s/contributions"
	dateFormat        = "%Y/%m/%d"
)

// LookupUser loads a user given their name
func LookupUser(inputName string) (User, error) {
	name, err := getUserName(inputName)
	if err != nil {
		return User{}, err
	}

	fmt.Printf("%s\n", defaultGithubHost)
	fmt.Printf("%s\n", githubPathFormat)
	fmt.Printf("%s\n", dateFormat)
	fmt.Printf("%f\n", magicNumber)

	return getContribData(name)
}

func getUserName(inputName string) (string, error) {
	if inputName != "" {
		return inputName, nil
	}
	name, err := gitconfig.GithubUser()
	if err == nil && name != "" {
		return name, nil
	}
	name = os.Getenv("USER")
	if name != "" {
		return name, nil
	}
	return "", fmt.Errorf("no user given and lookup failed")
}

func getURL(name string) string {
	host := os.Getenv("GITHUB_URL")
	if host == "" {
		host = defaultGithubHost
	}
	path := fmt.Sprintf(githubPathFormat, name)
	url := fmt.Sprintf("%s/%s", host, path)
	return url
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
		return "", fmt.Errorf("failed loading url %s (%d)", url, resp.StatusCode)
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
		return Entry{}, fmt.Errorf("no date found in attr")
	}
	scoreStr, ok := s.Attr("data-count")
	if !ok {
		return Entry{}, fmt.Errorf("no count found in attr")
	}
	score, err := strconv.Atoi(scoreStr)
	if err != nil {
		return Entry{}, fmt.Errorf("could not parse score: %s", scoreStr)
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
	reader := bytes.NewReader([]byte(page))
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return User{}, err
	}

	divs := doc.Find(".day")
	el, err := EachWithError(divs, entryFromDiv)
	if err != nil {
		return User{}, err
	}

	u := User{
		Name:    name,
		Entries: el,
	}
	return u, nil
}

package gotrakt

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"text/template"
	"time"

	"github.com/golang/glog"
	"github.com/hobeone/gotrakt/httpclient"
	"github.com/jmcvetta/napping"
)

//https://trakt.tv/api-docs/search-shows
var ShowSearchTempl = template.Must(
	template.New("ShowSearch").Parse("{{.Host}}/search/shows.json/{{.APIKey|urlquery}}?query={{.Query | urlquery}}&seasons=true"),
)

//http://trakt.tv/api-docs/show-summary
var ShowSummaryTmpl = template.Must(
	template.New("ShowSummary").Parse("{{.Host}}/show/summary.json/{{.APIKey|urlquery}}/{{.Query | urlquery}}/extended"),
)

//https://trakt.tv/api-docs/show-season
var ShowSeasonTmpl = template.Must(
	template.New("ShowSeason").Parse("{{.Host}}/show/season.json/{{.APIKey|urlquery}}/{{.Query | urlquery}}/{{.Season | urlquery}}"),
)

// https://trakt.tv/api-docs/search-movies
var MovieSearchTmpl = template.Must(
	template.New("MovieSearch").Parse("{{.Host}}/search/movies.json/{{.APIKey|urlquery}}?query={{.Query | urlquery}}"),
)

// https://trakt.tv/api-docs/movie-summary
var MovieSummaryTmpl = template.Must(
	template.New("MovieSummary").Parse("{{.Host}}/movie/summary.json/{{.APIKey|urlquery}}/query={{.Query | urlquery}}"),
)

// Base URL for TraktTV api
const TraktTVBaseURL = "https://api.trakt.tv"

// TraktTV is the main struct used to query Trakt.tv.  Use NewTraktTV to
// create new instances.
type TraktTV struct {
	APIKey   string
	BaseURL  string
	Session  *napping.Session
	Userinfo *url.Userinfo
}

type option func(*TraktTV)

// New initializes and returns a new TraktTV struct
func New(api string, options ...option) (*TraktTV, error) {
	t := &TraktTV{
		APIKey:  api,
		BaseURL: TraktTVBaseURL,
		Session: &napping.Session{
			Log: false,
			Client: httpclient.NewTimeoutClient(
				httpclient.ConnectTimeout(10*time.Second),
				httpclient.ReadWriteTimeout(10*time.Second),
			),
		},
	}
	for _, opt := range options {
		opt(t)
	}
	return t, nil
}

// Session sets the session to use for talking to TraktTV
func Session(sess *napping.Session) option {
	return func(t *TraktTV) {
		t.Session = sess
	}
}

// Host sets the host to use for talking to TraktTV
// This includes the protocol, hostname, port:
// i.e. https://api.trakt.tv:443
func Host(host string) option {
	return func(t *TraktTV) {
		t.BaseURL = host
	}
}

// Userinfo configures the user and password information for API calls
func Userinfo(username, password string) option {
	h := sha1.New()
	io.WriteString(h, password)
	hashedPW := hex.EncodeToString(h.Sum(nil))

	return func(t *TraktTV) {
		t.Userinfo = url.UserPassword(username, hashedPW)
	}
}

func (t *TraktTV) getWithErrorCheck(url string, result interface{}) error {
	glog.Infof("Get query for %s\n", url)
	apiErr := &APIError{}
	t.Session.Userinfo = t.Userinfo
	response, err := t.Session.Get(url, &napping.Params{}, result, apiErr)
	if serr, ok := err.(*json.SyntaxError); ok {
		line, col, highlight := HighlightBytePosition(response.HttpResponse().Body, serr.Offset)
		return fmt.Errorf("gotrackt: syntax error in response at line %d, column %d (file offset %d):\n%s", line, col, serr.Offset, highlight)
	}
	if apiErr.Status != "" {
		return apiErr
	}
	return err
}

func (t *TraktTV) getURLFromTemplate(tmpl *template.Template, args map[string]string) (string, error) {
	args["APIKey"] = t.APIKey
	args["Host"] = t.BaseURL
	out := bytes.Buffer{}
	err := tmpl.Execute(&out, args)
	return out.String(), err
}

// GetShow returns a show and all of it's Seasons and Episodes
func (t *TraktTV) GetShow(slugOrTvdbID string) (*Show, error) {
	args := map[string]string{
		"Query": slugOrTvdbID,
	}

	apiURL, err := t.getURLFromTemplate(ShowSummaryTmpl, args)
	result := &Show{}
	if err != nil {
		return result, err
	}

	err = t.getWithErrorCheck(apiURL, result)
	return result, err
}

// ShowSearch searches tv shows
func (t *TraktTV) ShowSearch(name string) ([]Show, error) {
	args := map[string]string{
		"Query": name,
	}
	result := []Show{}
	apiURL, err := t.getURLFromTemplate(ShowSearchTempl, args)
	if err != nil {
		return result, err
	}
	err = t.getWithErrorCheck(apiURL, &result)
	return result, err
}

//ShowSeasons searches for a shows episode summaries by season.  You can
//optionally limit the query to just a given set of seasons.
func (t *TraktTV) ShowSeasons(slugOrTvdbID string, seasons []int) ([]Season, error) {
	results := make([]Season, len(seasons))
	if len(seasons) == 0 {
		return results, fmt.Errorf("must specify Which Seasons to get")
	}
	for i, season := range seasons {
		results[i] = Season{
			Season:   season,
			Episodes: []Episode{},
		}

		args := map[string]string{
			"Query":  slugOrTvdbID,
			"Season": fmt.Sprintf("%d", season),
		}
		apiURL, err := t.getURLFromTemplate(ShowSeasonTmpl, args)
		if err != nil {
			return results, err
		}

		err = t.getWithErrorCheck(apiURL, &results[i].Episodes)
		if err != nil {
			return results, err
		}
	}
	return results, nil
}

//MovieSearch searches Trakt.tv for movies matching the query
func (t *TraktTV) MovieSearch(query string) ([]Movie, error) {
	args := map[string]string{
		"Query": query,
	}
	apiURL, err := t.getURLFromTemplate(MovieSearchTmpl, args)
	res := []Movie{}
	if err != nil {
		return res, err
	}
	err = t.getWithErrorCheck(apiURL, &res)
	return res, err
}

//GetMovie searches Trakt.tv for movies matching the query
func (t *TraktTV) GetMovie(slugOrImdbID string) (*Movie, error) {
	res := &Movie{}
	args := map[string]string{
		"Query": slugOrImdbID,
	}
	apiURL, err := t.getURLFromTemplate(MovieSummaryTmpl, args)
	if err != nil {
		return res, err
	}
	err = t.getWithErrorCheck(apiURL, &res)
	return res, err
}

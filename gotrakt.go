package gotrakt

import (
	"encoding/json"
	"fmt"

	"github.com/golang/glog"
	"github.com/hobeone/gotrakt/httpclient"
	"github.com/jmcvetta/napping"
)

//TODO: maybe change these to text.Template for named paramters?
//https://trakt.tv/api-docs/search-shows
const TvSearchURL = "/search/shows.json/%s?query=%s&limit=%d&seasons=true"

//http://trakt.tv/api-docs/show-summary
const TvSummaryURL = "/show/summary.json/%s/%s/extended"

//https://trakt.tv/api-docs/show-season
//We always use the TVDB ID for the query
const TvSeasonURL = "/show/season.json/%s/%d/%d"

// https://trakt.tv/api-docs/search-movies
const MovieSearchURL = "/search/movies.json/%s?query=%s&limit=%d"

// https://trakt.tv/api-docs/movie-summary
// We use the IMDB ID for all summary queries
const MovieSummaryURL = "/movie/summary.json/%s/%s"

// Base URL for TraktTV api
const TraktTVBaseURL = "https://api.trakt.tv"

// TraktTV is the main struct used to query Trakt.tv.  Use NewTraktTV to
// create new instances.
type TraktTV struct {
	APIKey  string
	BaseURL string
	Session *napping.Session
}

// New initializes and returns a new TraktTV struct
func New(api string, options ...func(*TraktTV) error) (*TraktTV, error) {
	t := TraktTV{
		APIKey:  api,
		BaseURL: TraktTVBaseURL,
		Session: &napping.Session{
			Client: httpclient.NewTimeoutClient(),
		},
	}
	return &t, nil
}

func (t *TraktTV) getWithErrorCheck(url string, result interface{}) error {
	glog.Infof("Get query for %s\n", url)
	response, err := t.Session.Get(url, &napping.Params{}, result, nil)
	if serr, ok := err.(*json.SyntaxError); ok {
		line, col, highlight := HighlightBytePosition(response.HttpResponse().Body, serr.Offset)
		return fmt.Errorf("gotrackt: syntax error in response at line %d, column %d (file offset %d):\n%s", line, col, serr.Offset, highlight)
	}
	return err
}

func (t *TraktTV) genURL(u string, a ...interface{}) string {
	fullURL := fmt.Sprintf("%s%s", t.BaseURL, u)
	fullURL = fmt.Sprintf(fullURL, a...)
	glog.Infof("Generated api url: %s", fullURL)
	return fullURL
}

// GetShow returns a show and all of it's Seasons and Episodes
func (t *TraktTV) GetShow(slugOrTvdbID string) (*Show, error) {
	s := t.genURL(TvSummaryURL, t.APIKey, slugOrTvdbID)
	result := &Show{}
	err := t.getWithErrorCheck(s, result)
	return result, err
}

// ShowSearch searches tv shows
func (t *TraktTV) ShowSearch(name string) ([]Show, error) {
	s := t.genURL(TvSearchURL, t.APIKey, name, 10)
	result := []Show{}
	err := t.getWithErrorCheck(s, &result)
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

		apiURL := t.genURL(TvSeasonURL, t.APIKey, slugOrTvdbID, season)

		glog.Infof("Query for %s Season %d: %s\n", slugOrTvdbID, season, apiURL)
		err := t.getWithErrorCheck(apiURL, results[i].Episodes)
		if err != nil {
			return results, err
		}
	}
	return results, nil
}

//MovieSearch searches Trakt.tv for movies matching the query
func (t *TraktTV) MovieSearch(query string) ([]Movie, error) {
	s := t.genURL(MovieSearchURL, t.APIKey, query, 10)
	res := []Movie{}
	err := t.getWithErrorCheck(s, &res)
	return res, err
}

//GetMovieByIMDB searches Trakt.tv for movies matching the query
func (t *TraktTV) GetMovieByIMDB(imdbID string) (*Movie, error) {
	s := t.genURL(MovieSummaryURL, t.APIKey, imdbID)
	res := &Movie{}
	err := t.getWithErrorCheck(s, res)
	return res, err
}

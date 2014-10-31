package gotrakt

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/golang/glog"
	"github.com/jmcvetta/napping"
)

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
}

// NewTraktTV initializes and returns a new TraktTV struct
func NewTraktTV(api string, options ...func(*TraktTV) error) (*TraktTV, error) {
	t := TraktTV{
		APIKey:  api,
		BaseURL: TraktTVBaseURL,
	}
	return &t, nil
}

func (t *TraktTV) genURL(u string, a ...interface{}) string {
	fullURL := fmt.Sprintf("%s%s", t.BaseURL, u)
	fullURL = fmt.Sprintf(fullURL, a...)
	glog.Infof("Generated api url: %s", fullURL)
	return fullURL
}

//ShowEpisodes searches for a shows episode summaries by season.  You can optionally limit the query to just a given set of seasons.
func (t *TraktTV) ShowEpisodes(slugOrTvdbID string, seasons []int) ([]Season, error) {
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
		_, err := napping.Get(apiURL, &napping.Params{}, &results[i].Episodes, nil)
		if err != nil {
			return results, err
		}
	}
	return results, nil
}

// ShowSummary returns a show and all of it's Seasons and Episodes
func (t *TraktTV) ShowSummary(slugOrTvdbID string) (*Show, error) {
	s := t.genURL(TvSummaryURL, t.APIKey, slugOrTvdbID)
	result := &Show{}
	glog.Infof("Query %s\n", s)
	_, err := napping.Get(s, &napping.Params{}, &result, nil)
	return result, err
}

// TvSearch searches tv shows
func (t *TraktTV) TvSearch(name string) ([]Show, error) {
	s := t.genURL(TvSearchURL, t.APIKey, name, 10)
	result := []Show{}
	glog.Infof("Query %s\n", s)
	_, err := napping.Get(s, &napping.Params{}, &result, nil)
	return result, err
}

//MovieSearch searches Trakt.tv for movies matching the query
func (t *TraktTV) MovieSearch(query string) ([]Movie, error) {
	s := t.genURL(MovieSearchURL, t.APIKey, query, 10)
	res := []Movie{}
	_, err := napping.Get(s, &napping.Params{}, &res, nil)
	return res, err
}

//GetMovieByIMDB searches Trakt.tv for movies matching the query
func (t *TraktTV) GetMovieByIMDB(imdbID string) (*Movie, error) {
	s := t.genURL(MovieSummaryURL, t.APIKey, imdbID)
	res := &Movie{}
	_, err := napping.Get(s, &napping.Params{}, res, nil)
	return res, err
}

// HighlightBytePosition takes a reader and the location in bytes of a parse
// error (for instance, from json.SyntaxError.Offset) and returns the line, column,
// and pretty-printed context around the error with an arrow indicating the exact
// position of the syntax error.
func HighlightBytePosition(f io.Reader, pos int64) (line, col int, highlight string) {
	line = 1
	br := bufio.NewReader(f)
	lastLine := ""
	thisLine := new(bytes.Buffer)
	for n := int64(0); n < pos; n++ {
		b, err := br.ReadByte()
		if err != nil {
			break
		}
		if b == '\n' {
			lastLine = thisLine.String()
			thisLine.Reset()
			line++
			col = 1
		} else {
			col++
			thisLine.WriteByte(b)
		}
	}
	if line > 1 {
		highlight += fmt.Sprintf("%5d: %s\n", line-1, lastLine)
	}
	highlight += fmt.Sprintf("%5d: %s\n", line, thisLine.String())
	highlight += fmt.Sprintf("%s^\n", strings.Repeat(" ", col+5))
	return
}

package gotrakt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestTvSearch(t *testing.T) {
	searchRes, err := ioutil.ReadFile("testdata/battlestar_tv_search.json")
	if err != nil {
		t.Fatalf("Error reading test data: %s", err)
	}

	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, string(searchRes))
			}))
	defer ts.Close()

	trakt, _ := New("testing")
	trakt.BaseURL = ts.URL

	term := "Battlestar+Galactica"
	res, err := trakt.ShowSearch(term)
	if err != nil {
		t.Fatalf("Error searching: %s", err)
	}
	if len(res) != 5 {
		t.Fatalf("Expecting 5 results, got %d", len(res))
	}
}

func TestTvSummary(t *testing.T) {
	f, err := os.Open("testdata/battlestar_tv_summary_extended.json")
	if err != nil {
		t.Fatalf("Error opening test data: %s", err)
	}
	searchresult, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("Error reading test data: %s", err)
	}

	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, string(searchresult))
			}))
	defer ts.Close()

	trakt, _ := New("anyapi")
	trakt.BaseURL = ts.URL
	tvshow, err := trakt.GetShow("battlestar-galactica-2003")

	if serr, ok := err.(*json.SyntaxError); ok {
		line, col, highlight := HighlightBytePosition(f, serr.Offset)
		t.Fatalf("Error (%s) at line %d, column %d (file offset %d):\n%s", err, line, col, serr.Offset, highlight)
	}

	if tvshow.Title != "Battlestar Galactica (2003)" {
		t.Fatalf("Expecting title of \"Battlestar Galactica (2003)\" got %s", tvshow.Title)
	}
}

func TestMovieSearch(t *testing.T) {
	f, err := os.Open("testdata/batman_movie_search_fmt.json")
	if err != nil {
		t.Fatalf("Error opening test data: %s", err)
	}
	searchresult, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("Error reading test data: %s", err)
	}
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, string(searchresult))
			}))
	defer ts.Close()

	trakt, _ := New("testing")
	trakt.BaseURL = ts.URL

	term := "batman"
	res, err := trakt.MovieSearch(term)
	if serr, ok := err.(*json.SyntaxError); ok {
		line, col, highlight := HighlightBytePosition(f, serr.Offset)
		t.Fatalf("Error (%s) at line %d, column %d (file offset %d):\n%s", err, line, col, serr.Offset, highlight)
	}

	if err != nil {
		t.Fatalf("Error getting Movie search: %s", err)
	}
	if len(res) != 30 {
		t.Fatalf("Didn't parse response correctly, should have gotten 1 record, got %d", len(res))
	}
	if res[0].ImdbID != "tt0096895" {
		t.Fatalf("Unexpected IMDB value parsed, got %s", res[0].ImdbID)
	}
}

func TestMovieSummary(t *testing.T) {
	f, err := os.Open("testdata/batman_movie_summary.json")
	if err != nil {
		t.Fatalf("Error opening test data: %s", err)
	}
	summaryData, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf("Error reading test data: %s", err)
	}

	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, string(summaryData))
			}))
	defer ts.Close()

	trakt, _ := New("testing")
	trakt.BaseURL = ts.URL

	m, err := trakt.GetMovieByIMDB("tt0133093")
	if serr, ok := err.(*json.SyntaxError); ok {
		line, col, highlight := HighlightBytePosition(f, serr.Offset)
		t.Fatalf("Error (%s) at line %d, column %d (file offset %d):\n%s", err, line, col, serr.Offset, highlight)
	}

	if err != nil {
		t.Fatalf("Error getting Movie search: %s", err)
	}

	if m.Title != "Batman" {
		t.Fatalf("Unexpected title: %s", m.Title)
	}
}

package gotrakt

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/jmcvetta/napping"
)

func TestOptionSetting(t *testing.T) {
	sess := &napping.Session{}
	sess.Params = &napping.Params{
		"testing": "true",
	}
	trakt, err := New("testingapi", Session(sess))
	if err != nil {
		t.Fatalf("Unexpected error when creating new TraktTV: %s", err)
	}
	p := *trakt.Session.Params
	if val, ok := p["testing"]; ok {
		if val != "true" {
			t.Fatalf("Expected value \"true\" got \"%s\"", val)
		}
	} else {
		t.Fatalf("Didn't find \"testing\" key in the session params")
	}
}

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
	f, err := ioutil.ReadFile("testdata/battlestar_tv_summary_extended.json")
	if err != nil {
		t.Fatalf("Error reading test data: %s", err)
	}

	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, string(f))
			}))
	defer ts.Close()

	trakt, _ := New("anyapi")
	trakt.BaseURL = ts.URL
	tvshow, err := trakt.GetShow("battlestar-galactica-2003")

	if tvshow.Title != "Battlestar Galactica (2003)" {
		t.Fatalf("Expecting title of \"Battlestar Galactica (2003)\" got %s", tvshow.Title)
	}
}

func TestShowSeasons(t *testing.T) {
	seasZero, err := ioutil.ReadFile("testdata/battlestar_tv_season_0.json")
	if err != nil {
		t.Fatalf("Error reading testdata: %s", err)
	}
	seasOne, err := ioutil.ReadFile("testdata/battlestar_tv_season_1.json")
	if err != nil {
		t.Fatalf("Error reading testdata: %s", err)
	}

	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				switch {
				case strings.HasSuffix(r.URL.String(), "0"):
					fmt.Fprintf(w, string(seasZero))
				case strings.HasSuffix(r.URL.String(), "1"):
					fmt.Fprintln(w, string(seasOne))
				default:
					fmt.Fprintf(w, "Unknown request")
				}
			}))
	defer ts.Close()

	trakt, err := New("testing", Host(ts.URL))
	if err != nil {
		t.Fatalf("Error creating TraktTV: %s", err)
	}

	seas, err := trakt.ShowSeasons("battlestar-galactica-2003", []int{0, 1})
	if err != nil {
		t.Fatalf("Error getting seasons: %s", err)
	}
	if len(seas) != 2 {
		t.Fatalf("Expected 2 seasons returned, got %d", len(seas))
	}
}

func TestMovieSearch(t *testing.T) {
	f, err := ioutil.ReadFile("testdata/batman_movie_search_fmt.json")
	if err != nil {
		t.Fatalf("Error reading test data: %s", err)
	}
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, string(f))
			}))
	defer ts.Close()

	trakt, _ := New("testing")
	trakt.BaseURL = ts.URL

	term := "batman"
	res, err := trakt.MovieSearch(term)

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

	m, err := trakt.GetMovie("tt0133093")

	if err != nil {
		t.Fatalf("Error getting Movie search: %s", err)
	}

	if m.Title != "Batman" {
		t.Fatalf("Unexpected title: %s", m.Title)
	}
}

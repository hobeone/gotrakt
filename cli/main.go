package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/golang/glog"
	"github.com/hobeone/gotrakt"
)

func main() {
	defer glog.Flush()

	flag.Set("logtostderr", "true")

	apiKey := flag.String("apikey", "", "Trakt.TV API KEY")

	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Printf("\n%s --apikey KEY movie|show searchterm\n", os.Args[0])
		fmt.Printf("\n**** Flags: ****\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	args := flag.Args()
	if *apiKey == "" {
		glog.Fatal("Must supply an --apikey flag.")
	}

	if len(args) < 2 {
		glog.Error("Missing required arguments")
		flag.Usage()
	}

	t, err := gotrakt.New(*apiKey)
	if err != nil {
		glog.Fatalf("Error creating trakt client: %s\n", err)
	}

	switch args[0] {
	case "movie":
		movies, err := t.MovieSearch(args[1])
		if err != nil {
			glog.Fatalf("Error searching for movies matching \"%s\": %s", args[1], err)
		}
		for i, m := range movies {
			fmt.Printf("[%d] -  %s", i, m.Title)
		}
	case "show":
		shows, err := t.ShowSearch(args[1])
		if err != nil {
			glog.Fatalf("Error searching for shows matching \"%s\": %s", args[1], err)
		}
		for i, s := range shows {
			fmt.Printf("[%d] -  %s", i, s.Title)
		}
	}
}

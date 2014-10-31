package gotrakt

// Show is the show result from Trakt
type Show struct {
	Title         string            `json:"title"`
	Year          int               `json:"year"`
	URL           string            `json:"url"`
	FirstAired    int               `json:"first_aired"`
	Country       string            `json:"country"`
	Overview      string            `json:"overview"`
	Runtime       int               `json:"runtime"`
	Network       string            `json:"network"`
	AirDay        string            `json:"air_day"`
	AirTime       string            `json:"air_time"`
	Certification string            `json:"certification"`
	ImdbID        string            `json:"imdb_id"`
	TvdbID        int               `json:"tvdb_id"`
	TvrageID      int               `json:"tvrage_id"`
	Ended         bool              `json:"ended"`
	Images        map[string]string `json:"images"`
	Genres        []string          `json:"genres"`
	Seasons       []Season          `json:"seasons"`
}

// Season is a containter for tv episodes
type Season struct {
	Season   int       `json:"season"`
	URL      string    `json:"url"`
	Poster   string    `json:"poster"`
	Episodes []Episode `json:"episodes"`
}

// Episode contains the information for a given Show Episode
type Episode struct {
	Season        int               `json:"season"`
	Episode       int               `json:"episode"`
	Number        int               `json:"number"`
	TvdbID        int               `json:"tvdb_id"`
	Title         string            `json:"title"`
	Overview      string            `json:"overview"`
	FirstAired    int               `json:"first_aired"`
	FirstAiredIso string            `json:"first_aired_iso"`
	FirstAiredUtc int               `json:"first_aired_utc"`
	URL           string            `json:"url"`
	Screen        string            `json:"screen"`
	Images        map[string]string `json:"images"`
	Ratings       Ratings           `json:"ratings"`
	//Not filled out as we don't do auth with the api
	Watched        bool `json:"watched"`
	InCollection   bool `json:"in_collection"`
	InWatchlist    bool `json:"in_watchlist"`
	Rating         bool `json:"rating"`
	RatingAdvanced int  `json:"rating_advanced"`
}

// Ratings represents the how the thing was rated by Trakt users
type Ratings struct {
	Percentage int `json:"percentage"`
	Votes      int `json:"votes"`
	Loved      int `json:"loved"`
	Hated      int `json:"hated"`
}

// Movie holds the result of a Movie search from Trakt
type Movie struct {
	Title         string `json:"title"`
	Year          int    `json:"year"`
	Released      int    `json:"released"`
	URL           string `json:"url"`
	Trailer       string `json:"trailer"`
	Runtime       int    `json:"runtime"`
	Tagline       string `json:"tagline"`
	Overview      string `json:"overview"`
	Certification string `json:"certification"`
	ImdbID        string `json:"imdb_id"`
	TmdbID        int    `json:"tmdb_id"`
	RtID          int    `json:"rt_id"`
	LastUpdated   int    `json:"last_updated"`
	Images        struct {
		Poster string `json:"poster"`
		Fanart string `json:"fanart"`
	} `json:"images"`
	Genres []string `json:"genres"`
	// It seems like if Age is unset it returns an empty string "" rather than 0,
	// setting these to interface{} values lets the decoder work properly.
	TopWatchers []struct {
		Plays     interface{} `json:"plays"`
		Username  string      `json:"username"`
		Protected bool        `json:"protected"`
		FullName  string      `json:"full_name"`
		Gender    string      `json:"gender"`
		Age       interface{} `json:"age"`
		Location  string      `json:"location"`
		About     string      `json:"about"`
		Joined    interface{} `json:"joined"`
		Avatar    string      `json:"avatar"`
		URL       string      `json:"url"`
	} `json:"top_watchers"`
	Ratings struct {
		Percentage int `json:"percentage"`
		Votes      int `json:"votes"`
		Loved      int `json:"loved"`
		Hated      int `json:"hated"`
	} `json:"ratings"`
	Stats struct {
		Watchers   int `json:"watchers"`
		Plays      int `json:"plays"`
		Scrobbles  int `json:"scrobbles"`
		Checkins   int `json:"checkins"`
		Collection int `json:"collection"`
	} `json:"stats"`
	People struct {
		Directors []struct {
			Name   string `json:"name"`
			Images struct {
				Headshot string `json:"headshot"`
			} `json:"images"`
		} `json:"directors"`
		Writers []struct {
			Name   string `json:"name"`
			Job    string `json:"job"`
			Images struct {
				Headshot string `json:"headshot"`
			} `json:"images"`
		} `json:"writers"`
		Producers []struct {
			Name      string `json:"name"`
			Executive bool   `json:"executive"`
			Images    struct {
				Headshot string `json:"headshot"`
			} `json:"images"`
		} `json:"producers"`
		Actors []struct {
			Name      string `json:"name"`
			Character string `json:"character"`
			Images    struct {
				Headshot string `json:"headshot"`
			} `json:"images"`
		} `json:"actors"`
	} `json:"people"`
	Watched        bool   `json:"watched"`
	Plays          int    `json:"plays"`
	Rating         string `json:"rating"`
	RatingAdvanced int    `json:"rating_advanced"`
	InWatchlist    bool   `json:"in_watchlist"`
	InCollection   bool   `json:"in_collection"`
}

gotrakt
=======

Golang API for Trakt.tv

Usage
=====

Get a Show's information:
```
show, err := t.GetShow("battlestar-galactica-2003)

show, err := t.GetShow(TvdbID)
```

Get information about particular seasons of a show:
```
seasons, err := t.ShowSeasons("battlestar-galactica-2003", []int{0,1})
```

Search for a show matching a particular string:
```
shows, err := t.ShowSearch("query")
```

ToDo
====
Authentication support
More API coverage

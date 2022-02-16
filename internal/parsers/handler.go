package parsers

import (
	"net/http"
	"net/url"
)

const (
	ParseBody = 1 + iota
	ParseElement
	ParseAllOfElement
)

// Webpage represents a single cached URL, affiliated with a grantmaker via ID, but can be a many-to-one relationship.
type Webpage struct {
	ID            int
	Name          string
	URL           string
	CollySelector string `db:"colly_selector"`
	UpdatedAt     string `db:"updated_at"`
	CachedBody    string `db:"body"`
	ParsedBody    string
	IsUpdated     bool
	HttpStatus    int
	ParserID      int `db:"parser_id"`
	Error         error
}


// ScrapeSites makes concurrent HTTP calls to each webpage and pulls down all the relevant text they contain.
func ScrapeSites(cache []Webpage) (liveData []Webpage, errors []Webpage) {
	ch := make(chan Webpage, len(cache))
	defer close(ch)
	for _, c := range cache {
		go processPage(c, ch)
	}

	var collection []Webpage
	for range cache {
		collection = append(collection, <-ch)
	}

	for _, w := range collection {
		if w.Error != nil {
			errors = append(errors, w)
			continue
		}
		liveData = append(liveData, w)
	}
	return liveData, errors
}

// processPage validates a URL, fetches the page, then parses it according to its module.
func processPage(w Webpage, ch chan Webpage) {
	if _, w.Error = url.Parse(w.URL); w.Error != nil {
		ch <- w
		return
	}

	w.headReq()
	if w.Error != nil || w.HttpStatus >= 300 {
		ch <- w
		return
	}

	w.parse()
	if w.Error != nil {
		ch <- w
		return
	}
	w.IsUpdated = w.ParsedBody != w.CachedBody
}

// headReq sends a HEAD request to the target URL, then records the HTTP status code and any errors.
func (w *Webpage) headReq() {
	var res *http.Response
	res, w.Error = http.Head(w.URL)
	w.HttpStatus = res.StatusCode
}

// parse loads the relevant parser for the page, then parses it and returns the db-friendly, cacheable string.
func (w *Webpage) parse() {
	w.ParsedBody, w.Error = loadParser(w.ParserID, w.URL, w.CollySelector)
}

// loadParser calls a different function depending on what type of page body is being parsed.
func loadParser(moduleID int, url string, selector string) (retValue string, err error) {
	switch moduleID {
	case ParseBody:
		retValue, err = parseBody(url)
	case ParseElement:
		retValue, err = parseElement(url, selector)
	case ParseAllOfElement:
		retValue, err = parseAllOfElement(url, selector)
	}
	return retValue, err
}

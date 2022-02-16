package main

import (
	"fmt"
	"grants_scraper/internal/parsers"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/jmoiron/sqlx"
	// sqlite3 driver
	_ "modernc.org/sqlite"
)

var db *sqlx.DB

// webpage represents a single cached URL, affiliated with a grantmaker via ID, but can be a many-to-one relationship.
type webpage struct {
	ID         int
	Name       string
	URL        string
	UpdatedAt  string `db:"updated_at"`
	CachedBody string `db:"body"`
	ParsedBody string
	IsUpdated  bool
	HttpStatus int
	ParserID   int `db:"parser_id"`
	Error      error
}

func main() {
	dbPath := os.Getenv("DB_PATH")
	db = setup(dbPath)

	urls, err := fetchCachedPages()
	if err != nil {
		log.Fatal("problem fetching cache:", err)
	}

	liveSiteData, erroredSites := scrapeSites(urls)
	for _, e := range erroredSites {
		log.Println("problem scraping url:", e.ID, e.URL, e.Error)
	}

	erroredSites = updateDatabase(liveSiteData)
	for _, e := range erroredSites {
		log.Println("problem updating database for url:", e.ID, e.URL, e.Error)
	}
}

// setup takes a path to a SQLite3 database and connects to it.
func setup(path string) *sqlx.DB {
	str := fmt.Sprintf("file:%s?cache=shared&mode=rwc", path)
	conn, err := sqlx.Connect("sqlite", str)
	if err != nil {
		log.Fatal("db connection failure:", err)
	}
	return conn
}

// fetchCachedPages pulls a list of URLs that have been cached in the database.
func fetchCachedPages() (cache []webpage, err error) {
	query := `
		select g.id, g.name, g.updated_at, c.url, c.parser_id, c.body 
		from (select id, name, updated_at from grantmakers where active = 1) g
		inner join (select grantmaker_id, url, parser_id, body from cached_pages where active = 1) c on c.grantmaker_id = g.id`
	err = db.Select(&cache, query)
	if err != nil {
		return
	}
	return
}

// scrapeSites makes concurrent HTTP calls to each webpage and pulls down all the relevant text they contain.
func scrapeSites(cache []webpage) (liveData []webpage, errors []webpage) {
	ch := make(chan webpage, len(cache))
	defer close(ch)
	for _, c := range cache {
		go processPage(c, ch)
	}

	var collection []webpage
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
func processPage(w webpage, ch chan webpage) {
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
func (w *webpage) headReq() {
	var res *http.Response
	res, w.Error = http.Head(w.URL)
	w.HttpStatus = res.StatusCode
}

// parse loads the relevant parser for the page, then parses it and returns the db-friendly, cacheable string.
func (w *webpage) parse() {
	w.ParsedBody, w.Error = parsers.LoadParser(w.ParserID, w.URL)
}

// updateDatabase records the scraper run in scrapes table, then updates the cached body if needed.
func updateDatabase(pages []webpage) (errors []webpage) {
	for _, p := range pages {
		query := `insert into scrapes (grantmaker_id) values ($1)`
		if _, err := db.Exec(query, p.ID); err != nil {
			errors = append(errors, p)
			continue
		}

		if p.IsUpdated {
			query = `
				update cached_pages 
				set updated_at = current_timestamp, body = $1 
				where grantmaker_id = $2 and url = $3`
			if _, err := db.Exec(query, p.ParsedBody, p.ID, p.URL); err != nil {
				errors = append(errors, p)
				continue
			}

		}
	}
	return
}

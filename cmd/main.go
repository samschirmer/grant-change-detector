package main

import (
	"fmt"
	parser "grants_scraper/internal/parsers"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	// sqlite3 driver
	_ "modernc.org/sqlite"
)

var db *sqlx.DB

func main() {
	dbPath := os.Getenv("DB_PATH")
	db = setup(dbPath)

	urls, err := fetchCachedPages()
	if err != nil {
		log.Fatal("problem fetching cache:", err)
	}

	liveSiteData, erroredSites := parser.ScrapeSites(urls)
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
func fetchCachedPages() (cache []parser.Webpage, err error) {
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

// updateDatabase records the scraper run in scrapes table, then updates the cached body if needed.
func updateDatabase(pages []parser.Webpage) (errors []parser.Webpage) {
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

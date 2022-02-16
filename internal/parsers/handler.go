package parsers

import (
	"fmt"
	"github.com/gocolly/colly"
)

const (
	BasicHtmlAll = 1 + iota
	BasicHtmlTargeted
)

// LoadParser calls a different function depending on what type of page body is being parsed.
func LoadParser(moduleID int, url string) (retValue string, err error) {
	switch moduleID {
	case BasicHtmlAll:
		retValue, err = parseAllHtml(url)
	case BasicHtmlTargeted:
	}
	return retValue, err
}

// parseAllHtml simply caches the text of all HTML elements on the page.
func parseAllHtml(page string) (toCache string, err error) {
	c := colly.NewCollector()

	c.OnHTML("body", func(e *colly.HTMLElement) {
		body := e
		toCache = body.Text
		fmt.Println(body.Text)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting", page)
	})

	err = c.Visit(page)
	return
}

package parsers

import (
	"fmt"
	"github.com/gocolly/colly"
)

// parseBody simply caches the text of all HTML elements in the body of the page.
func parseBody(page string) (toCache string, err error) {
	c := colly.NewCollector()
	c.OnHTML("body", func(e *colly.HTMLElement) {
		body := e
		toCache = body.Text
	})
	err = c.Visit(page)
	return
}

// parseElement takes a Colly-specific selector as an arg and returns the text of only that element.
func parseElement(page string, selector string) (toCache string, err error) {
	c := colly.NewCollector()
	c.OnHTML(selector, func(e *colly.HTMLElement) {
		toCache = e.Text
	})
	err = c.Visit(page)
	return
}

// parseAllOfElement takes a Colly-specific selector and returns the semicolon-delimited text of every instance of that element.
func parseAllOfElement(page string, selector string) (toCache string, err error) {
	var elements []colly.HTMLElement
	c := colly.NewCollector()
	c.OnHTML(selector, func(e *colly.HTMLElement) {
		elements = append(elements, *e)
	})
	err = c.Visit(page)
	for _, e := range elements {
		toCache = fmt.Sprintf("%v ; %v", toCache, e.Text)
	}
	return toCache, err
}

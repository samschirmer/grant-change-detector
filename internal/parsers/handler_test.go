package parsers

import (
	"log"
	"testing"
)

func Test_webpage_headReq(t *testing.T) {
	type fields struct {
		ID         int
		Name       string
		URL        string
		HttpStatus int
		Error      error
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "samschirmer.com",
			fields: fields{
				ID:   1,
				Name: "Sam Schirmer",
				URL:  "https://samschirmer.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Webpage{
				ID:         tt.fields.ID,
				Name:       tt.fields.Name,
				URL:        tt.fields.URL,
				HttpStatus: tt.fields.HttpStatus,
				Error:      tt.fields.Error,
			}
			w.headReq()
			if w.Error != nil {
				t.Fail()
			}
			log.Println(w.HttpStatus)
		})
	}
}

func Test_webpage_parse(t *testing.T) {
	type fields struct {
		ID            int
		Name          string
		URL           string
		CollySelector string
		HttpStatus    int
		ParserID      int
		Error         error
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "samschirmer.com",
			fields: fields{
				ID:       0,
				Name:     "Sam Schirmer",
				URL:      "https://samschirmer.com",
				ParserID: ParseBody,
			},
		},
		{
			name: "Doris Day Animal Foundation",
			fields: fields{
				ID:            1,
				Name:          "Doris Day Animal Foundation",
				URL:           "https://www.dorisdayanimalfoundation.org/grants",
				ParserID:      ParseElement,
				CollySelector: "div#main_content",
			},
		},
		{
			name: "Maddie's Fund",
			fields: fields{
				ID:            2,
				Name:          "Maddie's Fund",
				URL:           "https://www.maddiesfund.org/grant-opportunities.htm",
				ParserID:      ParseAllOfElement,
				CollySelector: "li",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Webpage{
				ID:            tt.fields.ID,
				Name:          tt.fields.Name,
				URL:           tt.fields.URL,
				CollySelector: tt.fields.CollySelector,
				HttpStatus:    tt.fields.HttpStatus,
				ParserID:      tt.fields.ParserID,
				Error:         tt.fields.Error,
			}
			w.parse()
			if w.Error != nil {
				log.Println(w.Error)
				t.Fail()
			}
			log.Println(w.ParsedBody)
		})
	}
}

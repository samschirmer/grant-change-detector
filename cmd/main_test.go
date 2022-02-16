package main

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
				ID:         1,
				Name:       "Sam Schirmer",
				URL:        "https://samschirmer.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webpage{
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
		ID         int
		Name       string
		URL        string
		ByteBody   []byte
		HttpStatus int
		ParserID   int
		Error      error
	}
	tests := []struct{
		name   string
		fields fields
	}{
		{
			name: "samschirmer.com",
			fields: fields{
				ID:         1,
				Name:       "Sam Schirmer",
				URL:        "https://samschirmer.com",
				ParserID:   1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webpage{
				ID:         tt.fields.ID,
				Name:       tt.fields.Name,
				URL:        tt.fields.URL,
				ByteBody:   tt.fields.ByteBody,
				HttpStatus: tt.fields.HttpStatus,
				ParserID:   tt.fields.ParserID,
				Error:      tt.fields.Error,
			}
			w.parse()
			if w.Error != nil {
				t.Fail()
			}
			log.Println(string(w.ByteBody))
		})
	}
}

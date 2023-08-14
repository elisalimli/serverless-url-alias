package domain

import "github.com/elisalimli/serverless-url-alias/internal/sheets"

type Domain struct {
	client sheets.Client
}

func NewDomain(client sheets.Client) *Domain {
	return &Domain{client: client}
}

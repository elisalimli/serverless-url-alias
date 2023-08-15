package api

import (
	"github.com/elisalimli/serverless-url-alias/domain"
	"github.com/elisalimli/serverless-url-alias/internal/sheets"
)

type Handlers struct {
	Domain    domain.Domain
	SheetAuth sheets.SheetAuth
}

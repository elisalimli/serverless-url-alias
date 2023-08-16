package domain

import (
	"net/url"
	"sync"
	"time"

	"github.com/elisalimli/serverless-url-alias/internal/initializers"
	"github.com/elisalimli/serverless-url-alias/internal/sheets"
)

type UrlMap map[string]*url.URL
type CachedUrlMap struct {
	sync.RWMutex
	LastUpdated   time.Time
	Ttl           time.Duration // time-to-live
	UrlMap        UrlMap
	SheetProvider SheetProvider
}

type Domain struct {
	DB *CachedUrlMap
}

func NewDomain(client sheets.Client) *Domain {
	return &Domain{DB: &CachedUrlMap{Ttl: time.Second * 10, UrlMap: make(UrlMap), SheetProvider: SheetProvider{GoogleSheetId: initializers.GoogleSheetId, GoogleSheetName: initializers.GoogleSheetName, client: client}}}
}

package domain

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	mySheets "github.com/elisalimli/serverless-url-alias/internal/sheets"
)

type SheetProvider struct {
	GoogleSheetId   string
	GoogleSheetName string
	client          mySheets.Client
}

func urlMap(in [][]interface{}) UrlMap {
	res := UrlMap{}
	// starting from first index to skip the headers
	for _, row := range in[1:] {
		// row should be in the follwing format : [alias, url]
		if len(row) < 2 {
			continue
		}

		alias, ok := row[0].(string)
		if !ok {
			log.Printf("warn: %v alias is invalid", alias)
			continue
		}
		value, ok := row[1].(string)
		if !ok || value == "" {
			log.Printf("warn: %v url value is invalid", value)
			continue
		}

		alias = strings.ToLower(alias)

		url, err := url.Parse(value)
		if err != nil {
			log.Printf("warn: %s=%s url is invalid", alias, value)
			continue
		}

		_, exists := res[alias]
		if exists {
			log.Printf("warn: duplicate alias %s, overwritting", alias)
		}
		res[alias] = url
	}
	return res
}

func (c *CachedUrlMap) Get() (UrlMap, error) {
	if time.Since(c.LastUpdated) > c.Ttl {
		err := c.Refresh()
		if err != nil {
			return nil, err
		}
	}
	return c.UrlMap, nil
}

func (c *CachedUrlMap) Refresh() error {
	fmt.Println("before lock")
	c.Lock()
	defer c.Unlock()

	if time.Since(c.LastUpdated) <= c.Ttl {
		return nil
	}
	fmt.Println("after lock - attempting to refresh")

	readRange := c.SheetProvider.GoogleSheetName + "!" + "A:B"

	data, err := c.SheetProvider.Query(readRange)
	if err != nil {
		return err
	}

	newUrlMap := urlMap(data)
	c.UrlMap = newUrlMap
	c.LastUpdated = time.Now()

	return nil

}

func (s SheetProvider) Query(readRange string) ([][]interface{}, error) {
	fmt.Println("querying sheet")

	resp, err := s.client.Service.Spreadsheets.Values.Get(s.GoogleSheetId, readRange).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		return nil, nil
	}

	return resp.Values, nil
}

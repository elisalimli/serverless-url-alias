package api

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func (h *Handlers) Redirect(w http.ResponseWriter, r *http.Request) {

	// Prints the names and majors of students in a sample spreadsheet:
	// https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit
	readRange := h.SheetAuth.GoogleSheetName + "!A:B"
	data, err := h.Domain.GetSpreadsheetData(h.SheetAuth.GoogleSheetId, readRange)

	if err != nil {
		log.Fatal(err)
	}

	urlMap := urlMap(data)
	segments := strings.Split(r.URL.Path, "/")
	fmt.Println("segments", r.URL.Path, segments)
	baseURL, discardedPaths := getRedirect(urlMap, segments)

	fmt.Println("test", baseURL, discardedPaths)
	redirectTo := prepRedirect(baseURL, discardedPaths, r.URL)
	if redirectTo == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "alias %v not found", r.URL.Path)
		return
	}
	http.Redirect(w, r, redirectTo.String(), http.StatusMovedPermanently)
}

func getRedirect(m map[string]*url.URL, segments []string) (*url.URL, []string) {
	discarded := []string{}
	for len(segments) > 0 {
		baseURL := checkRedirect(m, strings.Join(segments, "/"))
		log.Println("baseUrl", baseURL, strings.Join(segments, "/"))
		if baseURL != nil {
			return baseURL, discarded
		}
		discarded = append([]string{segments[len(segments)-1]}, discarded...)
		segments = segments[:len(segments)-1]
	}
	return nil, nil
}

func prepRedirect(baseURL *url.URL, discardedPaths []string, path *url.URL) *url.URL {
	values := path.Query()
	for key := range baseURL.Query() {
		values.Add(key, baseURL.Query().Get(key))
	}

	if !strings.HasSuffix(baseURL.Path, "/") {
		baseURL.Path += "/"
	}

	baseURL.Path = baseURL.Path + strings.Join(discardedPaths, "/")

	baseURL.RawQuery = values.Encode()

	return baseURL
}

func checkRedirect(m map[string]*url.URL, path string) *url.URL {
	path = strings.TrimPrefix(path, "/")
	return m[path]
}

/*
alias http://localhost:4000/ex/sub/?c=d
url https://example.com/?a=b
output https://example.com/sub/?a=b&c=d

url https://example.com/?a=b ~~> https://example.com/ ++
url query https://example.com/?a=b ~~> ?a=b
alias /sub/ ~~> /sub/ ++
alias query /sub/?c=d ~~> ?c=d



*/

func urlMap(in [][]interface{}) map[string]*url.URL {
	res := map[string]*url.URL{}
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

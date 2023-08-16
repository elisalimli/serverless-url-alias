package api

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/elisalimli/serverless-url-alias/domain"
)

func (h *Handlers) Redirect(w http.ResponseWriter, req *http.Request) {

	// Prints the names and majors of students in a sample spreadsheet:
	// https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit
	// data, err := h.Domain.GetSpreadsheetData(readRange)

	// if err != nil {
	// log.Fatal(err)
	// }
	if req.Body != nil {
		defer req.Body.Close()
	}

	urlMap, err := h.Domain.DB.Get()
	if err != nil {
		writeError(w, "couldn't get spread sheet data", http.StatusInternalServerError)
	}
	segments := strings.Split(req.URL.Path, "/")
	baseURL, discardedPaths := getRedirect(urlMap, segments)
	redirectTo := prepRedirect(baseURL, discardedPaths, req.URL)

	if redirectTo == nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "alias %v not found", req.URL.Path)
		return
	}

	http.Redirect(w, req, redirectTo.String(), http.StatusMovedPermanently)
}

func getRedirect(m domain.UrlMap, segments []string) (*url.URL, []string) {
	discarded := []string{}
	for len(segments) > 0 {
		baseURL := checkRedirect(m, strings.Join(segments, "/"))
		log.Println("redirecting = ", strings.Join(segments, "/"), "to=", baseURL)
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

func checkRedirect(m domain.UrlMap, path string) *url.URL {
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

package api

import (
	"fmt"
	"net/http"
)

func writeError(w http.ResponseWriter, msg string, code int, vals ...any) {
	w.WriteHeader(code)
	fmt.Fprintf(w, msg, vals)
}

package http

import (
	"fmt"
	"net/http"
)

// HealthHandler returns "OK" if the method is GET.
// It returns 405 Method Not Allowed for other methods.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprint(w, "OK")
}

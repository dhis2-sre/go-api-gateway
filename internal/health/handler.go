package health

import (
	"fmt"
	"net/http"
)

func Handler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = fmt.Fprintln(w, `{"status": "UP"}`)
}

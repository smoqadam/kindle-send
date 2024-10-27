package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/onurhanak/libgenapi"
)

func handleLibgenSearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		q := r.URL.Query().Get("q")
		fmt.Println(q)
		query := libgenapi.NewQuery("default", q, 25)

		query.Search()
		books := query.Results

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(books)
	}
}

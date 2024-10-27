package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/smoqadam/kindle-send/classifier"
	"github.com/smoqadam/kindle-send/handler"
	"github.com/smoqadam/kindle-send/types"
	"github.com/smoqadam/kindle-send/util"
)

func handleSend() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		urls, err := util.ParseURLs(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if len(urls) == 0 {
			http.Error(w, "No valid URLs provided", http.StatusBadRequest)
			return
		}

		responses := make([]types.ResponseItem, 0, len(urls))
		downloadRequests := classifier.Classify(urls)

		if len(downloadRequests) == 0 {
			http.Error(w, "No valid content found in provided URLs", http.StatusBadRequest)
			return
		}

		downloadedRequests := handler.Queue(downloadRequests)
		if len(downloadedRequests) == 0 {
			http.Error(w, "Failed to process any of the provided URLs", http.StatusInternalServerError)
			return
		}

		err = handler.Mail(downloadedRequests, 30)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error sending mail: %v", err), http.StatusInternalServerError)
			return
		}

		util.CyanBold.Printf("Downloaded %d files:\n", len(downloadedRequests))

		for idx, req := range downloadedRequests {
			fileInfo, err := os.Stat(req.Path)
			if err != nil {
				responses = append(responses, types.ResponseItem{
					URL:     req.Path,
					Success: false,
					Error:   fmt.Sprintf("Failed to get file info: %v", err),
				})
				continue
			}

			responses = append(responses, types.ResponseItem{
				URL:     fileInfo.Name(),
				Success: true,
			})
			util.Cyan.Printf("%d. %s\n", idx+1, fileInfo.Name())
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(responses); err != nil {
			util.Red.Printf("Error encoding response: %v\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

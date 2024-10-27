package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/nikhil1raghav/kindle-send/classifier"
	"github.com/nikhil1raghav/kindle-send/handler"
	"github.com/nikhil1raghav/kindle-send/types"
	"github.com/nikhil1raghav/kindle-send/util"
)

func handleSend() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		urls, err := util.ParseURLs(r)
		fmt.Println(urls)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		responses := make([]types.ResponseItem, len(urls))
		downloadRequests := classifier.Classify(urls)
		downloadedRequests := handler.Queue(downloadRequests)
		handler.Mail(downloadedRequests, 30)

		util.CyanBold.Printf("Downloaded %d files :\n", len(downloadRequests))
		for idx, req := range downloadedRequests {
			fileInfo, _ := os.Stat(req.Path)
			responses[idx] = types.ResponseItem{
				URL:     fileInfo.Name(),
				Success: true,
			}
			util.Cyan.Printf("%d. %s\n", idx+1, fileInfo.Name())
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responses)
	}
}

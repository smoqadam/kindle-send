package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/nikhil1raghav/kindle-send/classifier"
	"github.com/nikhil1raghav/kindle-send/config"
	"github.com/nikhil1raghav/kindle-send/handler"
	"github.com/nikhil1raghav/kindle-send/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().Int("port", 8080, "Port to run the server on")
}

type ResponseItem struct {
	URL     string `json:"url"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the web server with send and download endpoints",
	Long: `Starts a web server that provides two endpoints:
- GET /send?urls=url1,url2 : Processes URLs and returns results
- GET /download?urls=url1,url2 : Downloads content from provided URLs`,
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt("port")
		configPath, _ := cmd.Flags().GetString("config")

		_, err := config.Load(configPath)
		if err != nil {
			util.Red.Println("Error loading config:", err)
			return
		}

		http.HandleFunc("/send", handleSend())
		http.HandleFunc("/download", handleDownload())

		address := fmt.Sprintf(":%d", port)
		util.Green.Printf("Starting server on http://localhost%s\n", address)
		if err := http.ListenAndServe(address, nil); err != nil {
			util.Red.Printf("Server error: %v\n", err)
		}
	},
}

func parseURLs(r *http.Request) ([]string, error) {
	urls := r.URL.Query()["url[]"]
	if len(urls) == 0 {
		return nil, fmt.Errorf("no URLs provided")
	}

	// Remove empty strings and trim spaces
	validURLs := make([]string, 0)
	for _, url := range urls {
		if trimmedURL := strings.TrimSpace(url); trimmedURL != "" {
			validURLs = append(validURLs, trimmedURL)
		}
	}

	if len(validURLs) == 0 {
		return nil, fmt.Errorf("no valid URLs provided")
	}

	return validURLs, nil
}

func handleSend() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		urls, err := parseURLs(r)
		fmt.Println(urls)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		responses := make([]ResponseItem, len(urls))
		downloadRequests := classifier.Classify(urls)
		downloadedRequests := handler.Queue(downloadRequests)
		handler.Mail(downloadedRequests, 30)

		util.CyanBold.Printf("Downloaded %d files :\n", len(downloadRequests))
		for idx, req := range downloadedRequests {
			fileInfo, _ := os.Stat(req.Path)
			responses[idx] = ResponseItem{
				URL:     fileInfo.Name(),
				Success: true,
			}
			util.Cyan.Printf("%d. %s\n", idx+1, fileInfo.Name())
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responses)
	}
}

func handleDownload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		urls, err := parseURLs(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		responses := make([]ResponseItem, len(urls))
		downloadRequests := classifier.Classify(urls)
		downloadedRequests := handler.Queue(downloadRequests)

		util.CyanBold.Printf("Downloaded %d files :\n", len(downloadRequests))
		for idx, req := range downloadedRequests {
			fileInfo, _ := os.Stat(req.Path)
			responses[idx] = ResponseItem{
				URL:     fileInfo.Name(),
				Success: true,
			}
			util.Cyan.Printf("%d. %s\n", idx+1, fileInfo.Name())
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responses)
	}
}

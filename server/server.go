package server

import (
	"fmt"
	"net/http"

	"github.com/nikhil1raghav/kindle-send/util"
)

type Server struct {
}

func Start(port int) {
	http.HandleFunc("/send", handleSend())
	http.HandleFunc("/download", handleDownload())

	address := fmt.Sprintf(":%d", port)
	util.Green.Printf("Starting server on http://localhost%s\n", address)
	if err := http.ListenAndServe(address, nil); err != nil {
		util.Red.Printf("Server error: %v\n", err)
	}
}

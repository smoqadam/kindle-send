package handler

import (
	"fmt"

	"github.com/smoqadam/kindle-send/config"
	"github.com/smoqadam/kindle-send/epubgen"
	"github.com/smoqadam/kindle-send/mail"
	"github.com/smoqadam/kindle-send/types"
	"github.com/smoqadam/kindle-send/util"
)

func Queue(downloadRequests []types.Request) []types.Request {
	var processedRequests []types.Request
	for _, req := range downloadRequests {
		switch req.Type {
		case types.TypeFile:
			processedRequests = append(processedRequests, req)

		case types.TypeRemoteFile:
			path, err := util.DownloadFile(req.Path, config.GetInstance().StorePath)
			if err != nil {
				util.Red.Printf("SKIPPING %s: %s\n", req.Path, err)
			} else {
				processedRequests = append(processedRequests, types.NewRequest(path, types.TypeFile, nil))
				util.Green.Printf("Successfully downloaded: %s\n", path)
			}

		case types.TypeUrl:
			path, err := epubgen.Make([]string{req.Path}, "")
			if err != nil {
				util.Red.Printf("SKIPPING %s: %s\n", req.Path, err)
			} else {
				processedRequests = append(processedRequests, types.NewRequest(path, types.TypeFile, nil))
			}
		}
	}
	return processedRequests
}

func Mail(mailRequests []types.Request, timeout int) error {
	if len(mailRequests) == 0 {
		return fmt.Errorf("no files to send")
	}

	filePaths := make([]string, 0, len(mailRequests))
	for _, req := range mailRequests {
		if req.Path == "" {
			continue
		}
		filePaths = append(filePaths, req.Path)
	}

	if len(filePaths) == 0 {
		return fmt.Errorf("no valid file paths to send")
	}

	if timeout < 60 {
		timeout = config.DefaultTimeout
	}

	return mail.Send(filePaths, timeout)
}

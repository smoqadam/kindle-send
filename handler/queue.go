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
	if len(downloadRequests) == 0 {
		return nil
	}

	processedRequests := make([]types.Request, 0, len(downloadRequests))
	for _, req := range downloadRequests {
		if req.Path == "" {
			util.Red.Println("Skipping empty path")
			continue
		}

		switch req.Type {
		case types.TypeFile:
			processedRequests = append(processedRequests, req)
		case types.TypeRemoteFile:
			conf := config.GetInstance()
			path, err := util.DownloadFile(req.Path, conf.StorePath)
			if err != nil {
				util.Red.Printf("SKIPPING %s: %v\n", req.Path, err)
				continue
			}
			processedRequests = append(processedRequests, types.NewRequest(path, types.TypeFile, nil))
		case types.TypeUrl:
			path, err := epubgen.Make([]string{req.Path}, "")
			if err != nil {
				util.Red.Printf("SKIPPING %s: %v\n", req.Path, err)
				continue
			}
			processedRequests = append(processedRequests, types.NewRequest(path, types.TypeFile, nil))
		default:
			util.Red.Printf("Unknown type for %s: %v\n", req.Path, req.Type)
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

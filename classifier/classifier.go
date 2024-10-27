package classifier

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/smoqadam/kindle-send/types"
)

func isLibgen(u string) bool {
	if !isUrl(u) {
		return false
	}
	parsedURL, err := url.Parse(u)
	if err != nil {
		return false
	}

	libgenDomains := []string{
		"https://download.library.lol",
		// "libgen.is",
		// "libgen.rs",
		// "library.lol",
		// "gen.lib.rus.ec",
	}

	for _, domain := range libgenDomains {
		if strings.Contains(parsedURL.Host, domain) {
			return true
		}
	}
	return false
}

func isUrl(u string) bool {
	for _, proto := range []string{"http://", "https://"} {
		if strings.HasPrefix(u, proto) {
			return true
		}
	}
	return false
}

func isRemoteFile(u string) bool {
	if !isUrl(u) {
		return false
	}
	extension := strings.ToLower(filepath.Ext(u))
	for _, ext := range []string{".mobi", ".pdf", ".epub", ".azw3", ".txt"} {
		if extension == ext {
			return true
		}
	}
	return false
}

func isUrlFile(u string) bool {
	file, err := os.Open(u)
	if err != nil {
		return false
	}
	defer file.Close()

	buf := make([]byte, 1024)
	n, _ := file.Read(buf)
	content := string(buf[:n])
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, "http") {
			return false
		}
	}
	return true
}

func isLocalBook(u string) bool {
	extension := filepath.Ext(u)
	_, err := os.Stat(u)
	if err != nil {
		return false
	}
	for _, ext := range []string{".mobi", ".pdf", ".epub", ".azw3", ".txt"} {
		if extension == ext {
			return true
		}
	}
	return false
}

func processUrlFile(path string) []string {
	var urls []string
	file, err := os.Open(path)
	if err != nil {
		return urls
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return urls
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			urls = append(urls, line)
		}
	}
	return urls
}

func Classify(args []string) []types.Request {
	if len(args) == 0 {
		return nil
	}

	requests := make([]types.Request, 0, len(args))
	for _, arg := range args {
		if arg == "" {
			continue
		}

		// if isLibgen(arg) {
		// 	requests = append(requests, types.NewRequest(arg, types.TypeLibgen, nil))
		// } else
		if isRemoteFile(arg) {
			requests = append(requests, types.NewRequest(arg, types.TypeRemoteFile, nil))
		} else if isUrl(arg) {
			requests = append(requests, types.NewRequest(arg, types.TypeUrl, nil))
		} else if isUrlFile(arg) {
			urls := processUrlFile(arg)
			for _, url := range urls {
				if url == "" {
					continue
				}
				if isLibgen(url) {
					requests = append(requests, types.NewRequest(url, types.TypeLibgen, nil))
				} else if isRemoteFile(url) {
					requests = append(requests, types.NewRequest(url, types.TypeRemoteFile, nil))
				} else if isUrl(url) {
					requests = append(requests, types.NewRequest(url, types.TypeUrl, nil))
				}
			}
		} else if isLocalBook(arg) {
			requests = append(requests, types.NewRequest(arg, types.TypeFile, nil))
		}
	}
	fmt.Println(requests)
	return requests
}

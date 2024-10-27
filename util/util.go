package util

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
)

var Red = color.New(color.FgRed)
var RedBold = color.New(color.FgRed).Add(color.Bold)
var Cyan = color.New(color.FgCyan)
var CyanBold = color.New(color.FgCyan).Add(color.Bold)
var Green = color.New(color.FgGreen)
var GreenBold = color.New(color.FgGreen).Add(color.Bold)
var Magenta = color.New(color.FgMagenta)

func Scanline() string {
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return scanner.Text()
	}
	color.Red("\nInterrupted")
	os.Exit(1)
	return ""
}

// Scan input and trim
func ScanlineTrim() string {
	return strings.TrimSpace(Scanline())
}

// ExtractLinks extracts links from a file containing urls
func ExtractLinks(filename string) (links []string) {
	links = make([]string, 0)
	file, err := os.Open(filename)
	if err != nil {
		Red.Println("Error opening link file", err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			continue
		}
		links = append(links, scanner.Text())
	}
	return
}

func ParseURLs(r *http.Request) ([]string, error) {
	urls := r.URL.Query()["url[]"]
	if len(urls) == 0 {
		return nil, fmt.Errorf("no URLs provided")
	}

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

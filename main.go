package main

import (
	"bytes"
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func getFilenameFromURL(url *url.URL) string {
	return url.Hostname() + ".html"
}
func getFolderNameFromURL(url *url.URL) string {
	return url.Hostname() + "_content"
}

func isValidURL(url *url.URL) bool {
	return url.Scheme == "http" || url.Scheme == "https"
}

func main() {
	metadata := flag.Bool("metadata", false, "Enable metadata mode")
	// Parse command-line arguments
	flag.Parse()

	tail := flag.Args()
	urls := make([]*url.URL, 0)
	for _, str := range tail {
		parsedURL, err := url.Parse(str)
		if err != nil {
			fmt.Printf("Failed to parse URL: %s\n", err)
			return
		}
		if !isValidURL(parsedURL) {
			fmt.Printf("Invalid URL: %s\n", parsedURL)
			return
		}
		urls = append(urls, parsedURL)
	}
	downloadDomains(*metadata, urls)
}

func downloadDomains(metadata bool, urls []*url.URL) {
	for _, url := range urls {
		downloadDomain(metadata, url)
	}
}

func downloadDomain(metadata bool, url *url.URL) {

	filename := getFilenameFromURL(url)
	folderName := getFolderNameFromURL(url)

	err := os.Mkdir(folderName, 0755)
	if err != nil {
		if !os.IsExist(err) {
			fmt.Printf("Failed to create folder: %v\n", err)
			return
		}
	}

	var lastFetch *string = nil
	if metadata {
		fileInfo, err := os.Stat(filename)
		if err == nil {
			temp := fileInfo.ModTime().String()
			lastFetch = &temp
		}
	}

	// Send GET request to the URL
	response, err := http.Get(url.String())
	if err != nil {
		fmt.Printf("Failed to send GET request: %v\n", err)
		return
	}
	defer response.Body.Close()

	// Check if the response was successful (status code 200)
	if response.StatusCode != http.StatusOK {
		fmt.Printf("Request failed with status code %d\n", response.StatusCode)
		return
	}

	// Read the response body
	htmlContent, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		return
	}

	// Create a new file to save the downloaded content
	file, err := os.Create(filename) // Replace with the desired file name
	if err != nil {
		fmt.Printf("Failed to create file: %v\n", err)
		return
	}
	defer file.Close()

	//site: www.google.com
	//num_links: 35
	//images: 3
	//last_fetch: Tue Mar 16 2021 15:46 UTC
	data, newHTML := parseHTML(url, string(htmlContent))
	if metadata {
		fmt.Println("site:", url.Hostname())
		fmt.Println("num_links:", data.links)
		fmt.Println("images:", data.images)
		if lastFetch != nil {
			fmt.Println("last_fetch:", *lastFetch)
		} else {
			fmt.Println("last_fetch: N/A")
		}
	}

	_, err = file.WriteString(newHTML)

	if err != nil {
		fmt.Printf("Failed to download content: %v\n", err)
		return
	}
}

type MetaData struct {
	links  int
	images int
}

func parseHTML(url *url.URL, htmlContent string) (*MetaData, string) {
	metaData := &MetaData{}

	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		panic("Failed to parse HTML:" + err.Error())
	}

	traverseHTML(url, doc, metaData)

	htmlStringNew := renderNodeToString(doc)
	return metaData, htmlStringNew
}

func traverseHTML(url *url.URL, node *html.Node, data *MetaData) {
	if node.Type == html.ElementNode && node.Data == "a" {
		data.links++
	}
	if node.Type == html.ElementNode && node.Data == "img" {
		data.images++
	}
	if node.Type == html.ElementNode {
		newAttributes := make([]html.Attribute, 0)
		for _, attribute := range node.Attr {
			if attribute.Key == "src" || (node.Data == "link" && attribute.Key == "href") {
				newFileName := downloadAndSaveContent(url, attribute.Val)
				attribute.Val = newFileName
			}
			newAttributes = append(newAttributes, attribute)
		}
		node.Attr = newAttributes
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		traverseHTML(url, child, data)
	}
}

func downloadAndSaveContent(url *url.URL, downloadContentURL string) string {
	if downloadContentURL == "" {
		return ""
	}

	if strings.HasPrefix(downloadContentURL, "/") && !strings.HasPrefix(downloadContentURL, "//") {
		return downloadAndSaveContent(url, url.Scheme+"://"+url.Hostname()+downloadContentURL)
	}

	urlParsed, err := url.Parse(downloadContentURL)
	if err != nil {
		fmt.Printf("Failed to parse URL: %s\n", err)
		return ""
	}

	folderName := getFolderNameFromURL(url)
	dirPath := filepath.Dir(urlParsed.Path)
	if urlParsed.Path == "" {
		dirPath = ""
	}

	if dirPath == "/" {
		return ""
	}

	newFileName := folderName + urlParsed.Path
	if urlParsed.Path == "" {
		newFileName = folderName + downloadContentURL[1:]
	}

	newFolderName := folderName + dirPath

	err = os.MkdirAll(newFolderName, 0755)
	if err != nil {
		fmt.Printf("Failed to create folder: %v\n", err)
		return ""
	}
	finalDownloadUrl := downloadContentURL
	if strings.HasPrefix(downloadContentURL, "//") {
		finalDownloadUrl = url.Scheme + ":" + downloadContentURL
	}

	response, err := http.Get(finalDownloadUrl)
	if err != nil {
		fmt.Printf("Failed to send GET request: %v\n", err)
		return ""
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Printf("Request failed with status code %d URL: %s\n", response.StatusCode, finalDownloadUrl)
		return ""
	}

	file, err := os.Create(newFileName)
	if err != nil {
		fmt.Printf("Failed to create file: %v\n", err)
		return ""
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)

	if err != nil {
		fmt.Printf("Failed to download content: %v\n", err)
		return ""
	}
	return newFileName
}

func renderNodeToString(node *html.Node) string {
	var buf bytes.Buffer
	err := html.Render(&buf, node)
	if err != nil {
		panic("Failed to render HTML node: " + err.Error())
	}
	return buf.String()
}

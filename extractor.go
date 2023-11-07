package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

func main() {
	// Scanner to read from standard input
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		// Read the input URL and trim space to remove new line character
		inputURL := strings.TrimSpace(scanner.Text())

		// Make HTTP GET request
		resp, err := http.Get(inputURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching URL: %v\n", err)
			continue // Skip to the next input instead of exiting the program
		}

		// Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading response: %v\n", err)
			resp.Body.Close()
			continue
		}
		resp.Body.Close()

		// Compile the regular expression to find input tags and extract the name attribute
		inputTagRegex := regexp.MustCompile(`<input.*?>`)
		nameAttrRegex := regexp.MustCompile(`name=["']([^"']+)`)

		// Find all input tags in the content
		inputTags := inputTagRegex.FindAllString(string(body), -1)

		// Ensure the input URL ends with a slash
		if !strings.HasSuffix(inputURL, "/") {
			inputURL += "/"
		}

		// Create a URL object from the base URL
		parsedURL, err := url.Parse(inputURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing URL: %v\n", err)
			continue
		}

		// Query values container
		values := url.Values{}

		for _, tag := range inputTags {
			// Extract the name attribute from each input tag
			nameMatches := nameAttrRegex.FindStringSubmatch(tag)
			if len(nameMatches) > 1 {
				// Add the name attribute to the URL's query string with a placeholder value
				values.Set(nameMatches[1], "test")
			}
		}

		// Attach the query parameters to the URL
		parsedURL.RawQuery = values.Encode()

		// Print the constructed URL
		fmt.Println(parsedURL.String())
	}

	// Check for any errors that may have occurred during the scan
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading standard input: %v\n", err)
	}
}


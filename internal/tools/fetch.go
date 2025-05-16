package tools

import (
	"bytes"
	"io"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

var FetchUserAgent string

type FetchInput struct {
	URL string `json:"url" jsonschema_description:"The URL of the webpage to fetch."`
}

var multiWhitespaceRegex = regexp.MustCompile(`\s[\s]+`)

func htmlToText(content []byte) []byte {
	var text strings.Builder
	tokenizer := html.NewTokenizer(bytes.NewReader(content))
	previousStartTokenTest := tokenizer.Token()
	for {
		tt := tokenizer.Next()
		if tt == html.ErrorToken {
			break
		}
		if tt == html.StartTagToken {
			previousStartTokenTest = tokenizer.Token()
			continue
		}
		if tt == html.TextToken {
			if previousStartTokenTest.Data == "script" || previousStartTokenTest.Data == "style" {
				continue
			}
			nodeText := strings.TrimSpace(html.UnescapeString(string(tokenizer.Text())))
			if nodeText != "" {
				nodeText = multiWhitespaceRegex.ReplaceAllString(nodeText, " ")
				text.WriteString(nodeText)
				text.WriteRune(' ')
			}
		}
	}
	return []byte(strings.TrimSpace(text.String()))
}

func init() {
	Register(ToolDefinition{
		Name:        "fetch",
		Description: "Fetch the contents of a given webpage URL.",
		Parameters:  GenerateSchema[FetchInput](),
		Function: WithDecodedInput(func(val FetchInput) (string, error) {
			req, err := http.NewRequest("GET", val.URL, nil)
			if err != nil {
				return "", err
			}
			if FetchUserAgent != "" {
				req.Header.Set("User-Agent", FetchUserAgent)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return "", err
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return "", err
			}

			if strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
				body = htmlToText(body)
			}

			return string(body), nil
		}),
	})
}

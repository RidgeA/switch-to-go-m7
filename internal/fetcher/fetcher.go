package fetcher

import (
	"bytes"
	"context"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	TagAnchor  = []byte("a")
	TagHeader1 = []byte("h1")

	AttrHref = []byte("href")
	AttrId   = []byte("id")

	ValFirstHeading = []byte("firstHeading")
)

func Fetch(ctx context.Context, url string) (string, []string) {
	body, err := load(ctx, url)
	if err != nil {
		fmt.Printf("Failed to fetch url %s, error: %s\n", url, err)
		return "", nil
	}

	return parse(url, body)
}

func load(ctx context.Context, url string) ([]byte, error) {

	client := http.Client{}

	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make a request: %w", err)
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to make a request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

func parse(url string, data []byte) (string, []string) {
	tokenizer := html.NewTokenizer(bytes.NewReader(data))
	title := ""
	links := make([]string, 0)

loop:
	for {
		token := tokenizer.Next()
		switch token {
		case html.ErrorToken:
			if tokenizer.Err() != io.EOF {
				fmt.Printf("Decode error: %s\n", tokenizer.Err())
			}
			break loop
		case html.StartTagToken, html.SelfClosingTagToken:
			tag, _ := tokenizer.TagName()
			switch {
			case is(tag, TagAnchor):
				var attr, value []byte
				var next = true
				for ; next; {
					attr, value, next = tokenizer.TagAttr()
					if !is(attr, AttrHref) {
						continue
					}
					normalized := normalizeUrl(url, string(value))
					if normalized != "" && normalized != url {
						links = append(links, normalized)
					}
					break
				}
			case is(tag, TagHeader1):
				var attr, value []byte
				var next = true
				for ; next; {
					attr, value, next = tokenizer.TagAttr()
					if is(attr, AttrId) && is(value, ValFirstHeading) {
						title = lookupHeader(tokenizer)
					}
				}
			}

		}
	}

	return title, links
}

func is(tag, compare []byte) bool {
	return bytes.Compare(tag, compare) == 0
}

func lookupHeader(tokenizer *html.Tokenizer) string {
	for {
		token := tokenizer.Next()
		switch token {
		case html.TextToken:
			return string(tokenizer.Text())
		case html.EndTagToken:
			return ""
		}
	}
}

func normalizeUrl(baseUrl, link string) string {
	baseParsed, err := url.Parse(baseUrl)
	if err != nil {
		return ""
	}
	linkParsed, err := url.Parse(link)
	if err != nil {
		return ""
	}

	resolved := baseParsed.ResolveReference(linkParsed)

	resolved.Fragment = ""
	return resolved.String()
}

package pkg

import (
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// Crawl request a given Page from both domains and parses it for more pages.
func Crawl(p Page, httpClient *http.Client, m *Migration) {
	// @todo Make URLs part of Page to remove dependency on *Migration.
	oldURL, newURL := m.makeURLs(p)

	oldResponse, err := httpClient.Get(oldURL.String())
	if nil != err {
		panic(err)
	}
	defer oldResponse.Body.Close()
	p.OldResponse = oldResponse

	newResponse, err := httpClient.Get(newURL.String())
	if nil != err {
		panic(err)
	}
	newResponse.Body.Close()
	p.NewResponse = newResponse

	for _, path := range parse(oldResponse.Body) {
		p.Links = append(p.Links, path)
	}

	// @todo Inject m.Pages channel so that property can be made private.
	m.Pages <- p
}

func parse(body io.ReadCloser) []string {
	z := html.NewTokenizer(body)

	paths := []string{}

	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			return paths

		case tt == html.StartTagToken:
			t := z.Token()
			if t.Data == "a" {
				if href := parseHref(t); href != "" {
					paths = append(paths, href)
				}
			}
		}
	}
}

func parseHref(t html.Token) string {
	href := ""

	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			break
		}
	}

	if i := strings.Index(href, "#"); i >= 0 {
		href = href[0:i]
	}

	lenHref := len(href)

	// Cloudflare email obfuscation.
	if lenHref > 26 && "/cdn-cgi/l/email-protection" == href[0:27] {
		return ""
	}

	// @todo Handle absolute hrefs.
	if lenHref > 6 && "http://" == href[0:7] {
		return ""
	}
	if lenHref > 7 && "https://" == href[0:8] {
		return ""
	}
	if lenHref > 1 && "//" == href[0:2] {
		return ""
	}

	// @todo Handle absolute paths.
	if lenHref > 0 && "/" != href[0:1] {
		return ""
	}

	return href
}

package pkg

import (
	"net/http"
	"net/url"
)

// Migration represents a move between two Sites.
type Migration struct {
	oldSite *Site
	newSite *Site

	Pages    chan Page
	visiting map[string]Page
	visited  map[string]Page
}

// NewMigration instantiates a new Migration.
func NewMigration(oldURL, newURL *url.URL) *Migration {
	return &Migration{
		&Site{URL: oldURL},
		&Site{URL: newURL},
		make(chan Page, 1024),
		map[string]Page{},
		map[string]Page{},
	}
}

// Crawl initiates a comparison of the old and new Site.
func (m *Migration) Crawl(httpClient *http.Client, printer Printer, reporter Reporter) {
	m.Pages <- Page{Path: "/"}

mainloop:
	for {
		select {
		case page := <-m.Pages:
			if page.isCrawled() {
				printer(page)

				delete(m.visiting, page.Path)
				m.visited[page.Path] = page

				for _, path := range page.Links {
					m.Pages <- Page{Path: path}
				}
			} else {
				_, hasVisited := m.visited[page.Path]
				_, isVisiting := m.visiting[page.Path]

				if !hasVisited && !isVisiting {
					m.visiting[page.Path] = page
					go Crawl(page, httpClient, m)
				}
			}

		default:
			if m.isCrawled() {
				break mainloop
			}
		}
	}

	close(m.Pages)

	pages := []Page{}
	for _, page := range m.visited {
		pages = append(pages, page)
	}

	reporter(pages)
}

func (m *Migration) isCrawled() bool {
	return len(m.visiting) == 0
}

func (m *Migration) makeURLs(p Page) (*url.URL, *url.URL) {
	oldURL, _ := url.Parse(m.oldSite.URL.String())
	oldURL.Path = p.Path

	newURL, _ := url.Parse(m.newSite.URL.String())
	newURL.Path = p.Path

	return oldURL, newURL
}

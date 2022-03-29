package pkg

import "net/http"

// Page represents a page within a Site.
type Page struct {
	Path        string
	OldResponse *http.Response
	NewResponse *http.Response
	Links       []string
}

func (p *Page) isCrawled() bool {
	return p.OldResponse != nil && p.NewResponse != nil && p.OldResponse.StatusCode != 0 && p.NewResponse.StatusCode != 0
}

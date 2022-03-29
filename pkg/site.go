package pkg

import "net/url"

type Site struct {
	URL           *url.URL
	TrailingSlash bool
}

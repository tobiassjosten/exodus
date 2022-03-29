package main

import (
	"net/url"
	"os"

	"github.com/tobiassjosten/exodus/pkg"
)

func main() {
	oldOrigin, newOrigin := parseOrigins()

	httpClient := pkg.NewHTTPClient(pkg.RateLimit{
		// @todo Read this from a command line flag.
		// oldOrigin.Host: 2,
		newOrigin.Host: 1, // one request per second
	})

	site := pkg.NewMigration(oldOrigin, newOrigin)
	site.Crawl(httpClient, pkg.DefaultPrinter, pkg.DefaultReporter)
}

func parseOrigins() (oldOrigin, newOrigin *url.URL) {
	args := os.Args

	if len(args) != 3 {
		panic("usage: exodus <old origin> <new origin>")
	}

	oldOrigin = urlize(args[1])
	newOrigin = urlize(args[2])

	return
}

func urlize(origin string) *url.URL {
	u, err := url.Parse(origin)
	if nil != err {
		panic(err)
	}

	if "" == u.Scheme {
		u.Scheme = "http"
		// @todo Make a test request to check for HTTP->HTTPS redirectiom.
	}

	return u
}

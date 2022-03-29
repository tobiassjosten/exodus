package pkg

import "fmt"

type Reporter func([]Page)

var DefaultReporter = func(pages []Page) {
	fmt.Printf("Crawled %d pages\n", len(pages))
}

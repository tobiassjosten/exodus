package pkg

import (
	"fmt"
	"net/url"
	"strconv"
)

// Printer prints the results of a Page.
type Printer func(Page)

// DefaultPrinter is the default Printer.
var DefaultPrinter = func(page Page) {
	oldStatusCode := page.OldResponse.StatusCode
	oldColored := colorStatusCode(oldStatusCode)
	oldIsSuccessful := oldStatusCode >= 200 && oldStatusCode < 300

	newStatusCode := page.NewResponse.StatusCode
	newColored := colorStatusCode(newStatusCode)
	newIsSuccessful := newStatusCode >= 200 && newStatusCode < 300
	newIsRedirection := newStatusCode >= 300 && newStatusCode < 400

	switch {
	case oldIsSuccessful && newIsSuccessful:
		fmt.Printf("%s %s\n", newColored, page.Path)

	case oldStatusCode == newStatusCode:
		switch {
		case newIsRedirection:
			fmt.Printf("%s %s -> %s\n", newColored, page.Path, formatLocation(page))

		default:
			fmt.Printf("%s %s\n", newColored, page.Path)
		}

	default:
		switch {
		case oldIsSuccessful:
			fmt.Printf("%s %s\n", newColored, page.Path)

		default:
			fmt.Printf("%s %s %s\n", oldColored, newColored, page.Path)
		}
	}
}

func colorStatusCode(statusCode int) string {
	switch {
	case statusCode >= 200 && statusCode < 300:
		return fmt.Sprintf("\033[34m%s\033[0m", strconv.Itoa(statusCode))

	case statusCode >= 300 && statusCode < 400:
		return fmt.Sprintf("\033[33m%s\033[0m", strconv.Itoa(statusCode))

	case statusCode >= 400:
		return fmt.Sprintf("\033[31m%s\033[0m", strconv.Itoa(statusCode))

	default:
		return strconv.Itoa(statusCode)
	}
}

func formatLocation(page Page) string {
	location := page.NewResponse.Header.Get("Location")

	from := page.OldResponse.Request.URL
	to, err := url.Parse(location)
	if err != nil {
		return location
	}

	if from.Host != to.Host {
		return to.Host
	}

	return to.Path
}

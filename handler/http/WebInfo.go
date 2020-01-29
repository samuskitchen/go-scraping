package http

import (
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// GetTitleAndLogo Get the Logo and the title of the indicated address
func GetTitleAndLogo(address string) (string, string, error) {
	// Make HTTP GET request
	addressComplete := "https://www." + address + "/"
	response, err := http.Get(addressComplete)

	if err != nil {
		log.Println(err)
		return "", "", err
	}

	defer response.Body.Close()
	if response.StatusCode != 200 {
		log.Printf("status code error: %d %s", response.StatusCode, response.Status)
	}

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Println("Error loading HTTP response body. ", err)
		return "", "", err
	}

	pageTitle := document.Find("title").Contents().First().Text()
	var pageLogo string

	document.Find("link").Each(func(index int, element *goquery.Selection) {
		value, exists := element.Attr("type")
		if exists && validateType(value) {
			valueHref, _ := element.Attr("href")
			pageLogo = valueHref
		}
	})

	document.Find("meta").Each(func(index int, element *goquery.Selection) {
		value, exists := element.Attr("itemprop")
		if exists && "image" == value {
			valueContent, _ := element.Attr("content")
			pageLogo = valueContent
		}
	})

	if !strings.Contains(pageLogo, "https") {
		pageLogo = addressComplete + pageLogo
	}

	return pageTitle, pageLogo, err
}

func validateType(value string) bool {
	types := []string{"image/x-icon", "image/icon", "image/vnd.microsoft.icon", "image/svg+xml", "image/png", "image/jpg"}

	for _, element := range types {
		if element == value {
			return true
		}
	}

	return false
}

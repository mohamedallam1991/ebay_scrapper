package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func check(error error) {
	if error != nil {
		fmt.Println(error)
	}
}

func getHtml(url string) *http.Response {
	response, error := http.Get(url)
	check(error)

	if response.StatusCode > 400 {
		fmt.Printf("status code: ", response.StatusCode)
	}
	return response
}

func writeCsv(scrapedData []string) {
	fileName := "data.csv"
	file, error := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	check(error)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	error = writer.Write(scrapedData)
	check(error)
}
func scrapePageData(doc *goquery.Document) {
	doc.Find("ul.srp-results>li.s-item").Each(func(index int, item *goquery.Selection) {
		itemSpan := item.Find("a.s-item__link")

		title := strings.TrimSpace(itemSpan.Text())
		href, _ := itemSpan.Attr("href")
		price := strings.TrimSpace(item.Find("span.s-item__price").Text())

		item_subtitle := strings.TrimSpace(item.Find("div.s-item__subtitle").Text())
		pruchase_options := strings.TrimSpace(item.Find("span.s-item__purchase-options-with-icon").Text())
		// ebay_item_name := strings.TrimSpace(item.Find("span.SECONDARY_INFO").Text())
		// item_subtitle := strings.Trim(reviews_span, "product ratings")
		// item_subtitle := strings.TrimSpace(item_subtitle)

		reviews_span := strings.TrimSpace(item.Find("span.s-item__reviews-count").Text())
		reviews_count := strings.Trim(reviews_span, "product")
		reviews_count = strings.Trim(reviews_span, "ratings")
		reviews := strings.TrimSpace(reviews_count)

		location_span := strings.TrimSpace(item.Find("span.s-item__location").Text())
		location_from := strings.Trim(location_span, "from")
		location := strings.TrimSpace(location_from)

		scrapedData := []string{title, price, location, item_subtitle, reviews, pruchase_options, href}
		writeCsv(scrapedData)
	})
}

var urls = []string{}

func getUrls(url string) {
	response := getHtml(url)
	defer response.Body.Close()
	doc, error := goquery.NewDocumentFromReader(response.Body)
	check(error)

	urls = []string{url}
	// numbers := []int{1, 2, 3, 4}
	// numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for i := 3; i < 10; i++ {
		href, _ := doc.Find("nav.pagination>a.pagination__next").Attr("href")
		// if contains(urls, href) == true {
		// 	break
		// } else {
		// 	urls = append(urls, href)
		// }

		urls = append(urls, href)
	}

	// for i := 3; i < 5; i++ {
	// 	fmt.Printf("%+v\n", i)
	// 	lastUrl := urls[len(urls)-1]
	// 	withoutPage := lastUrl[:len(lastUrl)-1]
	// 	newUrl := withoutPage + string(i)
	// 	urls = append(urls, newUrl)
	// }

}

func contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func main() {
	url := "https://www.ebay.com/sch/i.html?_from=R40&_nkw=iphone&_sacat=0&_ipg=240"
	// url := "https://www.ebay.com/sch/i.html?_from=R40&_nkw=airpod+pro&_sacat=0&_ipg=240&_pgn=4"
	// getUrls(url)
	// for _, url := range urls {
	// 	// fmt.Printf("%+v\n", index)
	// 	fmt.Printf("%+v\n", url)

	// 	response := getHtml(url)

	// 	defer response.Body.Close()
	// 	doc, error := goquery.NewDocumentFromReader(response.Body)
	// 	check(error)
	// 	scrapePageData(doc)
	// }

	// urls := []string{"https://www.ebay.com/sch/i.html?_from=R40&_nkw=airpod+pro&_sacat=0&_ipg=240"}

	// check(error)
	// scrapePageData(doc)

	for {
		var previousUrl string
		response := getHtml(url)
		defer response.Body.Close()
		doc, error := goquery.NewDocumentFromReader(response.Body)
		check(error)
		scrapePageData(doc)
		href, _ := doc.Find("nav.pagination>a.pagination__next").Attr("href")
		fmt.Printf("%+v\n", previousUrl)
		fmt.Printf("%+v\n", href)

		if href == previousUrl {
			break
		} else {
			url = href
			previousUrl = href
		}
	}

}

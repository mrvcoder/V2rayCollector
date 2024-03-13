package collector

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/projectdiscovery/gologger"
	"net/http"
)

func loadMore(link string) *goquery.Document {
	req, _ := http.NewRequest("GET", link, nil)
	fmt.Println(link)
	resp, _ := client.Do(req)
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	return doc
}

func HttpRequest(url string) *http.Response {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		gologger.Fatal().Msg(fmt.Sprintf("Error When requesting to: %s Error : %s", url, err))
	}

	resp, err := client.Do(req)
	if err != nil {
		gologger.Fatal().Msg(err.Error())
	}
	return resp
}

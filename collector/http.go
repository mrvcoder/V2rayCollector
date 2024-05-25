package collector

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/projectdiscovery/gologger"
	"github.com/valyala/fasthttp"
)

func loadMore(link string) *goquery.Document {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(link)
	fmt.Println(link)
	resp := fasthttp.AcquireResponse()

	if err := client.Do(req, resp); err != nil {
		gologger.Fatal().Msg(err.Error())
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
	if err != nil {
		gologger.Fatal().Msg(err.Error())
	}

	return doc
}

func HttpRequest(url string) *fasthttp.Response {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)

	resp := fasthttp.AcquireResponse()

	if err := client.Do(req, resp); err != nil {
		gologger.Fatal().Msg(err.Error())
	}
	return resp
}

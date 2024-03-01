package collector

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strconv"
	"strings"
)

func loadMore(link string) *goquery.Document {
	req, _ := http.NewRequest("GET", link, nil)
	fmt.Println(link)
	resp, _ := client.Do(req)
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	return doc
}

// GetMessages تعریف GetMessages با یک کانال برای ارسال نتیجه
func GetMessages(length int, doc *goquery.Document, number string, channel string, resultChan chan *goquery.Document) {
	// فراخوانی GetMessages و انتقال نتیجه از طریق کانال
	resultChan <- GetMessagesHelper(length, doc, number, channel)
}

// GetMessagesHelper تعریف GetMessagesHelper برای انجام عملیات واقعی و بازگشت نتیجه
func GetMessagesHelper(length int, doc *goquery.Document, number string, channel string) *goquery.Document {
	// اجرای عملیات GetMessages به صورت معمول
	x := loadMore(channel + "?before=" + number)

	html2, _ := x.Html()
	reader2 := strings.NewReader(html2)
	doc2, _ := goquery.NewDocumentFromReader(reader2)

	// _, exist := doc.Find(".js-messages_more").Attr("href")
	doc.Find("body").AppendSelection(doc2.Find("body").Children())

	newDoc := goquery.NewDocumentFromNode(doc.Selection.Nodes[0])
	// fmt.Println(newDoc.Find(".js-messages_more").Attr("href"))
	messages := newDoc.Find(".js-widget_message_wrap").Length()

	fmt.Println(messages)
	if messages > length {
		return newDoc
	} else {
		num, _ := strconv.Atoi(number)
		n := num - 21
		if n > 0 {
			ns := strconv.Itoa(n)
			GetMessagesHelper(length, newDoc, ns, channel)
		} else {
			return newDoc
		}
	}

	return newDoc
}

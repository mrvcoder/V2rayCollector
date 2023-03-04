package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var client = &http.Client{}

func main() {

	channels := []string{
		"https://t.me/s/v2rayng_fa2",
		"https://t.me/s/v2rayng_org",
		"https://t.me/s/V2rayNGvpni",
		"https://t.me/s/custom_14",
		"https://t.me/s/v2rayNG_VPNN",
		"https://t.me/s/FreeV2rays{all_messages}",
	}

	configs := map[string]string{
		"vmess":  "",
		"trojan": "",
		"vless":  "",
		"ss":     "",
	}

	protocol := ""
	all_messages := false
	for i := 0; i < len(channels); i++ {
		if strings.Contains(channels[i], "{all_messages}") {
			all_messages = true
			channels[i] = strings.Split(channels[i], "{all_messages}")[0]
		}

		req, err := http.NewRequest("GET", channels[i], nil)
		if err != nil {
			log.Fatalf("Error When requesting to: %d Error : %s", channels[i], err)
		}

		resp, err1 := client.Do(req)
		if err1 != nil {
			log.Fatal(err1)
		}
		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		messages := doc.Find(".tgme_widget_message_wrap").Length()
		link, exist := doc.Find(".tme_messages_more").Attr("href")
		if messages < 150 && exist == true {
			number := strings.Split(link, "=")[1]
			doc = GetMessages(150, doc, number, channels[i])
		}

		if all_messages == true {
			doc.Find(".tgme_widget_message_text").Each(func(j int, s *goquery.Selection) {
				// For each item found, get the band and title
				code := s.Text()
				protocol = strings.Split(code, "://")[0]
				for proto, _ := range configs {
					if protocol == proto {
						configs[proto] += code + "\n"
					}
				}
			})
		} else {
			doc.Find("code").Each(func(j int, s *goquery.Selection) {
				// For each item found, get the band and title
				code := s.Text()
				protocol = strings.Split(code, "://")[0]
				for proto, _ := range configs {
					if protocol == proto {
						configs[proto] += code + "\n"
					}
				}
			})
		}

	}

	for proto, configcontent := range configs {
		WriteToFile(configcontent, proto+"_iran.txt")
	}

}

func WriteToFile(fileContent string, filePath string) {

	// Check if the file exists
	if _, err := os.Stat(filePath); err == nil {
		// If the file exists, clear its content
		err = ioutil.WriteFile(filePath, []byte{}, 0644)
		if err != nil {
			fmt.Println("Error clearing file:", err)
			return
		}
	} else if os.IsNotExist(err) {
		// If the file does not exist, create it
		_, err = os.Create(filePath)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
	} else {
		// If there was some other error, print it and return
		fmt.Println("Error checking file:", err)
		return
	}

	// Write the new content to the file
	err := ioutil.WriteFile(filePath, []byte(fileContent), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("File written successfully")
}

func load_more(link string) *goquery.Document {
	req, _ := http.NewRequest("GET", link, nil)
	fmt.Println(link)
	resp, _ := client.Do(req)
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	return doc
}

func GetMessages(length int, doc *goquery.Document, number string, channel string) *goquery.Document {
	x := load_more(channel + "?before=" + number)

	html1, _ := x.Html()
	html2, _ := doc.Html()

	combinedHtml := strings.Join([]string{html1, html2}, "")

	newDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(combinedHtml))

	messages := newDoc.Find(".tgme_widget_message_wrap").Length()

	fmt.Println(messages)
	if messages > length {
		return newDoc
	} else {
		num, _ := strconv.Atoi(number)
		n := num + 1
		ns := strconv.Itoa(n)
		GetMessages(length, newDoc, ns, channel)
	}

	return newDoc
}

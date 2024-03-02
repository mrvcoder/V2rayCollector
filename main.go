package main

import (
	"flag"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/jszwec/csvutil"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
)

var (
	client       = &http.Client{}
	max_messages = 100
	ConfigsNames = "@Vip_Security join us"
	configs      = map[string]string{
		"ss":     "",
		"vmess":  "",
		"trojan": "",
		"vless":  "",
		"mixed":  "",
	}
	myregex = map[string]string{
		"ss":     `(?m)^ss:\/\/.+(%3A%40|#)`,
		"vmess":  `(?m)^vmess:\/\/.+(%3A%40|#)`,
		"trojan": `(?m)^trojan:\/\/.+(%3A%40|#)`,
		"vless":  `(?m)^vless:\/\/.+(%3A%40|#)`,
	}
	sort = flag.Bool("sort", false, "sort from latest to oldest (default : false)")
)

type ChannelsType struct {
	URL             string `csv:"URL"`
	AllMessagesFlag bool   `csv:"AllMessagesFlag"`
}

func main() {

	gologger.DefaultLogger.SetMaxLevel(levels.LevelDebug)
	flag.Parse()

	file_data, _ := readFileContent("./channels.csv")
	var channels []ChannelsType
	if err := csvutil.Unmarshal([]byte(file_data), &channels); err != nil {
		gologger.Fatal().Msg("error: " + err.Error())
	}

	// loop through the channels lists
	for _, channel := range channels {

		// change url
		channel.URL = ChangeUrlToTelegramWebUrl(channel.URL)

		// get channel messgages
		resp := HttpRequest(channel.URL)
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		resp.Body.Close()

		if err != nil {
			gologger.Error().Msg(err.Error())
		}

		fmt.Println("---------------------------------------")
		gologger.Info().Msg("Crawling " + channel.URL)
		fmt.Println("")
		CrawlForV2ray(doc, channel.URL, channel.AllMessagesFlag)
		fmt.Println("")
		gologger.Info().Msg("Crawled " + channel.URL + " ! ")
		fmt.Println("---------------------------------------")
	}

	gologger.Info().Msg("Creating output files !")
	for proto, configcontent := range configs {

		lines := RemoveDuplicate(configcontent)

		if *sort {
			// 		reverse mode :
			lines_arr := strings.Split(configcontent, "\n")
			lines_arr = reverse(lines_arr)
			lines = strings.Join(lines_arr, "\n")
		}

		WriteToFile(lines, proto+"_iran.txt")

	}

	gologger.Info().Msg("All Done :D")

}

func CrawlForV2ray(doc *goquery.Document, channel_link string, HasAllMessagesFlag bool) {
	// here we are updating our DOM to include the x messages
	// in our DOM and then extract the messages from that DOM
	messages := doc.Find(".tgme_widget_message_wrap").Length()
	link, exist := doc.Find(".tgme_widget_message_wrap .js-widget_message").Last().Attr("data-post")

	if messages < max_messages && exist {
		number := strings.Split(link, "/")[1]
		doc = GetMessages(max_messages, doc, number, channel_link)
	}

	// extract v2ray based on message type and store configs at [configs] map
	if HasAllMessagesFlag {
		// get all messages and check for v2ray configs
		fmt.Println(doc.Find(".js-widget_message_wrap").Length())
		doc.Find(".tgme_widget_message_text").Each(func(j int, s *goquery.Selection) {
			// For each item found, get the band and title
			message_text, _ := s.Html()
			lines := strings.Split(message_text, "<br/>")
			for a := 0; a < len(lines); a++ {
				for _, regex_value := range myregex {
					re := regexp.MustCompile(regex_value)
					matches := re.FindStringSubmatch(lines[a])
					if len(matches) > 0 {
						// add extracted values to configs[]
						lines[a] = "\n" + matches[0] + ConfigsNames
						lines[a] = strings.TrimSpace(lines[a])
						for proto, _ := range configs {
							if strings.Contains(lines[a], proto) {
								configs["mixed"] += "\n" + lines[a] + "\n"
							}
						}
					}
				}

			}

		})
	} else {
		// get only messages that are inside code or pre tag and check for v2ray configs
		doc.Find("code,pre").Each(func(j int, s *goquery.Selection) {
			message_text, _ := s.Html()
			lines := strings.Split(message_text, "<br/>")
			for a := 0; a < len(lines); a++ {
				for proto_regex, regex_value := range myregex {

					re := regexp.MustCompile(regex_value)
					matches := re.FindStringSubmatch(lines[a])
					if len(matches) > 0 {
						lines[a] = "\n" + matches[0] + ConfigsNames
						// add extracted values to configs[]
						lines[a] = strings.TrimSpace(lines[a])
						configs[proto_regex] += "\n" + lines[a] + "\n"
					}
				}
			}
		})
	}
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

	html2, _ := x.Html()
	reader2 := strings.NewReader(html2)
	doc2, _ := goquery.NewDocumentFromReader(reader2)

	doc.Find("body").AppendSelection(doc2.Find("body").Children())

	newDoc := goquery.NewDocumentFromNode(doc.Selection.Nodes[0])
	messages := newDoc.Find(".js-widget_message_wrap").Length()

	if messages > length {
		return newDoc
	} else {
		num, _ := strconv.Atoi(number)
		n := num - 21
		if n > 0 {
			ns := strconv.Itoa(n)
			GetMessages(length, newDoc, ns, channel)
		} else {
			return newDoc
		}
	}

	return newDoc
}

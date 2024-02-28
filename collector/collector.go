// by MehrabSp

package collector

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var client = &http.Client{
	Transport: &http.Transport{
		// limit
		MaxIdleConns:    10,
		IdleConnTimeout: time.Second,
	},
}

// speed up
var wg sync.WaitGroup
var mutex sync.Mutex

// func createLinkSlice(linksString string, defaultChannelsConfig []string) []string {
// 	if linksString == "" {
// 		return defaultChannelsConfig
// 	}

// 	links := strings.Fields(linksString)
// 	linkSlice := make([]string, len(links))
// 	copy(linkSlice, links)

// 	return linkSlice
// }

// App
func GetConfigs(defaultChannelsConfig []string) string {
	// defaultChannelsConfig := createLinkSlice(linksString, defaultChannelsConfig)

	configs := map[string]string{
		"ss":     "",
		"vmess":  "",
		"trojan": "",
		"vless":  "",
		"mixed":  "",
	}
	myregex := map[string]string{
		"ss":     `(.{3})ss:\/\/`,
		"vmess":  `vmess:\/\/`,
		"trojan": `trojan:\/\/`,
		"vless":  `vless:\/\/`,
	}

	//protocol := ""
	// for _, channel := range channels {
	for i := 0; i < len(defaultChannelsConfig); i++ {
		wg.Add(1)
		go func(channel string) {
			defer wg.Done()
			all_messages := false
			if strings.Contains(channel, "{all_messages}") {
				all_messages = true
				channel = strings.Split(channel, "{all_messages}")[0]
			}

			req, err := http.NewRequest("GET", channel, nil)
			if err != nil {
				log.Fatalf("Error When requesting to: %s Error : %s", channel, err) //fixed
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
			link, exist := doc.Find(".tgme_widget_message_wrap .js-widget_message").Last().Attr("data-post")
			if messages < 100 && exist {
				number := strings.Split(link, "/")[1]
				fmt.Println(number)

				doc = GetMessages(100, doc, number, channel)
			}

			if all_messages {
				mutex.Lock()
				fmt.Println(doc.Find(".js-widget_message_wrap").Length())
				doc.Find(".tgme_widget_message_text").Each(func(j int, s *goquery.Selection) {
					// For each item found, get the band and title
					message_text := s.Text()
					lines := strings.Split(message_text, "\n")
					for a := 0; a < len(lines); a++ {
						for _, regex_value := range myregex {
							re := regexp.MustCompile(regex_value)
							lines[a] = re.ReplaceAllStringFunc(lines[a], func(match string) string {
								return "\n" + match
							})
						}
						for proto := range configs {
							if strings.Contains(lines[a], proto) {
								configs["mixed"] += "\n" + lines[a] + "\n"
							}
						}
					}

				})
				mutex.Unlock()
			} else {
				mutex.Lock()
				doc.Find("code,pre").Each(func(j int, s *goquery.Selection) {
					// For each item found, get the band and title
					message_text := s.Text()
					lines := strings.Split(message_text, "\n")
					for a := 0; a < len(lines); a++ {
						for proto_regex, regex_value := range myregex {
							re := regexp.MustCompile(regex_value)
							lines[a] = re.ReplaceAllStringFunc(lines[a], func(match string) string {
								if proto_regex == "ss" {
									if match[:3] == "vme" {
										return "\n" + match
									} else if match[:3] == "vle" {
										return "\n" + match
									} else {
										return "\n" + match
									}
								} else {
									return "\n" + match
								}
							})

							if len(strings.Split(lines[a], "\n")) > 1 {
								myconfigs := strings.Split(lines[a], "\n")
								for i := 0; i < len(myconfigs); i++ {
									if myconfigs[i] != "" {
										re := regexp.MustCompile(regex_value)
										myconfigs[i] = strings.ReplaceAll(myconfigs[i], " ", "")
										match := re.FindStringSubmatch(myconfigs[i])
										if len(match) >= 1 {
											if proto_regex == "ss" {
												if match[1][:3] == "vme" {
													configs["vmess"] += "\n" + myconfigs[i] + "\n"
												} else if match[1][:3] == "vle" {
													configs["vless"] += "\n" + myconfigs[i] + "\n"
												} else {
													configs["ss"] += "\n" + myconfigs[i][3:] + "\n"
												}
											} else {
												configs[proto_regex] += "\n" + myconfigs[i] + "\n"
											}
										}

									}

								}
							}
						}
					}
				})
				mutex.Unlock()
			}
		}(defaultChannelsConfig[i])
	}

	wg.Wait()

	var myString string
	// _ --> proto
	for _, configcontent := range configs {
		myString += RemoveDuplicate(configcontent)
	}
	return myString

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
			GetMessages(length, newDoc, ns, channel)
		} else {
			return newDoc
		}
	}

	return newDoc
}

func RemoveDuplicate(config string) string {
	lines := strings.Split(config, "\n")

	// Use a map to keep track of unique lines
	uniqueLines := make(map[string]bool)

	// Loop over lines and add unique lines to map
	for _, line := range lines {
		if len(line) > 0 {
			uniqueLines[line] = true
		}
	}

	// Join unique lines into a string
	uniqueString := strings.Join(getKeys(uniqueLines), "\n")

	return uniqueString
}

func getKeys(m map[string]bool) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

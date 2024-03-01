// by MehrabSp

package collector

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
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

// GetConfigs App
func GetConfigs(defaultChannelsConfig []string) map[string]string {
	// defaultChannelsConfig := createLinkSlice(linksString, defaultChannelsConfig)

	configs := map[string]string{
		"ss":     "",
		"vmess":  "",
		"trojan": "",
		"vless":  "",
		"mixed":  "",
	}
	regex := map[string]string{
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
			allMessages := false
			if strings.Contains(channel, "{all_messages}") {
				allMessages = true
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
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					log.Fatal(err)
				}
			}(resp.Body)

			doc, err := goquery.NewDocumentFromReader(resp.Body)

			// ساختن یک کانال برای ارتباط بین گوروتین و مین گوروتین
			resultChan := make(chan *goquery.Document)

			if err != nil {
				log.Fatal(err)
			}

			messages := doc.Find(".tgme_widget_message_wrap").Length()
			link, exist := doc.Find(".tgme_widget_message_wrap .js-widget_message").Last().Attr("data-post")
			if messages < 100 && exist {
				number := strings.Split(link, "/")[1]
				fmt.Println(number)

				// فراخوانی GetMessages به صورت ناهمزمان
				go GetMessages(100, doc, number, channel, resultChan)

				// دریافت نتیجه از کانال
				newDoc := <-resultChan
				// استفاده از newDoc برای کارهای بعدی
				doc = newDoc
				//doc = GetMessages(100, doc, number, channel)
			}

			if allMessages {
				mutex.Lock()
				fmt.Println(doc.Find(".js-widget_message_wrap").Length())
				doc.Find(".tgme_widget_message_text").Each(func(j int, s *goquery.Selection) {
					// For each item found, get the band and title
					messageText := s.Text()
					lines := strings.Split(messageText, "\n")
					for a := 0; a < len(lines); a++ {
						for _, regexValue := range regex {
							re := regexp.MustCompile(regexValue)
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
					messageText := s.Text()
					lines := strings.Split(messageText, "\n")
					for a := 0; a < len(lines); a++ {
						for protoRegex, regexValue := range regex {
							re := regexp.MustCompile(regexValue)
							lines[a] = re.ReplaceAllStringFunc(lines[a], func(match string) string {
								if protoRegex == "ss" {
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
								myConfigs := strings.Split(lines[a], "\n")
								for i := 0; i < len(myConfigs); i++ {
									if myConfigs[i] != "" {
										re := regexp.MustCompile(regexValue)
										myConfigs[i] = strings.ReplaceAll(myConfigs[i], " ", "")
										match := re.FindStringSubmatch(myConfigs[i])
										if len(match) >= 1 {
											if protoRegex == "ss" {
												if match[1][:3] == "vme" {
													configs["vmess"] += "\n" + myConfigs[i] + "\n"
												} else if match[1][:3] == "vle" {
													configs["vless"] += "\n" + myConfigs[i] + "\n"
												} else {
													configs["ss"] += "\n" + myConfigs[i][3:] + "\n"
												}
											} else {
												configs[protoRegex] += "\n" + myConfigs[i] + "\n"
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

	// var myString string
	// // _ --> proto
	// for _, configcontent := range configs {
	// 	myString += RemoveDuplicate(configcontent)
	// }
	// return myString

	// for proto, configcontent := range configs {
	// 	// 		reverse mode :
	// 	// 		lines := strings.Split(configcontent, "\n")
	// 	// 		reversed := reverse(lines)
	// 	// 		WriteToFile(strings.Join(reversed, "\n"), proto+"_iran.txt")
	// 	// 		simple mode :
	// 	WriteToFile(RemoveDuplicate(configcontent), proto+"_iran.txt")
	// }
	return configs
}

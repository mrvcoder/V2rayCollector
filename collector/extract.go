package collector

import (
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strings"
)

func ExtractConfig(Txt string, Tempconfigs []string) string {

	// filename can be "" or mixed
	for protoRegex, regexValue := range Myregex {
		re := regexp.MustCompile(regexValue)
		matches := re.FindStringSubmatch(Txt)
		extractedConfig := ""
		if len(matches) > 0 {
			if protoRegex == "ss" {
				Prefix := strings.Split(matches[0], "ss://")[0]
				if Prefix == "" {
					extractedConfig = "\n" + matches[0]
				} else if Prefix != "vle" { //  (Prefix != "vme" && Prefix != "") always true!
					d := strings.Split(matches[0], "ss://")
					extractedConfig = "\n" + "ss://" + d[1]
				}
			} else if protoRegex == "vmess" {
				extractedConfig = "\n" + matches[0]
			} else {
				extractedConfig = "\n" + matches[0]
			}

			Tempconfigs = append(Tempconfigs, extractedConfig)
			Txt = strings.ReplaceAll(Txt, matches[0], "")
			ExtractConfig(Txt, Tempconfigs)
		}
	}
	d := strings.Join(Tempconfigs, "\n")
	return d
}

func CrawlForV2ray(doc *goquery.Document, channelLink string, HasAllMessagesFlag bool) {
	// here we are updating our DOM to include the x messages
	// in our DOM and then extract the messages from that DOM
	messages := doc.Find(".tgme_widget_message_wrap").Length()
	link, exist := doc.Find(".tgme_widget_message_wrap .js-widget_message").Last().Attr("data-post")

	if messages < MaxMessages && exist {
		number := strings.Split(link, "/")[1]
		doc = GetMessages(MaxMessages, doc, number, channelLink)
	}

	// extract v2ray based on message type and store configs at [configs] map
	if HasAllMessagesFlag {
		// get all messages and check for v2ray configs
		doc.Find(".tgme_widget_message_text").Each(func(j int, s *goquery.Selection) {
			// For each item found, get the band and title
			messageText, _ := s.Html()
			str := strings.Replace(messageText, "<br/>", "\n", -1)
			doc, _ := goquery.NewDocumentFromReader(strings.NewReader(str))
			messageText = doc.Text()
			line := strings.TrimSpace(messageText)
			lines := strings.Split(line, "\n")
			for _, data := range lines {
				extractedConfigs := strings.Split(ExtractConfig(data, []string{}), "\n")
				for _, extractedConfig := range extractedConfigs {
					extractedConfig = strings.ReplaceAll(extractedConfig, " ", "")
					if extractedConfig != "" {

						// check if it is vmess or not
						re := regexp.MustCompile(Myregex["vmess"])
						matches := re.FindStringSubmatch(extractedConfig)

						if len(matches) > 0 {
							extractedConfig = EditVmessPs(extractedConfig, "mixed", false)
							if line != "" {
								Configs["mixed"] += extractedConfig + "\n"
							}
						} else {
							Configs["mixed"] += extractedConfig + "\n"
						}

					}
				}
			}
		})
	} else {
		// get only messages that are inside code or pre tag and check for v2ray configs
		doc.Find("code,pre").Each(func(j int, s *goquery.Selection) {
			messageText, _ := s.Html()
			str := strings.ReplaceAll(messageText, "<br/>", "\n")
			doc, _ := goquery.NewDocumentFromReader(strings.NewReader(str))
			messageText = doc.Text()
			line := strings.TrimSpace(messageText)
			lines := strings.Split(line, "\n")
			for _, data := range lines {
				extractedConfigs := strings.Split(ExtractConfig(data, []string{}), "\n")
				for protoRegex, regexValue := range Myregex {

					for _, extractedConfig := range extractedConfigs {

						re := regexp.MustCompile(regexValue)
						matches := re.FindStringSubmatch(extractedConfig)
						if len(matches) > 0 {
							extractedConfig = strings.ReplaceAll(extractedConfig, " ", "")
							if extractedConfig != "" {
								if protoRegex == "vmess" {
									extractedConfig = EditVmessPs(extractedConfig, protoRegex, false)
									if extractedConfig != "" {
										Configs[protoRegex] += extractedConfig + "\n"
									}
								} else if protoRegex == "ss" {
									Prefix := strings.Split(matches[0], "ss://")[0]
									if Prefix == "" {
										Configs[protoRegex] += extractedConfig + "\n"
									}
								} else {
									Configs[protoRegex] += extractedConfig + "\n"
								}

							}
						}

					}

				}
			}

		})
	}
}

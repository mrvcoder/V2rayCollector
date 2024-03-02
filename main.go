package main

import (
	"encoding/base64"
	"encoding/json"
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
	ConfigFileIds = map[string]int32{
		"ss":     0,
		"vmess":  0,
		"trojan": 0,
		"vless":  0,
		"mixed":  0,
	}
	myregex = map[string]string{
		"ss":     `(?m)...ss:\/\/.+?(%3A%40|#)`,
		"vmess":  `(?m)vmess:\/\/.+`,
		"trojan": `(?m)trojan:\/\/.+?(%3A%40|#)`,
		"vless":  `(?m)vless:\/\/.+?(%3A%40|#)`,
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

		fmt.Println(" ")
		fmt.Println(" ")
		fmt.Println("---------------------------------------")
		gologger.Info().Msg("Crawling " + channel.URL)
		CrawlForV2ray(doc, channel.URL, channel.AllMessagesFlag)
		gologger.Info().Msg("Crawled " + channel.URL + " ! ")
		fmt.Println("---------------------------------------")
		fmt.Println(" ")
		fmt.Println(" ")
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
		lines = strings.TrimSpace(lines)
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
		doc.Find(".tgme_widget_message_text").Each(func(j int, s *goquery.Selection) {
			// For each item found, get the band and title
			message_text, _ := s.Html()
			str := strings.Replace(message_text, "<br/>", "\n", -1)
			doc, _ := goquery.NewDocumentFromReader(strings.NewReader(str))
			message_text = doc.Text()
			line := strings.TrimSpace(message_text)
			lines := strings.Split(line, "\n")
			for _, data := range lines {
				extracted_configs := strings.Split(ExtractConfig(data, []string{}), "\n")

				for _, extractedConfig := range extracted_configs {
					extractedConfig = strings.ReplaceAll(extractedConfig, " ", "")
					if extractedConfig != "" {

						// check if it is vmess or not
						re := regexp.MustCompile(myregex["vmess"])
						matches := re.FindStringSubmatch(extractedConfig)

						if len(matches) > 0 {
							extractedConfig = EditVmessPs(extractedConfig, "mixed")
							if line != "" {
								configs["mixed"] += extractedConfig + "\n"
							}
						} else {
							ConfigFileIds["mixed"] += 1
							extractedConfig = extractedConfig + " - " + strconv.Itoa(int(ConfigFileIds["mixed"]))
							configs["mixed"] += extractedConfig + "\n"
						}

					}
				}
			}
		})
	} else {
		// get only messages that are inside code or pre tag and check for v2ray configs
		doc.Find("code,pre").Each(func(j int, s *goquery.Selection) {
			message_text, _ := s.Html()
			str := strings.ReplaceAll(message_text, "<br/>", "\n")
			doc, _ := goquery.NewDocumentFromReader(strings.NewReader(str))
			message_text = doc.Text()
			line := strings.TrimSpace(message_text)
			lines := strings.Split(line, "\n")
			for _, data := range lines {
				extracted_configs := strings.Split(ExtractConfig(data, []string{}), "\n")
				for proto_regex, regex_value := range myregex {

					for _, extractedConfig := range extracted_configs {

						re := regexp.MustCompile(regex_value)
						matches := re.FindStringSubmatch(extractedConfig)

						if len(matches) > 0 {
							extractedConfig = strings.TrimSpace(extractedConfig)
							extractedConfig = strings.ReplaceAll(extractedConfig, " ", "")
							if extractedConfig != "" {
								if proto_regex == "vmess" {
									extractedConfig = EditVmessPs(extractedConfig, proto_regex)
									if extractedConfig != "" {
										configs[proto_regex] += extractedConfig + "\n"
									}
								} else if proto_regex == "ss" {
									Prefix := strings.Split(matches[0], "ss://")[0]
									if Prefix == "" {
										ConfigFileIds[proto_regex] += 1
										configs[proto_regex] += extractedConfig + " - " + strconv.Itoa(int(ConfigFileIds[proto_regex])) + "\n"
									}
								} else {

									ConfigFileIds[proto_regex] += 1
									configs[proto_regex] += extractedConfig + " - " + strconv.Itoa(int(ConfigFileIds[proto_regex])) + "\n"
								}

							}
						}

					}

				}
			}

		})
	}
}

func ExtractConfig(Txt string, Tempconfigs []string) string {

	// filename can be "" or mixed
	for proto_regex, regex_value := range myregex {
		re := regexp.MustCompile(regex_value)
		matches := re.FindStringSubmatch(Txt)
		extracted_config := ""
		if len(matches) > 0 {
			if proto_regex == "ss" {
				Prefix := strings.Split(matches[0], "ss://")[0]
				if Prefix == "" {
					extracted_config = "\n" + matches[0] + ConfigsNames
				} else if Prefix != "vle" || Prefix != "vme" {
					d := strings.Split(matches[0], "ss://")
					extracted_config = "\n" + "ss://" + d[1] + ConfigsNames
				}
			} else if proto_regex == "vmess" {
				extracted_config = "\n" + matches[0]
			} else {
				extracted_config = "\n" + matches[0] + ConfigsNames
			}

			Tempconfigs = append(Tempconfigs, extracted_config)
			Txt = strings.ReplaceAll(Txt, matches[0], "")
			ExtractConfig(Txt, Tempconfigs)
		}
	}

	return strings.Join(Tempconfigs, "\n")
}

func EditVmessPs(config string, fileName string) string {
	// Decode the base64 string
	if config == "" {
		return ""
	}
	slice := strings.Split(config, "vmess://")
	if len(slice) > 0 {
		decodedBytes, err := base64.StdEncoding.DecodeString(slice[1])
		if err == nil {
			// Unmarshal JSON into a map
			var data map[string]interface{}
			err = json.Unmarshal(decodedBytes, &data)
			if err == nil {
				ConfigFileIds[fileName] += 1
				data["ps"] = ConfigsNames + " - " + strconv.Itoa(int(ConfigFileIds[fileName])) + "\n"
				// marshal JSON into a map
				jsonData, _ := json.Marshal(data)
				// Encode JSON to base64
				base64Encoded := base64.StdEncoding.EncodeToString(jsonData)

				return "vmess://" + base64Encoded
			}
		}
	}

	return ""
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

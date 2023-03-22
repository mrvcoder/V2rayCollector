package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
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
		"https://t.me/s/v2ray_outlineir",
		"https://t.me/s/v2_vmess",
		"https://t.me/s/FreeVlessVpn",
		"https://t.me/s/freeland8",
		"https://t.me/s/vmess_vless_v2rayng",
		"https://t.me/s/PrivateVPNs",
		"https://t.me/s/vmessiran",
		"https://t.me/s/Outline_Vpn",
		"https://t.me/s/vmessq",
		"https://t.me/s/WeePeeN",
		"https://t.me/s/V2rayNG3",
		"https://t.me/s/ShadowsocksM",
		"https://t.me/s/shadowsocksshop",
		"https://t.me/s/v2rayan",
		"https://t.me/ShadowSocks_s",
		"https://t.me/s/VmessProtocol",
		"https://t.me/s/napsternetv_config",
		"https://t.me/s/Easy_Free_VPN",
		"https://t.me/s/V2Ray_FreedomIran",
		"https://t.me/s/V2RAY_VMESS_free",
		"https://t.me/s/v2ray_for_free",
		"https://t.me/s/V2rayN_Free",
		"https://t.me/s/free4allVPN",
		"https://t.me/s/vpn_ocean",
		"https://t.me/s/configV2rayForFree",
		"https://t.me/s/FreeV2rays{all_messages}",
		"https://t.me/s/DigiV2ray",
		"https://t.me/s/v2rayNG_VPN",
		"https://t.me/s/freev2rayssr",
		"https://t.me/s/v2rayn_server",
		"https://t.me/s/Shadowlinkserverr",
		"https://t.me/s/iranvpnet",
		"https://t.me/s/vmess_iran",
		"https://t.me/s/mahsaamoon1",
		"https://t.me/s/V2RAY_NEW",
		"https://t.me/s/v2RayChannel",
		"https://t.me/s/configV2rayNG{all_messages}",
		"https://t.me/s/config_v2ray",
		"https://t.me/s/vpn_proxy_custom",
		"https://t.me/s/vpnmasi{all_messages}",
		"https://t.me/s/v2ray_custom",
		"https://t.me/s/VPNCUSTOMIZE",
		"https://t.me/s/HTTPCustomLand",
		"https://t.me/s/vpn_proxy_custom",
		"https://t.me/s/ViPVpn_v2ray",
		"https://t.me/s/FreeNet1500",
		"https://t.me/s/v2ray_ar{all_messages}",
		"https://t.me/s/beta_v2ray",
		"https://t.me/s/vip_vpn_2022",
		"https://t.me/s/FOX_VPN66",
		"https://t.me/s/VorTexIRN",
		"https://t.me/s/YtTe3la",
		"https://t.me/s/V2RayOxygen",
		"https://t.me/s/Network_442",
		"https://t.me/s/VPN_443",
		"https://t.me/s/v2rayng_v",
		"https://t.me/s/ultrasurf_12",
	}

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
	for i := 0; i < len(channels); i++ {
		all_messages := false
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
		link, exist := doc.Find(".tgme_widget_message_wrap .js-widget_message").Last().Attr("data-post")
		if messages < 100 && exist == true {
			number := strings.Split(link, "/")[1]
			fmt.Println(number)

			doc = GetMessages(100, doc, number, channels[i])
		}

		if all_messages {
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
					for proto, _ := range configs {
						if strings.Contains(lines[a], proto) {
							configs["mixed"] += "\n" + lines[a] + "\n"
						}
					}
				}

			})
		} else {
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
		}

	}

	for proto, configcontent := range configs {
		// 		reverse mode :
		// 		lines := strings.Split(configcontent, "\n")
		// 		reversed := reverse(lines)
		// 		WriteToFile(strings.Join(reversed, "\n"), proto+"_iran.txt")
		// 		simple mode :
		WriteToFile(RemoveDuplicate(configcontent), proto+"_iran.txt")
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

func reverse(lines []string) []string {
	for i := 0; i < len(lines)/2; i++ {
		j := len(lines) - i - 1
		lines[i], lines[j] = lines[j], lines[i]
	}
	return lines
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

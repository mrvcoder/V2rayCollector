package collector

import (
	"encoding/base64"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
)

func EditVmessPs(config string, fileName string, AddConfigName bool) string {
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
				if AddConfigName {
					ConfigFileIds[fileName] += 1
					data["ps"] = ConfigsNames + " - " + strconv.Itoa(int(ConfigFileIds[fileName])) + "\n"
				} else {
					data["ps"] = ""
				}

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

func GetMessages(length int, doc *goquery.Document, number string, channel string) *goquery.Document {
	x := loadMore(channel + "?before=" + number)

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

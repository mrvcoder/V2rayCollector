package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/mrvcoder/V2rayCollector/collector"

	"github.com/PuerkitoBio/goquery"
	"github.com/jszwec/csvutil"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
)

type ChannelsType struct {
	URL             string `csv:"URL"`
	AllMessagesFlag bool   `csv:"AllMessagesFlag"`
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) // تنظیم به تعداد هسته‌های سیستم

	gologger.DefaultLogger.SetMaxLevel(levels.LevelDebug)
	flag.Parse()

	fileData, _ := collector.ReadFileContent("channels.csv")
	var channels []ChannelsType
	if err := csvutil.Unmarshal([]byte(fileData), &channels); err != nil {
		gologger.Fatal().Msg("error: " + err.Error())
	}

	// loop through the channels lists
	for _, channel := range channels {

		// change url
		channel.URL = collector.ChangeUrlToTelegramWebUrl(channel.URL)
		log.Printf("Making request to: %s\n", channel.URL)

		// get channel messages
		resp := collector.HttpRequest(channel.URL)
		if resp.Body() == nil {
			gologger.Error().Msg("BodyStream is nil")
			continue
		}
		//defer func(resp *fasthttp.Response) {
		//	err := resp.CloseBodyStream()
		//	if err != nil {
		//		gologger.Error().Msg("BodyStream is nil")
		//		return
		//	}
		//}(resp) // Make sure to close the body stream

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
		if err != nil {
			gologger.Error().Msg(err.Error())
			return
		}
		// Use doc here

		fmt.Println(" ")
		fmt.Println(" ")
		fmt.Println("---------------------------------------")
		gologger.Info().Msg("Crawling " + channel.URL)
		collector.CrawlForV2ray(doc, channel.URL, channel.AllMessagesFlag)
		gologger.Info().Msg("Crawled " + channel.URL + " ! ")
		fmt.Println("---------------------------------------")
		fmt.Println(" ")
		fmt.Println(" ")
	}

	gologger.Info().Msg("Creating output files !")

	for proto, configcontent := range collector.Configs {
		lines := collector.RemoveDuplicate(configcontent)
		lines = AddConfigNames(lines, proto)
		if *collector.Sort {
			// 		from latest to oldest mode :
			linesArr := strings.Split(lines, "\n")
			linesArr = collector.Reverse(linesArr)
			lines = strings.Join(linesArr, "\n")
		} else {
			// 		from oldest to latest mode :
			linesArr := strings.Split(lines, "\n")
			linesArr = collector.Reverse(linesArr)
			linesArr = collector.Reverse(linesArr)
			lines = strings.Join(linesArr, "\n")
		}
		lines = strings.TrimSpace(lines)
		collector.WriteToFile(lines, proto+"_iran.txt")

	}

	gologger.Info().Msg("All Done :D")

}

func AddConfigNames(config string, configtype string) string {
	configs := strings.Split(config, "\n")
	newConfigs := ""
	for protoRegex, regexValue := range collector.Myregex {

		for _, extractedConfig := range configs {

			re := regexp.MustCompile(regexValue)
			matches := re.FindStringSubmatch(extractedConfig)
			if len(matches) > 0 {
				extractedConfig = strings.ReplaceAll(extractedConfig, " ", "")
				if extractedConfig != "" {
					if protoRegex == "vmess" {
						extractedConfig = collector.EditVmessPs(extractedConfig, configtype, true)
						if extractedConfig != "" {
							newConfigs += extractedConfig + "\n"
						}
					} else if protoRegex == "ss" {
						Prefix := strings.Split(matches[0], "ss://")[0]
						if Prefix == "" {
							collector.ConfigFileIds[configtype] += 1
							newConfigs += extractedConfig + collector.ConfigsNames + " - " + strconv.Itoa(int(collector.ConfigFileIds[configtype])) + "\n"
						}
					} else {

						collector.ConfigFileIds[configtype] += 1
						newConfigs += extractedConfig + collector.ConfigsNames + " - " + strconv.Itoa(int(collector.ConfigFileIds[configtype])) + "\n"
					}
				}
			}

		}
	}
	return newConfigs
}

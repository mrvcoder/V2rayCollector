package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	
	"github.com/jagger235711/V2rayCollector/collector"

	"github.com/PuerkitoBio/goquery"
	"github.com/jszwec/csvutil"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
)

type Config struct {
	Client        *http.Client
	MaxMessages   int
	ConfigsNames  string
	Configs       map[string]string
	ConfigFileIds map[string]int32
	RegexPatterns map[string]string
	Sort         *bool
}

var (
	maxMessages = 100
	client      *http.Client
	cfg         Config
)

func init() {
	// 初始化HTTP客户端配置
	transport := &http.Transport{
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 20 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	// 尝试从环境变量获取代理地址
	if proxyAddr := os.Getenv("HTTP_PROXY"); proxyAddr != "" {
		if proxyURL, err := url.Parse(proxyAddr); err == nil {
			transport.Proxy = http.ProxyURL(proxyURL)
			gologger.Info().Msgf("使用代理服务器: %s", proxyAddr)
		} else {
			gologger.Warning().Msgf("无效的代理地址: %s", proxyAddr)
		}
	} else {
		gologger.Info().Msg("未使用代理服务器")
	}

	client = &http.Client{
		Timeout:   60 * time.Second,
		Transport: transport,
	}

	cfg = Config{
		Client: &http.Client{
			Timeout:   60 * time.Second,
			Transport: transport,
		},
		MaxMessages:  100,
		ConfigsNames: "@Vip_Security join us",
		Configs: map[string]string{
			"ss":     "",
			"vmess":  "",
			"trojan": "",
			"vless":  "",
			"mixed":  "",
		},
		ConfigFileIds: map[string]int32{
			"ss":     0,
			"vmess":  0,
			"trojan": 0,
			"vless":  0,
			"mixed":  0,
		},
		RegexPatterns: map[string]string{
			"ss":     `(?m)(...ss:|^ss:)\/\/.+?(%3A%40|#)`,
			"vmess":  `(?m)vmess:\/\/.+`,
			"trojan": `(?m)trojan:\/\/.+?(%3A%40|#)`,
			"vless":  `(?m)vless:\/\/.+?(%3A%40|#)`,
		},
		Sort: flag.Bool("sort", false, "sort from latest to oldest (default : false)"),
	}
}

type ChannelsType struct {
	URL             string `csv:"URL"`
	AllMessagesFlag bool   `csv:"AllMessagesFlag"`
}

func main() {
	gologger.DefaultLogger.SetMaxLevel(levels.LevelDebug)
	flag.Parse()

	// 读取频道列表
	fileData, err := collector.ReadFileContent("channels.csv")
	if err != nil {
		gologger.Fatal().Msgf("读取channels.csv失败: %v", err)
	}

	var channels []ChannelsType
	if err = csvutil.Unmarshal([]byte(fileData), &channels); err != nil {
		gologger.Fatal().Msgf("解析channels.csv失败: %v", err)
	}

	// 使用worker pool并发处理
	var wg sync.WaitGroup
	ch := make(chan ChannelsType, len(channels))
	results := make(chan error, len(channels))

	// 启动worker
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for channel := range ch {
				if err := processChannel(channel); err != nil {
					results <- err
					return
				}
			}
		}()
	}

	// 分发任务
	for _, channel := range channels {
		ch <- channel
	}
	close(ch)

	// 等待所有worker完成
	wg.Wait()
	close(results)

	// 检查错误
	for err := range results {
		if err != nil {
			gologger.Fatal().Msgf("处理频道时出错: %v", err)
		}
	}

	gologger.Info().Msg("开始生成输出文件...")

	var fileWg sync.WaitGroup
	for proto, configcontent := range cfg.Configs {
		fileWg.Add(1)
		go func(proto, content string) {
			defer fileWg.Done()
			lines := collector.RemoveDuplicate(content)
			lines = AddConfigNames(lines, proto)
			if *cfg.Sort {
				linesArr := strings.Split(lines, "\n")
				linesArr = collector.Reverse(linesArr)
				lines = strings.Join(linesArr, "\n")
			}
			lines = strings.TrimSpace(lines)
			
			// 确保文件存在并清空内容
			// 确保results目录存在
			if err := os.MkdirAll("results", 0755); err != nil {
				gologger.Error().Msgf("创建results目录失败: %v", err)
				return
			}
			
			filePath := "results/" + proto + ".txt"
			
			// 如果文件存在则清空内容，不存在则创建
			file, err := os.Create(filePath)
			if err != nil {
				gologger.Error().Msgf("创建/清空文件 %s 失败: %v", filePath, err)
				return
			}
			file.Close()
			
			// 以追加模式打开文件写入内容
			file, err = os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				gologger.Error().Msgf("打开文件 %s 失败: %v", filePath, err)
				return
			}
			defer file.Close()
			
			if _, err := file.WriteString(lines); err != nil {
				gologger.Error().Msgf("写入文件 %s 失败: %v", filePath, err)
				return
			}
		}(proto, configcontent)
	}

	fileWg.Wait()
	gologger.Info().Msg("所有任务完成！")
}

func processChannel(channel ChannelsType) error {
	channel.URL = collector.ChangeUrlToTelegramWebUrl(channel.URL)
	gologger.Info().Msgf("开始爬取频道: %s", channel.URL)

	resp, err := HttpRequest(channel.URL)
	if err != nil {
		gologger.Error().Msgf("请求失败详情: URL=%s, 错误=%v", channel.URL, err)
		return fmt.Errorf("请求频道 %s 失败: %v", channel.URL, err)
	}
	defer resp.Body.Close()

	gologger.Debug().Msgf("请求成功: URL=%s, 状态码=%d", channel.URL, resp.StatusCode)
	
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		gologger.Error().Msgf("解析HTML失败: URL=%s, 错误=%v", channel.URL, err)
		return fmt.Errorf("解析频道 %s 失败: %v", channel.URL, err)
	}

	CrawlForV2ray(doc, channel.URL, channel.AllMessagesFlag)
	
	gologger.Info().Msgf("成功爬取频道: %s", channel.URL)
	return nil
}

func AddConfigNames(config string, configtype string) string {
	configs := strings.Split(config, "\n")
	newConfigs := ""
	for protoRegex, regexValue := range cfg.RegexPatterns {

		for _, extractedConfig := range configs {

			re := regexp.MustCompile(regexValue)
			matches := re.FindStringSubmatch(extractedConfig)
			if len(matches) > 0 {
				extractedConfig = strings.ReplaceAll(extractedConfig, " ", "")
				if extractedConfig != "" {
					if protoRegex == "vmess" {
						extractedConfig = EditVmessPs(extractedConfig, configtype, true)
						if extractedConfig != "" {
							newConfigs += extractedConfig + "\n"
						}
					} else if protoRegex == "ss" {
						Prefix := strings.Split(matches[0], "ss://")[0]
						if Prefix == "" {
	cfg.ConfigFileIds[configtype] += 1
	newConfigs += extractedConfig + cfg.ConfigsNames + " - " + strconv.Itoa(int(cfg.ConfigFileIds[configtype])) + "\n"
						}
					} else {

	cfg.ConfigFileIds[configtype] += 1
	newConfigs += extractedConfig + cfg.ConfigsNames + " - " + strconv.Itoa(int(cfg.ConfigFileIds[configtype])) + "\n"
					}
				}
			}

		}
	}
	return newConfigs
}

func CrawlForV2ray(doc *goquery.Document, channelLink string, HasAllMessagesFlag bool) {
	// here we are updating our DOM to include the x messages
	// in our DOM and then extract the messages from that DOM
	messages := doc.Find(".tgme_widget_message_wrap").Length()
	link, exist := doc.Find(".tgme_widget_message_wrap .js-widget_message").Last().Attr("data-post")

	if messages < maxMessages && exist {
		number := strings.Split(link, "/")[1]
		doc = GetMessages(maxMessages, doc, number, channelLink)
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
						re := regexp.MustCompile(cfg.RegexPatterns["vmess"])
						matches := re.FindStringSubmatch(extractedConfig)

						if len(matches) > 0 {
							extractedConfig = EditVmessPs(extractedConfig, "mixed", false)
							if line != "" {
								cfg.Configs["mixed"] += extractedConfig + "\n"
							}
						} else {
							cfg.Configs["mixed"] += extractedConfig + "\n"
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
				for protoRegex, regexValue := range cfg.RegexPatterns {

					for _, extractedConfig := range extractedConfigs {

						re := regexp.MustCompile(regexValue)
						matches := re.FindStringSubmatch(extractedConfig)
						if len(matches) > 0 {
							extractedConfig = strings.ReplaceAll(extractedConfig, " ", "")
							if extractedConfig != "" {
								if protoRegex == "vmess" {
									extractedConfig = EditVmessPs(extractedConfig, protoRegex, false)
									if extractedConfig != "" {
										cfg.Configs[protoRegex] += extractedConfig + "\n"
									}
								} else if protoRegex == "ss" {
									Prefix := strings.Split(matches[0], "ss://")[0]
									if Prefix == "" {
										cfg.Configs[protoRegex] += extractedConfig + "\n"
									}
								} else {

									cfg.Configs[protoRegex] += extractedConfig + "\n"
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
	for protoRegex, regexValue := range cfg.RegexPatterns {
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
	cfg.ConfigFileIds[fileName] += 1
	data["ps"] = cfg.ConfigsNames + " - " + strconv.Itoa(int(cfg.ConfigFileIds[fileName])) + "\n"
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

func loadMore(link string) *goquery.Document {
	gologger.Info().Msgf("正在加载更多消息: %s", link)
	
	var doc *goquery.Document
	var err error
	
	// 重试3次
	for i := 0; i < 3; i++ {
		req, err := http.NewRequest("GET", link, nil)
		if err != nil {
			gologger.Error().Msgf("创建请求失败: %v", err)
			continue
		}

		resp, err := client.Do(req)
		if err != nil {
			gologger.Warning().Msgf("请求失败，重试中... (尝试 %d/3)", i+1)
			time.Sleep(2 * time.Second)
			continue
		}
		defer resp.Body.Close()

		doc, err = goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			gologger.Warning().Msgf("解析HTML失败，重试中... (尝试 %d/3)", i+1)
			time.Sleep(2 * time.Second)
			continue
		}

		return doc
	}

	gologger.Error().Msgf("加载更多消息失败: %v", err)
	return nil
}

func HttpRequest(url string) (*http.Response, error) {
	var resp *http.Response
	var err error
	
	// 重试3次
	for i := 0; i < 3; i++ {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			gologger.Error().Msg(fmt.Sprintf("创建请求失败: %s 错误: %s", url, err))
			continue
		}
		
		// 添加浏览器headers
		userAgents := []string{
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Safari/605.1.15",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/115.0",
		}
		rand.Seed(time.Now().UnixNano())
		req.Header.Set("User-Agent", userAgents[rand.Intn(len(userAgents))])
		req.Header.Set("Referer", "https://web.telegram.org/")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
		req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Sec-Fetch-Dest", "document")
		req.Header.Set("Sec-Fetch-Mode", "navigate")
		req.Header.Set("Sec-Fetch-Site", "same-origin")
		req.Header.Set("Sec-Fetch-User", "?1")
		req.Header.Set("Upgrade-Insecure-Requests", "1")
		req.Header.Set("Cache-Control", "max-age=0")
		
		// 添加随机延迟
		time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)

		resp, err = client.Do(req)
		if err != nil {
			gologger.Warning().Msgf("请求失败，重试中... (尝试 %d/3) 错误: %v", i+1, err)
			time.Sleep(2 * time.Second)
			continue
		}

		// 检查响应状态码
		if resp.StatusCode != http.StatusOK {
			gologger.Warning().Msgf("请求 %s 返回非200状态码: %d", url, resp.StatusCode)
			resp.Body.Close()
			continue
		}

		// 检查响应内容长度
		if resp.ContentLength == 0 {
			gologger.Warning().Msgf("请求 %s 返回空响应", url)
			resp.Body.Close()
			continue
		}

		return resp, nil
	}

	gologger.Error().Msgf("请求 %s 失败: %v", url, err)
	return nil, fmt.Errorf("请求失败，请检查网络连接或目标网站是否可访问")
}

func GetMessages(length int, doc *goquery.Document, number string, channel string) *goquery.Document {
	gologger.Info().Msgf("正在获取更多消息，当前数量: %d", length)
	
	x := loadMore(channel + "?before=" + number)
	if x == nil {
		gologger.Error().Msg("加载更多消息失败")
		return doc
	}

	html2, err := x.Html()
	if err != nil {
		gologger.Error().Msgf("获取HTML内容失败: %v", err)
		return doc
	}

	reader2 := strings.NewReader(html2)
	doc2, err := goquery.NewDocumentFromReader(reader2)
	if err != nil {
		gologger.Error().Msgf("解析HTML失败: %v", err)
		return doc
	}

	doc.Find("body").AppendSelection(doc2.Find("body").Children())
	newDoc := goquery.NewDocumentFromNode(doc.Selection.Nodes[0])
	
	messages := newDoc.Find(".js-widget_message_wrap").Length()
	gologger.Info().Msgf("当前总消息数: %d", messages)

	if messages > length {
		return newDoc
	}

	num, err := strconv.Atoi(number)
	if err != nil {
		gologger.Error().Msgf("转换消息编号失败: %v", err)
		return newDoc
	}

	n := num - 21
	if n < 0 {
		return newDoc
	}
	return GetMessages(length, newDoc, strconv.Itoa(n), channel)
}

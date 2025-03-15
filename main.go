package main

import (
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
	"runtime"
	
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
	mu           sync.Mutex // 添加互斥锁
}

type FileLock struct {
	mu    *sync.Mutex
	cond  *sync.Cond
	count int
}

var (
	maxMessages = 100
	client      *http.Client
	cfg         Config
	fileLocks = map[string]*FileLock{
		"ss":        &FileLock{mu: &sync.Mutex{}, count: 0},
		"vmess":     &FileLock{mu: &sync.Mutex{}, count: 0},
		"trojan":    &FileLock{mu: &sync.Mutex{}, count: 0},
		"vless":     &FileLock{mu: &sync.Mutex{}, count: 0},
		"hysteria2": &FileLock{mu: &sync.Mutex{}, count: 0},
		"mixed":     &FileLock{mu: &sync.Mutex{}, count: 0},
	} // 为每种协议类型创建带等待队列的锁
)

func init() {
	// 初始化条件变量
	for _, lock := range fileLocks {
		lock.cond = sync.NewCond(lock.mu)
	}
}

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
		// ConfigsNames: "@Vip_Security join us",
		ConfigsNames: "@Jagger235711",
	Configs: map[string]string{
		"ss":        "",
		"vmess":     "",
		"trojan":    "",
		"vless":     "",
		"hysteria2": "",
		"mixed":     "",
	},
	ConfigFileIds: map[string]int32{
		"ss":        0,
		"vmess":     0,
		"trojan":    0,
		"vless":     0,
		"hysteria2": 0,
		"mixed":     0,
	},
	RegexPatterns: map[string]string{
		"ss":        `(?m)(...ss:|^ss:)\/\/.+?(%3A%40|#)`,
		"vmess":     `(?m)vmess:\/\/[A-Za-z0-9+/=]+`,  // 更严格的vmess匹配
		"trojan":    `(?m)trojan:\/\/.+?(%3A%40|#)`,
		"vless":     `(?m)vless:\/\/.+?(%3A%40|#)`,
		"hysteria2": `(?m)hysteria2:\/\/.+?(%3A%40|#)`,
	},
		Sort: flag.Bool("sort", false, "sort from latest to oldest (default : false)"),
	}
}

type ChannelsType struct {
	URL             string `csv:"URL"`
	AllMessagesFlag bool   `csv:"AllMessagesFlag"`
}

func main() {
	// 添加日志级别控制参数
	logLevel := flag.String("log", "info", "日志级别 (debug, info, warn, error, fatal)")
	flag.Parse()

	// 设置日志级别
	switch strings.ToLower(*logLevel) {
	case "debug":
		gologger.DefaultLogger.SetMaxLevel(levels.LevelDebug)
		gologger.Info().Msg("日志级别设置为：debug")
	case "info":
		gologger.DefaultLogger.SetMaxLevel(levels.LevelInfo)
		gologger.Info().Msg("日志级别设置为：info")
	case "warn":
		gologger.DefaultLogger.SetMaxLevel(levels.LevelWarning)
		gologger.Info().Msg("日志级别设置为：warn")
	case "error":
		gologger.DefaultLogger.SetMaxLevel(levels.LevelError)
		gologger.Info().Msg("日志级别设置为：error")
	case "fatal":
		gologger.DefaultLogger.SetMaxLevel(levels.LevelFatal)
		gologger.Info().Msg("日志级别设置为：fatal")
	default:
		gologger.DefaultLogger.SetMaxLevel(levels.LevelInfo)
		gologger.Info().Msg("使用默认日志级别：info")
	}

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

	// 根据CPU核心数动态设置worker数量
	numWorkers := runtime.NumCPU() * 2
	// 启动worker
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					gologger.Error().Msgf("worker panic: %v", r)
					results <- fmt.Errorf("worker panic: %v", r)
				}
			}()
			
			for channel := range ch {
				if err := processChannel(channel); err != nil {
					gologger.Error().Msgf("处理频道 %s 失败: %v", channel.URL, err)
					results <- err
					continue
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
	// 先写入具体协议文件
	for proto, configcontent := range cfg.Configs {
		if proto == "mixed" {
			continue
		}
		fileWg.Add(1)
		go func(proto, content string) {
			defer fileWg.Done()
			// 添加调试日志
			gologger.Debug().Msgf("准备写入 %s 配置，原始内容长度: %d", proto, len(content))
			
			lines := collector.RemoveDuplicate(content)
			gologger.Debug().Msgf("去重后内容长度: %d", len(lines))
			
			// lines = AddConfigNames(lines, proto)
			// gologger.Debug().Msgf("添加配置名称后内容长度: %d", len(lines))
			
			if *cfg.Sort {
				linesArr := strings.Split(lines, "\n")
				linesArr = collector.Reverse(linesArr)
				lines = strings.Join(linesArr, "\n")
				gologger.Debug().Msgf("排序后内容长度: %d", len(lines))
			}
			lines = strings.TrimSpace(lines)
			gologger.Debug().Msgf("最终内容长度: %d", len(lines))
			
	// 确保results目录存在并检查权限
	if err := os.MkdirAll("results", 0755); err != nil {
		gologger.Fatal().Msgf("创建results目录失败: %v", err)
	}
	
	// 提前创建mixed文件并检查写入权限
	mixedFilePath := "results/mixed.txt"
	if _, err := os.OpenFile(mixedFilePath, os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		gologger.Fatal().Msgf("无法创建mixed.txt文件: %v", err)
	}

			// 使用带等待队列的文件锁
			lock := fileLocks[proto]
			lock.mu.Lock()
			for lock.count > 0 {
				lock.cond.Wait()
			}
			lock.count++
			lock.mu.Unlock()
			
			defer func() {
				lock.mu.Lock()
				lock.count--
				lock.cond.Signal()
				lock.mu.Unlock()
			}()
			
			filePath := "results/" + proto + ".txt"
			if err := collector.WriteToFile(lines, filePath); err != nil {
				gologger.Error().Msgf("写入文件 %s 失败: %v", filePath, err)
				return
			}
			gologger.Info().Msgf("成功写入文件: %s", filePath)
			
			// 将内容追加到mixed.txt
			if proto != "mixed" {
				lock := fileLocks["mixed"]
				lock.mu.Lock()
				for lock.count > 0 {
					lock.cond.Wait()
				}
				lock.count++
				lock.mu.Unlock()
				
				if err := collector.AppendToFile(lines+"\n", "results/mixed.txt"); err != nil {
					gologger.Error().Msgf("追加内容到mixed.txt失败: %v", err)
				} else {
					gologger.Debug().Msgf("成功追加%s配置到mixed.txt", proto)
				}
				
				lock.mu.Lock()
				lock.count--
				lock.cond.Signal()
				lock.mu.Unlock()
			}
		}(proto, configcontent)
	}

	// 等待所有协议文件写入完成
	fileWg.Wait()

	// 确保results目录存在
	if err := os.MkdirAll("results", 0755); err != nil {
		gologger.Fatal().Msgf("创建results目录失败: %v", err)
	}

	// 检查mixed.txt文件状态
	mixedFilePath := "results/mixed.txt"
	if _, err := os.Stat(mixedFilePath); err == nil {
		gologger.Info().Msg("mixed.txt文件已存在")
	} else if os.IsNotExist(err) {
		// 创建空文件
		if err := collector.WriteToFile("", mixedFilePath); err != nil {
			gologger.Fatal().Msgf("创建mixed.txt文件失败: %v", err)
		}
		gologger.Info().Msg("已创建新的mixed.txt文件")
	} else {
		gologger.Fatal().Msgf("检查mixed.txt文件状态失败: %v", err)
	}
	
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

	if err := CrawlForV2ray(doc, channel.URL, channel.AllMessagesFlag); err != nil {
		gologger.Error().Msgf("爬取频道 %s 失败: %v", channel.URL, err)
		return fmt.Errorf("爬取频道 %s 失败: %v", channel.URL, err)
	}
	
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
						extractedConfig = collector.EditVmessPs(extractedConfig)
						if extractedConfig != "" {
							newConfigs += extractedConfig + "\n"
						}
					} else if protoRegex == "ss" {
						Prefix := strings.Split(matches[0], "ss://")[0]
						if Prefix != "vle" {
			cfg.mu.Lock()
			cfg.ConfigFileIds[configtype] += 1
			id := cfg.ConfigFileIds[configtype]
			cfg.mu.Unlock()
			
			cfg.mu.Lock()
			newConfigs += extractedConfig + cfg.ConfigsNames + " - " + strconv.Itoa(int(id)) + "\n"
			cfg.mu.Unlock()
						}
					} else if protoRegex == "hysteria2" {
						cfg.mu.Lock()
						cfg.ConfigFileIds[configtype] += 1
						id := cfg.ConfigFileIds[configtype]
						cfg.mu.Unlock()
						newConfigs += extractedConfig + cfg.ConfigsNames + " - " + strconv.Itoa(int(id)) + "\n"
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

func CrawlForV2ray(doc *goquery.Document, channelLink string, HasAllMessagesFlag bool) error {
	// 初始化配置计数
	initialCounts := make(map[string]int)
	for proto := range cfg.Configs {
		initialCounts[proto] = len(cfg.Configs[proto])
	}

	// 更新DOM以包含更多消息
	messages := doc.Find(".tgme_widget_message_wrap").Length()
	link, exist := doc.Find(".tgme_widget_message_wrap .js-widget_message").Last().Attr("data-post")

	if messages < maxMessages && exist {
		number := strings.Split(link, "/")[1]
		doc = collector.GetMessages(maxMessages, doc, number, channelLink)
	}

	if HasAllMessagesFlag {
		// 获取所有消息并检查v2ray配置
		doc.Find(".tgme_widget_message_text").Each(func(j int, s *goquery.Selection) {
			messageText, _ := s.Html()
			str := strings.ReplaceAll(messageText, "<br/>", "\n")
			doc, _ := goquery.NewDocumentFromReader(strings.NewReader(str))
			messageText = doc.Text()
			line := strings.TrimSpace(messageText)
			gologger.Debug().Msgf("原始消息内容: %s", line)
			
			lines := strings.Split(line, "\n")
			for _, data := range lines {
				data = strings.TrimSpace(data)
				if data == "" {
					continue
				}
				
				extractedConfigs := strings.Split(ExtractConfig(data, []string{}), "\n")
				gologger.Debug().Msgf("提取的配置: %v", extractedConfigs)
				
				for _, extractedConfig := range extractedConfigs {
					extractedConfig = strings.TrimSpace(extractedConfig)
					if extractedConfig == "" {
						continue
					}

					// 检查每种协议类型
					for proto, pattern := range cfg.RegexPatterns {
						re := regexp.MustCompile(pattern)
						if re.MatchString(extractedConfig) {
							gologger.Debug().Msgf("匹配到 %s 协议: %s", proto, extractedConfig)
							
							if proto == "vmess" {
								extractedConfig = collector.EditVmessPs(extractedConfig)
							}
							
							if extractedConfig != "" {
								cfg.Configs[proto] += extractedConfig + "\n"
								// 直接写入mixed文件
				lock := fileLocks["mixed"]
				lock.mu.Lock()
				for lock.count > 0 {
					lock.cond.Wait()
				}
				lock.count++
				lock.mu.Unlock()
				
				if err := collector.AppendToFile(extractedConfig+"\n", "results/mixed.txt"); err != nil {
					gologger.Error().Msgf("写入mixed.txt失败: %v", err)
				}
				
				lock.mu.Lock()
				lock.count--
				lock.cond.Signal()
				lock.mu.Unlock()
							}
							break
						}
					}
				}
			}
		})
	} else {
		// 优化后的配置提取逻辑
		doc.Find(".tgme_widget_message_text").Each(func(j int, s *goquery.Selection) {
			// 提取所有文本内容
			messageText := s.Text()
			lines := strings.Split(messageText, "\n")
			
			// 定义常见配置前缀
			configPrefixes := []string{
				"vmess://",
				"vless://", 
				"trojan://",
				"ss://",
				"hysteria2://",
			}
			
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" {
					continue
				}
				
				// 检查是否包含配置前缀
				hasConfig := false
				for _, prefix := range configPrefixes {
					if strings.Contains(line, prefix) {
						hasConfig = true
						break
					}
				}
				
				if !hasConfig {
					continue
				}
				
				// 提取配置
				extractedConfigs := strings.Split(ExtractConfig(line, []string{}), "\n")
				gologger.Debug().Msgf("提取的配置: %v", extractedConfigs)
				
				for _, extractedConfig := range extractedConfigs {
					extractedConfig = strings.TrimSpace(extractedConfig)
					if extractedConfig == "" {
						continue
					}

					// 检查每种协议类型
					for proto, pattern := range cfg.RegexPatterns {
						re := regexp.MustCompile(pattern)
						if re.MatchString(extractedConfig) {
							gologger.Debug().Msgf("匹配到 %s 协议: %s", proto, extractedConfig)
							
							if proto == "vmess" {
								extractedConfig = collector.EditVmessPs(extractedConfig)
							}
							
							if extractedConfig != "" {
								cfg.Configs[proto] += extractedConfig + "\n"
								// 确保mixed包含所有协议
								if proto != "mixed" {
									cfg.Configs["mixed"] += extractedConfig + "\n"
								}
							}
							break
						}
					}
				}
			}
		})
	}

	// 检查是否有新配置被添加
	finalCounts := make(map[string]int)
	for proto := range cfg.Configs {
		finalCounts[proto] = len(cfg.Configs[proto])
	}
	
	hasNewConfig := false
	for proto := range cfg.Configs {
		if finalCounts[proto] > initialCounts[proto] {
			hasNewConfig = true
			break
		}
	}
	
	if !hasNewConfig {
		return fmt.Errorf("没有找到新的配置")
	}
	
	return nil
}

func ExtractConfig(Txt string, Tempconfigs []string) string {
	for protoRegex, regexValue := range cfg.RegexPatterns {
		re := regexp.MustCompile(regexValue)
		matches := re.FindStringSubmatch(Txt)
		extractedConfig := ""
		if len(matches) > 0 {
			if protoRegex == "ss" {
				Prefix := strings.Split(matches[0], "ss://")[0]
				if Prefix == "" {
					extractedConfig = "\n" + matches[0]
				} else if Prefix != "vle" {
					d := strings.Split(matches[0], "ss://")
					extractedConfig = "\n" + "ss://" + d[1]
				}
			} else {
				extractedConfig = "\n" + matches[0]
			}

			Tempconfigs = append(Tempconfigs, extractedConfig)
			Txt = strings.ReplaceAll(Txt, matches[0], "")
			ExtractConfig(Txt, Tempconfigs)
		}
	}
	return strings.Join(Tempconfigs, "\n")
}

func loadMore(link string) *goquery.Document {
	gologger.Info().Msgf("正在加载更多消息: %s", link)
	
	var doc *goquery.Document
	var err error
	
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
	return doc
}

func HttpRequest(url string) (*http.Response, error) {
	for i := 0; i < 3; i++ {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			gologger.Error().Msg(fmt.Sprintf("创建请求失败: %s 错误: %s", url, err))
			continue
		}
		
		userAgents := []string{
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
		}

		req.Header.Set("User-Agent", userAgents[rand.Intn(len(userAgents))])
		resp, err := client.Do(req)
		if err == nil {
			return resp, nil
		}
		gologger.Warning().Msgf("请求失败，重试中... (尝试 %d/3)", i+1)
		time.Sleep(2 * time.Second)
	}
	return nil, fmt.Errorf("请求失败: %s", url)
}

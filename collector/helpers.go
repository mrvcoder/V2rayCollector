package collector

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func ChangeUrlToTelegramWebUrl(input string) string {
	if !strings.Contains(input, "/s/") {
		index := strings.Index(input, "/t.me/")
		if index != -1 {
			modifiedURL := input[:index+len("/t.me/")] + "s/" + input[index+len("/t.me/"):]
			return modifiedURL
		}
	}
	return input
}

func ReadFileContent(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func Reverse(lines []string) []string {
	for i := 0; i < len(lines)/2; i++ {
		j := len(lines) - i - 1
		lines[i], lines[j] = lines[j], lines[i]
	}
	return lines
}

func RemoveDuplicate(config string) string {
	lines := strings.Split(config, "\n")
	uniqueLines := make(map[string]bool)
	result := []string{}
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// 只对配置部分去重，忽略备注
		configPart := strings.Split(line, "#")[0]
		configPart = strings.TrimSpace(configPart)
		
		// 验证配置格式
		if !isValidConfig(configPart) {
			continue
		}
		
		if !uniqueLines[configPart] {
			uniqueLines[configPart] = true
			result = append(result, line)
		}
	}
	
	return strings.Join(result, "\n")
}

func isValidConfig(config string) bool {
	// 检查常见协议格式
	protocols := []string{"vmess://", "vless://", "ss://", "trojan://", "hysteria2://"}
	for _, protocol := range protocols {
		if strings.HasPrefix(config, protocol) {
			return true
		}
	}
	return false
}

func GetMessages(maxMessages int, doc *goquery.Document, number string, channelLink string) *goquery.Document {
	loadMoreURL := fmt.Sprintf("%s?before=%s", channelLink, number)
	newDoc := loadMore(loadMoreURL)
	if newDoc == nil {
		return doc
	}
	
	newMessages := newDoc.Find(".tgme_widget_message_wrap")
	doc.Find("body").AppendSelection(newMessages)
	
	if doc.Find(".tgme_widget_message_wrap").Length() < maxMessages {
		newNumber := newDoc.Find(".tgme_widget_message_wrap .js-widget_message").Last().AttrOr("data-post", "")
		if newNumber != "" {
			newNumber = strings.Split(newNumber, "/")[1]
			return GetMessages(maxMessages, doc, newNumber, channelLink)
		}
	}
	
	return doc
}

func WriteToFile(content string, filePath string) error {
	// 确保目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 按协议类型分类
	protocols := map[string][]string{
		"ss":       []string{},
		"vmess":    []string{},
		"vless":    []string{},
		"trojan":   []string{},
		"hysteria2": []string{},
	}

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !isValidConfig(line) {
			continue
		}

		// 根据协议类型分类
		switch {
		case strings.HasPrefix(line, "ss://"):
			protocols["ss"] = append(protocols["ss"], line)
		case strings.HasPrefix(line, "vmess://"):
			line = EditVmessPs(line)
			protocols["vmess"] = append(protocols["vmess"], line)
		case strings.HasPrefix(line, "vless://"):
			protocols["vless"] = append(protocols["vless"], line)
		case strings.HasPrefix(line, "trojan://"):
			protocols["trojan"] = append(protocols["trojan"], line)
		case strings.HasPrefix(line, "hysteria2://"):
			protocols["hysteria2"] = append(protocols["hysteria2"], line)
		}
	}

	// 写入对应文件
	for protocol, configs := range protocols {
		if len(configs) > 0 {
			outputPath := filepath.Join(dir, protocol+".txt")
			content := strings.Join(configs, "\n")
			if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
				return fmt.Errorf("写入%s文件失败: %v", protocol, err)
			}
			fmt.Printf("[INFO] 成功写入 %d 条%s配置到 %s\n", len(configs), protocol, outputPath)
		}
	}

	return nil
}

func AppendToFile(content string, filePath string) error {
	// 确保目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 以追加模式打开文件
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	// 写入内容
	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return nil
}

func EditVmessPs(vmess string) string {
	// 解析vmess配置
	config := strings.TrimPrefix(vmess, "vmess://")
	decoded, err := base64.StdEncoding.DecodeString(config)
	if err != nil {
		return vmess
	}

	// 解析JSON
	var vmessConfig map[string]interface{}
	if err := json.Unmarshal(decoded, &vmessConfig); err != nil {
		return vmess
	}

	// 保留原始PS值
	if ps, ok := vmessConfig["ps"].(string); ok {
		// 格式化配置名称
		vmessConfig["ps"] = fmt.Sprintf("%s | V2rayCollector", ps)
	}

	// 重新编码为JSON
	encoded, err := json.Marshal(vmessConfig)
	if err != nil {
		return vmess
	}

	// 返回新的vmess配置
	return fmt.Sprintf("vmess://%s", base64.StdEncoding.EncodeToString(encoded))
}

func loadMore(link string) *goquery.Document {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	for i := 0; i < 3; i++ {
		req, err := http.NewRequest("GET", link, nil)
		if err != nil {
			fmt.Printf("[WARN] 创建请求失败: %v (重试 %d/3)\n", err, i+1)
			continue
		}

		// 设置请求头
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
		req.Header.Set("Referer", "https://t.me/")
		req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("[WARN] 请求失败: %v (重试 %d/3)\n", err, i+1)
			time.Sleep(2 * time.Second)
			continue
		}
		defer resp.Body.Close()

		// 检查响应状态码
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("[WARN] 请求返回非200状态码: %d (重试 %d/3)\n", resp.StatusCode, i+1)
			time.Sleep(2 * time.Second)
			continue
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			fmt.Printf("[WARN] 解析HTML失败: %v (重试 %d/3)\n", err, i+1)
			time.Sleep(2 * time.Second)
			continue
		}

		fmt.Printf("[INFO] 成功加载更多消息 from %s\n", link)
		return doc
	}

	fmt.Printf("[ERROR] 加载更多消息失败 after 3次重试: %s\n", link)
	return nil
}

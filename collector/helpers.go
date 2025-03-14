package collector

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"slices"
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
	slices.Sort(lines)
	lines = slices.Compact(lines)
	return strings.Join(lines, "\n")
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

	// 使用临时文件实现原子写入
	tmpFile := filePath + ".tmp"
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("写入临时文件失败: %v", err)
	}

	// 重命名实现原子操作
	if err := os.Rename(tmpFile, filePath); err != nil {
		// 清理临时文件
		os.Remove(tmpFile)
		return fmt.Errorf("重命名文件失败: %v", err)
	}

	return nil
}

func loadMore(link string) *goquery.Document {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	for i := 0; i < 3; i++ {
		req, err := http.NewRequest("GET", link, nil)
		if err != nil {
			continue
		}

		resp, err := client.Do(req)
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}
		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}

		return doc
	}

	return nil
}

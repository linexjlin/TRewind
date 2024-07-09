package main

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
	"unicode"
)

func md5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

var baseDataRoot = "imported"

func saveToFile(category, id, position, input string) error {
	// 获取当前日期目录
	dateDir := fmt.Sprintf("%s/%s/%s", baseDataRoot, time.Now().Format("2006-01-02"), category)
	if _, err := os.Stat(dateDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dateDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
	}

	// 去除Windows不合法字符并取前20个rune字符
	var peekChars []rune
	for _, r := range input {
		if unicode.IsGraphic(r) && !isWindowsIllegalChar(r) && r != '\n' {
			peekChars = append(peekChars, r)
			if len(peekChars) == 20 {
				break
			}
		}
	}

	// 生成文件名
	fileName := fmt.Sprintf("%s/%s-%s-%s···.txt", dateDir, id, position, string(peekChars))

	// 创建文件并写入内容
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(input)
	return err
}

// 检查是否为Windows不合法字符
func isWindowsIllegalChar(r rune) bool {
	switch r {
	case '<', '>', ':', '"', '/', '\\', '|', '?', '*':
		return true
	default:
		return false
	}
}

func extractFilenameAndExtra(input string) (name, extra string) {
	for i, r := range input {
		if r == '\n' && len(name) != 0 {
			// left input is left
			extra = input[i+1:]
			break
		}
		name += string(r)
	}
	return name, extra
}

func checkSSL(address string) (bool, error) {
	conn, err := tls.DialWithDialer(&net.Dialer{
		Timeout: 2 * time.Second,
	}, "tcp", address, &tls.Config{
		InsecureSkipVerify: true,
	})

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "timeout") {
			return false, err
		} else {
			return false, nil
		}
	}
	defer conn.Close()
	return true, nil
}

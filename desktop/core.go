package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Core struct {
	apiAddress string
}

func NewCore(apiAddr string) *Core {
	c := Core{apiAddress: apiAddr}
	return &c
}

func (c *Core) uploadClipboard(collection, clipBoardContent, clipBoardType, importType string) {
	api := fmt.Sprintf("%s%s/%s/upload_clipboard", os.Getenv("API_SCHEME"), c.apiAddress, collection)
	// 创建一个 ClipBoard 结构体实例
	cb := struct {
		ClipBoardContent string `json:"clipBoardContent"`
		ClipBoardType    string `json:"clipBoardType"`
		ImportType       string `json:"importType"`
	}{
		ClipBoardContent: clipBoardContent,
		ClipBoardType:    clipBoardType,
		ImportType:       importType,
	}

	// 将结构体序列化为 JSON
	jsonData, err := json.Marshal(cb)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		return
	}

	// 创建一个 HTTP 客户端
	client := &http.Client{}

	// 创建一个 POST 请求
	req, err := http.NewRequest("POST", api, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error creating request:", err)
		return
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()
	// 打印响应状态码
	log.Println("Response status:", resp.Status)
}

func (c *Core) importToDocname(collection, text string) {
	id := md5Hash(text)
	saveToFile(collection, id, "name", text)
	c.uploadClipboard(collection, text, "Text", "DocName")
}

// only the fist line will be indexed, extra info will not be indexed
func (c *Core) importToExtra(collection, input string) {
	embText, _ := extractFilenameAndExtra(input)
	id := md5Hash(embText)
	saveToFile(collection, id, "extra", input)
	c.uploadClipboard(collection, input, "Text", "Extra")
}

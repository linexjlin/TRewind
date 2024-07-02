package main

import (
	"log"
	"time"

	"github.com/linexjlin/TRewind/chromaManager"
)

type Core struct {
	db *chromaManager.ChromaManager
}

func NewCore(db *chromaManager.ChromaManager) *Core {
	c := Core{db: db}
	c.init()
	return &c
}

func (c *Core) importClipboardText(collection, text string) {
	id := md5Hash(text)
	metadata := map[string]string{
		"update":   time.Now().Format("20060102150405"),
		"filename": text,
	}
	c.db.UpsertDoc(collection, text, id, metadata)
	log.Println("imported", text)
}

func (c *Core) init() {
	//
}

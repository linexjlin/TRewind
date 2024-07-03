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

func (c *Core) importToDocname(collection, text string) {
	id := md5Hash(text)
	metadata := map[string]string{
		"update":   time.Now().Format("20060102150405"),
		"filename": text,
	}
	saveToFile(collection, id, "name", text)
	c.db.UpsertDoc(collection, text, id, metadata)
	log.Println("imported", text)
}

// only the fist line will be indexed, extra info will not be indexed
func (c *Core) importToExtra(collection, input string) {
	embText, extra := extractFilenameAndExtra(input)
	id := md5Hash(embText)
	metadata := map[string]string{
		"update":   time.Now().Format("20060102150405"),
		"filename": embText,
		"extra":    extra,
	}
	saveToFile(collection, id, "extra", input)
	c.db.UpsertDoc(collection, embText, id, metadata)
	log.Println("imported as extra", input)
}

func (c *Core) delDocByName(collection, name string) {
	id := md5Hash(name)
	c.db.DeleteByID(collection, id)
	log.Println("del file", name)
}

func (c *Core) delDocByID(collection, id string) {
	c.db.DeleteByID(collection, id)
	log.Println("del ID", id)
}

func (c *Core) init() {
	//
}

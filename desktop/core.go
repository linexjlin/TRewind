package main

import (
	"log"

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

func (c *Core) importText(text string) {
	log.Println("import", text)
	//c.db.UpsertDoc()
}

func (c *Core) init() {
	//
}

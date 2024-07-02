package main

import "log"

type Core struct {
	//u           *UserCore
	//st *SysTray
}

func NewCore() *Core {
	c := Core{}
	c.init()
	return &c
}

func (c *Core) importText(text string) {
	log.Println("import", text)
}

func (c *Core) init() {
	//
}

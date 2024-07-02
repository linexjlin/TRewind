package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/linexjlin/TRewind/chromaManager"
)

var UText func(string) string
var UMenuText func(string) string

func initUText() {
	UMenuText = func(s string) string {
		return s
	}
}

func initUMenuText() {
	UMenuText = func(s string) string {
		return s
	}
}

var Version = "dev"

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	initUText()
	initUMenuText()
	docDB, err := chromaManager.NewChromaManager()
	if err != nil {
		panic(err)
	}
	core := NewCore(docDB)
	tray := NewSysTray(core)

	tray.Run()
}

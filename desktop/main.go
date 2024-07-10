package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/linexjlin/TRewind/apiServer"
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
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if err := godotenv.Load("env.txt"); err != nil {
		log.Println(err)
		log.Error("Error loading .env file")
		return
	}

	initUText()
	initUMenuText()

	backendAPI := ""

	if remoteAPIAddr := os.Getenv("REMOTE_API_ADDR"); remoteAPIAddr != "" {
		backendAPI = remoteAPIAddr
	}

	if listenAddr := os.Getenv("API_LISTEN_ADDR"); listenAddr != "" {
		docDB, err := chromaManager.NewChromaManager()
		if err != nil {
			panic(err)
		}
		apiServer := apiServer.NewServer(docDB)
		backendAPI = listenAddr

		go func() {
			log.Fatal(apiServer.ListenAndServe(listenAddr))
		}()
	}

	log.Println("using backend API:", backendAPI)

	if isSSL, _ := checkSSL(backendAPI); isSSL {
		os.Setenv("API_SCHEME", "https://")
	} else {
		os.Setenv("API_SCHEME", "http://")
	}

	core := NewCore(backendAPI)
	tray := NewSysTray(core)

	hideWinConsole()

	tray.Run()
}

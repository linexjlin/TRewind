package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/linexjlin/TRewind/apiServer"
	"github.com/linexjlin/TRewind/chromaManager"
)

var Version = "dev"

func main() {
	versionFlag := flag.Bool("version", false, "Print the version and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("Version: %s\n", Version)
		return
	}

	if err := godotenv.Load("env.txt"); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize DocDB (assuming ChromaManager has a constructor)
	docDB, _ := chromaManager.NewChromaManager()

	apiServer := apiServer.NewServer(docDB)

	listenAddr := os.Getenv("API_LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = "127.0.0.1:8601"
	}

	log.Fatal(apiServer.ListenAndServe(listenAddr))
}

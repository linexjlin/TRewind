package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/getlantern/systray"
	icon "github.com/linexjlin/systray-icons/tape"
	"github.com/skratchdot/open-golang/open"
	"golang.design/x/clipboard"
)

func NewSysTray(c *Core) *SysTray {
	defaultCollection := os.Getenv("DEFAULT_COLLECTION")
	if defaultCollection == "" {
		defaultCollection = "docs"
	}

	var collections = []string{}
	if collectionsStr := os.Getenv("COLLECTIONS"); collectionsStr != "" {
		collections = strings.Split(collectionsStr, ";")
	} else {
		collections = []string{}
	}

	serverAddr := os.Getenv("API_LISTEN_ADDR")
	if serverAddr == "" {
		serverAddr = "127.0.0.1:8601"
	}

	tray := SysTray{core: c, defaultCollection: defaultCollection, collections: collections, serverAddr: serverAddr}
	return &tray
}

type SysTray struct {
	core              *Core
	defaultCollection string // Changed to lowercase to make it unexported
	collections       []string
	serverAddr        string
}

func (st *SysTray) Run() {
	systray.Run(st.onReady, st.onExit)
}

func (st *SysTray) onExit() {
	fmt.Println("exit")
}

func (st *SysTray) onReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)

	mQuitOrig := systray.AddMenuItem(UMenuText("Exit"), UMenuText("Quit the whole app"))
	go func() {
		<-mQuitOrig.ClickedCh
		systray.Quit()
		fmt.Println("Finished quitting")
	}()

	mAbout := systray.AddMenuItem(UMenuText("About")+fmt.Sprintf(" ( %s )", Version), UMenuText("Open the project page"))
	go func() {
		for {
			<-mAbout.ClickedCh
			open.Start("https://github.com/linexjlin/TRewind")
		}
	}()

	var allCollections = []string{st.defaultCollection}
	allCollections = append(allCollections, st.collections...)
	for _, collection := range allCollections {
		systray.AddSeparator()
		systray.AddSeparator()
		mSearch := systray.AddMenuItem(fmt.Sprintf("[%s] %s", collection, UMenuText("Search")), UMenuText("Search"))
		go func() {
			for {
				<-mSearch.ClickedCh
				open.Start(fmt.Sprintf("http://%s/%s/recall/ui/", st.serverAddr, collection))
			}
		}()

		mAdd := systray.AddMenuItem(fmt.Sprintf("[%s] %s", collection, UMenuText("Add All")), UMenuText("Add new doc from clipboard"))
		go func() {
			for {
				<-mAdd.ClickedCh
				if err := clipboard.Init(); err != nil {
					log.Println(err)
				}
				clipboardText := string(clipboard.Read(clipboard.FmtText))
				log.Println("Got clipboardText", clipboardText)
				if len(clipboardText) > 0 {
					st.core.importToDocname(collection, clipboardText)
				}
			}
		}()

		mAddExtra := systray.AddMenuItem(fmt.Sprintf("[%s] %s", collection, UMenuText("Add First Line Only")), UMenuText("Add first line only to Extra from clipboard"))
		go func() {
			for {
				<-mAddExtra.ClickedCh
				if err := clipboard.Init(); err != nil {
					log.Println(err)
				}
				clipboardText := string(clipboard.Read(clipboard.FmtText))
				log.Println("Got clipboardText", clipboardText)
				if len(clipboardText) > 0 {
					st.core.importToExtra(collection, clipboardText)
				}
			}
		}()

		/*mDelByName := systray.AddMenuItem(fmt.Sprintf("[%s] %s", collection, UMenuText("Del by Name")), UMenuText("Del doc by Name from clipboard"))
		go func() {
			for {
				<-mDelByName.ClickedCh
				if err := clipboard.Init(); err != nil {
					log.Println(err)
				}
				clipboardText := string(clipboard.Read(clipboard.FmtText))
				log.Println("Got clipboardText", clipboardText)
				if len(clipboardText) > 0 {
					st.core.delDocByName(collection, clipboardText)
				}
			}
		}()

		mDelByID := systray.AddMenuItem(fmt.Sprintf("[%s] %s", collection, UMenuText("Del by ID")), UMenuText("Del doc by ID from clipboard"))
		go func() {
			for {
				<-mDelByID.ClickedCh
				if err := clipboard.Init(); err != nil {
					log.Println(err)
				}
				clipboardText := string(clipboard.Read(clipboard.FmtText))
				log.Println("Got clipboardText", clipboardText)
				if len(clipboardText) > 0 {
					st.core.delDocByID(collection, clipboardText)
				}
			}
		}()*/
	}
}

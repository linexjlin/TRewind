package main

import (
	"fmt"
	"log"

	"github.com/getlantern/systray"
	icon "github.com/linexjlin/systray-icons/tape"
	"github.com/skratchdot/open-golang/open"
	"golang.design/x/clipboard"
)

func NewSysTray(c *Core) *SysTray {
	tray := SysTray{core: c}
	return &tray
}

type SysTray struct {
	core              *Core
	updateHotKeyTitle func(string)
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
		fmt.Println("Requesting quit")
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

	systray.AddSeparator()

	mSearch := systray.AddMenuItem(UMenuText("Search"), UMenuText("Search"))
	go func() {
		for {
			<-mSearch.ClickedCh
			open.Start("http://127.0.0.1:8023/search")
		}
	}()

	mAdd := systray.AddMenuItem(UMenuText("Add"), UMenuText("Open the project page"))
	go func() {
		for {
			<-mAdd.ClickedCh
			if err := clipboard.Init(); err != nil {
				log.Println(err)
			}
			clipboardText := string(clipboard.Read(clipboard.FmtText))
			log.Println("Got clipboardText", clipboardText)
			if len(clipboardText) > 0 {
				st.core.importText(clipboardText)
			}
		}
	}()
}

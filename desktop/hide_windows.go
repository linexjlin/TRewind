package main

import "github.com/lxn/win"

func hideWinConsole() {
	win.ShowWindow(win.GetConsoleWindow(), win.SW_HIDE)
}

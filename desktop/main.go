package main

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
	initUText()
	initUMenuText()
	core := NewCore()
	tray := NewSysTray(core)

	tray.Run()
}

package main

import (
	"fwallet/wallet"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"log"
)

func main() {
	a := app.NewWithID("github.com/axengine/fyne_wallet")
	a.SetIcon(theme.FyneLogo())
	makeTray(a)
	logLifecycle(a)
	w := a.NewWindow("Fyne Wallet For Ethereum like")

	wl := wallet.NewFyneWallet(a, w)

	w.SetMainMenu(wl.MainMenu(a, w))
	w.SetContent(wl.Content(a, w))
	w.SetMaster()

	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}

func logLifecycle(a fyne.App) {
	a.Lifecycle().SetOnStarted(func() {
		log.Println("Lifecycle: Started")
	})
	a.Lifecycle().SetOnStopped(func() {
		log.Println("Lifecycle: Stopped")
	})
	a.Lifecycle().SetOnEnteredForeground(func() {
		log.Println("Lifecycle: Entered Foreground")
	})
	a.Lifecycle().SetOnExitedForeground(func() {
		log.Println("Lifecycle: Exited Foreground")
	})
}

func makeTray(a fyne.App) {
	if desk, ok := a.(desktop.App); ok {
		h := fyne.NewMenuItem("Hello", func() {})
		menu := fyne.NewMenu("Hello World", h)
		h.Action = func() {
			log.Println("System tray menu tapped")
			h.Label = "Welcome"
			menu.Refresh()
		}
		desk.SetSystemTrayMenu(menu)
	}
}

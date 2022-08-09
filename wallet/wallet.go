package wallet

import (
	"fmt"
	"fwallet/dao"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type function struct {
	name string
	fn   func(fyne.App, fyne.Window) fyne.CanvasObject
}

type FyneWallet struct {
	app    fyne.App
	window fyne.Window

	content     *fyne.Container
	navIndex    map[string][]string
	navHandlers map[string]function

	d *dao.Dao
}

func NewFyneWallet(a fyne.App, w fyne.Window) *FyneWallet {
	wl := &FyneWallet{
		app:         a,
		window:      w,
		content:     container.NewMax(),
		navIndex:    make(map[string][]string),
		navHandlers: make(map[string]function),
		d:           dao.New("fynewallet.db"),
	}
	wl.init()
	wl.handleMain(a, w)
	return wl
}

func (wl *FyneWallet) init() {
	wl.navIndex = make(map[string][]string)
	wl.navIndex[""] = []string{"main", "networks", "accounts", "assets"}

	wl.navHandlers[""] = function{
		name: "main",
		fn:   wl.handleMain,
	}
	wl.navHandlers["main"] = function{
		name: "main",
		fn:   wl.handleMain,
	}
	wl.navHandlers["networks"] = function{
		name: "networks",
		fn:   wl.handleNetworks,
	}
	wl.navHandlers["accounts"] = function{
		name: "accounts",
		fn:   wl.handleAccounts,
	}
	wl.navHandlers["assets"] = function{
		name: "assets",
		fn:   wl.handleAssets,
	}
}

func (wl *FyneWallet) handleMain(a fyne.App, w fyne.Window) fyne.CanvasObject {
	label := container.NewCenter(widget.NewLabel("main"))
	wl.content.Objects = []fyne.CanvasObject{label}
	wl.content.Refresh()
	return nil
}

func (wl *FyneWallet) handleNetworks(a fyne.App, w fyne.Window) fyne.CanvasObject {
	_, items, err := wl.d.ListNetwork(1, 10)
	if err != nil {
		label := container.NewCenter(widget.NewLabel(err.Error()))
		wl.content.Objects = []fyne.CanvasObject{label}
		wl.content.Refresh()
		return nil
	}

	tips := widget.NewLabel("Manage Network")

	//name := widget.NewEntry()
	//name.SetPlaceHolder("blockchain network name")
	//
	//chainId := widget.NewLabel("0")
	//rpc := widget.NewEntry()
	//rpc.SetPlaceHolder("rpc url")
	//
	//symbol := widget.NewEntry()
	//symbol.SetPlaceHolder("asset symbol")
	//
	//explorer := widget.NewEntry()
	//explorer.SetPlaceHolder("explorer url")
	//
	//form := &widget.Form{
	//	Items: []*widget.FormItem{
	//		{Text: "Name", Widget: name, HintText: "Your full name"},
	//		{Text: "ChainId", Widget: chainId, HintText: "A valid email address"},
	//		{Text: "Rpc", Widget: rpc, HintText: "A valid email address"},
	//		{Text: "Symbol", Widget: symbol, HintText: "A valid email address"},
	//		{Text: "Explorer", Widget: explorer, HintText: "A valid email address"},
	//	},
	//	OnCancel: func() {
	//		fmt.Println("Cancelled")
	//	},
	//	OnSubmit: func() {
	//		fmt.Println("Form submitted")
	//		fyne.CurrentApp().SendNotification(&fyne.Notification{
	//			Title: "Form for: " + name.Text,
	//			//Content: largeText.Text,
	//			Content: "xxxxxxxxxxxxxxx",
	//		})
	//	},
	//}
	form := &widget.Form{
		OnCancel: func() {
			fmt.Println("Cancelled")
		},
		OnSubmit: func() {
		},
	}
	vbox := container.NewVBox(tips, form)

	list := widget.NewList(
		func() int {
			return len(items)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(items[id].Name)
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		tips.SetText(items[id].Name)

		name := widget.NewEntry()
		name.SetText(items[id].Name)
		rpc := widget.NewEntry()
		rpc.SetText(items[id].Rpc)
		chainId := widget.NewLabel(fmt.Sprintf("%d", items[id].ChainId))
		symbol := widget.NewEntry()
		symbol.SetText(items[id].Symbol)
		explorer := widget.NewEntry()
		explorer.SetText(items[id].Explorer)

		form.Append("name", name)
		form.Append("rpc", rpc)
		form.Append("chainId", chainId)
		form.Append("symbol", symbol)
		form.Append("explorer", explorer)
	}
	list.OnUnselected = func(id widget.ListItemID) {
		tips.SetText("Select An Item From The List")
	}

	c := container.NewHSplit(list, container.NewMax(vbox))
	c.Offset = 0.2
	wl.content.Objects = []fyne.CanvasObject{c}
	wl.content.Refresh()
	return nil

}

func (wl *FyneWallet) handleAccounts(a fyne.App, w fyne.Window) fyne.CanvasObject {
	label := container.NewCenter(widget.NewLabel("accounts manage"))
	wl.content.Objects = []fyne.CanvasObject{label}
	wl.content.Refresh()
	return nil
}

func (wl *FyneWallet) handleAssets(a fyne.App, w fyne.Window) fyne.CanvasObject {
	label := container.NewCenter(widget.NewLabel("assets manage"))
	wl.content.Objects = []fyne.CanvasObject{label}
	wl.content.Refresh()
	return nil
}

func (wl *FyneWallet) MainMenu(a fyne.App, w fyne.Window) *fyne.MainMenu {
	return nil
}

func (wl *FyneWallet) Content(a fyne.App, w fyne.Window) fyne.CanvasObject {
	c := container.NewHSplit(wl.nav(a, w), wl.content)
	c.Offset = 0.2
	return c
}

// 导航栏
func (wl *FyneWallet) nav(a fyne.App, w fyne.Window) fyne.CanvasObject {
	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return wl.navIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := wl.navIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Default name")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			if t, ok := wl.navHandlers[uid]; ok {
				obj.(*widget.Label).SetText(t.name)
			}
		},
		OnSelected: func(uid string) {
			if t, ok := wl.navHandlers[uid]; ok {
				t.fn(wl.app, wl.window)
			}
		},
	}
	return container.NewBorder(nil, nil, nil, nil, tree)
}

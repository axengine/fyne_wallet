package wallet

import (
	"fmt"
	"fwallet/dao"
	"fwallet/model"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"strconv"
	"strings"
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
	var (
		page, size = 1, 10
	)
	_, items, err := wl.d.ListNetwork(page, size)
	if err != nil {
		wl.showErr(err)
		return nil
	}

	tips := widget.NewLabel("Manage Network")

	networkId := widget.NewEntry()
	networkId.Disable()
	name := widget.NewEntry()
	rpc := widget.NewEntry()
	chainId := widget.NewEntry()
	symbol := widget.NewEntry()
	explorer := widget.NewEntry()

	form := &widget.Form{}
	form.Append("Id", networkId)
	form.Append("Name", name)
	form.Append("RPC", rpc)
	form.Append("ChainId", chainId)
	form.Append("Symbol", symbol)
	form.Append("Explorer", explorer)

	delButton := widget.NewButton("DEL", func() {
		d := dialog.NewConfirm("Sure", "deleted", func(b bool) {
			if !b {
				return
			}
			id, _ := strconv.ParseInt(networkId.Text, 10, 64)
			if err := wl.d.Delete(nil, id, &model.Network{}); err != nil {
				wl.showErr(err)
				return
			}
		}, w)
		d.Show()
	})

	updateButton := widget.NewButton("UPDATE", func() {
		d := dialog.NewConfirm("Sure", "updated", func(b bool) {
			if !b {
				return
			}
			id, _ := strconv.ParseInt(networkId.Text, 10, 64)

			cid, err := strconv.ParseInt(strings.TrimSpace(chainId.Text), 10, 64)
			if err != nil {
				wl.showErr(err)
				return
			}

			if err := wl.d.Update(nil, id, &model.Network{
				Name:     strings.TrimSpace(name.Text),
				Rpc:      strings.TrimSpace(rpc.Text),
				ChainId:  cid,
				Symbol:   strings.TrimSpace(symbol.Text),
				Explorer: strings.TrimSpace(explorer.Text),
			}); err != nil {
				wl.showErr(err)
				return
			}
		}, w)
		d.Show()
	})

	vbox := container.NewVBox(tips, form, container.NewHSplit(delButton, updateButton))

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
		tips.SetText(tips.Text + ":" + items[id].Name)
		networkId.SetText(fmt.Sprintf("%d", items[id].Id))
		name.SetText(items[id].Name)
		rpc.SetText(items[id].Rpc)
		chainId.SetText(fmt.Sprintf("%d", items[id].ChainId))
		symbol.SetText(items[id].Symbol)
		explorer.SetText(items[id].Explorer)
	}
	list.OnUnselected = func(id widget.ListItemID) {
		tips.SetText("Please Add Network")
	}
	if len(items) > 0 {
		list.Select(0)
	}

	addButton := widget.NewButton("Add", func() {
		var (
			name, rpc, chainId, symbol, explorer = widget.NewEntry(), widget.NewEntry(), widget.NewEntry(), widget.NewEntry(), widget.NewEntry()
		)
		items := []*widget.FormItem{
			widget.NewFormItem("Name", name),
			widget.NewFormItem("RPC", rpc),
			widget.NewFormItem("ChainId", chainId),
			widget.NewFormItem("Symbol", symbol),
			widget.NewFormItem("Explorer", explorer),
		}

		dialog.ShowForm("Add Network", "Add", "Cancel", items, func(b bool) {
			if !b {
				return
			}

			cid, err := strconv.ParseInt(strings.TrimSpace(chainId.Text), 10, 64)
			if err != nil {
				wl.showErr(err)
				return
			}
			if err := wl.d.Insert(nil, &model.Network{
				Name:     strings.TrimSpace(name.Text),
				Rpc:      strings.TrimSpace(rpc.Text),
				ChainId:  cid,
				Symbol:   strings.TrimSpace(symbol.Text),
				Explorer: strings.TrimSpace(explorer.Text),
			}); err != nil {
				wl.showErr(err)
				return
			}
		}, wl.window)
	})

	left := container.NewVSplit(list, addButton)
	left.SetOffset(0.9)

	c := container.NewHSplit(left, container.NewMax(vbox))
	c.Offset = 0.2
	wl.showContent(c)
	return nil
}

func (wl *FyneWallet) showContent(c fyne.CanvasObject) {
	wl.content.Objects = []fyne.CanvasObject{c}
	wl.content.Refresh()
}

func (wl *FyneWallet) showErr(err error) {
	dialog.NewError(err, wl.window).Show()
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

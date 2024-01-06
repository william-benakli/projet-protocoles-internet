package UI

import (
	"encoding/hex"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"io"
	"net"
	"os"
	"path/filepath"
	. "projet-protocoles-internet/Tools"
	"projet-protocoles-internet/restpeer"
	"projet-protocoles-internet/udppeer/arbre"
)

import . "projet-protocoles-internet/udppeer"

var users []string
var userClicked string
var path string = "tmp/peers/"
var oldpath string = "tmp/peers/"

//https://developer.fyne.io/started/ doc

func InitPage() {
	InitRoot()
	LoginPage()
}

var windows fyne.Window

func init() {
	a := app.New()
	windows = a.NewWindow("PEERS | PROJET INTERNET ")
	windows.Resize(fyne.NewSize(700, 400))
	windows.CenterOnScreen()
	windows.FixedSize()

}

func LoginPage() {

	labelWelcome := widget.NewLabel("Bienvenue sur le projet PEERS, Connectez vous pour pouvoir continuer")

	//labelBienvenue.Resize()

	input := widget.NewEntry()
	input.SetPlaceHolder("Votre pseudo ")

	login := widget.NewButton("Connexion", func() {
		if len(input.Text) == 0 {
			Name = "0000HEHEH"
		} else {
			Name = "0000" + input.Text
		}
		ServeurPeer, _ := restpeer.GetMasterAddresse(ClientRestAPI, "https://jch.irif.fr:8443/peers/jch.irif.fr/addresses")
		SendUdpRequest(ConnUDP, NewRequestUDPExtension(GetGlobalID(), HelloRequest, int16(len(Name)), []byte(Name)), ServeurPeer.ListOfAddresses[0], " FIRST CONNEXION JULIUZS")
		PageMain()
	})

	boxLogin := container.New(layout.NewVBoxLayout(), labelWelcome, input, layout.NewSpacer(), login)
	content := container.New(layout.NewCenterLayout(), boxLogin)

	windows.SetContent(content)
	windows.ShowAndRun()
}

func PageMain() {

	windows.CenterOnScreen()
	windows.FixedSize()

	labelWelcome := widget.NewLabel("Bienvenue sur le projet PEERS " + Name[4:])
	contentCenter := container.New(layout.NewCenterLayout(), labelWelcome)

	tabs := container.NewAppTabs(
		container.NewTabItem("PEERS", getListUserGraphic()),
		container.NewTabItem("MES FICHIERS", widget.NewLabel("Vos fichiers")), //uploadeFileGraphic()),
	)

	boxLogin := container.NewBorder(contentCenter, nil, nil, nil, tabs)

	tabs.SetTabLocation(container.TabLocationLeading)
	windows.SetContent(boxLogin)
}

func getListUserGraphic() *fyne.Container {

	listPeerName := widget.NewList(
		func() int {
			return len(users)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Element")
		},
		func(i int, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(users[i])
		},
	)

	listPeerName.OnSelected = func(index int) {
		if index >= 0 && index < len(users) {
			userClicked = users[index]
			IP_ADRESS = restpeer.GetAdrFromNamePeers(userClicked)
			PageUser()
		}
	}

	butonRefresh := widget.NewButton("Rafraichir les pairs", func() {
		users = restpeer.GetRestPeerNames(ClientRestAPI)
		listPeerName.Refresh()
	})

	footer := container.NewVBox(
		butonRefresh,
	)
	listPeerName.Refresh()

	c := container.New(layout.NewBorderLayout(nil, footer, nil, nil), footer, listPeerName)
	return c
}

func PageUser() {

	var label = widget.NewLabel(" Menu de " + userClicked)

	retour := widget.NewButton("Retour", func() {
		path = "tmp/peers/"
		PageMain()
	})

	cRetourLabel := container.New(layout.NewBorderLayout(nil, nil, retour, nil), retour, label)

	PublicKey := widget.NewButton("Envoyer PublicKey", func() {
		rq := NewRequestUDPExtension(GetGlobalID()+1, PublicKeyReply, 0, make([]byte, 0))
		SendUdpRequest(ConnUDP, rq, IP_ADRESS, "MAIN")
	})

	Root := widget.NewButton("Envoyer Root", func() {
		rq := NewRequestUDPExtension(GetGlobalID()+1, RootRequest, int16(len(GetRacine().HashReceive)), GetRacine().HashReceive)
		SendUdpRequest(ConnUDP, rq, IP_ADRESS, "MAIN")
	})

	noOP := widget.NewButton("Envoyer NoOp", func() {
		rq := NewRequestUDPExtension(GetGlobalID()+1, NoOp, 0, []byte(""))
		SendUdpRequest(ConnUDP, rq, IP_ADRESS, "MAIN")
	})

	Hello := widget.NewButton("Envoyer Hello", func() {
		rq := NewRequestUDPExtension(GetGlobalID()+1, HelloRequest, int16(len(Name)), []byte(Name))
		SendUdpRequest(ConnUDP, rq, IP_ADRESS, "MAIN")
	})

	butonDownload := widget.NewButton("Telecharger", func() {
		fmt.Println("telechargement en cours... : ", userClicked)
		downloadFile(ConnUDP)
	})

	butonDownloadFileOnDisk := widget.NewButton("Mettre à jour les fichiers", func() {
		fmt.Println("telechargement en cours ", userClicked)
		arbre.BuildImage(GetRoot(), "tmp/peers/"+userClicked)
	})

	printTree := widget.NewButton("Afficher arbre", func() {
		arbre.AfficherArbre(GetRoot(), 0)
	})

	fileList := widget.NewList(
		func() int {
			files, _ := os.ReadDir(path)
			return len(files)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			files, _ := os.ReadDir(path)
			o.(*widget.Label).SetText(files[i].Name())
		},
	)

	fileList.OnSelected = func(id widget.ListItemID) {
		files, _ := os.ReadDir(path)
		file := files[id]
		if file.IsDir() {
			oldpath = path
			path += file.Name()
		}
		fileList.Refresh()
	}

	retourArbo := widget.NewButton("...", func() {
		path = oldpath
		oldpath = filepath.Dir(path)
		fileList.Refresh()
	})

	lefter := container.NewVBox(
		Hello,
		Root,
		PublicKey,
		noOP,
		layout.NewSpacer(),
		butonDownload,
		printTree,
		butonDownloadFileOnDisk,
	)
	header := container.NewVBox(
		cRetourLabel,
	)

	cfileListe := container.New(layout.NewBorderLayout(retourArbo, nil, nil, nil), retourArbo, fileList)
	c := container.New(layout.NewBorderLayout(header, nil, lefter, nil), header, lefter, cfileListe)
	windows.SetContent(c)
}

func uploadeFileGraphic() *fyne.Container {
	path = "tmp/users/"

	fileList := widget.NewList(
		func() int {
			files, _ := os.ReadDir(path)
			return len(files)
		},
		func() fyne.CanvasObject {
			return container.NewVBox(
				widget.NewLabel("Name"),
				widget.NewButton("Supprimer le fichier", nil),
			)
		},

		func(id widget.ListItemID, object fyne.CanvasObject) {
			files, _ := os.ReadDir(path)
			container := object.(*fyne.Container)
			button := container.Objects[0].(*widget.Button)
			label := container.Objects[1].(*widget.Label)
			button.SetText("Supprimer le fichier")
			label.SetText(files[id].Name())
		},
	)

	fileList.OnSelected = func(id widget.ListItemID) {
		files, _ := os.ReadDir(path)
		file := files[id]
		if file.IsDir() {
			oldpath = path
			path += file.Name()
		}
		fileList.Refresh()
	}

	retourArbo := widget.NewButton("...", func() {
		path = oldpath
		oldpath = filepath.Dir(path)
		fileList.Refresh()
	})

	cfileListe := container.New(layout.NewBorderLayout(retourArbo, nil, nil, nil), retourArbo, fileList)
	c := container.New(layout.NewBorderLayout(nil, nil, nil, nil), cfileListe)
	return c
}

func downloadFile(connexion *net.UDPConn) {

	//client -> root ->
	get, err := ClientRestAPI.Get("https://jch.irif.fr:8443/peers/" + userClicked + "/root")
	if err != nil {
		return
	}
	rootKey, err := io.ReadAll(get.Body)
	if err != nil {
		return
	}

	fmt.Println(hex.EncodeToString(rootKey) + "####################### ROOT KEYYYYYYYYY")

	requestDatum := NewRequestUDPExtension(GetGlobalID(), GetDatumRequest, int16(len(rootKey)), rootKey)
	go SendUdpRequest(connexion, requestDatum, IP_ADRESS, "DATUM")
}

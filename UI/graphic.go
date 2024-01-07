package UI

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"io"
	"net"
	"os"
	"path/filepath"
	. "projet-protocoles-internet/Tools"
	"projet-protocoles-internet/restpeer"
	"projet-protocoles-internet/udppeer/arbre"
	"projet-protocoles-internet/udppeer/cryptographie"
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
			Name = "HEHEH"
		} else {
			Name = input.Text
		}
		byteName := make([]byte, 4)
		byteName[0] = 0
		byteName[1] = 0
		byteName[2] = 0
		byteName[3] = 0
		byteName = append(byteName, []byte(Name)...)

		ServeurPeer, _ := restpeer.GetMasterAddresse(ClientRestAPI, "https://jch.irif.fr:8443/peers/jch.irif.fr/addresses")
		SendUdpRequest(ConnUDP, NewRequestUDPExtensionSigned(GetGlobalID(), HelloRequest, int16(len(byteName)), byteName), ServeurPeer.ListOfAddresses[0], " FIRST CONNEXION JULIUZS")
		PageMain()
		go MaintainConnexion(ConnUDP, ServeurPeer)
	})

	boxLogin := container.New(layout.NewVBoxLayout(), labelWelcome, input, layout.NewSpacer(), login)
	content := container.New(layout.NewCenterLayout(), boxLogin)

	windows.SetContent(content)
	windows.ShowAndRun()
}

func PageMain() {

	windows.CenterOnScreen()
	windows.FixedSize()

	labelWelcome := widget.NewLabel("Bienvenue sur le projet PEERS " + Name)
	contentCenter := container.New(layout.NewCenterLayout(), labelWelcome)

	OpenFile := widget.NewButton("Ajouter des fichiers à mon pair", func() {
		dialog.ShowFileOpen(func(uri fyne.URIReadCloser, err error) {
			if uri != nil && err == nil {
				fmt.Println("Chemin du dossier sélectionné :", uri.URI().Path())
				go func() {
					destination := "tmp/user/" + filepath.Base(uri.URI().Path())
					copyFile(uri.URI().Path(), destination)
				}()
			}
		}, windows)
	})

	LoadFile := widget.NewButton("Recharger mes fichiers", func() {
		Racine, _ = arbre.ParcourirRepertoire2("tmp/user/")
		arbre.AfficherArbre(GetRacine(), 0)
	})

	bottomButton := container.New(layout.NewGridLayoutWithColumns(2), OpenFile, LoadFile)
	boxLogin := container.NewBorder(contentCenter, bottomButton, nil, nil, getListUserGraphic(), bottomButton)

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
		rq := NewRequestUDPExtensionSigned(GetGlobalID()+1, PublicKeyReply, 64, cryptographie.FormateKey()) // On utilise la fonction FormateKey
		SendUdpRequest(ConnUDP, rq, IP_ADRESS, "MAIN")
	})

	Root := widget.NewButton("Envoyer Root", func() {
		rq := NewRequestUDPExtensionSigned(GetGlobalID()+1, RootRequest, int16(len(GetRacine().HashReceive)), GetRacine().HashReceive)
		SendUdpRequest(ConnUDP, rq, IP_ADRESS, "MAIN")
	})

	noOP := widget.NewButton("Envoyer NoOp", func() {
		rq := NewRequestUDPExtension(GetGlobalID()+1, NoOp, 0, []byte(""))
		SendUdpRequest(ConnUDP, rq, IP_ADRESS, "MAIN")
	})

	Hello := widget.NewButton("Envoyer Hello", func() {
		byteName := make([]byte, 4)
		byteName[0] = 0
		byteName[1] = 0
		byteName[2] = 0
		byteName[3] = 0
		byteName = append(byteName, []byte(Name)...)
		rq := NewRequestUDPExtensionSigned(GetGlobalID()+1, HelloRequest, int16(len(byteName)), byteName)
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

func downloadFile(connexion *net.UDPConn) {

	get, err := ClientRestAPI.Get("https://jch.irif.fr:8443/peers/" + userClicked + "/root")
	if err != nil {
		return
	}
	rootKey, err := io.ReadAll(get.Body)
	if err != nil {
		return
	}

	requestDatum := NewRequestUDPExtension(GetGlobalID(), GetDatumRequest, int16(len(rootKey)), rootKey)
	go SendUdpRequest(connexion, requestDatum, IP_ADRESS, "DATUM")
}

func copyFile(src, dst string) {
	srcFile, err := os.Open(src)
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier source :", err)
		return
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		fmt.Println("Erreur lors de la création du fichier de destination :", err)
		return
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		fmt.Println("Erreur lors de la copie :", err)
		return
	}

	fmt.Println("Fichier copié avec succès :", dst)
}

package UI

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	. "projet-protocoles-internet/Tools"
	"projet-protocoles-internet/restpeer"
	"projet-protocoles-internet/udppeer/arbre"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

import . "projet-protocoles-internet/udppeer"

var listOfPeers restpeer.ListOfPeers
var users []string
var userClicked string
var arborescence string

//https://developer.fyne.io/started/ doc

//TODO Faire un composant peer
//TODO Ajouter autant de composant qu'il y'a des peers
//TODO faire que les peers soit clickable avec action dessus
//TODO faire un terminal de debug

func InitPage() {

	a := app.New()
	w := a.NewWindow("PEER | PROJET INTERNET ")
	w.Resize(fyne.NewSize(800, 600))
	var label = widget.NewLabel(arborescence)
	w.SetContent(widget.NewLabel("Hello World!"))

	butonRefresh := widget.NewButton("Rafraichir les pairs", func() {
		users = restpeer.GetRestPeerNames(ClientRestAPI)
		w.Resize(fyne.NewSize(801, 600))
	})

	butonDownload := widget.NewButton("Telecharger", func() {
		fmt.Println("telechargement en cours... : ", userClicked)
		downloadFile(ConnUDP)
		arborescence = "coucou bg\ncooooooc" // TODO DOIT RECUPERER LA REEL ARBORECENCE
		label.SetText(arborescence)
	})

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
			fmt.Println(IP_ADRESS, " SELECTED ADDRESSE")
			SendUdpRequest(ConnUDP, NewRequestUDPExtension(GetGlobalID(), HelloRequest, int16(len(Name)), []byte(Name)), IP_ADRESS, "MAIN")
		}
	}

	butonDownloadFileOnDisk := widget.NewButton("Mettre à jour les fichiers", func() {
		fmt.Println("telechargement en cours ", userClicked)
		arbre.BuildImage(GetRoot(), "tmp/peers/"+userClicked)
	})

	arbreUpdate := widget.NewButton("Mettre à jour mon arbre", func() {
		fmt.Println("Mise à jour de l'arbre ", userClicked)
		InitRoot()
	})

	header := container.NewVBox(
		butonRefresh,
	)
	lefter := container.NewVBox(
		butonDownload,
		butonDownloadFileOnDisk,
		arbreUpdate,
	)
	footer := container.NewVBox(
		label,
	)
	c := container.New(layout.NewBorderLayout(header, footer, lefter, nil), header, footer, lefter, listPeerName)
	w.SetContent(c)
	w.ShowAndRun()
}

func getListUserGraphic(w fyne.Window) {
	//renvoyer la liste des pairs
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

	requestDatum := NewRequestUDPExtension(rand.Int31(), GetDatumRequest, int16(len(rootKey)), rootKey)
	go SendUdpRequest(connexion, requestDatum, IP_ADRESS, "DATUM")
}

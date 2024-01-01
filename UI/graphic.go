package UI

import (
	"fmt"
	"net/http"
	"projet-protocoles-internet/restpeer"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var listOfPeers restpeer.ListOfPeers
var users []string
var userClicked string
var arborescence string

//https://developer.fyne.io/started/ doc

//TODO Faire un composant peer
//TODO Ajouter autant de composant qu'il y'a des peers
//TODO faire que les peers soit clickable avec action dessus
//TODO faire un terminal de debug

func InitPage(client *http.Client) {
	a := app.New()
	w := a.NewWindow("Hello World")
	w.Resize(fyne.NewSize(800, 600))
	var label = widget.NewLabel(arborescence)
	w.SetContent(widget.NewLabel("Hello World!"))
	butonRefresh := widget.NewButton("Rafraichir les pairs", func() {
		users = restpeer.GetRestPeerNames(client)
		w.Resize(fyne.NewSize(801, 600))
	})
	butonDownload := widget.NewButton("Telecharger", func() {
		fmt.Println("telechargement en cours... : ", userClicked)
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
			fmt.Println("Cliqué sur :", users[index])
			userClicked = users[index]
		}
	}

	header := container.NewVBox(
		butonRefresh,
	)
	lefter := container.NewVBox(
		butonDownload,
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

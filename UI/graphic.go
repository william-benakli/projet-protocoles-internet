package UI

import (
	"net/http"
	"projet-protocoles-internet/restpeer"
)

var listOfPeers restpeer.ListOfPeers

//https://developer.fyne.io/started/ doc

//TODO Faire un composant peer
//TODO Ajouter autant de composant qu'il y'a des peers
//TODO faire que les peers soit clickable avec action dessus
//TODO faire un terminal de debug

func InitPage(client *http.Client) {
	/*	a := app.New()
		w := a.NewWindow("Hello World")

		w.SetContent(widget.NewLabel("Hello World!"))
		content := widget.NewButton("Rafraichir les pairs", func() {
			log.Println("/* Update restCommands ")
		})

		w.SetContent(content)
		w.ShowAndRun()
	*/}

func getListUserGraphic() {
	//renvoyer la liste des pairs
}

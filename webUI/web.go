package webUI

import (
	"fmt"
	"log"
	"net/http"
	"projet-protocoles-internet/restpeer"
	"text/template"
)

var listOfPeers restpeer.ListOfPeers
var clientPeers *http.Client

/* Tuto suivi: 	https://www.youtube.com/watch?v=Gybqs9kJKZU&ab_channel=DrVipinClasses */
var templates *template.Template

type Data struct {
	DivContent string
}

func init() {
	templates = template.Must(template.ParseFiles("webUI/html/index.html"))
}

func SetupPage(client *http.Client) {
	clientPeers = client
	listOfPeers = restpeer.GetListOfPeers(clientPeers, restpeer.GetRestPeerNames(clientPeers))
	http.HandleFunc("/", peersPage)
	http.HandleFunc("/peersList", peersList)

	err := http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", nil)
	log.Fatal("ListenAndServe: ", err)
}

func peersPage(w http.ResponseWriter, r *http.Request) {

	PeersList := getListUserGraphic()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	data := Data{
		DivContent: PeersList,
	}

	err := templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		return
	}

}

func peersList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprint(w, getListUserGraphic())
}

// Cette fonction renvoie chaque utilisateurs dans une div séparée
func getListUserGraphic() string {
	listOfPeers = restpeer.GetListOfPeers(clientPeers, restpeer.GetRestPeerNames(clientPeers))

	PeersList := ""
	for i := 0; i < len(listOfPeers.ListOfPeers); i++ {
		addressesPeers := ""
		peer := listOfPeers.ListOfPeers[i]

		for y := 0; y < len(peer.ListOfAddresses); y++ {
			addressesPeers += peer.ListOfAddresses[y] + "<br>"
		}

		PeersList += " <div class='peer_contains' > <h4>" + peer.NameUser + "</h4> <p> Port:" + peer.Port + "</p><p> Addresses: </p>" + addressesPeers + "<br></div>"
	}
	return PeersList
}

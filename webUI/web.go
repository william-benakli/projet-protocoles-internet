package webUI

import (
	"log"
	"net/http"
	"projet-protocoles-internet/restpeer"
	"text/template"
)

var listOfPeers restpeer.ListOfPeers

/* Tuto suivi: 	https://www.youtube.com/watch?v=Gybqs9kJKZU&ab_channel=DrVipinClasses */
var templates *template.Template

type Data struct {
	DivContent string
}

func init() {
	templates = template.Must(template.ParseFiles("webUI/html/index.html"))
}

func SetupPage(client *http.Client) {
	listOfPeers = restpeer.GetListOfPeers(client, restpeer.GetRestPeerNames(client))
	http.HandleFunc("/", peersPage)
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

	if r.Method != "HEAD" && r.Method != "GET" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

}

// Cette fonction renvoie chaque utilisateurs dans une div séparée
func getListUserGraphic() string {
	PeersList := ""
	for i := 0; i < len(listOfPeers.ListOfPeers); i++ {
		addressesPeers := ""
		peer := listOfPeers.ListOfPeers[i]

		for y := 0; y < len(peer.ListOfAddresses); y++ {
			addressesPeers += peer.ListOfAddresses[y] + " "
		}

		PeersList += " <div class='peer_contains' > <h4>" + peer.NameUser + "</h4> <p> Addresses: " + addressesPeers + " Port: " + peer.Port + " </p> </div>"
	}
	return PeersList
}

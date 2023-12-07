package webUI

import (
	"fmt"
	"log"
	"net/http"
	"projet-protocoles-internet/restpeer"
)

var listOfPeers restpeer.ListOfPeers

func SetupPage(client *http.Client) {
	http.HandleFunc("/", peersPage)
	listOfPeers = restpeer.GetListOfPeers(client, restpeer.GetRestPeerNames(client))
	err := http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", nil)
	log.Fatal("ListenAndServe: ", err)
}

func peersPage(w http.ResponseWriter, r *http.Request) {

	peersList := ""
	for i := 0; i < len(listOfPeers.ListOfPeers); i++ {
		addressesPeers := ""
		peer := listOfPeers.ListOfPeers[i]

		for y := 0; y < len(peer.ListOfAddresses); y++ {
			addressesPeers += peer.ListOfAddresses[y] + " "
		}

		peersList += "Name: " + peer.NameUser + " Addresses: " + addressesPeers + " Port: " + peer.Port + "</br>"
	}

	data := `<!DOCTYPE html> <html> <head></head> <body> 
		<div style="display: flex; flex-direction: rows; border: solid;">
			<div style="border: solid;"> 
			<h1> Pair connectés: </h1 
			<p>` + peersList + `</p> 
			</div> 
			<div style="border: solid;">
				<h1> Racine </h1 
			</div>
		</div> 
		</body></html>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Method != "HEAD" && r.Method != "GET" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, data)

}

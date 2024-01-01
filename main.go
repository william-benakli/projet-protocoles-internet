package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"projet-protocoles-internet/restpeer"
	"projet-protocoles-internet/udppeer"
	"time"
)

func main() {

	fmt.Println("Lancement du programme")

	/* DEBUT Client pour REST API */
	transport := &*http.DefaultTransport.(*http.Transport)
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{
		Transport: transport,
		Timeout:   50 * time.Second,
	}
	/* FIN  Client pour REST API */

	/*
		Preparation des dossiers
	*/

	if err := os.MkdirAll("tmp/peers/", os.ModePerm); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll("tmp/user/", os.ModePerm); err != nil {
		log.Fatal(err)
	}

	udppeer.InitId()

	var PeerSingleton restpeer.PeersUser
	PeerSingleton.NameUser = "CharlyWilly"
	PeerSingleton.NameLen = int16(len(PeerSingleton.NameUser))

	fmt.Println("Connexion REST API terminée")

	ServeurPeer, err := restpeer.GetMasterAddresse(client, "https://jch.irif.fr:8443/peers/jch.irif.fr/addresses")
	connUdp, err := net.ListenUDP("udp", &net.UDPAddr{})
	channel := make(chan []byte)

	startClient(channel, connUdp, ServeurPeer)

	/* Lancement de UI Thread Principal */
	//UI.InitPage(client)

	//Si tu veux tester un autre thread lancer UI avec go
	//comme go UI.SetupPage(client)

	if err != nil {
		fmt.Println("Erreur lors de la création de la connexion UDP :", err)
		return
	}

}

func startClient(channel chan []byte, connUdp *net.UDPConn, ServeurPeer restpeer.PeersUser) {
	//Tout d'abord on écoute
	go udppeer.ListenActive(connUdp, channel)

	//on envoie Hello
	request, err := udppeer.SendUdpRequest(connUdp, udppeer.GetRequet(udppeer.HelloRequest, udppeer.GetGlobalID()), ServeurPeer.ListOfAddresses[0]+":"+ServeurPeer.Port, "MAIN")
	if err != nil {
		return
	}
	if request {
		//si tout c bien passé on envoie la suite des requetes et on reste connecté au serveur
		udppeer.SendUDPPacketFromResponse(connUdp, channel)
		go udppeer.MaintainConnexion(connUdp, ServeurPeer)

	} else {
		fmt.Println("La requête Hello n'a pas été envoyé...")
		fmt.Println("Fin du programme")
	}
}

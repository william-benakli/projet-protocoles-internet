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

var ConnUdp *net.UDPConn

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
	channel := make(chan []byte)
	ConnUdp, _ = net.ListenUDP("udp", &net.UDPAddr{})
	startClient(channel, ServeurPeer)

	/* Lancement de UI Thread Principal */
	//UI.InitPage(client)

	//Si tu veux tester un autre thread lancer UI avec go
	//comme go UI.SetupPage(client)

	if err != nil {
		fmt.Println("Erreur lors de la création de la connexion UDP :", err)
		return
	}

}

func startClient(channel chan []byte, ServeurPeer restpeer.PeersUser) {
	//Tout d'abord on écoute
	go udppeer.ListenActive(ConnUdp, channel)

	//on envoie Hello
	request, err := udppeer.SendUdpRequest(ConnUdp, udppeer.GetRequet(udppeer.HelloRequest, udppeer.GetGlobalID()), "81.194.27.155:8443", "MAIN")
	if err != nil {
		return
	}
	if request {
		//si tout c bien passé on envoie la suite des requetes et on reste connecté au serveur
		udppeer.SendUDPPacketFromResponse(ConnUdp, channel)
		go udppeer.MaintainConnexion(ConnUdp, ServeurPeer)

	} else {
		fmt.Println("La requête Hello n'a pas été envoyé...")
		fmt.Println("Fin du programme")
	}
}

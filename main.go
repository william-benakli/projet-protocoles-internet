package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	. "projet-protocoles-internet/Tools"
	"projet-protocoles-internet/UI"
	"projet-protocoles-internet/restpeer"
	"projet-protocoles-internet/udppeer"
	"sync"
	"time"
)

import . "projet-protocoles-internet/udppeer"

func main() {

	fmt.Println("Lancement du programme")

	/* DEBUT Client pour REST API */
	transport := &*http.DefaultTransport.(*http.Transport)
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	ClientRestAPI = &http.Client{
		Transport: transport,
		Timeout:   50 * time.Second,
	}

	if err := os.MkdirAll("tmp/peers/", os.ModePerm); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll("tmp/user/", os.ModePerm); err != nil {
		log.Fatal(err)
	}

	ShowDebug()

	udppeer.InitId()
	fmt.Println("Connexion REST API terminée")

	ServeurPeer, err := restpeer.GetMasterAddresse(ClientRestAPI, "https://jch.irif.fr:8443/peers/jch.irif.fr/addresses")
	channel := make(chan RequestUDPExtension)

	go startClient(channel, ConnUDP, ServeurPeer)
	UI.InitPage()

	if err != nil {
		fmt.Println("Erreur lors de la création de la connexion UDP :", err)
		return
	}

}

func startClient(channel chan udppeer.RequestUDPExtension, connUDP *net.UDPConn, ServeurPeer restpeer.PeersUser) {
	RequestTimes = sync.Map{}
	go ListenActive(connUDP, channel)
	//on envoie Hello
	go RemissionPaquets(connUDP, IP_ADRESS)
	go SendUDPPacketFromResponse(connUDP, channel)
	go MaintainConnexion(connUDP, ServeurPeer)

}

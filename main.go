package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
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

	udppeer.InitId()

	fmt.Println("Connexion REST API terminée")

	var peerSingleton restpeer.PeersUser
	peerSingleton.NameUser = "CharlyWilly"
	peerSingleton.NameLen = int16(len(peerSingleton.NameUser))

	ServeurPeer, err := restpeer.GetMasterAddresse(client, "https://jch.irif.fr:8443/peers/jch.irif.fr/addresses")
	connUdp, err := net.ListenUDP("udp", &net.UDPAddr{})
	channel := make(chan []byte)

	fmt.Println("Préparation UDP terminée")
	fmt.Println("Lancement des threads")

	startClient(channel, connUdp, ServeurPeer)

	if err != nil {
		fmt.Println("Erreur lors de la création de la connexion UDP :", err)
		return
	}

}

func startClient(channel chan []byte, connUdp *net.UDPConn, ServeurPeer restpeer.PeersUser) {

	request, err := udppeer.SendUdpRequest(connUdp, udppeer.GetRequet(udppeer.HelloRequest, udppeer.GetGlobalID()), ServeurPeer.AddressIpv4+":"+ServeurPeer.Port)
	if err != nil {
		return
	}
	if !request {
		fmt.Println("Erreur premier Hello n'a pas été envoyé fin du programme")
	} else {
		udppeer.ListenActive(connUdp, channel)
		go udppeer.MaintainConnexion(connUdp, ServeurPeer)
		go udppeer.SendUDPPacketFromResponse(connUdp, channel)
	}
}

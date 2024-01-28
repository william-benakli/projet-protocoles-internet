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
	"projet-protocoles-internet/udppeer"
	"projet-protocoles-internet/udppeer/cryptographie"
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

	cryptographie.PrivateKey = cryptographie.GeneratePrivateKey()
	cryptographie.PublicKey = cryptographie.GetPublicKey(cryptographie.PrivateKey)

	if err := os.MkdirAll("tmp/peers/", os.ModePerm); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll("tmp/user/", os.ModePerm); err != nil {
		log.Fatal(err)
	}
	//ShowDebug()
	udppeer.InitId()
	fmt.Println("Connexion REST API terminée")
	channel := make(chan RequestUDPExtension)
	go startClient(channel, ConnUDP)
	UI.InitPage()

}

func startClient(channel chan udppeer.RequestUDPExtension, connUDP *net.UDPConn) {
	RequestTimes = sync.Map{}
	go ListenActive(connUDP, channel)
	go RemissionPaquets(connUDP, IP_ADRESS)
	go SendUDPPacketFromResponse(connUDP, channel)
}

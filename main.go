package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"projet-protocoles-internet/UI"
	"projet-protocoles-internet/restpeer"
	"projet-protocoles-internet/udppeer"
	"projet-protocoles-internet/udppeer/Tools"
	"sync"
	"time"
)

import . "projet-protocoles-internet/udppeer"

import . "projet-protocoles-internet/udppeer/Tools"

var name = "PROUTE"

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

	Tools.ShowDebug()

	udppeer.InitId()
	fmt.Println("Connexion REST API terminée")

	ServeurPeer, err := restpeer.GetMasterAddresse(ClientRestAPI, "https://jch.irif.fr:8443/peers/jch.irif.fr/addresses")
	channel := make(chan RequestUDPExtension)
	connUDP, _ := net.ListenUDP("udp", &net.UDPAddr{})
	go startClient(channel, connUDP, ServeurPeer)

	UI.InitPage()

	/*for {
		var i int
		fmt.Print("Send Hello: ")
		fmt.Scan(&i)

		if i == 1 {
			SendUdpRequest(connUDP, NewRequestUDPExtension(GetGlobalID(), HelloRequest, int16(len(name)), []byte(name)), IP_ADRESS, "MAIN")

		} else if i == 2 {
			arbre.BuildImage(GetRoot(), "tmp/peers/juliuz/")
			arbre.AfficherArbre(GetRoot(), 0)

		} else {
			InitRoot()
		}
	}*/

	//UI.InitPage(client)
	/*if i == 0 {
		arbre.AfficherArbre(udppeer.GetRoot(), 0)

		//	cp := -1

		/*for {
			result := udppeer.VerifieNotEmpty(udppeer.GetRoot())
			fmt.Println("///////// cp", cp, " result", result)
			if result > 15 {
				udppeer.Remplir(udppeer.GetRoot(), connUDP)
				fmt.Println("----------------- ")
				cp = result
				time.Sleep(time.Millisecond * 10)
				if result > 1000 {
					arbre.AfficherArbre(udppeer.GetRoot(), 0)
					break
				}
			} else {
				break
			}
			//	arbre.AfficherArbre(udppeer.GetRoot(), 0)

		}

		fmt.Println(" ")
		fmt.Println(" ")
		fmt.Println(" ")
		fmt.Println(" ")
		fmt.Print("Attendre : ")
		fmt.Scan(&i)
		arbre.AfficherArbre(udppeer.GetRoot(), 0)
		arbre.BuildImage(udppeer.GetRoot())

	}

	/* Lancement de UI Thread Principal
	*/
	//UI.InitPage(client)

	//Si tu veux tester un autre thread lancer UI avec go
	//comme go UI.SetupPage(client)

	if err != nil {
		fmt.Println("Erreur lors de la création de la connexion UDP :", err)
		return
	}

}

func startClient(channel chan udppeer.RequestUDPExtension, connUDP *net.UDPConn, ServeurPeer restpeer.PeersUser) {
	//udppeer.LastPaquets = udppeer.LastPaquetsMutex{Paquets: make(map[int32]udppeer.RequestTime)}
	//Tout d'abord on écoute
	RequestTimes = sync.Map{}
	go ListenActive(connUDP, channel)
	//on envoie Hello
	go RemissionPaquets(connUDP, IP_ADRESS)
	go SendUDPPacketFromResponse(connUDP, channel)
	go MaintainConnexion(connUDP, ServeurPeer)

}

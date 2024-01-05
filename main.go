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
	"projet-protocoles-internet/udppeer/arbre"
	"projet-protocoles-internet/udppeer/cryptographie"
	"projet-protocoles-internet/utils"
	"sync"
	"time"
)

var name = "PROUTE"

func main() {
	cryptographie.PrivateKey = cryptographie.GeneratePrivateKey()
	cryptographie.PublicKey = cryptographie.GetPublicKey(cryptographie.PrivateKey)
	fmt.Println("Lancement du programme")

	/* DEBUT Client pour REST API */
	transport := &*http.DefaultTransport.(*http.Transport)
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{
		Transport: transport,
		Timeout:   50 * time.Second,
	}
	//UI.InitPage(client)

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
	fmt.Println("Connexion REST API terminée")

	ServeurPeer, err := restpeer.GetMasterAddresse(client, "https://jch.irif.fr:8443/peers/jch.irif.fr/addresses")
	channel := make(chan udppeer.RequestUDPExtension)
	connUDP, _ := net.ListenUDP("udp", &net.UDPAddr{})
	go startClient(channel, connUDP, ServeurPeer)

	var i int
	fmt.Print("Type a number: ")
	fmt.Scan(&i)

	if i == 0 {
		arbre.ParcoursRec(udppeer.GetRoot())
		arbre.AfficherArbre(udppeer.GetRoot(), 0)
		arbre.BuildImage(udppeer.GetRoot())
	}
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
	udppeer.RequestTimes = sync.Map{}

	go udppeer.ListenActive(connUDP, channel)

	//on envoie Hello
	go udppeer.RemissionPaquets(connUDP, udppeer.IP_ADRESS)
	request, err := udppeer.SendUdpRequest(connUDP, udppeer.NewRequestUDPExtensionSigned(udppeer.GetGlobalID(), udppeer.HelloRequest, int16(len(utils.NameUser)), []byte(utils.NameUser)), udppeer.IP_ADRESS, "MAIN")
	if err != nil {
		return
	}
	if request {

		udppeer.SendUDPPacketFromResponse(connUDP, channel)
		go udppeer.MaintainConnexion(connUDP, ServeurPeer)

	} else {
		fmt.Println("La requête Hello n'a pas été envoyé...")
		fmt.Println("Fin du programme")
	}
}

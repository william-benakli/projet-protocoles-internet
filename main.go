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

	/* Client test pour REST API */

	transport := &*http.DefaultTransport.(*http.Transport)
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{
		Transport: transport,
		Timeout:   50 * time.Second,
	}

	fmt.Println("-------- appel au fonction ------------ ")
	fmt.Println(restpeer.SendRestPeerNames(client))
	fmt.Println(restpeer.SendRestPeerAdresses(client, "jch.irif.fr"))
	fmt.Println("-------- appelle a getListOfPeers ------------ ")

	listeOfPeer := restpeer.GetListOfPeers(client, restpeer.SendRestPeerNames(client))
	fmt.Println(listeOfPeer)
	/* Client test pour REST API */
	udppeer.SendHello(listeOfPeer.ListOfPeers[0].AddressIpv4 + ":" + listeOfPeer.ListOfPeers[0].Port) // need to give IP+":"+port
	/*
		tableau := listenActive()
		fmt.Println(tableau)*/
	fmt.Println("Fin du programme")

}

func listenActive() []string {
	connUdp, err := net.ListenUDP("udppeer", &net.UDPAddr{})
	tableau := make([]string, 5)

	for compteur := 0; compteur < 10; compteur++ {
		if err != nil {
			fmt.Println("erreur listen not working")
		}
		maxRequest := make([]byte, 32*7)
		n, _, err := connUdp.ReadFromUDP(maxRequest)
		if n != len(maxRequest) {
			fmt.Println("Pas toutes les bits")
		}
		if err != nil {
			fmt.Println("Erreur ReadFromUDP")
		}
		tableau = append(tableau, string(maxRequest))
	}
	return tableau
}

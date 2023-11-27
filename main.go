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

var tableau []string

func main() {

	/* Client test pour REST API */
	transport := &*http.DefaultTransport.(*http.Transport)
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{
		Transport: transport,
		Timeout:   50 * time.Second,
	}
	udppeer.InitId()
	fmt.Println("-------- appel au fonction ------------ ")
	fmt.Println(restpeer.SendRestPeerNames(client))
	fmt.Println(restpeer.SendRestPeerAdresses(client, "jch.irif.fr"))
	fmt.Println("-------- appelle a getListOfPeers ------------ ")

	listeOfPeer := restpeer.GetListOfPeers(client, restpeer.SendRestPeerNames(client))
	fmt.Println(listeOfPeer)
	/* Client test pour REST API */
	channel := make(chan []byte)

	connUdp, err := net.ListenUDP("udp", &net.UDPAddr{})
	if err != nil {
		fmt.Println("Erreur lors de la création de la connexion UDP :", err)
		return
	}

	go listenActive(channel, connUdp)

	go udppeer.SendHello(connUdp, listeOfPeer.ListOfPeers[0].AddressIpv4+":"+listeOfPeer.ListOfPeers[0].Port) // need to give IP+":"+port

	for {
		bytesReceive, ok := <-channel // Receiving a message from the channel
		if !ok {
			fmt.Println("Channel closed. Exiting receiver.")
			return
		}
		receiveStruct := udppeer.ByteToStruct(bytesReceive)
		fmt.Println("Received ID :", receiveStruct.Id)
		fmt.Println("Received TYPE :", receiveStruct.Type)
		fmt.Println("Received NAME :", receiveStruct.Name)
		fmt.Println("Received LENGTH :", receiveStruct.Length)
		fmt.Println("Received EXTENSION :", receiveStruct.Extensions)
		fmt.Println("Received SIGNATURE:", receiveStruct.Signature)

	}
}

func listenActive(ch chan []byte, connUdp *net.UDPConn) {

	maxRequest := make([]byte, 32)
	for {
		n, _, err := connUdp.ReadFromUDP(maxRequest)
		if err != nil {
			fmt.Println("Erreur lors de la lecture UDP :", err)
			return
		}
		if n != len(maxRequest) {
			fmt.Println("Pas tous les bits lus")
		}

		ch <- maxRequest[:n]
	}
}

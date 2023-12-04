package udppeer

import (
	"fmt"
	"net"
	"projet-protocoles-internet/restpeer"
	"time"
)

func SendUDPPacket() {

}

func SendUDPPacketFromResponse(connUdp *net.UDPConn, channel chan []byte) {
	for {
		fmt.Println("SendUDPPacketFromResponse ")

		bytesReceive, ok := <-channel
		if !ok {
			fmt.Println("Channel closed. Exiting receiver.")
			return
		}

		if bytesReceive == nil {
			fmt.Println("Error closed. Exiting receiver.")
		}

		receiveStruct := ByteToStruct(bytesReceive)
		PrintRequest(receiveStruct)

		/* Ici gerer le cas d'erreur */

		//addressesOfPeer := restpeer.GetAdrFromNamePeers(receiveStruct.Name)
		//TODO changer "81.194.27.155:8443" par l'adresses du pair avec qui on discute
		//pour ça le recuperer dans une la liste ou le faire passer par la structure
		var request bool
		var err error

		switch receiveStruct.Type {

		case PublicKeyRequest:
			fmt.Println("Envoie PublicKeyRequest")
			request, err = SendUdpRequest(connUdp, GetRequet(PublicKeyReply, receiveStruct.Id), "81.194.27.155:8443")
		case RootRequest:
			fmt.Println("Envoie RootReply ")
			request, err = SendUdpRequest(connUdp, GetRequet(RootReply, globalID), "81.194.27.155:8443")

		case HelloReply:
			request, err = SendUdpRequest(connUdp, GetRequet(HelloRequest, globalID), "81.194.27.155:8443")
		case HelloRequest:
			request, err = SendUdpRequest(connUdp, GetRequet(HelloReply, globalID), "81.194.27.155:8443")

		case GetDatumRequest:
			request, err = SendUdpRequest(connUdp, GetRequet(Datum, globalID), "81.194.27.155:8443")

		}

		if err != nil {
			fmt.Println("Il y'a une erreur")
		}

		if request {
			fmt.Println("Requete envoyé avec succes")
		} else {
			fmt.Println("Erreur requete echec")
		}

	}
}

func MaintainConnexion(connUdp *net.UDPConn, ServeurPeer restpeer.PeersUser) {
	for tick := range time.Tick(25 * time.Second) {
		fmt.Println("MaintainConnexion : Envoie de hello")
		_, err := SendUdpRequest(connUdp, GetRequet(HelloRequest, globalID), string(ServeurPeer.AddressIpv4+":"+ServeurPeer.Port))
		if err != nil {
			return
		}
		fmt.Println(tick, " Envoie pour maintenir la connexion ")
	}

}

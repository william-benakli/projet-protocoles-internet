package udppeer

import (
	"fmt"
	"net"
	"projet-protocoles-internet/restpeer"
	"time"
)

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
		PrintRequest(receiveStruct, "RECEIVED") // Pour le debugage

		/* Ici gerer le cas d'erreur */

		//addressesOfPeer := restpeer.GetAdrFromNamePeers(receiveStruct.Name)
		//TODO changer "81.194.27.155:8443" par l'adresses du pair avec qui on discute
		//pour ça le recuperer dans une la liste ou le faire passer par la structure
		var request bool
		var err error

		switch receiveStruct.Type {

		case PublicKeyRequest:
			fmt.Println("Envoie PublicKeyRequest")
			request, err = SendUdpRequest(connUdp, GetRequet(PublicKeyReply, receiveStruct.Id), "81.194.27.155:8443", "PublicKeyRequest")

		case RootRequest:
			fmt.Println("Envoie RootReply ")
			request, err = SendUdpRequest(connUdp, GetRequet(RootReply, globalID), "81.194.27.155:8443", "ROOT Request")

			fmt.Println(string(receiveStruct.Name[1:32]), " master on a rien modif")

			hash := make([]byte, 32)
			copy(hash, receiveStruct.Name[1:32])
			requestDatum := NewRequestUDPExtension(globalID, GetDatumRequest, receiveStruct.Length, 0, hash)
			_, err = SendUdpRequest(connUdp, requestDatum, "81.194.27.155:8443", "DATUM")

		case HelloReply:
			//request, err = SendUdpRequest(connUdp, GetRequet(HelloRequest, globalID), "81.194.27.155:8443", "HelloReply")
			//Enrengistrer la pair en mémoire pendant au moins 180secondes
		case HelloRequest:
			request, err = SendUdpRequest(connUdp, GetRequet(HelloReply, globalID), "81.194.27.155:8443", "HelloRequest")

		case Datum:
			fmt.Println("Datum ")
			fmt.Println(string(receiveStruct.Name))

		case NoDatum:
			fmt.Println("No datum")

		case NoOp:
			fmt.Println("No op ignoré")

		default:
			continue

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
	for tick := range time.Tick(28 * time.Second) {
		_, err := SendUdpRequest(connUdp, GetRequet(HelloRequest, globalID), string(ServeurPeer.ListOfAddresses[0]+":"+ServeurPeer.Port), "MaintainConnexion")
		if err != nil {
			return
		}
		fmt.Println(tick, "maintien de la connexion avec le serveur")
	}

}

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

		addressesOfPeer := restpeer.GetAdrFromNamePeers(receiveStruct.Name)
		var request bool
		var err error

		switch receiveStruct.Type {

		case HelloReply:
			request, err = SendUdpRequest(connUdp, GetRequet(HelloRequest, globalID), addressesOfPeer)
		case HelloRequest:
			request, err = SendUdpRequest(connUdp, GetRequet(HelloReply, globalID), addressesOfPeer)
		case GetDatumRequest:
			request, err = SendUdpRequest(connUdp, GetRequet(Datum, globalID), addressesOfPeer)
		case PublicKeyRequest:
			request, err = SendUdpRequest(connUdp, GetRequet(PublicKeyReply, globalID), addressesOfPeer)
		case RootRequest:
			request, err = SendUdpRequest(connUdp, GetRequet(RootReply, globalID), addressesOfPeer)
		}

		if err != nil {
			return
		}

		if request {
			fmt.Println("Requete envoyé avec succes")
		} else {
			fmt.Println("Erreur requete echec")
		}

	}
}

func MaintainConnexion(connUdp *net.UDPConn, ServeurPeer restpeer.PeersUser) {
	for tick := range time.Tick(30 * time.Second) {
		_, err := SendUdpRequest(connUdp, GetRequet(HelloRequest, globalID), string(ServeurPeer.AddressIpv4+":"+ServeurPeer.Port))
		if err != nil {
			return
		}
		fmt.Println(tick)
	}
}

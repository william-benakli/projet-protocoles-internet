package udppeer

import (
	"fmt"
	"net"
	"projet-protocoles-internet/restpeer"
	"strings"
	"time"
)

//var racine Noeud
//var current *Noeud

func SendUDPPacketFromResponse(connUdp *net.UDPConn, channel chan []byte) {
	for {
		//	fmt.Println("SendUDPPacketFromResponse ")

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
		//var request bool
		var err error

		switch receiveStruct.Type {

		case PublicKeyRequest:
			fmt.Println("Envoie PublicKeyRequest")
			_, err = SendUdpRequest(connUdp, GetRequet(PublicKeyReply, receiveStruct.Id), "81.194.27.155:8443", "PublicKeyRequest")

		case RootRequest:
			fmt.Println("Envoie RootReply ")
			//_, err = SendUdpRequest(connUdp, GetRequet(RootReply, globalID), "81.194.27.155:8443", "ROOT Request")

			//hash := make([]byte, 32)
			//copy(hash, receiveStruct.Body[0:32])

			requestDatum := NewRequestUDPExtension(globalID, GetDatumRequest, int16(len(receiveStruct.Body)), receiveStruct.Body)
			_, err = SendUdpRequest(connUdp, requestDatum, "81.194.27.155:8443", "DATUM")

		case HelloReply:
			//request, err = SendUdpRequest(connUdp, GetRequet(HelloRequest, globalID), "81.194.27.155:8443", "HelloReply")
			//Enrengistrer la pair en mémoire pendant au moins 180secondes
		case Datum:

			fmt.Println("Datum ")
			fmt.Println("Hash :  ", receiveStruct.Body[0:31])
			fmt.Println("Type file: ", receiveStruct.Body[32])

			typeFormat := receiveStruct.Body[32]

			if typeFormat == 2 {

				nbFils := (receiveStruct.Length - 33) / 64
				for i := 0; i < int(nbFils); i++ {
					start_name := 33 + i*64
					//map[Body[start_name: start_name+32]] = Body[start_name: start_name+32]
					fmt.Println(removeEmpty(string(receiveStruct.Body[start_name : start_name+32])))
					//fmt.Println("Hash", removeEmpty(string(receiveStruct.Body[start_name+32:start_name+64])))
					requestDatum := NewRequestUDPExtension(globalID, GetDatumRequest, int16(len(receiveStruct.Body[start_name+32:start_name+64])), receiveStruct.Body[start_name+32:start_name+64])
					_, err = SendUdpRequest(connUdp, requestDatum, "81.194.27.155:8443", "DATUM")
				}

			} else if typeFormat == 1 {
				//Big file donc pour télécharger aussi
			} else if typeFormat == 0 {

				//TODO télécharger

			} else {
				fmt.Println("Cas non traitable")
			}

		case NoDatum:
			fmt.Println("No datum")
			//fmt.Println(string(receiveStruct.Body))

		case NoOp:
			fmt.Println("No op ignoré")

		case GetDatumRequest:
			//TODO Pour l'instant répondre NoDatum
			//TODO Ici que l'on va envoyé les fichiers de l'arbre de merkel

		default:
			continue

		}

		fmt.Println("                ")
		if err != nil {
			fmt.Println("Il y'a une erreur")
		}

		/*
			if request {
				fmt.Println("Requete envoyé avec succes")
			} else {
				fmt.Println("Erreur requete echec")
			}
		*/
	}
}

func removeEmpty(stringBody string) string {
	nullIndex := strings.IndexByte(stringBody, '\000')
	if nullIndex == -1 {
		return stringBody
	}
	return stringBody[:nullIndex]
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

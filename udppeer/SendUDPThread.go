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

// TODO faire plutot ID-IP car plusieurs ID peuvent être les memes
var listIdDejaVu []int32

var mapOfDirectoryHash = make(map[string]string)
var mapOfBigFileHash = make(map[string][]string)
var mapOfNameInt = make(map[string]int)
var mapOfChunckHash = make(map[string][]byte)

var compteur_rep int
var compteur_chunck int
var compteur_bf int

//TODO pourquoi ça boucle

func construireFileFromMap(nameFile string) {
	var byteOfFile []byte
	for nameFileOfBigFile := range mapOfBigFileHash {
		//fmt.Println("on entre 1", nameFileOfBigFile, nameFile)
		if nameFileOfBigFile == nameFile {
			//fmt.Println("on entre 2")

			listeOfHash := mapOfBigFileHash[nameFileOfBigFile]
			//fmt.Println(listeOfHash)

			for i := 0; i < len(listeOfHash); i++ {
				//fmt.Println("listeOfHash", i)
				//fmt.Println("mapOfChunckHash[listeOfHash[i]]", len(mapOfChunckHash[listeOfHash[i]]))

				for y := 0; y < len(mapOfChunckHash[listeOfHash[i]]); y++ {
					//fmt.Println(mapOfChunckHash[listeOfHash[i]][y], "mapOfChunckHash[listeOfHash[i]][y]")
					byteOfFile = append(byteOfFile, mapOfChunckHash[listeOfHash[i]][y])
				}
				//	fmt.Println(byteOfFile)
			}

		}
	}
	//fmt.Println(byteOfFile)
	//	_ = os.WriteFile(nameFile, byteOfFile, 0644)

}

func SendUDPPacketFromResponse(connUdp *net.UDPConn, channel chan []byte) {
	compteur_chunck = 0
	compteur_bf = 0
	compteur_rep = 0
	for {
		//fmt.Println("SendUDPPacketFromResponse ")
		fmt.Println("compteur_chunck", compteur_chunck, "compteur_bf", compteur_bf, "compteur_rep", compteur_rep)
		fmt.Println("mapOfBigFileHash", len(mapOfBigFileHash), "mapOfNameInt", len(mapOfNameInt), "mapOfChunckHash", len(mapOfChunckHash), "mapOfDirectoryHash", len(mapOfDirectoryHash))

		/*

			for elem := range mapOfDirectoryHash {
				fmt.Println(elem)
			}
			for nameFile := range mapOfBigFileHash {

				for nameFileOfMapOfInt := range mapOfNameInt {
					fmt.Println("-------------- lenght", mapOfNameInt[nameFileOfMapOfInt], len(mapOfBigFileHash[nameFile]))
					if mapOfNameInt[nameFileOfMapOfInt] == len(mapOfBigFileHash[nameFile]) {
						construireFileFromMap(nameFile)
					}
				}
			}
		*/
		//	fmt.Println("Avant buffer")
		bytesReceive, ok := <-channel
		receiveStruct := ByteToStruct(bytesReceive)

		if !ok {
			fmt.Println("Channel closed. Exiting receiver.")
			return
		}

		if bytesReceive == nil {
			fmt.Println("Error closed. Exiting receiver.")
		}

		//En cas de réemission de message on ignore le packet
		if containedList(listIdDejaVu, receiveStruct.Id) {
			continue
		}

		listIdDejaVu = append(listIdDejaVu, receiveStruct.Id)

		PrintRequest(receiveStruct, "RECEIVED") // Pour le debugage

		/* Ici gerer le cas d'erreur */

		//addressesOfPeer := restpeer.GetAdrFromNamePeers(receiveStruct.Name)
		//TODO changer "81.194.27.155:8443" par l'adresses du pair avec qui on discute
		//pour ça le recuperer dans une la liste ou le faire passer par la structure
		//var request bool
		var err error

		switch receiveStruct.Type {

		case PublicKeyRequest:
			fmt.Println("Envoie PublicKeyReply")
			_, err = SendUdpRequest(connUdp, GetRequet(PublicKeyReply, receiveStruct.Id), "81.194.27.155:8443", "PublicKeyRequest")

		case RootRequest:
			fmt.Println("Envoie RootReply ")
			_, err = SendUdpRequest(connUdp, GetRequet(RootReply, globalID), "81.194.27.155:8443", "ROOT Request")
			requestDatum := NewRequestUDPExtension(globalID, GetDatumRequest, int16(len(receiveStruct.Body)), receiveStruct.Body)
			_, err = SendUdpRequest(connUdp, requestDatum, "81.194.27.155:8443", "DATUM")

		case HelloReply:
			fmt.Println("Hello")

			//_, err = SendUdpRequest(connUdp, GetRequet(HelloRequest, globalID), "81.194.27.155:8443", "HelloReply")
			//Enrengistrer la pair en mémoire pendant au moins 180secondes
		case Datum:
			fmt.Println("Hash :  ", receiveStruct.Body[0:31])
			fmt.Println("Type file: ", receiveStruct.Body[32])
			typeFormat := receiveStruct.Body[32]

			if typeFormat == 2 {
				compteur_rep = compteur_rep + 1

				go directory(receiveStruct, connUdp)
			} else if typeFormat == 1 {
				compteur_bf = compteur_bf + 1

				bigFile(receiveStruct, connUdp)
			} else if typeFormat == 0 {
				compteur_chunck = compteur_chunck + 1
				chuck(receiveStruct, connUdp)
			} else {
				fmt.Println("Cas non traitable")
			}
		case NoDatum:
			fmt.Println("No datum")
			//fmt.Println(string(receiveStruct.Body))

		case NoOp:
			fmt.Println("No op ignoré")

		case GetDatumRequest:
			fmt.Println("GetDatum ")

		case Error:
			fmt.Println("Error ----------------- ")

			//TODO Pour l'instant répondre NoDatum
			//TODO Ici que l'on va envoyé les fichiers de l'arbre de merkel

		}

		fmt.Println("                ")
		if err != nil {
			fmt.Println("Il y'a une erreur")
		}
	}
}

func removeEmpty(stringBody string) string {
	nullIndex := strings.IndexByte(stringBody, '\000')
	if nullIndex == -1 {
		return stringBody
	}
	return stringBody[:nullIndex]
}

func directory(receiveStruct RequestUDPExtension, connUdp *net.UDPConn) {
	nbFils := (receiveStruct.Length - 33) / 64
	for i := 0; i < int(nbFils); i++ {
		start_name := 33 + i*64
		//map[Body[start_name: start_name+32]] = Body[start_name: start_name+32]
		fmt.Println(removeEmpty(string(receiveStruct.Body[start_name:start_name+32])), "/")
		//fmt.Println("Hash", removeEmpty(string(receiveStruct.Body[start_name+32:start_name+64])))
		mapOfDirectoryHash[string(receiveStruct.Body[start_name:start_name+32])] = string(receiveStruct.Body[start_name+32 : start_name+64])
		requestDatum := NewRequestUDPExtension(globalID, GetDatumRequest, int16(len(receiveStruct.Body[start_name+32:start_name+64])), receiveStruct.Body[start_name+32:start_name+64])
		_, _ = SendUdpRequest(connUdp, requestDatum, "81.194.27.155:8443", "DATUM")
	}
}

func bigFile(receiveStruct RequestUDPExtension, connUdp *net.UDPConn) {

	nbFils := (receiveStruct.Length - 1) / 32
	var listHashString []string

	for i := 0; i < int(nbFils); i++ {
		start_name := 1 + i*32
		listHashString = append(listHashString, string(receiveStruct.Body[start_name:start_name+32]))
		requestDatum := NewRequestUDPExtension(globalID, GetDatumRequest, int16(len(receiveStruct.Body[start_name:start_name+32])), receiveStruct.Body[start_name:start_name+32])
		_, _ = SendUdpRequest(connUdp, requestDatum, "81.194.27.155:8443", "DATUM")
	}

	for nomFichier := range mapOfDirectoryHash {
		//if mapOfDirectoryHash[nomFichier] == string(receiveStruct.Body[0:31]) {
		mapOfBigFileHash[nomFichier] = listHashString
		mapOfNameInt[nomFichier] = int(nbFils)
		//}
	}

	fmt.Println("BIF FILE")
}

func chuck(receiveStruct RequestUDPExtension, connUdp *net.UDPConn) {
	mapOfChunckHash[string(receiveStruct.Body[0:31])] = receiveStruct.Body[31:]
	fmt.Println("################## CHUNK #####################")
}
func containedList(listId []int32, id int32) bool {
	for i := 0; i < len(listId); i++ {
		if listId[i] == id {
			return true
		}
	}
	return false

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

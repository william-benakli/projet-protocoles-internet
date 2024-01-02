package udppeer

import (
	"encoding/hex"
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

// TODO verifier quand on recoit un helloreply si on est deja en communication avec juliuz avec l'historique
type Noeud struct {
	//	HashCalculate []byte
	Type        int8
	HashReceive []byte
	Data        []byte

	Fils []*Noeud
}

func afficherArbre(noeud *Noeud, niveau int) {
	if noeud == nil {
		return
	}
	if noeud.Type == 0 {
		return
	}

	indent := ""
	for i := 0; i < niveau; i++ {
		indent += "  "
	}

	hashStr := hex.EncodeToString(noeud.HashReceive)
	dataStr := string(noeud.Data)

	fmt.Printf("%sNoeud : Type %d Hash: %.5s, Data: %s/\n", indent, noeud.Type, hashStr, dataStr)
	for _, enfant := range noeud.Fils {
		afficherArbre(enfant, niveau+1)
	}
}

func changeDataFromHash(root *Noeud, hashATrouver []byte, newData []byte) bool {
	var queue []*Noeud
	queue = append(queue, root)

	for len(queue) > 0 {
		currentNode := queue[0]
		queue = queue[1:]

		if compareHashes(currentNode.HashReceive, hashATrouver) {
			currentNode.Data = newData
			return true // Retourne vrai si les données ont été modifiées
		}

		for _, child := range currentNode.Fils {
			queue = append(queue, child)
		}
	}

	return false // Retourne faux si aucun nœud avec le hash spécifié n'est trouvé
}

func setType(root *Noeud, hashATrouver []byte, typeFile int8) bool {
	var queue []*Noeud
	queue = append(queue, root)

	for len(queue) > 0 {
		currentNode := queue[0]
		queue = queue[1:]

		if compareHashes(currentNode.HashReceive, hashATrouver) {
			currentNode.Type = typeFile
			return true // Retourne vrai si les données ont été modifiées
		}

		for _, child := range currentNode.Fils {
			queue = append(queue, child)
		}
	}

	return false // Retourne faux si aucun nœud avec le hash spécifié n'est trouvé
}

func addNodeFromHash(root *Noeud, hash []byte, noeudToAdd *Noeud) {
	var queue []*Noeud
	queue = append(queue, root)

	for len(queue) > 0 {
		currentNode := queue[0]
		queue = queue[1:]

		// Vérifie si le nœud actuel a le hash recherché
		//fmt.Println(hex.EncodeToString(currentNode.HashReceive), " avec ", hex.EncodeToString(hash))

		if compareHashes(currentNode.HashReceive, hash) {
			currentNode.Fils = append(currentNode.Fils, noeudToAdd)
			fmt.Println("Hash trouvé j'ajoute")
			return
		}

		// Ajoute les fils du nœud actuel à la file
		for _, child := range currentNode.Fils {
			queue = append(queue, child)
		}
	}

	fmt.Println("Aucun Hash j'ajoute au noeud racine")
	root.HashReceive = hash
	root.Fils = append(root.Fils, noeudToAdd)
}

// compareHashes compare deux slices de bytes (hashes).
func compareHashes(hash1, hash2 []byte) bool {
	if len(hash1) != len(hash2) {
		return false
	}
	for i, b := range hash1 {
		if b != hash2[i] {
			return false
		}
	}
	return true
}

var root Noeud

func SendUDPPacketFromResponse(connUdp *net.UDPConn, channel chan []byte) {
	root.Type = 2
	afficherArbre(&root, 0)
	for {
		//fmt.Println("SendUDPPacketFromResponse ")
		//fmt.Println("compteur_chunck", compteur_chunck, "compteur_bf", compteur_bf, "compteur_rep", compteur_rep)
		//fmt.Println("mapOfBigFileHash", len(mapOfBigFileHash), "mapOfNameInt", len(mapOfNameInt), "mapOfChunckHash", len(mapOfChunckHash), "mapOfDirectoryHash", len(mapOfDirectoryHash))

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
			//fmt.Println("Hash :  ", receiveStruct.Body[0:31])
			fmt.Println("Type file: ", receiveStruct.Body[32])
			typeFormat := receiveStruct.Body[32]

			if typeFormat == 2 {

				nbFils := (receiveStruct.Length - 33) / 64

				for i := 0; i < int(nbFils); i++ {
					start_name := 33 + i*64
					//fmt.Println(removeEmpty(string(receiveStruct.Body[start_name:start_name+32])), "/")

					//					fmt.Println(hex.EncodeToString(receiveStruct.Body[0:32]), "avec ", hex.EncodeToString(receiveStruct.Body[start_name+32:start_name+64]))
					fils := &Noeud{Fils: make([]*Noeud, 0), HashReceive: receiveStruct.Body[start_name+32 : start_name+64], Data: receiveStruct.Body[start_name : start_name+32]}
					addNodeFromHash(&root, receiveStruct.Body[0:32], fils)

					//	fmt.Println("Hash fils ajouté: ", hex.EncodeToString(receiveStruct.Body[start_name+32:start_name+64]))
					requestDatum := NewRequestUDPExtension(globalID, GetDatumRequest, int16(len(receiveStruct.Body[start_name+32:start_name+64])), receiveStruct.Body[start_name+32:start_name+64])
					_, _ = SendUdpRequest(connUdp, requestDatum, "81.194.27.155:8443", "DATUM")
				}
				setType(&root, receiveStruct.Body[0:32], 2)

			} else if typeFormat == 1 {

				nbFils := (receiveStruct.Length - 1) / 32

				for i := 0; i < int(nbFils); i++ {
					start_name := 1 + i*32

					//fmt.Println("Hash fils recu: ", hex.EncodeToString(receiveStruct.Body[start_name:start_name+32]))

					fils := &Noeud{Fils: make([]*Noeud, 0), HashReceive: receiveStruct.Body[start_name : start_name+32], Data: make([]byte, 0), Type: 0}
					addNodeFromHash(&root, receiveStruct.Body[0:32], fils)

					requestDatum := NewRequestUDPExtension(globalID, GetDatumRequest, int16(len(receiveStruct.Body[start_name:start_name+32])), receiveStruct.Body[start_name:start_name+32])
					_, _ = SendUdpRequest(connUdp, requestDatum, "81.194.27.155:8443", "DATUM")
				}
				setType(&root, receiveStruct.Body[0:32], 1)

			} else if typeFormat == 0 {

				fmt.Println(changeDataFromHash(&root, receiveStruct.Body[0:32], receiveStruct.Body[31:]))
				afficherArbre(&root, 0)
			} else {
				fmt.Println("Cas non traitable")
			}
		case NoDatum:
			fmt.Println("No datum")
			//log.Fatal("eerr")
			//fmt.Println(string(receiveStruct.Body))

			//TODO Pour l'instant répondre NoDatum
			//TODO Ici que l'on va envoyé les fichiers de l'arbre de merkel

		case NoOp:
			fmt.Println("No op ignoré")

		case GetDatumRequest:
			fmt.Println("GetDatum ")

		case Error:
			fmt.Println("Error ----------------- ")

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

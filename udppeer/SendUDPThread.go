package udppeer

import (
	"encoding/hex"
	"fmt"
	"net"
	"os"
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
	NAME        string
	Data        []byte
	Fils        []*Noeud
}

const IP_ADRESS = "81.194.27.155:8443"

/*
func ParcoursRec(noeud *Noeud) error {

		switch noeud.Type {
		case 1: // bigfile
			var donnéesComplètes []byte
			for _, fils := range noeud.Fils {
				donnéesComplètes = append(donnéesComplètes, fils.Data...)
			}
			return os.WriteFile("tmp/peers/"+noeud.NAME, donnéesComplètes, 0644)
		case 2: // directory
			for _, fils := range noeud.Fils {
				if err := ParcoursRec(fils); err != nil {
					return err
				}
			}
		default: // chunk
			return os.WriteFile("tmp/peers/"+noeud.NAME, noeud.Data, 0644)
		}

		return nil
	}
*/
func buildImage(root *Noeud) {

	var queue []*Noeud
	queue = append(queue, root)

	for len(queue) > 0 {
		currentNode := queue[0]
		queue = queue[1:]

		//fmt.Println(removeEmpty(string(currentNode.Data)))

		if currentNode.Type == 0 {
			if len(currentNode.NAME) > 0 {
				os.WriteFile("tmp/peers/"+currentNode.NAME, currentNode.Data, 0644)
			}
		} else if currentNode.Type == 1 {
			bytetab := make([]byte, 0)

			for i := 0; i < len(currentNode.Fils); i++ {
				for j := 0; j < len(currentNode.Fils[i].Fils); j++ {
					for k := 0; k < len(currentNode.Fils[i].Fils[j].Data); k++ {
						bytetab = append(bytetab, currentNode.Fils[i].Fils[j].Data[k])
					}
				}
			}
			os.WriteFile("tmp/peers/"+currentNode.NAME, bytetab, 0644)

		}

		/*
			if removeEmpty(currentNode.NAME) == "README" {
				os.WriteFile("README", currentNode.Data, 0644)
				fmt.Println("READEME ???")
			} else if removeEmpty(string(currentNode.NAME)) == "horse.jpg" {
				bytetab := make([]byte, 0)

				for i := 0; i < len(currentNode.Fils); i++ {
					for j := 0; j < len(currentNode.Fils[i].Fils); j++ {
						for k := 0; k < len(currentNode.Fils[i].Fils[j].Data); k++ {
							bytetab = append(bytetab, currentNode.Fils[i].Fils[j].Data[k])
						}
					}
				}
				//	fmt.Printf("%.100s,dataaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa ", string(bytetab))

				os.WriteFile("horse.jpg", bytetab, 0644)
				fmt.Println("image ???")
			}*/

		for _, child := range currentNode.Fils {
			queue = append(queue, child)
		}
	}

}

func afficherArbre(noeud *Noeud, niveau int) {
	if noeud == nil {
		return
	}
	/*	if noeud.Type == 0 {
			return
		}
	*/
	indent := ""
	for i := 0; i < niveau; i++ {
		indent += "  "
	}

	hashStr := hex.EncodeToString(noeud.HashReceive)
	//dataStr := string(noeud.Data)
	/*if len(noeud.Data) == 0 && noeud.Type == 0 {
		return
	}*/
	fmt.Printf("%sNoeud : Type %d Fils: %d Hash: %.5s, Name: %s, Data: %d\n", indent, noeud.Type, len(noeud.Fils), hashStr, noeud.NAME, len(noeud.Data))

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
			copybyte := make([]byte, len(newData))
			copy(copybyte, newData)
			currentNode.Data = copybyte

			currentNode.Type = 0
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
			currentNode.Type = 1
			//fmt.Println("Hash trouvé j'ajoute")
			return
		}

		// Ajoute les fils du nœud actuel à la file
		for _, child := range currentNode.Fils {
			queue = append(queue, child)
		}
	}

	//fmt.Println("Aucun Hash j'ajoute au noeud racine")
	root.HashReceive = hash
	root.Fils = append(root.Fils, noeudToAdd)
}

// compareHashes compare deux slices de bytes (hashes).
func compareHashes(hash1, hash2 []byte) bool {
	if len(hash1) != len(hash2) {
		//fmt.Println("lenght hash  ", len(hash1), " != ", len(hash2))
		//fmt.Println(hex.EncodeToString(hash1), " ", hex.EncodeToString(hash2))
		return false
	}
	for i, b := range hash1 {
		if b != hash2[i] {
			//fmt.Println("different")
			//fmt.Println(hex.EncodeToString(hash1), " ", hex.EncodeToString(hash2))
			return false
		}
	}
	return true
}

var root Noeud
var lastRequest RequestUDPExtension

//TODO faire un téléchargement chunck en paralallelle

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
			_, err = SendUdpRequest(connUdp, GetRequet(PublicKeyReply, receiveStruct.Id), IP_ADRESS, "PublicKeyRequest")

		case RootRequest:
			fmt.Println("Envoie RootReply ")
			_, err = SendUdpRequest(connUdp, GetRequet(RootReply, globalID), IP_ADRESS, "ROOT Request")
			requestDatum := NewRequestUDPExtension(globalID, GetDatumRequest, int16(len(receiveStruct.Body)), receiveStruct.Body)
			_, err = SendUdpRequest(connUdp, requestDatum, IP_ADRESS, "DATUM")

		case RootReply:
			fmt.Println("Envoie RootReply ")

			requestDatum := NewRequestUDPExtension(globalID, GetDatumRequest, int16(len(receiveStruct.Body)), receiveStruct.Body)
			_, err = SendUdpRequest(connUdp, requestDatum, IP_ADRESS, "DATUM")
			//	lastId = globalID
			lastRequest = requestDatum

		case HelloReply:
			fmt.Println("Hello")
			_, err = SendUdpRequest(connUdp, GetRequet(RootRequest, globalID), IP_ADRESS, "ROOT Request")

			//_, err = SendUdpRequest(connUdp, GetRequet(HelloRequest, globalID), "81.194.27.155:8443", "HelloReply")
			//Enrengistrer la pair en mémoire pendant au moins 180secondes
		case Datum:

			//ici

			//	if receiveStruct.Id == lastId {
			go datumTree(connUdp, receiveStruct)
			//	} else {
			_, err = SendUdpRequest(connUdp, lastRequest, IP_ADRESS, "DATUM")
		//	}

		//TODO verifier que les hash sont bon
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

func datumTree(connUdp *net.UDPConn, receiveStruct RequestUDPExtension) {
	typeFormat := receiveStruct.Body[32]

	if typeFormat == 2 {

		/*
			D'abord verifier le hash du directory
		*/

		//		hashReceive := receiveStruct.Body[0:31]

		//		hashVerifie :=

		/*for i := 0; i < int(nbFils); i++ {
			start_name := 33 + i*64
			receiveStruct.Body[start_name+32 : start_name+64]
		}*/

		nbFils := (receiveStruct.Length - 33) / 64

		for i := 0; i < int(nbFils); i++ {
			start_name := 33 + i*64
			//fmt.Println(removeEmpty(string(receiveStruct.Body[start_name : start_name+32])))

			if removeEmpty(string(receiveStruct.Body[start_name:start_name+32])) != "videos" {

				fils := &Noeud{Fils: make([]*Noeud, 0), HashReceive: receiveStruct.Body[start_name+32 : start_name+64], Data: make([]byte, 0), NAME: removeEmpty(string(receiveStruct.Body[start_name : start_name+32]))}
				addNodeFromHash(&root, receiveStruct.Body[0:32], fils)

				//	fmt.Println("Hash fils ajouté: ", hex.EncodeToString(receiveStruct.Body[start_name+32:start_name+64]))
				requestDatum := NewRequestUDPExtension(globalID, GetDatumRequest, int16(len(receiveStruct.Body[start_name+32:start_name+64])), receiveStruct.Body[start_name+32:start_name+64])
				_, _ = SendUdpRequest(connUdp, requestDatum, IP_ADRESS, "DATUM")
				//lastId = globalID
				lastRequest = requestDatum

			}
		}
		setType(&root, receiveStruct.Body[0:32], 2)

	} else if typeFormat == 1 {

		nbFils := (receiveStruct.Length - 33) / 32
		fmt.Println(nbFils, "##########################")

		for i := 0; i < int(nbFils); i++ {
			start_name := 33 + i*32

			//fmt.Println("Hash fils recu: ", hex.EncodeToString(receiveStruct.Body[start_name:start_name+32]))

			fils := &Noeud{Fils: make([]*Noeud, 0), HashReceive: receiveStruct.Body[start_name : start_name+32], Data: make([]byte, 0)}
			addNodeFromHash(&root, receiveStruct.Body[0:32], fils)

			requestDatum := NewRequestUDPExtension(globalID, GetDatumRequest, int16(len(receiveStruct.Body[start_name:start_name+32])), receiveStruct.Body[start_name:start_name+32])
			_, _ = SendUdpRequest(connUdp, requestDatum, IP_ADRESS, "DATUM")
			//lastId = globalID
			lastRequest = requestDatum

		}
		setType(&root, receiveStruct.Body[0:32], 1)

	} else if typeFormat == 0 {
		//	fmt.Println(receiveStruct.Body[31:])

		//fmt.Println(receiveStruct.Length, "data du prof", string(receiveStruct.Body[33:]))
		if changeDataFromHash(&root, receiveStruct.Body[0:32], receiveStruct.Body[33:]) == false {
			//requestDatum := NewRequestUDPExtension(globalID, GetDatumRequest, int16(len(receiveStruct.Body[start_name:start_name+32])), receiveStruct.Body[start_name:start_name+32])
			//_, _ = SendUdpRequest(connUdp, requestDatum, IP_ADRESS, "DATUM")
		}
		//fmt.Println(len(receiveStruct.Body[33:]), "lenght chunck")
		//afficherArbre(&root, 0)
		buildImage(&root)
		/*err := ParcoursRec(&root)
		if err != nil {
			fmt.Println(err, "aucune génération")
		}*/
	} else {
		fmt.Println("Cas non traitable")
	}
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

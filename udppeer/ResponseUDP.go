package udppeer

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"
)

import . "projet-protocoles-internet/udppeer/arbre"

const REMISSION = 3
const TempsRemissionMiliSeconde = 2

func receiveResponse(connUdp *net.UDPConn, receiveStruct RequestUDPExtension) {

	//	if contains(LastPaquets, receiveStruct.Id) {
	switch receiveStruct.Type {

	case HelloReply: //
		fmt.Println("HELLO REPLY RECU")
		/*hasher := sha256.New()
		hash := hasher.Sum(nil)
		_, _ = SendUdpRequest(connUdp, NewRequestUDPExtension(receiveStruct.Id, RootRequest, int16(len(hash)), hash), IP_ADRESS, "ROOT Request")
		*/
	case PublicKeyReply:
		/* stocker la cle crypto */

	case RootReply: //
		fmt.Println("root reply")

		/*body := make([]byte, 0)
		body = append(body, racine.HashReceive...)

		rootRequest := NewRequestUDPExtension(globalID, RootRequest, int16(len(body)), body)
		_, _ = SendUdpRequest(connUdp, rootRequest, IP_ADRESS, "ROOT Request")
		*/
		requestDatum := NewRequestUDPExtension(rand.Int31(), GetDatumRequest, int16(len(receiveStruct.Body)), receiveStruct.Body)
		_, _ = SendUdpRequest(connUdp, requestDatum, IP_ADRESS, "DATUM")

	case Datum:
		datumTree(connUdp, receiveStruct)
		time.Sleep(time.Millisecond * 20)

	case NoDatum:
		log.Fatal("no datum")
		//TODO à changer
	}

	//} else {
	if receiveStruct.Type != NoOp {
		/*	go sendBackOfExpo(connUdp)
					} else {
			//			fmt.Println("Noop supprimé")
		*/
	}
	//sendBackOfExpo(receiveStruct.Id, connUdp)

	//	}

}

func sendBackOfExpo(storeID int32, connUdp *net.UDPConn) {

	/*
		found := false
		nMilliSeconde := time.Duration(2) * time.Millisecond
		time.Sleep(nMilliSeconde)

		for i := 0; i < len(LastPaquets); i++ {

			for i := 0; i < REMISSION; i++ {
				time.Sleep(nMilliSeconde)

				_, err := SendUdpRequest(connUdp, LastPaquets[storeID], IP_ADRESS, GetName(LastPaquets[storeID].Type))
				if err != nil {
					return
				}

				if !contains(LastPaquets, storeID) {
					found = true
					break
				}
				nMilliSeconde *= TempsRemissionMiliSeconde
			}

			if !found {
				delete(LastPaquets, storeID)

			}
			found = false

		}*/

}

func contains(lastPaquets map[int32]RequestUDPExtension, id int32) bool {
	for idMap := range lastPaquets {
		if idMap == id {
			delete(lastPaquets, id)
			return true
		}
	}
	return false
}

func datumTree(connUdp *net.UDPConn, receiveStruct RequestUDPExtension) {
	typeFormat := receiveStruct.Body[32]

	if typeFormat == 2 { //
		//	fmt.Println("DIR ##########################")

		nbFils := (receiveStruct.Length - 33) / 64

		//	fils := &Noeud{Fils: make([]*Noeud, 0), HashReceive: receiveStruct.Body[start_name+32 : start_name+64], Data: make([]byte, 0), NAME: removeEmpty(string(receiveStruct.Body[start_name : start_name+32]))}
		//	go AddNodeFromHash(&root, receiveStruct.Body[0:32], fils)

		for i := 0; i < int(nbFils); i++ {
			start_name := 33 + i*64
			//	//fmt.Println(removeEmpty(string(receiveStruct.Body[start_name : start_name+32])))

			if removeEmpty(string(receiveStruct.Body[start_name:start_name+32])) != "videos" { // || removeEmpty(string(receiveStruct.Body[start_name:start_name+32])) == "r.mp4" {

				fils := &Noeud{Fils: make([]*Noeud, 0), HashReceive: receiveStruct.Body[start_name+32 : start_name+64], Data: make([]byte, 0), NAME: removeEmpty(string(receiveStruct.Body[start_name : start_name+32]))}
				go AddNodeFromHash(&root, receiveStruct.Body[0:32], fils)

				requestDatum := NewRequestUDPExtension(rand.Int31(), GetDatumRequest, int16(len(receiveStruct.Body[start_name+32:start_name+64])), receiveStruct.Body[start_name+32:start_name+64])
				SendUdpRequest(connUdp, requestDatum, IP_ADRESS, "DATUM")

			}
		}
		SetType(&root, receiveStruct.Body[0:32], 2)

	} else if typeFormat == 1 {

		nbFils := (receiveStruct.Length - 33) / 32 //
		fmt.Println(nbFils, "BIG FILE ##########################")

		for i := 0; i < int(nbFils); i++ {
			startName := 33 + i*32

			//	//fmt.Println("Hash fils recu: ", hex.EncodeToString(receiveStruct.Body[start_name:start_name+32]))

			fils := &Noeud{Fils: make([]*Noeud, 0), HashReceive: receiveStruct.Body[startName : startName+32], Data: make([]byte, 0)}
			go AddNodeFromHash(&root, receiveStruct.Body[0:32], fils)

			requestDatum := NewRequestUDPExtension(rand.Int31(), GetDatumRequest, int16(len(receiveStruct.Body[startName:startName+32])), receiveStruct.Body[startName:startName+32])
			SendUdpRequest(connUdp, requestDatum, IP_ADRESS, "DATUM")

		}

		SetType(&root, receiveStruct.Body[0:32], 1)

	} else if typeFormat == 0 { //
		fmt.Println("CHUNKC ##########################")

		////	fmt.Println(receiveStruct.Body[31:])
		if !ChangeDataFromHash(&root, receiveStruct.Body[0:32], receiveStruct.Body[33:]) {
			fmt.Println("NOT FOUND")
		}
		//requestDatum := NewRequestUDPExtension(globalID, GetDatumRequest, int16(len(receiveStruct.Body[start_name:start_name+32])), receiveStruct.Body[start_name:start_name+32])
		//_, _ = SendUdpRequest(connUdp, requestDatum, IP_ADRESS, "DATUM")
	} else { //
		fmt.Println("Cas non traitable")
	}
}

/*
func Remplir(root *Noeud, connUdp *net.UDPConn) bool {
	var queue []*Noeud
	queue = append(queue, root)

	for len(queue) > 0 {
		currentNode := queue[0]
		queue = queue[1:]

		if currentNode.Type == 0 {
			////	fmt.Println(len(currentNode.Data))
			if len(currentNode.Data) == 0 {
				//	fmt.Printf("%.5s %s", hex.EncodeToString(currentNode.HashReceive), "        ")

				requestDatum := NewRequestUDPExtension(rand.Int31(), GetDatumRequest, int16(len(currentNode.HashReceive)), currentNode.HashReceive)
				SendUdpRequest(connUdp, requestDatum, IP_ADRESS, "DATUM")
			}
		}

		for _, child := range currentNode.Fils {
			queue = append(queue, child)
		}
	}

	return false // Retourne faux si aucun nœud avec le hash spécifié n'est trouvé
}

func VerifieNotEmpty(root *Noeud) int {
	cp := 0
	var queue []*Noeud
	queue = append(queue, root)

	for len(queue) > 0 {
		currentNode := queue[0]
		queue = queue[1:]

		if currentNode.Type == 1 {
			cp = cp + 1

		}
		//	if currentNode.Type == 0 {
		////	fmt.Println(len(currentNode.Data))
		//		if len(currentNode.Data) == 0 {
		//	//fmt.Println("test")
		//		if currentNode.NAME != "empty.txt" {
		//		}
		//	}
		//}

		for _, child := range currentNode.Fils {
			queue = append(queue, child)
		}
	}

	return cp
}
*/

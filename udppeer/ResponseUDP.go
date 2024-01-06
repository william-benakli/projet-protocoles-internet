package udppeer

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/sha3"
	"net"
	. "projet-protocoles-internet/Tools"
	"time"
)

import . "projet-protocoles-internet/udppeer/arbre"

func receiveResponse(connUdp *net.UDPConn, receiveStruct RequestUDPExtension) {

	//	if contains(LastPaquets, receiveStruct.Id) {
	switch receiveStruct.Type {

	case HelloReply: //
		fmt.Println("HELLO REPLY RECU")
	case PublicKeyReply:
		/* stocker la cle crypto */

	case RootReply: //
		fmt.Println("root reply")

		/*body := make([]byte, 0)
		body = append(body, racine.HashReceive...)

		rootRequest := NewRequestUDPExtension(globalID, RootRequest, int16(len(body)), body)
		_, _ = SendUdpRequest(connUdp, rootRequest, IP_ADRESS, "ROOT Request")

		requestDatum := NewRequestUDPExtension(rand.Int31(), GetDatumRequest, int16(len(receiveStruct.Body)), receiveStruct.Body)
		SendUdpRequest(connUdp, requestDatum, IP_ADRESS, "DATUM")
		*/
	case Datum:
		go datumTree(connUdp, receiveStruct)
		time.Sleep(time.Millisecond * 20)

	case NoDatum:
		fmt.Println("NO DATUM")
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

func datumTree(connUdp *net.UDPConn, receiveStruct RequestUDPExtension) {
	typeFormat := receiveStruct.Body[32]

	/* Verication que la hash est correct */
	hashFromRequete := receiveStruct.Body[0:32]
	hashCalculate := sha3.Sum256(receiveStruct.Body[32:])

	if CompareHashes(hashFromRequete, hashCalculate[:]) == false {
		PrintDebug("Hash incorrect")
		requestDatum := NewRequestUDPExtension(globalID, Error, int16(len("Bad datum")), []byte("Bad datum"))
		SendUdpRequest(connUdp, requestDatum, IP_ADRESS, "DATUM")
		return
	}

	if typeFormat == 2 { //
		//	fmt.Println("DIR ##########################")

		nbFils := (receiveStruct.Length - 33) / 64

		//	fils := &Noeud{Fils: make([]*Noeud, 0), HashReceive: receiveStruct.Body[start_name+32 : start_name+64], Data: make([]byte, 0), NAME: removeEmpty(string(receiveStruct.Body[start_name : start_name+32]))}
		//	go AddNodeFromHash(&root, receiveStruct.Body[0:32], fils)

		for i := 0; i < int(nbFils); i++ {
			start_name := 33 + i*64
			//	//fmt.Println(removeEmpty(string(receiveStruct.Body[start_name : start_name+32])))

			//if removeEmpty(string(receiveStruct.Body[start_name:start_name+32])) != "videos" { // || removeEmpty(string(receiveStruct.Body[start_name:start_name+32])) == "images" {

			fils := &Noeud{Fils: make([]*Noeud, 0), HashReceive: receiveStruct.Body[start_name+32 : start_name+64], Data: make([]byte, 0), NAME: removeEmpty(string(receiveStruct.Body[start_name : start_name+32])), ID: i}
			go AddNodeFromHash(&root, receiveStruct.Body[0:32], fils)

			fmt.Println("Envoie du datum")
			fmt.Println("Hash bigfile / chunck / rep #############", hex.EncodeToString(receiveStruct.Body[start_name+32:start_name+64]))

			requestDatum := NewRequestUDPExtension(globalID, GetDatumRequest, int16(len(receiveStruct.Body[start_name+32:start_name+64])), receiveStruct.Body[start_name+32:start_name+64])
			SendUdpRequest(connUdp, requestDatum, IP_ADRESS, "DATUM")

			//}
		}
		go SetType(&root, receiveStruct.Body[0:32], 2)

	} else if typeFormat == 1 {

		nbFils := (receiveStruct.Length - 33) / 32 //
		fmt.Println(nbFils, "BIG FILE ##########################")

		for i := 0; i < int(nbFils); i++ {
			startName := 33 + i*32

			//	//fmt.Println("Hash fils recu: ", hex.EncodeToString(receiveStruct.Body[start_name:start_name+32]))

			fils := &Noeud{Fils: make([]*Noeud, 0), HashReceive: receiveStruct.Body[startName : startName+32], Data: make([]byte, 0), ID: i}
			go AddNodeFromHash(&root, receiveStruct.Body[0:32], fils)

			fmt.Println("Envoie du datum")
			fmt.Println("Hash chunck #############", hex.EncodeToString(receiveStruct.Body[startName:startName+32]))

			requestDatum := NewRequestUDPExtension(globalID, GetDatumRequest, int16(len(receiveStruct.Body[startName:startName+32])), receiveStruct.Body[startName:startName+32])
			SendUdpRequest(connUdp, requestDatum, IP_ADRESS, "DATUM")

		}

		go SetType(&root, receiveStruct.Body[0:32], 1)

	} else if typeFormat == 0 { //
		fmt.Println("CHUNKC ##########################")
		if !ChangeDataFromHash(&root, receiveStruct.Body[0:32], receiveStruct.Body[33:]) {
			fmt.Println("NOT FOUND")
		}
	} else { //
		fmt.Println("Cas non traitable")
	}
}

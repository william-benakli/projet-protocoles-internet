package udppeer

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	. "projet-protocoles-internet/Tools"
	"projet-protocoles-internet/restpeer"
	"projet-protocoles-internet/udppeer/cryptographie"
	"time"
)

import . "projet-protocoles-internet/udppeer/arbre"

func receiveResponse(connUdp *net.UDPConn, receiveStruct RequestUDPExtension) {

	//	if contains(LastPaquets, receiveStruct.Id) {
	switch receiveStruct.Type {

	case HelloReply: //

		fmt.Println(string(receiveStruct.Body), " vide ?")

		rep := restpeer.GetPublicKey(ClientRestAPI, RemoveEmpty(string(receiveStruct.Body)))

		if rep == 200 {
			if !cryptographie.VerifyHash(receiveStruct.Body, receiveStruct.Signature) {
				requestTOSend := requestErrorReply(receiveStruct, "Bad signature")
				go SendUdpRequest(connUdp, requestTOSend, IP_ADRESS_SEND, GetName(requestTOSend.Type))
			}
		}

		fmt.Println("HELLO REPLY RECU")
	case PublicKeyReply:
		fmt.Println("Public KEY reply bien recu")
		return
	case RootReply: //
		fmt.Println("Root reply recu")
		return
	case Datum:
		go datumTree(connUdp, receiveStruct)
		time.Sleep(time.Millisecond * 100)

	case NoDatum:
		fmt.Println("NO DATUM")
	}

}

func datumTree(connUdp *net.UDPConn, receiveStruct RequestUDPExtension) {
	typeFormat := receiveStruct.Body[32]

	fmt.Println(typeFormat, string(receiveStruct.Body[32:]))

	/* Verication que la hash est correct */
	hashFromRequete := receiveStruct.Body[0:32]
	hashCalculate := sha256.Sum256(receiveStruct.Body[32:])

	fmt.Println(hex.EncodeToString(hashFromRequete), " ", hex.EncodeToString(hashCalculate[:]))
	fmt.Println(receiveStruct.Body)

	if !CompareHashes(hashFromRequete, hashCalculate[:]) {
		PrintDebug("Hash incorrect")
		requestDatum := NewRequestUDPExtension(globalID, Error, int16(len("Bad datum")), []byte("Bad datum"))
		SendUdpRequest(connUdp, requestDatum, IP_ADRESS, "DATUM")
		return
	}

	if typeFormat == 2 { //

		nbFils := (receiveStruct.Length - 33) / 64

		for i := 0; i < int(nbFils); i++ {
			start_name := 33 + i*64

			//if removeEmpty(string(receiveStruct.Body[start_name:start_name+32])) != "videos" { // || removeEmpty(string(receiveStruct.Body[start_name:start_name+32])) == "images" {

			fils := &Noeud{Fils: make([]*Noeud, 0), HashReceive: receiveStruct.Body[start_name+32 : start_name+64], Data: make([]byte, 0), NAME: RemoveEmpty(string(receiveStruct.Body[start_name : start_name+32])), ID: i}
			go AddNodeFromHash(&root, receiveStruct.Body[0:32], fils)

			fmt.Println("Envoie du datum")
			fmt.Println("Hash bigfile / chunck / rep ############# ", hex.EncodeToString(receiveStruct.Body[start_name+32:start_name+64]))

			requestDatum := NewRequestUDPExtension(globalID, GetDatumRequest, int16(len(receiveStruct.Body[start_name+32:start_name+64])), receiveStruct.Body[start_name+32:start_name+64])
			SendUdpRequest(connUdp, requestDatum, IP_ADRESS, "DATUM")
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

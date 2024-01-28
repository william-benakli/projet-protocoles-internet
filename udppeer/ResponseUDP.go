package udppeer

import (
	"crypto/sha256"
	"fmt"
	"net"
	. "projet-protocoles-internet/Tools"
	"projet-protocoles-internet/restpeer"
	"projet-protocoles-internet/udppeer/cryptographie"
	"time"
)

import . "projet-protocoles-internet/udppeer/arbre"

var DownloadFileRestant = 0

func receiveResponse(connUdp *net.UDPConn, receiveStruct RequestUDPExtension) {

	switch receiveStruct.Type {

	case HelloReply:
		rep := restpeer.GetPublicKey(ClientRestAPI, RemoveEmpty(string(receiveStruct.Body)))
		if rep == 200 {
			if len(receiveStruct.Signature) > 0 {
				if !cryptographie.VerifyHash(receiveStruct.Body, receiveStruct.Signature) {
					requestTOSend := requestErrorReply(receiveStruct, "Bad signature")
					go SendUdpRequest(connUdp, requestTOSend, IP_ADRESS_SEND, GetName(requestTOSend.Type))
				}
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

	/* Verication que la hash est correct */
	hashFromRequete := receiveStruct.Body[0:32]
	hashCalculate := sha256.Sum256(receiveStruct.Body[32:])

	if !CompareHashes(hashFromRequete, hashCalculate[:]) {
		PrintDebug("Hash incorrect")
		requestDatum := NewRequestUDPExtension(globalID+1, Error, int16(len("Bad datum")), []byte("Bad datum"))
		SendUdpRequest(connUdp, requestDatum, IP_ADRESS, "DATUM")
		return
	}

	if typeFormat == 2 { //
		SetTypeRec(&Root, receiveStruct.Body[0:32], 2)
		nbFils := (receiveStruct.Length - 33) / 64

		for i := 0; i < int(nbFils); i++ {
			start_name := 33 + i*64
			fils := &Noeud{Fils: make([]*Noeud, 0), HashReceive: receiveStruct.Body[start_name+32 : start_name+64], Data: make([]byte, 0), NAME: RemoveEmpty(string(receiveStruct.Body[start_name : start_name+32])), ID: i}
			AddNodeFromHashRec(&Root, receiveStruct.Body[0:32], fils)
		}

		for i := 0; i < int(nbFils); i++ {
			start_name := 33 + i*64
			requestDatum := NewRequestUDPExtension(globalID+1, GetDatumRequest, int16(len(receiveStruct.Body[start_name+32:start_name+64])), receiveStruct.Body[start_name+32:start_name+64])
			SendUdpRequest(connUdp, requestDatum, IP_ADRESS, "DATUM")
		}

	} else if typeFormat == 1 {
		SetTypeRec(&Root, receiveStruct.Body[0:32], 1)
		nbFils := (receiveStruct.Length - 33) / 32 //
		DownloadFileRestant += int(nbFils)

		for i := 0; i < int(nbFils); i++ {
			startName := 33 + i*32
			fils := &Noeud{Fils: make([]*Noeud, 0), HashReceive: receiveStruct.Body[startName : startName+32], Data: make([]byte, 0), ID: i}
			AddNodeFromHashRec(&Root, receiveStruct.Body[0:32], fils)
		}

		for i := 0; i < int(nbFils); i++ {
			startName := 33 + i*32
			requestDatum := NewRequestUDPExtension(globalID+1, GetDatumRequest, int16(len(receiveStruct.Body[startName:startName+32])), receiveStruct.Body[startName:startName+32])
			SendUdpRequest(connUdp, requestDatum, IP_ADRESS, "DATUM")
		}

	} else if typeFormat == 0 { //
		if len(receiveStruct.Body[33:]) > 0 {
			ChangeDataFromHashRec(&Root, receiveStruct.Body[0:32], receiveStruct.Body[33:])
			DownloadFileRestant -= 1
			fmt.Println("Chunck restants avant la fin du téléchargement :", DownloadFileRestant)
		}
	} else {
		fmt.Println("Cas non traitable")
	}
}

func CheckChunck(noeud *Noeud) {

	if noeud.Type == ChunkType {
		if len(noeud.Data) == 0 && len(noeud.NAME) == 0 {
			requestDatum := NewRequestUDPExtension(GetGlobalID()+1, GetDatumRequest, int16(len(noeud.HashReceive)), noeud.HashReceive)
			SendUdpRequest(ConnUDP, requestDatum, IP_ADRESS, "DATUM")
		}
	}

	for _, child := range noeud.Fils {
		CheckChunck(child)
	}
}

package udppeer

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	. "projet-protocoles-internet/Tools"
	"projet-protocoles-internet/udppeer/arbre"
	"sort"
)

func receiveRequest(connexion *net.UDPConn, receiveStruct RequestUDPExtension) {

	var requestTOSend RequestUDPExtension

	switch receiveStruct.Type {

	case HelloRequest:
		requestTOSend = requestHelloReply(receiveStruct)
	case PublicKeyRequest:
		requestTOSend = requestPublicKeyReply(receiveStruct)
	case RootRequest:
		requestTOSend = requestRootReply(receiveStruct) /* Sa racine */
	case GetDatumRequest:
		requestGetDatumReply(connexion, receiveStruct)
		return
	case NoOp:
		fmt.Println("No OP")
		return
	case Error:
		fmt.Print("Paquet Error: ", string(receiveStruct.Body))
		return
	}

	go SendUdpRequest(connexion, requestTOSend, IP_ADRESS_SEND, GetName(requestTOSend.Type))
}

func requestHelloReply(receiveStruct RequestUDPExtension) RequestUDPExtension {
	return NewRequestUDPExtension(receiveStruct.Id, HelloReply, int16(len([]byte(Name))), []byte(Name))
}

func requestPublicKeyReply(receiveStruct RequestUDPExtension) RequestUDPExtension {
	return NewRequestUDPExtension(receiveStruct.Id, PublicKeyReply, 0, []byte(""))
}

func requestRootReply(receiveStruct RequestUDPExtension) RequestUDPExtension {
	return NewRequestUDPExtension(receiveStruct.Id, RootReply, int16(len(GetRacine().HashReceive)), GetRacine().HashReceive)
}

func requestGetDatumReply(connexion *net.UDPConn, receiveStruct RequestUDPExtension) {

	hashGetDatum := receiveStruct.Body

	fmt.Println("TEST 0 ")

	currentNode := arbre.GetHashDFS(GetRacine(), hashGetDatum)

	fmt.Println("TEST 1 ")
	if currentNode == nil {
		requestDatum := NewRequestUDPExtension(receiveStruct.Id, NoDatum, 0, make([]byte, 0))
		go SendUdpRequest(connexion, requestDatum, IP_ADRESS_SEND, "NO DATUM")
		return
	}

	fmt.Println(currentNode.NAME, currentNode.Type, currentNode.ID)

	fmt.Println("TEST 2 ")

	body := make([]byte, 0)
	body = append(body, currentNode.HashReceive...)

	if currentNode.Type == ChunkType {
		/* BODY [hash, type, data] */
		body = append(body, ChunkType)
		body = append(body, currentNode.Data...)
		fmt.Println("TEST 3 CHUNCK ")

	} else if currentNode.Type == BigFileType {

		body = append(body, BigFileType)

		sort.Sort(arbre.ByID(currentNode.Fils))

		for i := 0; i < len(currentNode.Fils); i++ {
			body = append(body, currentNode.Fils[i].HashReceive...)
		}

		fmt.Println("TEST 4 BIGFILE")

	} else if currentNode.Type == DirectoryType {

		fmt.Println("TEST 5 DIR")
		body = append(body, DirectoryType)

		sort.Sort(arbre.ByID(currentNode.Fils))

		for i := 0; i < len(currentNode.Fils); i++ {
			var arr [32]byte
			copy(arr[:], currentNode.Fils[i].NAME)
			body = append(body, arr[:]...)
			body = append(body, currentNode.Fils[i].HashReceive...)
		}
		fmt.Println("TEST 7 DIR")

	} else {
		error := "type non defini "
		requestDatum := NewRequestUDPExtension(receiveStruct.Id, Error, int16(len(error)), []byte(error))
		SendUdpRequest(connexion, requestDatum, IP_ADRESS_SEND, "NO DATUM")
		return
	}

	hashbody := sha256.Sum256(body[32:])
	fmt.Println(hex.EncodeToString(hashGetDatum), " ", hex.EncodeToString(hashbody[:]))

	requestDatum := NewRequestUDPExtension(receiveStruct.Id, Datum, int16(len(body)), body)
	SendUdpRequest(connexion, requestDatum, IP_ADRESS_SEND, "DATUM")
}

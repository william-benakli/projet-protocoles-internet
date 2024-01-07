package udppeer

import (
	"fmt"
	"net"
	. "projet-protocoles-internet/Tools"
	"projet-protocoles-internet/restpeer"
	"projet-protocoles-internet/udppeer/arbre"
	"projet-protocoles-internet/udppeer/cryptographie"
	"sort"
)

func receiveRequest(connexion *net.UDPConn, receiveStruct RequestUDPExtension) {

	var requestTOSend RequestUDPExtension

	switch receiveStruct.Type {

	case HelloRequest:

		rep := restpeer.GetPublicKey(ClientRestAPI, string(receiveStruct.Body))

		if rep == 200 {
			if (len(receiveStruct.Signature)) > 0 {
				if cryptographie.VerifyHash(receiveStruct.Body, receiveStruct.Signature) {
					requestTOSend = requestHelloReply(receiveStruct)
				} else {
					requestTOSend = requestErrorReply(receiveStruct, "Bad signature")
				}
			}
		} else if rep == 204 {
			requestTOSend = requestHelloReply(receiveStruct)
		} else {
			requestTOSend = requestHelloReply(receiveStruct)
			//requestTOSend = requestErrorReply(receiveStruct, "Pair inconnu")
		}
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

	fmt.Println(IP_ADRESS_SEND)

	go SendUdpRequest(connexion, requestTOSend, IP_ADRESS_SEND, GetName(requestTOSend.Type))
}

func requestHelloReply(receiveStruct RequestUDPExtension) RequestUDPExtension {
	return NewRequestUDPExtension(receiveStruct.Id, HelloReply, int16(len([]byte(Name))), []byte(Name))
}

func requestErrorReply(receiveStruct RequestUDPExtension, message string) RequestUDPExtension {
	return NewRequestUDPExtension(receiveStruct.Id, Error, int16(len([]byte(message))), []byte(message))
}

func requestPublicKeyReply(receiveStruct RequestUDPExtension) RequestUDPExtension {
	return NewRequestUDPExtensionSigned(receiveStruct.Id, PublicKeyReply, 64, cryptographie.FormateKey()) // On utilise la fonction FormateKey
}

func requestRootReply(receiveStruct RequestUDPExtension) RequestUDPExtension {
	return NewRequestUDPExtensionSigned(receiveStruct.Id, RootReply, int16(len(GetRacine().HashReceive)), GetRacine().HashReceive)
}

func requestGetDatumReply(connexion *net.UDPConn, receiveStruct RequestUDPExtension) {
	hashGetDatum := receiveStruct.Body
	currentNode := arbre.GetHashDFS(GetRacine(), hashGetDatum)
	if currentNode == nil {
		requestDatum := NewRequestUDPExtension(receiveStruct.Id, NoDatum, 0, make([]byte, 0))
		go SendUdpRequest(connexion, requestDatum, IP_ADRESS_SEND, "NO DATUM")
		return
	}

	fmt.Println(currentNode.NAME, currentNode.Type, currentNode.ID)

	body := make([]byte, 0)
	body = append(body, currentNode.HashReceive...)

	if currentNode.Type == ChunkType {
		/* BODY [hash, type, data] */
		body = append(body, ChunkType)
		body = append(body, currentNode.Data...)
	} else if currentNode.Type == BigFileType {

		body = append(body, BigFileType)
		sort.Sort(arbre.ByID(currentNode.Fils))
		for i := 0; i < len(currentNode.Fils); i++ {
			body = append(body, currentNode.Fils[i].HashReceive...)
		}
	} else if currentNode.Type == DirectoryType {
		body = append(body, DirectoryType)

		sort.Sort(arbre.ByID(currentNode.Fils))

		for i := 0; i < len(currentNode.Fils); i++ {
			var arr [32]byte
			copy(arr[:], currentNode.Fils[i].NAME)
			body = append(body, arr[:]...)
			body = append(body, currentNode.Fils[i].HashReceive...)
		}
	} else {
		error := "type non defini "
		requestDatum := NewRequestUDPExtension(receiveStruct.Id, Error, int16(len(error)), []byte(error))
		SendUdpRequest(connexion, requestDatum, IP_ADRESS_SEND, "NO DATUM")
		return
	}
	requestDatum := NewRequestUDPExtension(receiveStruct.Id, Datum, int16(len(body)), body)
	SendUdpRequest(connexion, requestDatum, IP_ADRESS_SEND, "DATUM")
}

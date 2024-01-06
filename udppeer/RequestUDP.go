package udppeer

import (
	"encoding/hex"
	"fmt"
	"net"
	. "projet-protocoles-internet/Tools"
	"projet-protocoles-internet/restpeer"
	"projet-protocoles-internet/udppeer/arbre"
)

func receiveRequest(connexion *net.UDPConn, receiveStruct RequestUDPExtension) {

	var requestTOSend RequestUDPExtension

	switch receiveStruct.Type {

	case HelloRequest:
		/* Son nom -> son ip */
		nom := string(receiveStruct.Body)
		IP_ADRESS_SEND = restpeer.GetAdrFromNamePeers(nom)
		if len(IP_ADRESS_SEND) > 0 {
			requestTOSend = requestHelloReply(receiveStruct)
		} else {
			fmt.Println("IP INCONNUE HelloRequest")
			return
		}
	case PublicKeyRequest:
		requestTOSend = requestPublicKeyReply(receiveStruct)
	case RootRequest:
		requestTOSend = requestRootReply(receiveStruct) /* Sa racine */
	case GetDatumRequest:
		requestGetDatumReply(connexion, receiveStruct)
	case NoOp:
		fmt.Println("No OP")
	case Error:
		fmt.Print("Paquet Error: ", string(receiveStruct.Body))
	}

	//restpeer.GetAdrFromNamePeers(requ)
	go SendUdpRequest(connexion, requestTOSend, IP_ADRESS, GetName(requestTOSend.Type))
}

func requestHelloReply(receiveStruct RequestUDPExtension) RequestUDPExtension {
	return NewRequestUDPExtension(receiveStruct.Id, HelloReply, int16(len(receiveStruct.Body)), receiveStruct.Body)
}

func requestPublicKeyReply(receiveStruct RequestUDPExtension) RequestUDPExtension {
	return NewRequestUDPExtension(receiveStruct.Id, PublicKeyReply, 0, []byte(""))
}

func requestRootReply(receiveStruct RequestUDPExtension) RequestUDPExtension {
	fmt.Println(GetRacine().HashReceive)
	return NewRequestUDPExtension(receiveStruct.Id, RootReply, int16(len(GetRacine().HashReceive)), GetRacine().HashReceive)
}

func requestGetDatumReply(connexion *net.UDPConn, receiveStruct RequestUDPExtension) {

	hashGetDatum := make([]byte, 0)
	copy(hashGetDatum, receiveStruct.Body[32:])

	fmt.Println("TEST 0 ")

	currentNode := getNoeudFromHash(hashGetDatum)

	fmt.Println(currentNode.NAME, currentNode.Type, currentNode.ID)

	fmt.Println("TEST 1 ")
	if currentNode == nil {
		requestDatum := NewRequestUDPExtension(globalID, NoDatum, 0, make([]byte, 0))
		go SendUdpRequest(connexion, requestDatum, IP_ADRESS, "NO DATUM")
		return
	}
	fmt.Println("TEST 2 ")

	body := make([]byte, 0)
	body = append(body, currentNode.HashReceive...)

	if currentNode.Type == ChunkType {
		/* BODY [hash, type, data] */
		hashCalculate := currentNode.HashReceive //sha256.Sum256(receiveStruct.Body[33:])

		fmt.Printf("%5.s\n", hex.EncodeToString(hashCalculate))

		if !arbre.CompareHashes(hashCalculate[:], hashGetDatum) {
			fmt.Println("PAS BON HASH 0")
			return
		}

		body = append(body, 0)
		body = append(body, currentNode.Data...)

		fmt.Println("TEST 3 CHUNCK ")

	} else if currentNode.Type == BigFileType {

		hashCalculate := currentNode.HashReceive

		if !arbre.CompareHashes(hashCalculate[:], hashGetDatum) {
			//requestDatum := NewRequestUDPExtension(globalID, Datum, int16(len(hashGetDatum)), hashGetDatum) // TODO constante
			//_, _ = SendUdpRequest(connexion, requestDatum, IP_ADRESS, "DATUM")
			fmt.Println("PAS BON HASH BF")
			return
		}

		body = append(body, 1)
		body = append(body, currentNode.Data...)

		for i := 0; i < len(currentNode.Fils); i++ {
			body = append(body, currentNode.Fils[i].HashReceive...)
		}
		fmt.Println("TEST 4 BIGFILE")

	} else if currentNode.Type == DirectoryType {

		fmt.Println("TEST 5 DIR")

		hashCalculate := currentNode.HashReceive

		if !arbre.CompareHashes(hashCalculate, hashGetDatum) {
			//requestDatum := NewRequestUDPExtension(globalID, Datum, int16(len(hashGetDatum)), hashGetDatum) // TODO constante
			//_, _ = SendUdpRequest(connexion, requestDatum, IP_ADRESS, "DATUM")
			fmt.Println("PAS BON HASH DIR")

			return
		}
		fmt.Println("TEST 6 DIR")

		body = append(body, 2)
		body = append(body, currentNode.Data...)

		for i := 0; i < len(currentNode.Fils); i++ {
			var arr [32]byte
			copy(arr[:], []byte(currentNode.Fils[i].NAME))
			body = append(body, arr[:]...)
			body = append(body, currentNode.Fils[i].HashReceive...)
		}
		fmt.Println("TEST 7 DIR")

	} else {
		error := "type non defini "
		requestDatum := NewRequestUDPExtension(globalID, Error, int16(len(error)), []byte(error))
		SendUdpRequest(connexion, requestDatum, IP_ADRESS, "NO DATUM")
	}

	fmt.Println("Envoie du datum")
	requestDatum := NewRequestUDPExtension(globalID, Datum, int16(len(body)), body)
	SendUdpRequest(connexion, requestDatum, IP_ADRESS, "DATUM")
}

package udppeer

import (
	"crypto/sha256"
	"fmt"
	"net"
	"projet-protocoles-internet/udppeer/arbre"
)

func receiveRequest(connexion *net.UDPConn, receiveStruct RequestUDPExtension) {

	var requestTOSend RequestUDPExtension

	switch receiveStruct.Type {

	case HelloRequest:
		requestTOSend = requestHelloReply(receiveStruct)
	case PublicKeyRequest:
		requestTOSend = requestPublicKeyReply(receiveStruct)
	case RootRequest:
		requestTOSend = requestRootReply(receiveStruct)
		requestDatum := NewRequestUDPExtension(globalID, GetDatumRequest, int16(len(receiveStruct.Body)), receiveStruct.Body)
		_, _ = SendUdpRequest(connexion, requestDatum, IP_ADRESS, "DATUM")
	case GetDatumRequest:
		requestGetDatumReply(connexion, receiveStruct)
	case NoOp:
		fmt.Println("No OP")
	case Error:
		fmt.Print(string(receiveStruct.Body))
	}

	_, err := SendUdpRequest(connexion, requestTOSend, IP_ADRESS, GetName(requestTOSend.Type))

	if err != nil {
		fmt.Println("sendUDP Failed RequestUDP.go ", err)
	}

}

func requestHelloReply(receiveStruct RequestUDPExtension) RequestUDPExtension {
	return NewRequestUDPExtension(receiveStruct.Id, HelloReply, int16(len(receiveStruct.Body)), receiveStruct.Body)
}

func requestPublicKeyReply(receiveStruct RequestUDPExtension) RequestUDPExtension {
	return NewRequestUDPExtension(receiveStruct.Id, PublicKeyReply, 0, []byte(""))
}

func requestRootReply(receiveStruct RequestUDPExtension) RequestUDPExtension {
	hasher := sha256.New()
	hash := hasher.Sum(nil)
	return NewRequestUDPExtension(receiveStruct.Id, RootReply, int16(len(hash)), hash)

}

func requestGetDatumReply(connexion *net.UDPConn, receiveStruct RequestUDPExtension) {
	hashGetDatum := receiveStruct.Body[32:64]
	currentNode := getNoeudFromHash(hashGetDatum)

	if currentNode == nil {
		requestDatum := NewRequestUDPExtension(globalID, NoDatum, 0, make([]byte, 0))
		_, _ = SendUdpRequest(connexion, requestDatum, IP_ADRESS, "NO DATUM")
		return
	}

	body := make([]byte, 0)
	body = append(body, currentNode.HashReceive...)

	if currentNode.Type == 0 {
		/* BODY [hash, type, data] */
		hashCalculate := sha256.Sum256(receiveStruct.Body[33:])

		if !arbre.CompareHashes(hashCalculate[:], hashGetDatum) {
			fmt.Println("PAS BON HASH 0")
			return
		}

		body = append(body, 0)
		body = append(body, currentNode.Data...)

	} else if currentNode.Type == 1 {

		hashCalculate := sha256.Sum256(receiveStruct.Body[33:])

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
	} else if currentNode.Type == 2 {

		hashCalculate := make([]byte, 0)

		sha := sha256.Sum256(hashCalculate[:])
		if !arbre.CompareHashes(sha[:], hashGetDatum) {
			//requestDatum := NewRequestUDPExtension(globalID, Datum, int16(len(hashGetDatum)), hashGetDatum) // TODO constante
			//_, _ = SendUdpRequest(connexion, requestDatum, IP_ADRESS, "DATUM")
			fmt.Println("PAS BON HASH DIR")

			return
		}

		body = append(body, 2)
		body = append(body, currentNode.Data...)

		for i := 0; i < len(currentNode.Fils); i++ {
			var arr [32]byte
			copy(arr[:], []byte(currentNode.Fils[i].NAME))
			body = append(body, arr[:]...)
			body = append(body, currentNode.Fils[i].HashReceive...)
		}
	} else {
		error := "type non defini "
		requestDatum := NewRequestUDPExtension(globalID, Error, int16(len(error)), []byte(error))
		_, _ = SendUdpRequest(connexion, requestDatum, IP_ADRESS, "NO DATUM")
	}

	requestDatum := NewRequestUDPExtension(globalID, Datum, int16(len(body)+32), body) // TODO constante
	_, _ = SendUdpRequest(connexion, requestDatum, IP_ADRESS, "DATUM")
}

package udppeer

import (
	"fmt"
	"net"
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
	case GetDatumRequest:
		/*//peut envoyer datum
		sendDatumReply()
		//ou
		sendNoDatum()


		déjà calculer notre arbre
		[doc]
			[images]
			[videos]

		doc -> images -> videos

		parcourir notre arbre, si on trouve le hash lui envoyer si c'est directory, bigfile, chuck
		sinon lui envoyer no datum

		*/
	case Error:
		fmt.Print(string(receiveStruct.Body))
	}

	_, err := SendUdpRequest(connexion, requestTOSend, IP_ADRESS, GetName(requestTOSend.Type))

	if err != nil {
		fmt.Println("sendUDP Failed RequestUDP.go ")
	}

}

func requestHelloReply(receiveStruct RequestUDPExtension) RequestUDPExtension {
	return NewRequestUDPExtension(receiveStruct.Id, HelloReply, int16(len(receiveStruct.Body)), receiveStruct.Body)
}

func requestPublicKeyReply(receiveStruct RequestUDPExtension) RequestUDPExtension {
	return NewRequestUDPExtension(receiveStruct.Id, PublicKeyReply, 0, []byte(""))
}

func requestRootReply(receiveStruct RequestUDPExtension) RequestUDPExtension {
	return NewRequestUDPExtension(receiveStruct.Id, RootReply, 0, []byte(""))

}

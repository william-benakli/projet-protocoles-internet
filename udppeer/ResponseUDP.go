package udppeer

import (
	"fmt"
	"net"
	"time"
)

const REMISSION = 3
const TempsRemissionMiliSeconde = 2

var lastPaquets map[int32]RequestUDPExtension

func receiveResponse(receiveStruct RequestUDPExtension, connUdp *net.UDPConn) {

	if contains(lastPaquets, receiveStruct.Id) {
		switch receiveStruct.Type {

		case HelloReply:

		}

	} else {
		if receiveStruct.Type != NoOp {
			/*	go sendBackOfExpo(connUdp)
				} else {
					fmt.Println("Noop supprimé")
			*/
		}

	}

}

func sendBackOfExpo(storeID int32, connUdp *net.UDPConn) {

	found := false
	nMilliSeconde := time.Duration(2) * time.Millisecond
	for i := 0; i < len(lastPaquets); i++ {

		for i := 0; i < REMISSION; i++ {
			err := connUdp.SetReadDeadline(time.Now().Add(nMilliSeconde))
			if err != nil {
				fmt.Println("Erreur envoier read udp")
			}

			_, err = SendUdpRequest(connUdp, lastRequest, IP_ADRESS, GetName(lastRequest.Type))
			if err != nil {
				return
			}
			if !contains(lastPaquets, storeID) {
				found = true
				break
			}
			nMilliSeconde *= TempsRemissionMiliSeconde
		}

		if !found {
			delete(lastPaquets, storeID)

		}
		found = false

	}

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

func getRequeteFromMap(lastPaquets map[int32]RequestUDPExtension, id int32) RequestUDPExtension {
	return lastPaquets[id]
}

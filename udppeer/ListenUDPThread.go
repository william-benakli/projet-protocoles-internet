package udppeer

import (
	"fmt"
	"net"
)

import . "projet-protocoles-internet/Tools"

func ListenActive(connUdp *net.UDPConn, ch chan RequestUDPExtension) {

	for {
		maxRequest := make([]byte, 1024+34+50)
		_, from, err := connUdp.ReadFromUDP(maxRequest)

		IP_ADRESS_SEND = fmt.Sprintf("%s:%d", from.IP.String(), from.Port)

		if err != nil {
			fmt.Println("Erreur lors de la lecture UDP :", err)
		} else {
			receiveStruct := ByteToStruct(maxRequest)

			if receiveStruct.Type >= 128 {
				RequestTimes.Delete(receiveStruct.Id)
			}

			ch <- receiveStruct

		}
	}
}

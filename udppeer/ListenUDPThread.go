package udppeer

import (
	"fmt"
	"net"
)

func ListenActive(connUdp *net.UDPConn, ch chan RequestUDPExtension) {

	for {
		maxRequest := make([]byte, 1024+34+50)
		_, _, err := connUdp.ReadFromUDP(maxRequest)
		//fmt.Println(maxRequest)
		if err != nil {
			fmt.Println("Erreur lors de la lecture UDP :", err)
		} else {
			receiveStruct := ByteToStruct(maxRequest)
			ch <- receiveStruct
		}
	}
}

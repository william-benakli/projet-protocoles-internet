package udppeer

import (
	"fmt"
	"net"
)

func ListenActive(connUdp *net.UDPConn, ch chan []byte) {

	maxRequest := make([]byte, 1024*32)
	for {
		_, _, err := connUdp.ReadFromUDP(maxRequest)

		if err != nil {
			fmt.Println("Erreur lors de la lecture UDP :", err)
			ch <- nil
		}

		ch <- maxRequest
	}
}

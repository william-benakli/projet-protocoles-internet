package udppeer

import (
	"fmt"
	"net"
)

func ListenActive(connUdp *net.UDPConn, ch chan []byte) {

	maxRequest := make([]byte, 64)
	for {
		n, _, err := connUdp.ReadFromUDP(maxRequest)

		if err != nil {
			fmt.Println("Erreur lors de la lecture UDP :", err)
			ch <- nil
		}

		if n != len(maxRequest) {
			fmt.Println("Pas tous les bits lus")
		}

		ch <- maxRequest
	}
}

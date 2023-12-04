package udppeer

import (
	"fmt"
	"net"
	"time"
)

func ListenActive(connUdp *net.UDPConn, ch chan []byte) {

	maxRequest := make([]byte, 32)
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

		time.Sleep(1 * time.Second)
		// ch <- maxRequest
		//ch <- "Données reçues : " + string(maxRequest[:n])

	}
}

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

			if receiveStruct.Type >= 128 {
				//	fmt.Println(LastPaquets.Paquets)
				//	fmt.Println("On supprime ", receiveStruct.Id)
				LastPaquets.mutex.Lock()
				delete(LastPaquets.Paquets, receiveStruct.Id)
				LastPaquets.mutex.Unlock()
				//	fmt.Println(LastPaquets.Paquets)
			}

			ch <- receiveStruct

		}
	}
}

func deleteFromPaquet(channelDelete chan RequestUDPExtension) {
	receiveStruct := <-channelDelete
	//	LastPaquets.mutex.Lock()
	delete(LastPaquets.Paquets, receiveStruct.Id)
	fmt.Println("On supprime ", receiveStruct.Id)
	//	LastPaquets.mutex.Unlock()
}

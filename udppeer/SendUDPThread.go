package udppeer

import (
	"fmt"
	"golang.org/x/exp/rand"
	"net"
	"projet-protocoles-internet/udppeer/arbre"
	"strings"
	"sync"
	"time"
)

type RequestTime struct {
	TIME    int64
	REQUEST RequestUDPExtension //	time.Now().UnixMilli()
}

// TODO faire plutot ID-IP car plusieurs ID peuvent être les memes
var listIdDejaVu []int32

var LastPaquets LastPaquetsMutex

type LastPaquetsMutex struct {
	mutex   sync.RWMutex
	Paquets map[int32]RequestTime
	//Paquets sync.Map
}

var requestTimes sync.Map

var dernierId int32

/*
SI j'envoie un packet j'ajoute à la liste
500ms -> renvoie
je retire de la map
*/
var root arbre.Noeud

func GetRoot() *arbre.Noeud {
	return &root
}

var LastRequest RequestUDPExtension //	time.Now().UnixMilli()

const IP_ADRESS = "81.194.27.155:8443"

func SendUDPPacketFromResponse(connUdp *net.UDPConn, channel chan RequestUDPExtension) {

	for {

		fmt.Println("Remission paquet ", len(LastPaquets.Paquets))

		receiveStruct, ok := <-channel

		//PrintRequest(receiveStruct, "RECEIVED -- before")
		if !ok {
			fmt.Println("Channel closed. Exiting receiver.")
			return
		}
		if containedList(listIdDejaVu, receiveStruct.Id) {
			continue
		}

		listIdDejaVu = append(listIdDejaVu, receiveStruct.Id)

		if receiveStruct.Type < 128 {
			receiveRequest(connUdp, receiveStruct)
		} else {
			receiveResponse(connUdp, receiveStruct)
		}

		//for tick := range time.Tick(2 * time.Second) {
		//		LastPaquets.mutex.Lock()

		//	LastPaquets.mutex.Unlock()
		//}

	}

}

func removeEmpty(stringBody string) string {
	nullIndex := strings.IndexByte(stringBody, '\000')
	if nullIndex == -1 {
		return stringBody
	}
	return stringBody[:nullIndex]
}

func containedList(listId []int32, id int32) bool {
	for i := 0; i < len(listId); i++ {
		if listId[i] == id {
			return true
		}
	}
	return false

}

func RemissionPaquets(connUdp *net.UDPConn, adressPort string) {
	for _ = range time.Tick(4 * time.Second) {
		for paquet := range LastPaquets.Paquets {
			//fmt.Println(time.Now().UnixMilli() - LastPaquets.Paquets[paquet].TIME)
			LastPaquets.mutex.RLock()

			if (time.Now().UnixMilli() - LastPaquets.Paquets[paquet].TIME) > 1000 {

				//	delete(LastPaquets.Paquets, LastPaquets.Paquets[paquet].REQUEST.Id)
				requestDatum := NewRequestUDPExtension(LastPaquets.Paquets[paquet].REQUEST.Id, LastPaquets.Paquets[paquet].REQUEST.Type, int16(len(LastPaquets.Paquets[paquet].REQUEST.Body)), LastPaquets.Paquets[paquet].REQUEST.Body)
				go SendUdpRequest(connUdp, requestDatum, IP_ADRESS, "remissionPaquets ")

				fmt.Println("Remission paquet ", len(LastPaquets.Paquets))
				time.Sleep(time.Duration(int(rand.Int63n(50))) * time.Millisecond)
			}
			LastPaquets.mutex.RUnlock()

		}

		//	LastPaquets.mutex.Unlock()
	}
}

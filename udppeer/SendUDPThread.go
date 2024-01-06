package udppeer

import (
	"fmt"
	"golang.org/x/exp/rand"
	"net"
	"projet-protocoles-internet/udppeer/arbre"
	"sync"
	"time"
)

import . "projet-protocoles-internet/Tools"

type RequestTime struct {
	TIME    int64
	REQUEST RequestUDPExtension //	time.Now().UnixMilli()
}

var RequestTimes sync.Map

var listIdDejaVu []int32
var root arbre.Noeud

func GetRoot() *arbre.Noeud {
	return &root
}

func SendUDPPacketFromResponse(connUdp *net.UDPConn, channel chan RequestUDPExtension) {

	for {

		receiveStruct, ok := <-channel

		PrintRequest(receiveStruct, " RECEIVED "+string(ReceiveCounter))

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

	}

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
	for _ = range time.Tick(5 * time.Second) {
		RequestTimes.Range(func(key, value interface{}) bool {

			if requestTime, ok := value.(RequestTime); ok {
				if (time.Now().UnixMilli() - requestTime.TIME) > 7000 {
					RequestTimes.Delete(key)
					if requestTime.REQUEST.Type > 0 {
						requestDatum := NewRequestUDPExtension(globalID, requestTime.REQUEST.Type, int16(len(requestTime.REQUEST.Body)), requestTime.REQUEST.Body)
						SendUdpRequest(connUdp, requestDatum, IP_ADRESS, "remissionPaquets ")
						time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
					}
				}
			}
			return true
		})
	}
}

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

var listIdDejaVu []int32

var RequestTimes sync.Map
var root arbre.Noeud

func GetRoot() *arbre.Noeud {
	return &root
}

var IP_ADRESS = "81.194.27.155:8443"

func SendUDPPacketFromResponse(connUdp *net.UDPConn, channel chan RequestUDPExtension) {

	for {

		receiveStruct, ok := <-channel

		PrintRequest(receiveStruct, " RECEIVED ")

		if !ok {
			fmt.Println("Channel closed. Exiting receiver.")
			return
		}
		if containedList(listIdDejaVu, receiveStruct.Id) {
			continue
		}

		fmt.Println(" LISTE ")
		listIdDejaVu = append(listIdDejaVu, receiveStruct.Id)

		if receiveStruct.Type < 128 {
			receiveRequest(connUdp, receiveStruct)
		} else {
			receiveResponse(connUdp, receiveStruct)
		}

		fmt.Println(" SORTIE ")

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
	for _ = range time.Tick(5 * time.Second) {
		RequestTimes.Range(func(key, value interface{}) bool {

			if requestTime, ok := value.(RequestTime); ok {
				if (time.Now().UnixMilli() - requestTime.TIME) > 7000 {
					RequestTimes.Delete(key)
					requestDatum := NewRequestUDPExtension(requestTime.REQUEST.Id, requestTime.REQUEST.Type, int16(len(requestTime.REQUEST.Body)), requestTime.REQUEST.Body)
					SendUdpRequest(connUdp, requestDatum, IP_ADRESS, "remissionPaquets ")
					fmt.Println("Remission paquet ")
					time.Sleep(time.Duration(int(rand.Int63n(50))) * time.Millisecond)
				}
			}
			return true
		})
	}
}

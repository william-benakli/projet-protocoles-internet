package udppeer

import (
	"fmt"
	"net"
	"projet-protocoles-internet/udppeer/arbre"
	"strings"
)

// TODO faire plutot ID-IP car plusieurs ID peuvent être les memes
var listIdDejaVu []int32
var LastPaquets map[int32]RequestUDPExtension

var root arbre.Noeud

func GetRoot() *arbre.Noeud {
	return &root
}

const IP_ADRESS = "81.194.27.155:8443"

func SendUDPPacketFromResponse(connUDP *net.UDPConn, channel chan RequestUDPExtension) {

	for {
		receiveStruct, ok := <-channel

		PrintRequest(receiveStruct, "RECEIVED -- before")

		if !ok {
			fmt.Println("Channel closed. Exiting receiver.")
			return
		}
		if containedList(listIdDejaVu, receiveStruct.Id) {
			continue
		}

		listIdDejaVu = append(listIdDejaVu, receiveStruct.Id)

		if receiveStruct.Type < 128 {
			receiveRequest(connUDP, receiveStruct)
		} else {
			receiveResponse(connUDP, receiveStruct)
		}
		fmt.Println("---------------------")
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

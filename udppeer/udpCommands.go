package udppeer

import (
	"fmt"
	"net"
)

/*
	UdpCommands
*/

var globalID int32

func InitId() {
	globalID = 5432345
}

func SendHello(connUdp *net.UDPConn, adresse string) {

	/* Packer hello  */
	helloUpdStruct := RequestUDPExtension{
		Id:         globalID,
		Type:       2,
		Length:     1,
		Extensions: 0,
		Name:       "Aaaaabbb",
	}
	globalID += 1
	isSend, err := SendUdpRequest(connUdp, helloUpdStruct, adresse)

	if err != nil {
		fmt.Print("Erreur SendUdpRequest", string(err.Error()))
	}
	if isSend {
		fmt.Println("Packet envoyé ")
	}

}

func sendPublicKey() {

}

func sendRoot() {

}

func sendGetDatum() {

}

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

	namePeer := "Hello!"
	helloUpdStruct := RequestUDPExtension{
		Id:         globalID,
		Type:       2,
		Length:     int16(len(namePeer)),
		Extensions: 0,
		Name:       []byte(namePeer),
	}
	fmt.Println(helloUpdStruct.Length, "##################")
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

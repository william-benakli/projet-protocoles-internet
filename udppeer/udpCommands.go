package udppeer

import "fmt"

/*
	UdpCommands
*/

var globalID int32

func InitId() {
	globalID = 3
}

func SendHello(port string) {

	/* Packer hello  */
	helloUpdStruct := RequestUDPExtension{
		Id:     globalID,
		Type:   2,
		Length: 0,
		Body:   make([]byte, 0),
		Name:   "ChachaBG",
	}
	globalID += 1
	isSend, err := SendUdpRequest(helloUpdStruct, port)

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

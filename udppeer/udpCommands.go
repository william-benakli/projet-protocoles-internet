package udppeer

import "fmt"

/*
	UdpCommands
*/

func SendHello() {

	/* Packer hello  */
	helloUpdStruct := RequestUDPExtension{
		Id:     43234,
		Type:   2,
		Length: 0,
		Body:   make([]byte, 0),
	}
	isSend, err := SendUdpRequest(helloUpdStruct, "")

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

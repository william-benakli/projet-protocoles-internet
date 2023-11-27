package udp

import "fmt"

/*
	UdpCommands
*/

func SendHello() {
	helloUpdStruct := RequestUDP{
		Id:         0,
		Type:       2,
		Length:     0,
		Body:       make([]byte, 0),
		Name:       make([]byte, 0),
		Extensions: make([]byte, 10),
		Signature:  make([]byte, 10),
	}
	response, err := SendUdpRequest(helloUpdStruct)

	if err != nil {
		fmt.Printf("Echec de envoie du message Hello", err)
	}

	if response.Id == 129 {
		fmt.Println("Enrengistrement terminé")
	}
}

func sendPublicKey() {

}

func sendRoot() {

}

func sendGetDatum() {

}

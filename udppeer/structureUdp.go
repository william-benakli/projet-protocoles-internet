package udppeer

import (
	"fmt"
	"log"
	"net"
)

type RequestUDP struct {
	Id     int32
	Type   int8
	Length int16
	Body   []byte
}

type RequestUDPExtension struct {
	Id         int32
	Type       int8
	Length     int16
	Body       []byte
	Extensions []byte
	Name       []byte
	Signature  []byte
}

func byteToStruct(bytes []byte) RequestUDPExtension {

	requestUdp := RequestUDPExtension{
		Id:         int32(bytes[0])*(1<<24) + (int32(bytes[1]) * (1 << 16)) + (int32(bytes[2]) * (1 << 8)) + int32(bytes[3]),
		Type:       int8(bytes[4]),
		Body:       make([]byte, 100),
		Extensions: nil,
		Name:       nil,
		Signature:  nil,
	}

	for i := 0; i < int(requestUdp.Length); i++ {
		requestUdp.Body[i] = bytes[i+7]
	}
	return requestUdp
}

func structToBytes(message RequestUDPExtension) []byte {
	msg := make([]byte, 7+message.Length)
	msg[0] = byte(message.Id >> 24)
	msg[1] = byte(message.Id >> 16)
	msg[2] = byte(message.Id >> 8)
	msg[3] = byte(message.Id)
	msg[4] = byte(message.Type)
	msg[5] = byte(message.Length * 256)
	msg[6] = byte(message.Length)
	for i := 0; i < int(message.Length); i++ {
		msg[7+i] = message.Body[i]
	}
	return msg
}

func SendUdpRequest(RequestUDP RequestUDPExtension, adressPort string) (bool, error) {
	structTobytes := structToBytes(RequestUDP)
	udpAddr, err := net.ResolveUDPAddr("udp", adressPort)
	if err != nil {
		fmt.Println("ResolveUDPAddr error : ")
		log.Fatal(err)
	}
	connUdp, err := net.ListenUDP("udp", &net.UDPAddr{})
	if err != nil {
		fmt.Println("ListenUDP error : ")
		log.Fatal(err)
	}
	count, err := connUdp.WriteToUDP(structTobytes, udpAddr)
	if err != nil {
		fmt.Println("WriteToUDP error : ")
		log.Fatal(err)
	}

	// verifier que le nbr caracter envoyé = taille structure
	return count == len(structTobytes), err
}

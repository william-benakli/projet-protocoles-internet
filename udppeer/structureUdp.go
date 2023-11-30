package udppeer

import (
	"fmt"
	"log"
	"net"
)

type RequestUDPExtension struct {
	Id         int32
	Type       uint8 // care changement a verifier!
	Length     int16
	Extensions int32
	Name       []byte
	//Signature  int8
}

/*
func byteToStruct(bytes []byte) RequestUDPExtension {

	requestUdp := RequestUDPExtension{
		Id:         int32(bytes[0])*(1<<24) + (int32(bytes[1]) * (1 << 16)) + (int32(bytes[2]) * (1 << 8)) + int32(bytes[3]),
		Type:       int8(bytes[4]),
		Extensions: int32(bytes[5])*(1<<24) + (int32(bytes[6]) * (1 << 16)) + (int32(bytes[7]) * (1 << 8)) + int32(bytes[8]),
		Name:       string(bytes[8]),
		Signature:  string(bytes[9]),
	}

	return requestUdp
}
*/
/*
func structToBytes(requete RequestUDPExtension) []byte {
	msg := make([]byte, 50)
	msg[0] = byte(requete.Id >> 24)
	msg[1] = byte(requete.Id >> 16)
	msg[2] = byte(requete.Id >> 8)
	msg[3] = byte(requete.Id)
	msg[4] = byte(requete.Type)
	msg[5] = byte(requete.Length * 256)
	msg[6] = byte(requete.Length)

	//copy(msg[11:], []byte(message.Name))
	return msg
}*/

func ByteToStruct(bytes []byte) RequestUDPExtension {
	result := RequestUDPExtension{}
	result.Id = int32(bytes[0])*(1<<24) + (int32(bytes[1]) * (1 << 16)) + (int32(bytes[2]) * (1 << 8)) + int32(bytes[3])
	result.Type = uint8(bytes[4])
	result.Length = int16(bytes[5])*(1<<8) + int16(bytes[6])
	result.Extensions = int32(bytes[7])*(1<<24) + (int32(bytes[8]) * (1 << 16)) + (int32(bytes[9]) * (1 << 8)) + int32(bytes[10])

	result.Name = make([]byte, result.Length)
	for i := 0; i < int(result.Length); i++ {
		result.Name[i] = bytes[11+i]
	}
	return result
}

func StructToBytes(requete RequestUDPExtension) []byte {
	lenBuffer := 4 + 1 + 2 + 4 + requete.Length
	buffer := make([]byte, lenBuffer)

	buffer[0] = byte(requete.Id >> 24)
	buffer[1] = byte(requete.Id >> 16)
	buffer[2] = byte(requete.Id >> 8)
	buffer[3] = byte(requete.Id)

	buffer[4] = byte(requete.Type)
	buffer[5] = byte(requete.Length >> 8)
	buffer[6] = byte(requete.Length)
	buffer[7] = byte(requete.Extensions >> 24)
	buffer[8] = byte(requete.Extensions >> 16)
	buffer[9] = byte(requete.Extensions >> 8)
	buffer[10] = byte(requete.Extensions)
	//copy(buffer[11:], requete.Name)

	for i := 0; i < int(requete.Length); i++ {
		buffer[11+i] = requete.Name[i]
	}
	fmt.Println(buffer, " BUFERRRRR PRINT")
	return buffer
}

func SendUdpRequest(connUdp *net.UDPConn, RequestUDP RequestUDPExtension, adressPort string) (bool, error) {
	structToBytes := StructToBytes(RequestUDP)

	fmt.Println("Bytes envoyés ")

	receiveStruct := ByteToStruct(structToBytes)
	fmt.Println("Received ID :", receiveStruct.Id)
	fmt.Println("Received TYPE :", receiveStruct.Type)
	fmt.Println("Received NAME :", string(receiveStruct.Name))
	fmt.Println("Received LENGTH :", receiveStruct.Length)
	fmt.Println("Received EXTENSION :", receiveStruct.Extensions)

	udpAddr, err := net.ResolveUDPAddr("udp", adressPort)
	if err != nil {
		fmt.Println("ResolveUDPAddr error : ")
		log.Fatal(err)
	}
	if err != nil {
		fmt.Println("ListenUDP error : ")
		log.Fatal(err)
	}

	count, err := connUdp.WriteToUDP(structToBytes, udpAddr)
	if err != nil {
		fmt.Println("WriteToUDP error : ")
		log.Fatal(err)
	}
	fmt.Println(structToBytes)
	// verifier que le nbr caracter envoyé = taille structure
	return count == len(structToBytes), err
}

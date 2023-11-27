package udppeer

import (
	"bytes"
	"encoding/binary"
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
	Extensions int32
	Name       string
	Signature  string
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
	result.Id = int32(binary.LittleEndian.Uint32(bytes[0:4]))
	result.Type = int8(bytes[4])
	result.Length = int16(binary.LittleEndian.Uint16(bytes[5:7]))
	result.Extensions = int32(binary.LittleEndian.Uint32(bytes[7:11]))
	result.Name = string(bytes[11:15])
	result.Signature = string(bytes[15:19])
	return result
}

func StructToBytes(requete RequestUDPExtension) ([]byte, error) {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.LittleEndian, requete.Id)
	err = binary.Write(buffer, binary.LittleEndian, requete.Type)
	err = binary.Write(buffer, binary.LittleEndian, requete.Length)
	err = binary.Write(buffer, binary.LittleEndian, requete.Extensions)
	err = binary.Write(buffer, binary.LittleEndian, []byte(requete.Name))
	err = binary.Write(buffer, binary.LittleEndian, []byte(requete.Signature))
	return buffer.Bytes(), err
}

func SendUdpRequest(connUdp *net.UDPConn, RequestUDP RequestUDPExtension, adressPort string) (bool, error) {
	fmt.Println(adressPort)
	structTobytes, err := StructToBytes(RequestUDP)
	udpAddr, err := net.ResolveUDPAddr("udp", adressPort)
	if err != nil {
		fmt.Println("ResolveUDPAddr error : ")
		log.Fatal(err)
	}
	if err != nil {
		fmt.Println("ListenUDP error : ")
		log.Fatal(err)
	}

	count, err := connUdp.WriteToUDP(structTobytes, udpAddr)
	if err != nil {
		fmt.Println("WriteToUDP error : ")
		log.Fatal(err)
	}
	fmt.Println(structTobytes)
	// verifier que le nbr caracter envoyé = taille structure
	return count == len(structTobytes), err
}

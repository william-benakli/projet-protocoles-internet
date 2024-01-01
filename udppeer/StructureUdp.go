package udppeer

import (
	"fmt"
	"net"
)

type RequestUDPExtension struct {
	Id     int32
	Type   uint8 // care changement a verifier!
	Length int16
	Body   []byte
	//Signature  int8
}

func NewRequestUDPExtension(id int32, typeVal uint8, length int16, body []byte) RequestUDPExtension {
	return RequestUDPExtension{
		Id:     id,
		Type:   typeVal,
		Length: length,
		Body:   body,
		//Name:       name,
		//Extensions: extensions,
	}
}

var globalID int32

func InitId() {
	globalID = 5432345
}

func GetGlobalID() int32 {
	return globalID
}

func ByteToStruct(bytes []byte) RequestUDPExtension {
	result := RequestUDPExtension{}
	result.Id = int32(bytes[0])*(1<<24) + (int32(bytes[1]) * (1 << 16)) + (int32(bytes[2]) * (1 << 8)) + int32(bytes[3])
	result.Type = bytes[4]
	result.Length = int16(bytes[5])*(1<<8) + int16(bytes[6])
	result.Body = make([]byte, result.Length)

	for i := 0; i < int(result.Length); i++ {
		result.Body[i] = bytes[7+i]
	}
	return result
}

// StructToBytes
// Cette fonction renvoie un tableau de bytes à partir d'une structure
// param: RequestUDPExtension, une structure
// return: un tableau de bytes
func StructToBytes(message RequestUDPExtension) []byte {
	buffer := make([]byte, 7+message.Length)
	buffer[0] = byte(message.Id >> 24)
	buffer[1] = byte(message.Id >> 16)
	buffer[2] = byte(message.Id >> 8)
	buffer[3] = byte(message.Id)
	buffer[4] = message.Type
	buffer[5] = byte(message.Length >> 8)
	buffer[6] = byte(message.Length)
	for i := 0; i < int(message.Length); i++ {
		buffer[7+i] = message.Body[i]
	}
	return buffer
}

func SendUdpRequest(connUdp *net.UDPConn, RequestUDP RequestUDPExtension, adressPort string, from string) (bool, error) {
	globalID += 1
	structToBytes := StructToBytes(RequestUDP)
	//PrintRequest(ByteToStruct(structToBytes), "SEND "+from) // Pour le debugage
	udpAddr, err := net.ResolveUDPAddr("udp", adressPort)
	count, err := connUdp.WriteToUDP(structToBytes, udpAddr)
	// verifier que le nbr caracter envoyé = taille structure
	return count == len(structToBytes), err // gestion d'erreur plus tard
}

func PrintRequest(requestUdp RequestUDPExtension, status string) {
	fmt.Println("                 ", status)
	fmt.Println("ID :", requestUdp.Id)
	fmt.Println("TYPE :", GetName(requestUdp.Type))
	//fmt.Println("NAME :", string(requestUdp.Body))
	fmt.Println("LENGTH :", requestUdp.Length)
	fmt.Println("                 ")

}

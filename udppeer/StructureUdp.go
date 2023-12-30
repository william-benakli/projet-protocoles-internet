package udppeer

import (
	"fmt"
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

func NewRequestUDPExtension(id int32, typeVal uint8, length int16, extensions int32, name []byte) RequestUDPExtension {
	return RequestUDPExtension{
		Id:         id,
		Type:       typeVal,
		Length:     length,
		Extensions: extensions,
		Name:       name,
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
	result.Extensions = int32(bytes[7])*(1<<24) + (int32(bytes[8]) * (1 << 16)) + (int32(bytes[9]) * (1 << 8)) + int32(bytes[10])

	result.Name = make([]byte, result.Length)
	for i := 0; i < int(result.Length); i++ {
		result.Name[i] = bytes[11+i]
	}
	return result
}

// StructToBytes
// Cette fonction renvoie un tableau de bytes à partir d'une structure
// param: RequestUDPExtension, une structure
// return: un tableau de bytes
func StructToBytes(requete RequestUDPExtension) []byte {
	lenBuffer := 4 + 1 + 2 + 4 + requete.Length
	buffer := make([]byte, lenBuffer)

	buffer[0] = byte(requete.Id >> 24)
	buffer[1] = byte(requete.Id >> 16)
	buffer[2] = byte(requete.Id >> 8)
	buffer[3] = byte(requete.Id)

	buffer[4] = requete.Type
	buffer[5] = byte(requete.Length >> 8)
	buffer[6] = byte(requete.Length)
	buffer[7] = byte(requete.Extensions >> 24)
	buffer[8] = byte(requete.Extensions >> 16)
	buffer[9] = byte(requete.Extensions >> 8)
	buffer[10] = byte(requete.Extensions)
	copy(buffer[11:], requete.Name)

	return buffer
}

func SendUdpRequest(connUdp *net.UDPConn, RequestUDP RequestUDPExtension, adressPort string, from string) (bool, error) {
	globalID += 1
	structToBytes := StructToBytes(RequestUDP)
	PrintRequest(ByteToStruct(structToBytes), "SEND "+from) // Pour le debugage
	udpAddr, err := net.ResolveUDPAddr("udp", adressPort)
	count, err := connUdp.WriteToUDP(structToBytes, udpAddr)
	// verifier que le nbr caracter envoyé = taille structure
	return count == len(structToBytes), err // gestion d'erreur plus tard
}

func PrintRequest(requestUdp RequestUDPExtension, status string) {
	fmt.Println("---------- ", status)
	fmt.Println("ID :", requestUdp.Id)
	fmt.Println("TYPE :", requestUdp.Type)
	fmt.Println("NAME :", string(requestUdp.Name))
	fmt.Println("LENGTH :", requestUdp.Length)
	fmt.Println("EXTENSION :", requestUdp.Extensions)
}

package udppeer

import (
	"fmt"
	"math/rand"
	"net"
	. "projet-protocoles-internet/Tools"
	"projet-protocoles-internet/restpeer"
	"projet-protocoles-internet/udppeer/cryptographie"
	"time"
)

type RequestUDPExtension struct {
	Id        int32
	Type      uint8 // care changement a verifier!
	Length    int16
	Body      []byte
	Signature []byte
}

func NewRequestUDPExtension(id int32, typeVal uint8, length int16, body []byte) RequestUDPExtension {
	return RequestUDPExtension{
		Id:     id,
		Type:   typeVal,
		Length: length,
		Body:   body,
	}
}

func NewRequestUDPExtensionSigned(id int32, typeVal uint8, length int16, body []byte) RequestUDPExtension {
	fmt.Println(length, "LENGHT aaaaaaaaaaa")
	buffer := make([]byte, 7+int(len(body)))
	buffer[0] = byte(id >> 24)
	buffer[1] = byte(id >> 16)
	buffer[2] = byte(id >> 8)
	buffer[3] = byte(id)
	buffer[4] = typeVal
	buffer[5] = byte(length >> 8)
	buffer[6] = byte(length)
	for i := 0; i < len(body); i++ {
		buffer[7+i] = body[i]
	}
	return RequestUDPExtension{
		Id:        id,
		Type:      typeVal,
		Length:    length,
		Body:      body,
		Signature: cryptographie.Encrypted(buffer),
	}

}

var globalID int32

func InitId() {
	globalID = rand.Int31n(130450)
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
	buffer := make([]byte, 7+int(message.Length)+len(message.Signature))
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
	if len(message.Signature) != 0 {
		for i := 0; i < 64; i++ {
			buffer[7+int(message.Length)+i] = message.Signature[i]
		}
	}
	return buffer
}

var SendCounter int = 0
var ReceiveCounter int = 0

func SendUdpRequest(connUdp *net.UDPConn, RequestUDP RequestUDPExtension, adressPort string, from string) {
	globalID += 1
	structToBytes := StructToBytes(RequestUDP)
	udpAddr, _ := net.ResolveUDPAddr("udp", adressPort)

	time.Sleep(time.Millisecond * 50)

	_, _ = connUdp.WriteToUDP(structToBytes, udpAddr)

	if RequestUDP.Type < 128 && RequestUDP.Type != 0 && RequestUDP.Type != 1 {
		var TimeRequestUDP RequestTime
		TimeRequestUDP.REQUEST = RequestUDP
		TimeRequestUDP.TIME = time.Now().UnixMilli()
		RequestTimes.Store(RequestUDP.Id, TimeRequestUDP)
	}

	PrintRequest(RequestUDP, "SEND: "+string(rune(SendCounter))+from)
	SendCounter += 1
}

func MaintainConnexion(connUdp *net.UDPConn, ServeurPeer restpeer.PeersUser) {
	for tick := range time.Tick(30 * time.Second) {
		byteName := make([]byte, 4)
		byteName[0] = 0
		byteName[1] = 0
		byteName[2] = 0
		byteName[3] = 0
		byteName = append(byteName, []byte(Name)...)
		SendUdpRequest(connUdp, NewRequestUDPExtensionSigned(GetGlobalID(), HelloRequest, int16(len(byteName)), byteName), ServeurPeer.ListOfAddresses[0], "MaintainConnexion")
		fmt.Println(tick, "maintien de la connexion avec le serveur")
	}
}

func PrintRequest(requestUdp RequestUDPExtension, status string) {
	fmt.Println("                 ", status)
	fmt.Println("ID :", requestUdp.Id)
	fmt.Println("TYPE :", GetName(requestUdp.Type), "(", requestUdp.Type, ")")
	fmt.Printf("NAME : %s %d\n", string(requestUdp.Body), len(requestUdp.Body))
	fmt.Println("LENGTH :", requestUdp.Length)
	fmt.Println("                 ")

}

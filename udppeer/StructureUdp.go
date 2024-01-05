package udppeer

import (
	"fmt"
	"net"
	"projet-protocoles-internet/restpeer"
	"time"
)

import . "projet-protocoles-internet/udppeer/Tools"

type RequestUDPExtension struct {
	Id     int32
	Type   uint8 // care changement a verifier!
	Length int16
	Body   []byte
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

func SendUdpRequest(connUdp *net.UDPConn, RequestUDP RequestUDPExtension, adressPort string, from string) {
	globalID += 1
	structToBytes := StructToBytes(RequestUDP)
	udpAddr, _ := net.ResolveUDPAddr("udp", adressPort)

	time.Sleep(time.Millisecond * 50)

	_, _ = connUdp.WriteToUDP(structToBytes, udpAddr)

	if RequestUDP.Type < 128 {
		var TimeRequestUDP RequestTime
		TimeRequestUDP.REQUEST = RequestUDP
		TimeRequestUDP.TIME = time.Now().UnixMilli()
		RequestTimes.Store(RequestUDP.Id, TimeRequestUDP)

	}

	PrintRequest(RequestUDP, "SEND")
}

func MaintainConnexion(connUdp *net.UDPConn, ServeurPeer restpeer.PeersUser) {
	for tick := range time.Tick(28 * time.Second) {

		//SendUdpRequest(connUdp, NewRequestUDPExtension(), string(ServeurPeer.ListOfAddresses[0]+":"+ServeurPeer.Port), "MaintainConnexion")
		//if err != nil {
		//	return
		//}
		fmt.Println(tick, "maintien de la connexion avec le serveur")
	}

}

var countReceive int

func PrintRequest(requestUdp RequestUDPExtension, status string) {
	countReceive += 1
	fmt.Println("                 ", status, countReceive)
	fmt.Println("ID :", requestUdp.Id)
	fmt.Println("TYPE :", GetName(requestUdp.Type), "(", requestUdp.Type, ")")
	fmt.Printf("NAME : %.10s %d\n", string(requestUdp.Body), len(requestUdp.Body))
	fmt.Println("LENGTH :", requestUdp.Length)
	fmt.Println("                 ")

}

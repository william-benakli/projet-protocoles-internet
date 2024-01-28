package udppeer

import (
	"math/rand"
	"projet-protocoles-internet/udppeer/cryptographie"
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

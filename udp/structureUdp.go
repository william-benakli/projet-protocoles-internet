package udp

type RequestUDP struct {
	Id         int32
	Type       int8
	Length     int16
	Body       []byte
	Extensions []byte
	Name       []byte
	Signature  []byte
}

func byteToStruct(bytes []byte) RequestUDP {
	udpMsg := RequestUDP{}
	udpMsg.Id += int32(bytes[0])*(1<<24) + (int32(bytes[1]) * (1 << 16)) + (int32(bytes[2]) * (1 << 8)) + int32(bytes[3])
	udpMsg.Type = int8(bytes[4])
	udpMsg.Length = int16(int(bytes[5])*256 + int(bytes[6]))
	udpMsg.Body = make([]byte, udpMsg.Length)

	for i := 0; i < int(udpMsg.Length); i++ {
		udpMsg.Body[i] = bytes[i+7]
	}
	return udpMsg
}

func structToBytes(message RequestUDP) []byte {
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

func SendUdpRequest(RequestUDP RequestUDP) (RequestUDP, error) {
	//structTobytes := structToBytes(RequestUDP)

	/*
		send(structTobytes)
	*/

	//messageStruct := byteToStruct(bytesReponseToStruct)
	//return messageStruct
	udpMsg := RequestUDP{}
	return udpMsg, nil
}

package udppeer

/*
	Commans constants
*/

// Constante des Requetes envoyer
const (
	HelloRequest     uint8 = 2
	PublicKeyRequest uint8 = 3
	RootRequest      uint8 = 4
	GetDatumRequest  uint8 = 5
	HelloReply       uint8 = 129
	PublicKeyReply   uint8 = 130
	RootReply        uint8 = 131
	Datum            uint8 = 132
	NoDatum          uint8 = 133
	NoOp             uint8 = 0
)

func GetRequet(codeCommands uint8, ID int32) RequestUDPExtension {

	name := []byte("Proute")
	lenName := int16(len(name)) + 4

	switch codeCommands {

	case HelloRequest:
		return NewRequestUDPExtension(ID, HelloRequest, lenName, 0, name)
	case HelloReply:
		return NewRequestUDPExtension(ID, HelloReply, lenName, 0, name)

	case PublicKeyRequest:
		return NewRequestUDPExtension(ID, PublicKeyRequest, 0, 0, []byte(""))
	case PublicKeyReply:
		return NewRequestUDPExtension(ID, PublicKeyReply, 0, 0, []byte(""))

	case RootRequest:
		return NewRequestUDPExtension(ID, RootRequest, lenName, 0, name)
	case RootReply:
		return NewRequestUDPExtension(ID, RootReply, lenName, 0, name)

	case GetDatumRequest:
		return NewRequestUDPExtension(ID, GetDatumRequest, lenName, 0, name)

	case Datum:
		return NewRequestUDPExtension(ID, Datum, lenName, 0, name)
	case NoDatum:
		return NewRequestUDPExtension(ID, NoDatum, lenName, 0, name)

	case NoOp:
		return NewRequestUDPExtension(ID, NoOp, lenName, 0, name)

	default:
		return NewRequestUDPExtension(ID, HelloRequest, lenName, 0, name)
	}
}

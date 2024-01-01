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
	Error            uint8 = 1
)

func GetRequet(codeCommands uint8, ID int32) RequestUDPExtension {

	name := "Proute"
	lenName := int16(len(name) + 4)

	namebyte := make([]byte, lenName)

	for i := 0; i < len(namebyte); i++ {
		namebyte[i] = '0'
	}

	copy(namebyte, name)

	switch codeCommands {

	case HelloRequest:

		return NewRequestUDPExtension(ID, HelloRequest, lenName, namebyte)
	case HelloReply:

		return NewRequestUDPExtension(ID, HelloReply, lenName, namebyte)
	case PublicKeyRequest:
		return NewRequestUDPExtension(ID, PublicKeyRequest, 0, []byte(""))
	case PublicKeyReply:
		return NewRequestUDPExtension(ID, PublicKeyReply, 0, []byte(""))
	case RootRequest:
		return NewRequestUDPExtension(ID, RootRequest, lenName, namebyte)
	case RootReply:
		return NewRequestUDPExtension(ID, RootReply, 0, []byte(""))
	case GetDatumRequest:

		return NewRequestUDPExtension(ID, GetDatumRequest, lenName, namebyte)
	case Datum:
		return NewRequestUDPExtension(ID, Datum, lenName, namebyte)
	case NoDatum:
		return NewRequestUDPExtension(ID, NoDatum, lenName, namebyte)

	case NoOp:
		return NewRequestUDPExtension(ID, NoOp, lenName, namebyte)

	default:
		return NewRequestUDPExtension(ID, HelloRequest, lenName, namebyte)
	}
}

func GetName(codeCommands uint8) string {

	switch codeCommands {
	case HelloRequest:
		return "HelloRequest"
	case HelloReply:
		return "HelloReply"
	case PublicKeyRequest:

		return "PublicKeyRequest"
	case PublicKeyReply:

		return "PublicKeyReply"
	case RootRequest:

		return "RootRequest"
	case RootReply:

		return "RootReply"
	case GetDatumRequest:

		return "GetDatumRequest"
	case Datum:

		return "Datum"
	case NoDatum:

		return "NoDatum"
	case NoOp:
		return "NoOp"
	case Error:
		return "Error"
	default:

		return "UNKNOW"
	}
}

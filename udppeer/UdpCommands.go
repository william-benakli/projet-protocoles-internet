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
	ErrorReply       uint8 = 128
)

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

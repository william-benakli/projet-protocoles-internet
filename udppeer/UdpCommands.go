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
)

func GetRequet(codeCommands uint8, ID int32) RequestUDPExtension {
	switch codeCommands {
	// Requete d'envoie
	case HelloRequest:
		return NewRequestUDPExtension(ID, HelloRequest, int16(len("BOulangerPatissierEtFiereDeLetre")), 0, []byte("BOulangerPatissierEtFiereDeLetre"))
	case PublicKeyRequest:
		return NewRequestUDPExtension(ID, PublicKeyRequest, int16(len("BOulangerPatissierEtFiereDeLetre")), 0, []byte("BOulangerPatissierEtFiereDeLetre"))
	case RootRequest:
		return NewRequestUDPExtension(ID, RootRequest, int16(len("BOulangerPatissierEtFiereDeLetre")), 0, []byte("BOulangerPatissierEtFiereDeLetre"))
	case GetDatumRequest:
		return NewRequestUDPExtension(ID, GetDatumRequest, int16(len("BOulangerPatissierEtFiereDeLetre")), 0, []byte("BOulangerPatissierEtFiereDeLetre"))
		// Requete reponse
	case HelloReply:
		return NewRequestUDPExtension(ID, HelloReply, int16(len("BOulangerPatissierEtFiereDeLetre")), 0, []byte("BOulangerPatissierEtFiereDeLetre"))
	case PublicKeyReply:
		return NewRequestUDPExtension(ID, PublicKeyReply, int16(len("BOulangerPatissierEtFiereDeLetre")), 0, []byte("BOulangerPatissierEtFiereDeLetre"))
	case RootReply:
		return NewRequestUDPExtension(ID, RootReply, int16(len("BOulangerPatissierEtFiereDeLetre")), 0, []byte("BOulangerPatissierEtFiereDeLetre"))
	case Datum:
		return NewRequestUDPExtension(ID, Datum, int16(len("BOulangerPatissierEtFiereDeLetre")), 0, []byte("BOulangerPatissierEtFiereDeLetre"))
	case NoDatum:
		return NewRequestUDPExtension(ID, NoDatum, int16(len("BOulangerPatissierEtFiereDeLetre")), 0, []byte("BOulangerPatissierEtFiereDeLetre"))
	default:
		return NewRequestUDPExtension(ID, HelloRequest, int16(len("BOulangerPatissierEtFiereDeLetre")), 0, []byte("BOulangerPatissierEtFiereDeLetre"))
	}
}

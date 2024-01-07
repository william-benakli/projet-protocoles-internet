package Tools

import (
	"math/rand"
	"net/http"
)

/* EDIT VALUE */

var ClientRestAPI *http.Client
var ConnUDP = SetListen(rand.Intn(45000) + 3000)

var Name = "MOOON"
var IP_ADRESS = "81.194.27.155:8443"
var IP_ADRESS_SEND = "81.194.27.155:8443"

var WantBigFile bool = true

/* EDIT VALUE */

/* REQUEST */
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

/* NOEUD */
const (
	DirectoryType = 2
	BigFileType   = 1
	ChunkType     = 0
)

/* NOEUD */

/* CONSTANT VALUE */

const REMISSION = 3
const TempsRemissionMiliSeconde = 2
const ChunkSize = 1024

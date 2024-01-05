package utils

import "net/http"

const NameUser = "    BoulangerPatissierEtFiereDeLetre"

var Client *http.Client

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

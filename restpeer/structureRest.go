package restpeer

type ListOfPeers struct {
	ListOfPeers []PeersUser
	Length      int16
}

type PeersUser struct {
	NameUser    string
	AddressIpv6 string
	AddressIpv4 string
	Port        string
}

package restpeer

type ListOfPeers struct {
	listOfPeers []PeersUser
	Length      int16
}

type PeersUser struct {
	nameUser    string
	addressIpv6 string
	addressIpv4 string
	port        string
}

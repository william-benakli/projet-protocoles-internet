package restpeer

type ListOfPeers struct {
	ListOfPeers []PeersUser
	Length      int16
}

type PeersUser struct {
	NameUser        string
	NameLen         int16
	ListOfAddresses []string
}

package restpeer

import (
	"fmt"
	"io"
	"log"
	"net/http"
	. "projet-protocoles-internet/Tools"
	"projet-protocoles-internet/udppeer/cryptographie"
	"strings"
)

/* Private fonction */
func getPeerStructFromStringTab(name string, userPeer []string) PeersUser {
	var peer PeersUser

	for i := 0; i < len(userPeer); i++ {
		if strings.HasPrefix(userPeer[i], "[") {
			port := strings.Split(userPeer[i], ":")[7]
			//IPV6
			addressIp := strings.Replace(strings.Replace(strings.Split(userPeer[i], "]:")[0], "[", "", -1), "]", "", -1) + ":" + port
			peer.ListOfAddresses = append(peer.ListOfAddresses, addressIp)
		} else {
			//IPV4
			addressIp := userPeer[i]
			peer.ListOfAddresses = append(peer.ListOfAddresses, addressIp)
		}
	}
	peer.NameUser = name
	return peer

}

/* public fonction */

func GetListOfPeers(client *http.Client, peersTableau []string) ListOfPeers {
	var listOfPeers ListOfPeers
	for i := 0; i < len(peersTableau); i++ {
		infoPeers := SendRestPeerAdresses(client, peersTableau[i])
		if !strings.Contains(infoPeers[0], "404") {
			listOfPeers.ListOfPeers = append(listOfPeers.ListOfPeers, getPeerStructFromStringTab(peersTableau[i], infoPeers))
		} else {
			fmt.Println(peersTableau[i], " Auncune ip")
		}
	}

	return listOfPeers
}

func GetRestPeerNames(client *http.Client) []string {
	resp, err := client.Get("https://jch.irif.fr:8443/peers/")

	if err != nil {

		log.Fatal("client fail to get peer Names ")
	}
	if resp.Body != nil {
		body, readIo := io.ReadAll(resp.Body)
		if readIo != nil {
			log.Fatal("io failed")
		}
		/* Gerer le cas d'erreur */
		return strings.Split(string(body), "\n")
	} else {
		fmt.Println("Erreur aucune nom de peers ")
		return nil
	}
}

func GetMasterAddresse(client *http.Client, url string) (PeersUser, error) {
	resp, err := client.Get(url)

	if err != nil {
		log.Fatal("client fail to get peer Names ")
	}
	body, readIo := io.ReadAll(resp.Body)
	if readIo != nil {
		log.Fatal("io failed")
	}

	var pair PeersUser
	ipv4Ipv6 := strings.Split(string(body), "\n")

	for i := 0; i < len(ipv4Ipv6); i++ {
		if strings.HasPrefix(ipv4Ipv6[i], "[") {
			port := strings.Split(ipv4Ipv6[i], ":")[7]
			//IPV6
			addressIp := strings.Replace(strings.Replace(strings.Split(ipv4Ipv6[i], "]:")[0], "[", "", -1), "]", "", -1) + ":" + port
			pair.ListOfAddresses = append(pair.ListOfAddresses, addressIp)
		} else {
			//IPV4
			addressIp := ipv4Ipv6[i]
			pair.ListOfAddresses = append(pair.ListOfAddresses, addressIp)
		}
	}

	return pair, err
}

func SendRestPeerAdresses(client *http.Client, namePeer string) []string {
	resp, err := client.Get("https://jch.irif.fr:8443/peers/" + RemoveEmpty(namePeer) + "/addresses")

	if err != nil {
		log.Fatal("client fail to get peer Names ")
	}
	if resp.Body != nil {
		body, readIo := io.ReadAll(resp.Body)
		if readIo != nil {
			log.Fatal("io failed")
		}

		/* Gerer le cas d'erreur */
		return strings.Split(string(body), "\n")
	} else {
		fmt.Println("Erreur aucune peers presente")
		return nil
	}
}

func GetAdrFromNamePeers(userName string) string {

	listPeers := GetListOfPeers(ClientRestAPI, GetRestPeerNames(ClientRestAPI))

	for i := range listPeers.ListOfPeers {
		if listPeers.ListOfPeers[i].NameUser == userName {
			return listPeers.ListOfPeers[i].ListOfAddresses[0]
		}
	}
	user, _ := GetMasterAddresse(ClientRestAPI, "https://jch.irif.fr:8443/peers/jch.irif.fr/addresses")
	return user.ListOfAddresses[0]
}

func GetPublicKey(client *http.Client, name string) int {
	resp, err := client.Get("https://jch.irif.fr:8443/peers/" + RemoveEmpty(name) + "/key")

	if err != nil {
		log.Fatal("getPublicKey")
	}
	if resp.Body != nil {
		if resp.StatusCode == 200 {
			body, readIo := io.ReadAll(resp.Body)
			if readIo != nil {
				log.Fatal("io failed")
			}
			cryptographie.OtherPublicKey = cryptographie.UnFormateKey(body)
			return 200
		} else if resp.StatusCode == 204 {
			fmt.Println("no key")
			return 204
		} else {
			fmt.Println("pair inconnu")
			return 404
		}
	} else {
		fmt.Println("Erreur aucune nom de peers ")
	}
	return 404
}

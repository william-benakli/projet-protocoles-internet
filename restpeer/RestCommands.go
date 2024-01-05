package restpeer

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"projet-protocoles-internet/udppeer/cryptographie"
	"strings"
)

/* Private fonction */
func getPeerStructFromStringTab(name string, userPeer []string) PeersUser {
	var peer PeersUser
	peer.Port = strings.Split(userPeer[0], ":")[1]

	for i := 0; i < len(userPeer); i++ {
		if strings.HasPrefix(userPeer[i], "[") {
			//IPV6
			addressIp := strings.Replace(strings.Replace(strings.Split(userPeer[i], "]:")[0], "[", "", -1), "]", "", -1)
			peer.ListOfAddresses = append(peer.ListOfAddresses, addressIp)
		} else {
			//IPV4
			addressIp := strings.Split(userPeer[i], ":")[0]
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
		if len(peersTableau[i]) < 16 && !strings.ContainsAny(peersTableau[i], " ") {

			infoPeers := SendRestPeerAdresses(client, peersTableau[i])
			if len(infoPeers) > 1 && !strings.Contains(infoPeers[0], "404") {
				fmt.Println(infoPeers[0], strings.Contains(infoPeers[0], "404 page not found "), " -------------- ")
				listOfPeers.ListOfPeers = append(listOfPeers.ListOfPeers, getPeerStructFromStringTab(peersTableau[i], infoPeers))
			} else {
				fmt.Println(peersTableau[i], " a été ignoré car il a aucune ip associé")
			}
		} else {
			fmt.Println(peersTableau[i], " a été ignoré car il comporte des espaces ou est trop long")
		}
	}
	return listOfPeers
}

func GetPublicKey(client *http.Client, name string) {
	resp, err := client.Get("https://jch.irif.fr:8443/peers/" + name + "/key")

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
		} else if resp.StatusCode == 204 {
			fmt.Println("no key")
		} else {
			fmt.Println("pair inconnu")
		}
	} else {
		fmt.Println("Erreur aucune nom de peers ")
	}
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
	pair.Port = strings.Split(ipv4Ipv6[0], ":")[1]

	for i := 0; i < len(ipv4Ipv6); i++ {
		if strings.HasPrefix(ipv4Ipv6[i], "[") {
			//IPV6
			addressIp := strings.Replace(strings.Replace(strings.Split(ipv4Ipv6[i], "]:")[0], "[", "", -1), "]", "", -1)
			pair.ListOfAddresses = append(pair.ListOfAddresses, addressIp)
		} else {
			//IPV4
			addressIp := strings.Split(ipv4Ipv6[i], ":")[0]
			pair.ListOfAddresses = append(pair.ListOfAddresses, addressIp)
		}
	}

	return pair, err
}

func SendRestPeerAdresses(client *http.Client, namePeer string) []string {
	resp, err := client.Get("https://jch.irif.fr:8443/peers/" + namePeer + "/addresses")

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

// TODO: à finir
func GetAdrFromNamePeers(userName []byte) string {
	var user PeersUser

	return user.ListOfAddresses[0] + ":" + user.Port
}

package restpeer

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func getPeerStructFromStringTab(name string, userPeer []string) PeersUser {
	var peer PeersUser
	addressIpv4 := strings.Split(userPeer[0], ":")[0]
	addressIpv6 := strings.Replace(strings.Replace(strings.Split(userPeer[1], ":")[0], "[", "", -1), "]", "", -1)
	port := strings.Split(userPeer[0], ":")[1]
	fmt.Println(userPeer[0], " aaaa ")
	peer.AddressIpv4 = addressIpv4
	peer.AddressIpv6 = addressIpv6
	peer.Port = port
	peer.NameUser = name
	return peer

}

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

/* Rendre ça générique */

func SendRestPeerNames(client *http.Client) []string {
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
	pair.AddressIpv4 = strings.Split(ipv4Ipv6[0], ":")[0]
	adrv6 := strings.Split(ipv4Ipv6[1], "]")[0]
	pair.AddressIpv6 = strings.ReplaceAll(adrv6, "[", "")

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

func GetAdrFromNamePeers(userName []byte) string {
	var user PeersUser

	return user.AddressIpv4 + ":" + user.Port
}

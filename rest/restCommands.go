package rest

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

/* Rendre ça générique */

func sendRequestPeerNames(client http.Client) {
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
		fmt.Print(string(body))
	}
}

func sendRequestPeerAdresses(client http.Client, namePeer string) {
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
		fmt.Print(string(body))
	}
}

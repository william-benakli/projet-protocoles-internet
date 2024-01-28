package udppeer

import (
	"fmt"
	"golang.org/x/exp/rand"
	"net"
	"projet-protocoles-internet/restpeer"
	"sync"
	"time"
)

import . "projet-protocoles-internet/Tools"

type RequestTime struct {
	TIME    int64
	REQUEST RequestUDPExtension //	time.Now().UnixMilli()
}

var RequestTimes sync.Map
var listIdDejaVu []int32        // Evite les remissions
var listOfIDFromRequest []int32 //s'assure que la reponse envoyer fait suite à une demande
var Tentative map[int32]int     // S'assurer que si une pair crash on arrete de lui parler au bout de n tentatives

func init() {
	Tentative = make(map[int32]int)
}

func SendUDPPacketFromResponse(connUdp *net.UDPConn, channel chan RequestUDPExtension) {

	for {

		receiveStruct, ok := <-channel

		PrintRequest(receiveStruct, " RECEIVED ")

		if !ok {
			fmt.Println("Channel closed. Exiting receiver.")
			return
		}

		//Permet de ignorer les rémissions
		if containedList(listIdDejaVu, receiveStruct.Id) {
			continue
		}

		listIdDejaVu = append(listIdDejaVu, receiveStruct.Id)

		if receiveStruct.Type < 128 {
			receiveRequest(connUdp, receiveStruct)
		} else {
			//List qui verifie que la reponse fait suite à une demande
			if containedList(listOfIDFromRequest, receiveStruct.Id) {
				go receiveResponse(connUdp, receiveStruct)
			}
		}
	}

}

func containedList(listId []int32, id int32) bool {
	for i := 0; i < len(listId); i++ {
		if listId[i] == id {
			return true
		}
	}
	return false
}

func RemissionPaquets(connUdp *net.UDPConn, adressPort string) {
	for _ = range time.Tick(5 * time.Second) {
		RequestTimes.Range(func(key, value interface{}) bool {

			if requestTime, ok := value.(RequestTime); ok {
				if (time.Now().UnixMilli() - requestTime.TIME) > 7000 {
					RequestTimes.Delete(key)
					if requestTime.REQUEST.Type > 0 {

						if Tentative[requestTime.REQUEST.Id] <= MaxTentative {
							//On envoie avec la meme id la requete
							requestDatum := NewRequestUDPExtension(requestTime.REQUEST.Id, requestTime.REQUEST.Type, int16(len(requestTime.REQUEST.Body)), requestTime.REQUEST.Body)
							SendUdpRequest(connUdp, requestDatum, IP_ADRESS, "remissionPaquets ")
							time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
							Tentative[requestTime.REQUEST.Id] = Tentative[requestTime.REQUEST.Id] + 1
						} else {
							PrintDebug("Max tentative atteint pour le paquet suivant :")
							PrintRequest(requestTime.REQUEST, "Remission de paquet")
							PrintDebug("-----------------")
						}

					}
				}
			}
			return true
		})
	}
}

var ReceiveCounter int = 0

func SendUdpRequest(connUdp *net.UDPConn, RequestUDP RequestUDPExtension, adressPort string, from string) {
	globalID += 1
	structToBytes := StructToBytes(RequestUDP)
	udpAddr, _ := net.ResolveUDPAddr("udp", adressPort)

	time.Sleep(time.Millisecond * 50)

	_, _ = connUdp.WriteToUDP(structToBytes, udpAddr)

	//On ajoute l'id aux id demander
	listOfIDFromRequest = append(listOfIDFromRequest, RequestUDP.Id)

	if RequestUDP.Type < 128 && RequestUDP.Type != 0 && RequestUDP.Type != 1 {
		var TimeRequestUDP RequestTime
		TimeRequestUDP.REQUEST = RequestUDP
		TimeRequestUDP.TIME = time.Now().UnixMilli()
		RequestTimes.Store(RequestUDP.Id, TimeRequestUDP)
	}

	PrintRequest(RequestUDP, "SEND: "+from)
}

func MaintainConnexion(connUdp *net.UDPConn, ServeurPeer restpeer.PeersUser) {
	for tick := range time.Tick(30 * time.Second) {
		byteName := make([]byte, 4)
		byteName[0] = 0
		byteName[1] = 0
		byteName[2] = 0
		byteName[3] = 0
		byteName = append(byteName, []byte(Name)...)
		SendUdpRequest(connUdp, NewRequestUDPExtensionSigned(GetGlobalID(), HelloRequest, int16(len(byteName)), byteName), ServeurPeer.ListOfAddresses[0], "MaintainConnexion")
		fmt.Println(tick, "maintien de la connexion avec le serveur")
	}
}

func PrintRequest(requestUdp RequestUDPExtension, status string) {
	if DebugPrint == true {
		fmt.Println("                 ", status)
		fmt.Println("ID :", requestUdp.Id)
		fmt.Println("TYPE :", GetName(requestUdp.Type), "(", requestUdp.Type, ")")
		fmt.Printf("BODY : %.20s %d\n", requestUdp.Body, len(requestUdp.Body))
		fmt.Println("LENGTH :", requestUdp.Length)
		fmt.Println("                 ")
	}
}

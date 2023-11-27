package projet_protocoles_internet

import (
	"fmt"
	"udp"
)

func main() {
	SendHello()
	listpeer := udp.SendRestListPeers()

	/*
		 Ecouter udp
			swtich sur l'element
	*/
	fmt.Println("Fin du programme")
}

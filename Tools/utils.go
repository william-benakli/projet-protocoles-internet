package Tools

import (
	"fmt"
	"net"
	"os"
	"strings"
)

/* DEBUG */

var debugPrint bool = true

func PrintDebug(messages ...any) {
	if debugPrint {
		fmt.Println("# DEBUG # ")
		for message := range messages {
			if &message != nil {
				fmt.Println(messages[message])
			}
		}
	}
}

func HideDebug() {
	debugPrint = false
}

func ShowDebug() {
	debugPrint = true
}

func SetListen(port int) *net.UDPConn {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: port})
	if err != nil {
		println("Listen failed:", err.Error())
		os.Exit(1)
	}

	return conn
}

func RemoveEmpty(stringBody string) string {
	nullIndex := strings.IndexByte(stringBody, '\000')
	if nullIndex == -1 {
		return stringBody
	}
	return stringBody[:nullIndex]
}

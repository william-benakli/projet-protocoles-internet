package Tools

import (
	"fmt"
	"net"
	"os"
	"strings"
)

/* DEBUG */

var DebugPrint bool = false

func PrintDebug(messages ...any) {
	if DebugPrint {
		fmt.Println("# DEBUG # ")
		for message := range messages {
			if &message != nil {
				fmt.Println(messages[message])
			}
		}
	}
}

func HideDebug() {
	DebugPrint = false
}

func ShowDebug() {
	DebugPrint = true
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

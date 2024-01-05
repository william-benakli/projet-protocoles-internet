package Tools

import "fmt"

/* DEBUG */

var debugPrint bool = false

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

package webUI

import (
	"fmt"
	"log"
	"net/http"
)

func SetupPage(client *http.Client) {
	http.HandleFunc("/", peersPage)
	err := http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", nil)
	log.Fatal("ListenAndServe: ", err)
}

func peersPage(w http.ResponseWriter, r *http.Request) {
	const data = `<!DOCTYPE html> <html> <head></head><body> <p>Bonjour !</p> </body></html>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Method != "HEAD" && r.Method != "GET" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, data)

}

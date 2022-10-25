// The executable program need in main package.
package main

import (
	// for Fprintf.
	"fmt"
	// for ResponseWriter, Request, HandleFunc, ListenAndServe.
	"net/http"
)

// Handler function.
func handler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Hello World, %s!", request.URL.Path[1:])
}

// main function.
func main() {
	// Subscribe handler.
	http.HandleFunc("/", handler)
	// Listen for port 443.
	http.ListenAndServeTLS("", "/usr/local/bin/go/https2/certificate.crt", "/usr/local/bin/go/https2/private.key", nil)
}

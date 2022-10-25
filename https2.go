// The executable program need in main package.
package main

import (
	// for Fprintf.
	"fmt"
	// for ResponseWriter, Request, HandleFunc, ListenAndServe.
	"net/http"
	// for FuncForPC
	"reflect"
	"runtime"
	"github.com/julienschmidt/httprouter"
)

//func hello(w http.ResponseWriter, r *http.Request) {
func hello(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", p.ByName("name"))
	//fmt.Fprintf(w, "Hello")
}

func world(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "World")
}

func log(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
		fmt.Println("Handler function called - " + name)
		h(w,r)
	}
}

func nothing(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("nothing")
		h(w,r)
	}
}

// Handler function.
//type MyHandler struct{}
//
//func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	//fmt.Fprintf(w, "Hello World, %s!", r.URL.Path[1:])
//	fmt.Fprintf(w, "Hello World!")
//}

// main function.
func main() {
	// Subscribe handler.
	// use multiplexer
	mux := httprouter.New()
	mux.GET("/hello/:name", hello)

	server := http.Server{
		Addr: "",
		Handler: mux,
	}

	//http.HandleFunc("/hello", nothing(log(hello)))
	//http.HandleFunc("/world", log(world))

	// Listen for port 443.
	server.ListenAndServeTLS("/usr/local/bin/go/https2/certificate.crt", "/usr/local/bin/go/https2/private.key")

//	// no use multiplexer
//	handler := MyHandler{}
//	server := http.Server{
//		Addr:    "",
//		Handler: &handler,
//	}
//	server.ListenAndServeTLS("/usr/local/bin/go/https2/certificate.crt", "/usr/local/bin/go/https2/private.key")
}

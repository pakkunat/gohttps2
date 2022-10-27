// The executable program need in main package.
package main

import (
	// for Fprintf.
	"fmt"
	// for ResponseWriter, Request, HandleFunc, ListenAndServe.
	"net/http"
	// for FuncForPC
	//"reflect"
	//"runtime"
	// for ReadAll
	"io/ioutil"
	// for Marshal
	"encoding/json"
	// for Unix
	"time"
	// for EncodeString, DecodeString
	"encoding/base64"
)

type Post struct {
	User    string
	Threads []string
}

func headers(w http.ResponseWriter, r *http.Request) {
	h := r.Header
	fmt.Fprintln(w, h)
}

func body(w http.ResponseWriter, r *http.Request) {
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	fmt.Fprintln(w, string(body))
}

func process(w http.ResponseWriter, r *http.Request) {
	//r.ParseForm()
	r.ParseMultipartForm(1024)
	fileHeader := r.MultipartForm.File["uploaded"][0]
	file, err := fileHeader.Open()
	if err == nil {
		data, err := ioutil.ReadAll(file)
		if err == nil {
			fmt.Fprintln(w, string(data))
		}
	}
	//fmt.Fprintln(w, "(1)", r.FormValue("hello"))
	//fmt.Fprintln(w, "(2)", r.PostFormValue("hello"))
	//fmt.Fprintln(w, "(3)", r.PostForm)
	//fmt.Fprintln(w, "(4)", r.MultipartForm)
}

func writeExample(w http.ResponseWriter, r *http.Request) {
	str := `<html>
<head><title>Go Web Programming</title></head>
<body><h1>Hello World</h1></body>
</html>`
	w.Write([]byte(str))
}

func writeHeaderExample(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(501)
	fmt.Fprintln(w, "NOT exist Such a service")
}

func headerExample(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Location", "https://www.google.com")
	w.WriteHeader(302)
}

func jsonExample(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	post := &Post{
		User:    "higaki koji",
		Threads: []string{"first", "second", "third"},
	}
	json, _ := json.Marshal(post)
	w.Write(json)
}

func setCookie(w http.ResponseWriter, r *http.Request) {
	c1 := http.Cookie{
		Name:     "first_cookie",
		Value:    "Go Web Programming",
		HttpOnly: true,
	}
	c2 := http.Cookie{
		Name:     "second_cookie",
		Value:    "Manning Publications Co",
		HttpOnly: true,
	}

	http.SetCookie(w, &c1)
	http.SetCookie(w, &c2)
}

func getCookie(w http.ResponseWriter, r *http.Request) {
	c1, err := r.Cookie("first_cookie")
	if err != nil {
		fmt.Fprintln(w, "Cannot get the first cookie")
	}
	cs := r.Cookies()
	fmt.Fprintln(w, c1)
	fmt.Fprintln(w, cs)
}

func setMessage(w http.ResponseWriter, r *http.Request) {
	msg := []byte("Hello World")
	c := http.Cookie{
		Name:  "flash",
		Value: base64.URLEncoding.EncodeToString(msg),
	}
	http.SetCookie(w, &c)
}

func showMessage(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("flash")
	if err != nil {
		if err == http.ErrNoCookie {
			fmt.Fprintln(w, "No messsages")
		}
	} else {
		rc := http.Cookie{
			Name:    "flash",
			MaxAge:  -1,
			Expires: time.Unix(1, 0),
		}
		http.SetCookie(w, &rc)
		val, _ := base64.URLEncoding.DecodeString(c.Value)
		fmt.Fprintln(w, string(val))
	}
}

//func log(h http.HandlerFunc) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		name := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
//		fmt.Println("Handler function called - " + name)
//		h(w,r)
//	}
//}

// main function.
func main() {
	// Subscribe handler.
	server := http.Server{
		Addr: "",
	}

	http.HandleFunc("/headers", headers)
	http.HandleFunc("/body", body)
	http.HandleFunc("/process", process)
	http.HandleFunc("/write", writeExample)
	http.HandleFunc("/writeheader", writeHeaderExample)
	http.HandleFunc("/redirect", headerExample)
	http.HandleFunc("/json", jsonExample)
	http.HandleFunc("/set_cookie", setCookie)
	http.HandleFunc("/get_cookie", getCookie)
	http.HandleFunc("/set_message", setMessage)
	http.HandleFunc("/show_message", showMessage)

	// Listen for port 443.
	server.ListenAndServeTLS("/usr/local/bin/go/https2/certificate.crt", "/usr/local/bin/go/https2/private.key")
}

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
	// for template.ParseFiles
	"html/template"
	// for rand.Seed, rand.Intn
	"math/rand"
	// for os.Create
	"os"
	// for csv.NewWriter, csv.NewReader, Write, ReadAll
	"encoding/csv"
	// for strconv.Itoa, strconv.ParseInt
	"strconv"
	// for bytes.NewBuffer
	"bytes"
	// for gob.NewEncoder, gob.NewDecoder
	"encoding/gob"
	// for sql.Open, Query, QueryRow
	"database/sql"
	// for
	_ "github.com/lib/pq"
)

type Post struct {
	User    string
	Threads []string
}

type Data struct {
	Id      int
	Content string
	Author  string
}

var PostById map[int]*Data
var PostsByAuthor map[string][]*Data
var Db *sql.DB

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

func store(post Data) {
	PostById[post.Id] = &post
	PostsByAuthor[post.Author] = append(PostsByAuthor[post.Author], &post)
}

//func log(h http.HandlerFunc) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		name := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
//		fmt.Println("Handler function called - " + name)
//		h(w,r)
//	}
//}

func tmpl(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("tmpl.html")
	rand.Seed(time.Now().Unix())
	t.Execute(w, rand.Intn(10) > 5)
}

func iterator(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("tmpl.html")
	dayOfWeek := []string{"mon", "tue", "wed", "thu", "fri", "sat", "sun"}
	t.Execute(w, dayOfWeek)
}

func assign(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("tmpl.html")
	t.Execute(w, "hello")
}

func include(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("t1.html", "t2.html")
	t.Execute(w, "Hello world!")
}

func layout(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().Unix())
	var t *template.Template
	if rand.Intn(10) > 5 {
		t, _ = template.ParseFiles("layout.html", "red_hello.html")
	} else {
		t, _ = template.ParseFiles("layout.html")
		//t, _ = template.ParseFiles("layout.html", "blue_hello.html")
	}
	t.ExecuteTemplate(w, "layout", "")
}

func newStore(data interface{}, filename string) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(data)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(filename, buffer.Bytes(), 0600)
	if err != nil {
		panic(err)
	}
}

func load(data interface{}, filename string) {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	buffer := bytes.NewBuffer(raw)
	dec := gob.NewDecoder(buffer)
	err = dec.Decode(data)
	if err != nil {
		panic(err)
	}
}

func init() {
	var err error
	Db, err = sql.Open("postgres", "user=gwp dbname=gwp password=gwp sslmode=disable")
	if err != nil {
		fmt.Println("db init error")
		return
	}
}

func Posts(limit int) (posts []Data, err error) {
	rows, err := Db.Query("select id, content, author from posts limit $1", limit)
	if err != nil {
		return
	}
	for rows.Next() {
		post := Data{}
		err = rows.Scan(&post.Id, &post.Content, &post.Author)
		if err != nil {
			return
		}
		posts = append(posts, post)
	}
	rows.Close()
	return
}

func GetPost(id int) (post Data, err error) {
	post = Data{}
	err = Db.QueryRow("select id, content, author from posts where id = $1", id).Scan(&post.Id, &post.Content, &post.Author)
	if err != nil {
		panic(err)
	}
	return
}

func (post *Data) Create() (err error) {
	statement := "insert into posts (content, author) values ($1, $2) returning id"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		fmt.Println("error Create: Db.Prepare")
		panic(err)
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(post.Content, post.Author).Scan(&post.Id)
	if err != nil {
		fmt.Println("error Create: stmt.QueryRow")
		panic(err)
	}
	return
}

func (post *Data) Update() (err error) {
	_, err = Db.Exec("update posts set content = $2, author = $3 where id = $1", post.Id, post.Content, post.Author)
	return
}

func (post *Data) Delete() (err error) {
	_, err = Db.Exec("delete from posts where id = $1", post.Id)
	return
}

// main function.
func main() {
	post := Data{Content: "Hello World!", Author: "higaki koji"}

	fmt.Println(post)
	post.Create()
	fmt.Println(post)

	readPost, _ := GetPost(post.Id)
	fmt.Println(readPost)

	readPost.Content = "Bonjour Monde!"
	readPost.Author = "Pierre"
	readPost.Update()

	postss, _ := Posts(10)
	fmt.Println(postss)

	readPost.Delete()

	post = Data{Id: 1, Content: "Hello World!", Author: "higaki koji"}
	newStore(post, "post1")
	var postRead Data
	load(&postRead, "post1")
	//fmt.Println(postRead)

	csvFile, err := os.Create("posts.csv")
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	allPosts := []Data{
		Data{Id: 1, Content: "Hello World!", Author: "higaki koji"},
		Data{Id: 2, Content: "Bonjour Mondel", Author: "Pierre"},
		Data{Id: 3, Content: "Hola Mundo!", Author: "Pedro"},
		Data{Id: 4, Content: "Greetings Earthings!", Author: "higaki koji"},
	}

	writer := csv.NewWriter(csvFile)
	for _, post := range allPosts {
		line := []string{strconv.Itoa(post.Id), post.Content, post.Author}
		err := writer.Write(line)
		if err != nil {
			panic(err)
		}
	}
	writer.Flush()

	file, err := os.Open("posts.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	record, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	var posts []Data
	for _, item := range record {
		id, _ := strconv.ParseInt(item[0], 0, 0)
		post := Data{Id: int(id), Content: item[1], Author: item[2]}
		posts = append(posts, post)
	}
	//fmt.Println(posts[0].Id)
	//fmt.Println(posts[0].Content)
	//fmt.Println(posts[0].Author)

	data := []byte("Hello World!\n")
	err = ioutil.WriteFile("data2", data, 0644)
	if err != nil {
		panic(err)
	}
	read1, _ := ioutil.ReadFile("data1")
	fmt.Print(string(read1))
	file1, _ := os.Create("data2")
	defer file1.Close()

	bytes, _ := file1.Write(data)
	fmt.Printf("Wrote %d bytes to file\n", bytes)

	file2, _ := os.Open("data2")
	defer file2.Close()
	read2 := make([]byte, len(data))
	bytes, _ = file2.Read(read2)
	//fmt.Printf("Read %d bytes from file\n", bytes)
	//fmt.Println(string(read2))

	PostById = make(map[int]*Data)
	PostsByAuthor = make(map[string][]*Data)

	post1 := Data{Id: 1, Content: "Hello World!", Author: "higaki koji"}
	post2 := Data{Id: 2, Content: "Bonjour Mondel", Author: "Pierre"}
	post3 := Data{Id: 3, Content: "Hola Mundo!", Author: "Pedro"}
	post4 := Data{Id: 4, Content: "Greetings Earthings!", Author: "higaki koji"}

	store(post1)
	store(post2)
	store(post3)
	store(post4)

	//fmt.Println(PostById[1])
	//fmt.Println(PostById[2])

	for _, post := range PostsByAuthor["higaki koji"] {
		fmt.Println(post)
	}
	for _, post := range PostsByAuthor["Pedro"] {
		fmt.Println(post)
	}

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
	http.HandleFunc("/template", tmpl)
	http.HandleFunc("/iterator", iterator)
	http.HandleFunc("/assign", assign)
	http.HandleFunc("/include", include)
	http.HandleFunc("/layout", layout)

	// Listen for port 443.
	server.ListenAndServeTLS("./certificate.crt", "private.key")
}

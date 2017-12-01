package main

import (
	"net/http"
	"time"
	//"fmt"
	"html/template"
	"./dbwork"
	"database/sql"
	//"./requests"
	"github.com/davecgh/go-spew/spew"
	"reflect"
)

var Server = &http.Server {
	Addr:	":80",
	ReadTimeout: 10 * time.Second,
	WriteTimeout: 10 * time.Second,
}

func main() {
	http.HandleFunc("/", HandleIndex)
	http.HandleFunc("/login", Login)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.HandleFunc("/select", SelectRequest)

	err := Server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func SelectRequest(w http.ResponseWriter, r *http.Request) {
	PrepareArgs(r)
	DB := dbwork.Connect()
	MC := &MainController{ DataBase : DB, UsrGroup : "User" }
	out := reflect.ValueOf(MC).MethodByName("Req1").Call([]reflect.Value{reflect.ValueOf(r)})
	spew.Dump(out)
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	    expiration := time.Now().Add(365 * 24 * time.Hour)
    cookie := http.Cookie{Name: "username", Value: "astaxie", Expires: expiration}
     http.SetCookie(w, &cookie)
	t := template.Must(template.ParseFiles("public/index.html"))
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, "Hello world!")
}

func Login(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("public/index.html"))
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, "Hello world!")
}

func PrepareArgs(r *http.Request) {
	//res := requests.Req1(DB)
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
}

type MainController struct {
	DataBase    *sql.DB
	UsrGroup    string
}
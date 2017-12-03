package main

import (
	"net/http"
	"time"
	"fmt"
	"html/template"
	"./dbwork"
	"database/sql"
	//"./requests"
	//"github.com/davecgh/go-spew/spew"
	"reflect"
	"encoding/json"
)

var Server = &http.Server {
	Addr:	":80",
	ReadTimeout: 10 * time.Second,
	WriteTimeout: 10 * time.Second,
}

func main() {
	http.HandleFunc("/", HandleIndex)
	http.HandleFunc("/index", HandleIndex)
	http.HandleFunc("/login", Login)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.HandleFunc("/select", SelectRequest)

	err := Server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func SelectRequest(w http.ResponseWriter, r *http.Request) {
	PrepareArgs(r)			//подготовка аргументо из запроса
	DB := dbwork.Connect()	//коннект к базе, до этого нужно сделать авторизацию
	MC := &MainController{ DataBase : DB, UsrGroup : "User" }	//создаю объект с контроллером, в котором хранятся коннект к бд и пользователь, опцианально, мб уберу
	out := reflect.ValueOf(MC).MethodByName(r.URL.Query().Get("action")).Call([]reflect.Value{reflect.ValueOf(r)})	//выполнение самого запроса, пришлось немного с рефлектом поебаца
	if out[1].IsNil() {
		fmt.Fprintf(w, "no data")
		return
	}
	tmplFile := Templates[r.URL.Query().Get("action")]
	t := template.Must(template.ParseFiles("public/" + tmplFile))
	output := PrepareForOut(out[1])
	if err := t.Execute(w, output); err != nil {
		panic(err)
	}
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

func PrepareForOut(res reflect.Value) []map[string]interface{} {
	outData := make([]map[string]interface{}, 0)
	switch reflect.TypeOf(res.Interface()).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(res.Interface())
		for i := 0; i < s.Len(); i++ {
			m := map[string]interface{}{}
			data, err := json.Marshal(s.Index(i).Interface())
			if err != nil {
				panic(err)
			}
			if err := json.Unmarshal(data, &m); err != nil {
				panic(err)
			}
			outData = append(outData, m)
		}
	}
	return outData
}

var Templates = map[string]string {
	"Req1" : "output1.html",
	"Req2" : "output2.html",
}
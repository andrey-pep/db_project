package main

import (
	"net/http"
	"time"
	//"fmt"
	"html/template"
	"./dbwork"
	"database/sql"
	//"./requests"
//	"github.com/davecgh/go-spew/spew"
	"reflect"
	"encoding/json"
)

var Server = &http.Server {
	Addr:	":80",
	ReadTimeout: 10 * time.Second,
	WriteTimeout: 10 * time.Second,
}

func main() {
	http.HandleFunc("/", Login)
	http.HandleFunc("/index", HandleIndex)
	http.HandleFunc("/login", Login)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.HandleFunc("/select", SelectRequest)
	http.HandleFunc("/authorization", CheckUser)

	err := Server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func CheckUser(w http.ResponseWriter, r *http.Request) {
	PrepareArgs(r)
	err, DB := dbwork.Connect("nonauth", "password")
	if err != nil {
		t := template.Must(template.ParseFiles("public/error.html"))
		if err := t.Execute(w, ""); err != nil {
			panic(err)
		}
		return	
	}
	MC := &MainController{ DataBase : DB, UsrGroup: "nonauth"}
	out := reflect.ValueOf(MC).MethodByName("SelectUser").Call([]reflect.Value{reflect.ValueOf(r)})
	if !out[0].IsNil() {
		return
	}
	if MC.UsrGroup == "nonauth" {
		t := template.Must(template.ParseFiles("public/noauth.html"))
		if err := t.Execute(w, ""); err != nil {
			t := template.Must(template.ParseFiles("public/error.html"))
			if err := t.Execute(w, ""); err != nil {
				panic(err)
			}
		return
	}
		return
	}
	t := template.Must(template.ParseFiles("public/index.html"))
	w.Header().Set("Content-Type", "text/html")
	if err := t.Execute(w, ""); err != nil {
		t := template.Must(template.ParseFiles("public/error.html"))
		if err := t.Execute(w, ""); err != nil {
			panic(err)
		}
		return
	}
}

func SelectRequest(w http.ResponseWriter, r *http.Request) {
	PrepareArgs(r)			//подготовка аргументо из запроса
	err, DB := dbwork.Connect("root", "ghjybr7")	//коннект к базе, до этого нужно сделать авторизацию
	if err != nil {
			t := template.Must(template.ParseFiles("public/error.html"))
		if err := t.Execute(w, ""); err != nil {
			panic(err)
		}
		return
	}
	MC := &MainController{ DataBase : DB, UsrGroup : "teacher"}	//создаю объект с контроллером, в котором хранятся коннект к бд и пользователь, опцианально, мб уберу
	out := reflect.ValueOf(MC).MethodByName(r.URL.Query().Get("action")).Call([]reflect.Value{reflect.ValueOf(r)})	//выполнение самого запроса, пришлось немного с рефлектом поебаца
	if !out[0].IsNil() {
		return
	}	
	if reflect.ValueOf(out[1].Interface()).Len() == 0 {
		t := template.Must(template.ParseFiles("public/nodata.html"))
		if err := t.Execute(w, ""); err != nil {
			panic(err)
		}
		return
	}
	tmplFile := Templates[r.URL.Query().Get("action")]
	t := template.Must(template.ParseFiles("public/" + tmplFile))
	output := PrepareForOut(out[1])
	w.Header().Set("Content-Type", "text/html")
	if err := t.Execute(w, output); err != nil {
		t := template.Must(template.ParseFiles("public/error.html"))
		if err := t.Execute(w, ""); err != nil {
			panic(err)
		}
		return
	}
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	expiration := time.Now().Add(365 * 24 * time.Hour)
    cookie := http.Cookie{Name: "username", Value: "astaxie", Expires: expiration}
    http.SetCookie(w, &cookie)
	t := template.Must(template.ParseFiles("public/index.html"))
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, "")
}

func Login(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("public/login.html"))
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, "")
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
	Login        string
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
package main

import (
	"net/http"
	"time"
	//"fmt"
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
	http.HandleFunc("/", Login)
	http.HandleFunc("/index", HandleIndex)
	http.HandleFunc("/login", Login)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.HandleFunc("/select", SelectRequest)
	http.HandleFunc("/authorization", Enter)

	err := Server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func CheckCookies(cArr []*http.Cookie,MC *MainController) bool {
	for _, c := range cArr {
		if c.Name == "ln" {
			MC.Login = c.Value
		}
		if c.Name == "ug" {
			MC.UsrGroup = c.Value
		}
		if c.Name == "ps" {
			MC.UsrPass = c.Value
		}
	}
	return true
}

func CheckUser(w http.ResponseWriter, r *http.Request, MC *MainController) bool {
	CheckCookies(r.Cookies(), MC)
	if (MC.Login != "" && MC.Login != "" && MC.UsrPass != "") {
		err, exists := MC.CheckIfUserExists()
		if err != nil {
			t := template.Must(template.ParseFiles("public/error.html"))
			if err := t.Execute(w, ""); err != nil {
				panic(err)
			}
		}
		if exists {
			return true
		} else {
			t := template.Must(template.ParseFiles("public/nono.html"))
			if err := t.Execute(w, ""); err != nil {
				panic(err)
			}
		}
	}
	return false
}

func Enter (w http.ResponseWriter, r *http.Request) {
	PrepareArgs(r)
	MC := &MainController{UsrGroup: "nonauth"}
	err, DB := dbwork.Connect(MC.UsrGroup)
	MC.DataBase = DB
	if err != nil {
		t := template.Must(template.ParseFiles("public/error.html"))
		if err := t.Execute(w, ""); err != nil {
			panic(err)
		}
		return
	}
	err = MC.SelectUser(r)
	if err != nil {
		t := template.Must(template.ParseFiles("public/error.html"))
		if err := t.Execute(w, ""); err != nil {
			panic(err)
		}
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
	expiration := time.Now().Add(365 * 24 * time.Hour)
    cookie1 := http.Cookie{Name: "ln", Value: MC.Login, Expires: expiration}
	cookie2 := http.Cookie{Name: "ud", Value: MC.UsrGroup, Expires: expiration}
	cookie3 := http.Cookie{Name: "ps", Value: MC.UsrPass, Expires: expiration}
    http.SetCookie(w, &cookie1)
    http.SetCookie(w, &cookie2)
    http.SetCookie(w, &cookie3)
	w.Header().Set("Content-Type", "text/html")
	if err := t.Execute(w, ""); err != nil {
		t := template.Must(template.ParseFiles("public/error.html"))
		if err := t.Execute(w, ""); err != nil {
			panic(err)
		}
		return
	}
	return
}

func SelectRequest(w http.ResponseWriter, r *http.Request) {
	PrepareArgs(r)			//подготовка аргументо из запроса
	MC := &MainController{UsrGroup: "nonauth"}
	if !CheckUser(w, r, MC) {
		http.Redirect(w, r, "login", http.StatusMovedPermanently)
	}
	err, DB := dbwork.Connect("teacher")	//коннект к базе, до этого нужно сделать авторизацию
	if (err != nil || reflect.ValueOf(DB).IsNil()) {
		t := template.Must(template.ParseFiles("public/error.html"))
		if err := t.Execute(w, ""); err != nil {
			panic(err)
		}
		return
	}
	MC.DataBase = DB
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
	MC := &MainController{UsrGroup: "nonauth"}
	if CheckUser(w, r, MC) {
		t := template.Must(template.ParseFiles("public/index.html"))
		w.Header().Set("Content-Type", "text/html")
		t.Execute(w, "")
	} else {
		http.Redirect(w, r, "login", http.StatusMovedPermanently)
	}
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
	Login       string
	UsrPass     string
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
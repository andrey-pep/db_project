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
	"crypto/md5"
	"encoding/hex"
	"io"
	"./models"
)

var Server = &http.Server {
	Addr:	":80",
	ReadTimeout: 10 * time.Second,
	WriteTimeout: 10 * time.Second,
}

var globalSessions, _ = NewManager("gosessionid",3600)

func main() {

	http.HandleFunc("/", Login)
	http.HandleFunc("/index", HandleIndex)
	http.HandleFunc("/login", Login)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.HandleFunc("/select", SelectRequest)
	http.HandleFunc("/authorization", Enter)
	http.HandleFunc("/logout", Exit)
	http.HandleFunc("/check", CheckOtchet)
	http.HandleFunc("/proc", Procedure)
	http.HandleFunc("/check_tech", IsTech)
	http.HandleFunc("/inmarks", InsertMarks)
	http.HandleFunc("/marks", NewMarks)
	http.HandleFunc("/error", Error)
	err := Server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func Error (w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("public/error.html", "public/helper.html"))
	if err := t.Execute(w, ""); err != nil {
		panic(err)
	}
	return
}

func NewMarks(w http.ResponseWriter, r *http.Request) {
	MC := &MainController{UsrGroup: "nonauth"}
	err, DB := dbwork.Connect(MC.UsrGroup)
	MC.DataBase = DB
	if CheckUser(w, r, MC) == false {
		http.Redirect(w, r, "login", http.StatusSeeOther)
		return
	}
	err, DB = dbwork.Connect(MC.UsrGroup)	//коннект к базе, до этого нужно сделать авторизацию
	MC.DataBase = DB

	if (err != nil || reflect.ValueOf(DB).IsNil()) {
		t := template.Must(template.ParseFiles("public/error.html", "public/helper.html"))
		if err := t.Execute(w, ""); err != nil {
			panic(err)
		}
		return
	}

	var u []models.MarksUpdate
	p := make([]byte, r.ContentLength)
	_, _ = r.Body.Read(p)

	err = json.Unmarshal(p, &u)
	if err != nil {
		t := template.Must(template.ParseFiles("public/error.html", "public/helper.html"))
		if err := t.Execute(w, ""); err != nil {
			panic(err)
		}
		return
	}
	affected, err := MC.MarksInput(u)
	if (affected != len(u) || err != nil) {
		t := template.Must(template.ParseFiles("public/error.html", "public/helper.html"))
		if err := t.Execute(w, ""); err != nil {
			panic(err)
		}
		return	
	}
	t := template.Must(template.ParseFiles("public/itsok.html", "public/helper.html"))
	if err := t.Execute(w, ""); err != nil {
		t := template.Must(template.ParseFiles("public/error.html", "public/helper.html"))
		if err := t.Execute(w, ""); err != nil {
			panic(err)
		}
		return
	}
	return

}

func InsertMarks(w http.ResponseWriter, r *http.Request) {
	PrepareArgs(r)			//подготовка аргументо из запроса
	MC := &MainController{UsrGroup: "nonauth"}
	err, DB := dbwork.Connect(MC.UsrGroup)
	MC.DataBase = DB
	if CheckUser(w, r, MC) == false {
		http.Redirect(w, r, "login", http.StatusSeeOther)
		return
	}
	err, DB = dbwork.Connect(MC.UsrGroup)	//коннект к базе, до этого нужно сделать авторизацию
	MC.DataBase = DB

	if (err != nil || reflect.ValueOf(DB).IsNil()) {
		t := template.Must(template.ParseFiles("public/error.html", "public/helper.html"))
		if err := t.Execute(w, ""); err != nil {
			panic(err)
		}
		return
	}

	err, students := MC.SelectStudents(r)
	if (err != nil || reflect.ValueOf(DB).IsNil()) {
		t := template.Must(template.ParseFiles("public/error.html", "public/helper.html"))
		if err := t.Execute(w, ""); err != nil {
			panic(err)
		}
		return
	}

	t := template.Must(template.ParseFiles("public/students.html", "public/helper.html"))
	output := PrepareForOut(reflect.ValueOf(students))
	w.Header().Set("Content-Type", "text/html")
	if err := t.Execute(w, output); err != nil {
		t := template.Must(template.ParseFiles("public/error.html", "public/helper.html"))
		if err := t.Execute(w, ""); err != nil {
			panic(err)
		}
		return
	}
	return

}

func IsTech(w http.ResponseWriter, r *http.Request) {
	PrepareArgs(r)			//подготовка аргументо из запроса
	MC := &MainController{UsrGroup: "nonauth"}
	_, DB := dbwork.Connect(MC.UsrGroup)
	MC.DataBase = DB
	if CheckUser(w, r, MC) == false {
		http.Redirect(w, r, "login", http.StatusSeeOther)
		return
	}
	resp := CheckResp{false}
	if MC.UsrGroup == "root" {
		resp.Status = true
	}
	js, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func Procedure(w http.ResponseWriter, r *http.Request) {
	PrepareArgs(r)			//подготовка аргументо из запроса
	MC := &MainController{UsrGroup: "nonauth"}
	err, DB := dbwork.Connect(MC.UsrGroup)
	MC.DataBase = DB
	if CheckUser(w, r, MC) == false {
		http.Redirect(w, r, "login", http.StatusSeeOther)
		return
	}
	err, DB = dbwork.Connect(MC.UsrGroup)	//коннект к базе, до этого нужно сделать авторизацию
	MC.DataBase = DB
	err, out := MC.Otchet()
	if err != nil {
		t := template.Must(template.ParseFiles("public/error.html", "public/helper.html"))
		if err := t.Execute(w, ""); err != nil {
			panic(err)
		}
		return
	}
	if len(out) > 1 {
		io.WriteString(w, "уже выполнялась")
		return
	} else {
		status := MC.MakeOtchet(r)
		if status != nil {
			t := template.Must(template.ParseFiles("public/error.html", "public/helper.html"))
			if err := t.Execute(w, ""); err != nil {
				panic(err)
			}
		return	
		}
	}
}

func CheckOtchet(w http.ResponseWriter, r *http.Request) {
	PrepareArgs(r)			//подготовка аргументо из запроса
	MC := &MainController{UsrGroup: "nonauth"}
	err, DB := dbwork.Connect(MC.UsrGroup)
	MC.DataBase = DB
	if CheckUser(w, r, MC) == false {
		http.Redirect(w, r, "login", http.StatusSeeOther)
		return
	}
	err, DB = dbwork.Connect(MC.UsrGroup)	//коннект к базе, до этого нужно сделать авторизацию
	MC.DataBase = DB
	err, out := MC.Otchet()
	if err != nil {
		t := template.Must(template.ParseFiles("public/error.html", "public/helper.html"))
		if err := t.Execute(w, ""); err != nil {
			panic(err)
		}
		return
	}
	if len(out) < 1 {
		//io.WriteString(w, "1\n")
		t := template.Must(template.ParseFiles("public/nodata.html", "public/helper.html"))
		if err := t.Execute(w, ""); err != nil {
			t := template.Must(template.ParseFiles("public/error.html"))
			if err := t.Execute(w, ""); err != nil {
				panic(err)
			}
		}
		return
	} else {
		t := template.Must(template.ParseFiles("/public/otchet.html", "public/helper.html"))
		if err := t.Execute(w, out); err != nil {
			t := template.Must(template.ParseFiles("public/error.html", "public/helper.html"))
			if err := t.Execute(w, ""); err != nil {
				panic(err)
			}
		}
		return
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
	sid, err := r.Cookie("sid")
	if err != nil {
		return false
	}
	session := globalSessions.CheckSession(sid.Value)
	if session != nil {
		MC.Login = session.GetValue("login").(string)
		MC.UsrGroup = session.GetValue("userGroup").(string)
		MC.UsrPass = session.GetValue("password").(string)
		return true
	} else {
		return false
	}
}

func Enter (w http.ResponseWriter, r *http.Request) {
	PrepareArgs(r)
	MC := &MainController{UsrGroup: "nonauth"}
	err, DB := dbwork.Connect(MC.UsrGroup)
	MC.DataBase = DB
	if err != nil {
		t := template.Must(template.ParseFiles("public/error.html", "public/helper.html"))
		if err := t.Execute(w, ""); err != nil {
			panic(err)
		}
		return
	}
	err = MC.SelectUser(r)
	if err != nil {
		t := template.Must(template.ParseFiles("public/error.html", "public/helper.html"))
		if err := t.Execute(w, ""); err != nil {
			panic(err)
		}
		return
	}
	if MC.UsrGroup == "nonauth" {
		t := template.Must(template.ParseFiles("public/noauth.html", "public/helper.html"))
		if err := t.Execute(w, ""); err != nil {
			t := template.Must(template.ParseFiles("public/error.html", "public/helper.html"))
			if err := t.Execute(w, ""); err != nil {
				panic(err)
			}
			return
		}
		return
	}
	t := template.Must(template.ParseFiles("public/index.html", "public/helper.html"))
	expiration := time.Now().Add(365 * 24 * time.Hour)
	hasher := md5.New()
	hasher.Write([]byte(r.URL.Query().Get("password")))
	hashedPass := hex.EncodeToString(hasher.Sum(nil))
    cookie1 := http.Cookie{Name: "ln", Value: MC.Login, Expires: expiration}
	cookie2 := http.Cookie{Name: "ud", Value: MC.UsrGroup, Expires: expiration}
	cookie3 := http.Cookie{Name: "ps", Value: hashedPass, Expires: expiration}
    http.SetCookie(w, &cookie1)
    http.SetCookie(w, &cookie2)
    http.SetCookie(w, &cookie3)
    sid := globalSessions.SessionInit()
    session := globalSessions.Sessions[sid]
    session.SetValue(MC.Login, "login") 
    session.SetValue(MC.UsrGroup, "userGroup")
    session.SetValue(hashedPass, "password")
    cookie4 := http.Cookie{Name: "sid", Value: sid, Expires: expiration}
    http.SetCookie(w, &cookie4)
	w.Header().Set("Content-Type", "text/html")
	if err := t.Execute(w, ""); err != nil {
		t := template.Must(template.ParseFiles("public/error.html", "public/helper.html"))
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
	err, DB := dbwork.Connect(MC.UsrGroup)
	MC.DataBase = DB
	if CheckUser(w, r, MC) == false {
		http.Redirect(w, r, "login", http.StatusSeeOther)
		return
	}
	err, DB = dbwork.Connect(MC.UsrGroup)	//коннект к базе, до этого нужно сделать авторизацию
	MC.DataBase = DB

	if (err != nil || reflect.ValueOf(DB).IsNil()) {
		t := template.Must(template.ParseFiles("public/error.html", "public/helper.html"))
		if err := t.Execute(w, ""); err != nil {
			panic(err)
		}
		return
	}

	out := reflect.ValueOf(MC).MethodByName(r.URL.Query().Get("action")).Call([]reflect.Value{reflect.ValueOf(r)})	//выполнение самого запроса, пришлось немного с рефлектом поебаца
	if !out[0].IsNil() {
		t := template.Must(template.ParseFiles("public/error.html", "public/helper.html"))
		if err := t.Execute(w, ""); err != nil {
			panic(err)
		}
		return
	}
	if reflect.ValueOf(out[1].Interface()).Len() == 0 {
		t := template.Must(template.ParseFiles("public/nodata.html", "public/helper.html"))
		if err := t.Execute(w, ""); err != nil {
			panic(err)
		}
		return
	}
	tmplFile := Templates[r.URL.Query().Get("action")]
	t := template.Must(template.ParseFiles("public/" + tmplFile, "public/helper.html"))
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
	err, DB := dbwork.Connect(MC.UsrGroup)
	if (err != nil || reflect.ValueOf(DB).IsNil()) {
		t := template.Must(template.ParseFiles("public/error.html", "public/helper.html"))
		if err := t.Execute(w, ""); err != nil {
			panic(err)
		}
		return
	}
	MC.DataBase = DB
	if !CheckUser(w, r, MC) {
		t := template.Must(template.ParseFiles("public/index.html", "public/helper.html"))
		w.Header().Set("Content-Type", "text/html")
		t.Execute(w, "")
	} else {
		http.Redirect(w, r, "login", http.StatusSeeOther)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	MC := &MainController{UsrGroup: "nonauth"}
	err, DB := dbwork.Connect(MC.UsrGroup)
	if (err != nil || reflect.ValueOf(DB).IsNil()) {
		t := template.Must(template.ParseFiles("public/error.html", "public/helper.html"))
		if err := t.Execute(w, ""); err != nil {
			panic(err)
		}
		return
	}
	MC.DataBase = DB
	if CheckUser(w, r, MC) {
		t := template.Must(template.ParseFiles("public/index.html", "public/helper.html"))
		w.Header().Set("Content-Type", "text/html")
		t.Execute(w, "")
	} else {
		t := template.Must(template.ParseFiles("public/login.html"))
		w.Header().Set("Content-Type", "text/html")
		t.Execute(w, "")
	}
}

func Exit(w http.ResponseWriter, r *http.Request) {
		expiration := time.Now().Add(0)
    	sid, _ := r.Cookie("sid")
    	globalSessions.DeleteSession(sid.Value)
    	cookie1 := http.Cookie{Name: "sid", Value: "", Expires: expiration}
    	http.SetCookie(w, &cookie1)
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
	"Req1" : 			"output1.html",
	"Req2" : 			"output2.html",
	"Req3" : 			"output3.html",
	"Req4" : 			"output4.html",
	"Req5" : 			"output5.html",
	"Req6" : 			"output6.html",
	"SelectStudents" : 	"students.html",
}

type CheckResp struct {
	Status bool
}
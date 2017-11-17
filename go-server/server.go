package main

import (
	"net/http"
	"time"
	"fmt"
)

func main() {
	Server := &http.Server {
		Addr:	":8080",
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	http.HandleFunc("/", HandleIndex)
	http.HandleFunc("/login", Login)

	Server.ListenAndServe()
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	    expiration := time.Now().Add(365 * 24 * time.Hour)
    cookie := http.Cookie{Name: "username", Value: "astaxie", Expires: expiration}
     http.SetCookie(w, &cookie)
	fmt.Fprintf(w, "Hello wrold!")
	Cookies := r.Cookies()
	for _, c := range Cookies {
		fmt.Fprintf(w, c.String())
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	
}
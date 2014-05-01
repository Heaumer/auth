package main

import (
	"github.com/gorilla/securecookie"
	"io/ioutil"
	"log"
	"runtime"
	"net/http"
)

// LogError calls log.Printf on error, and adds location in source code
func LogError(err error) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		log.Printf("%s:%d : %s", file, line, err)
	} else {
		log.Println(err)
	}
}

// LogHttp log error and sends it to browser
func LogHttp(w http.ResponseWriter, err error) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		log.Printf("%s:%d : %s", file, line, err)
	} else {
		log.Println(err)
	}
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

// LogFatal calls log.Fatalf on error, and adds location in source code
func LogFatal(err error) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		log.Fatalf("%s:%d : %s", file, line, err)
	} else {
		log.Fatal(err)
	}
	
}

// WriteFiles write the files it's given as argument to w
func writeFiles(w http.ResponseWriter, files ...string) error {
	for _, file := range files {
		b, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		w.Write(b)
	}
	return nil
}

func ok(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("ok"))
}
func ko(w http.ResponseWriter) {
	w.Write([]byte("ko"))
}

var hashKey = []byte(securecookie.GenerateRandomKey(32))
var blockKey = []byte(securecookie.GenerateRandomKey(32))
var s = securecookie.New(hashKey, blockKey)

func SetToken(w http.ResponseWriter, token string) error {
	encoded, err := s.Encode("auth-token", token)
	if err != nil { return err }

	cookie := &http.Cookie{
		Name:	"auth-token",
		Value:	encoded,
		Path:	"/",
	}
	http.SetCookie(w, cookie)

	return nil
}

func UnsetToken(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:	"auth-token",
		Value:	"",
		Path:	"/",
		MaxAge:	-1,
	}
	http.SetCookie(w, cookie)
}

func GetToken(r *http.Request) (token string, err error) {
	cookie, err := r.Cookie("auth-token")
	if err == nil {
		err = s .Decode("auth-token", cookie.Value, &token)
	}

	return
}

func GetNavbar(r *http.Request) string {
	navbar := "templates/navbar.html"

	// if connected
	token, err := GetToken(r)
	if err == nil && CheckToken(&Token{ Auth.Key, token }) {
		navbar = "templates/navbar2.html"
		if IsAdmin(token) {
			navbar = "templates/navbar3.html"
		}
	}

	return navbar
}

func ACheckToken(token string) bool {
	return CheckToken(&Token{ Auth.Key, token })
}
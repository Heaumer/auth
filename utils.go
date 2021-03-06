package main

import (
	"./cookie"
	"io/ioutil"
	"math/rand"
	"net/http"
)

const (
	alnum = "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ123456789"
)

// Generate random string of n bytes
func randomString(n int) string {
	buf := make([]byte, n)

	for i := 0; i < n; i++ {
		buf[i] = alnum[rand.Intn(len(alnum))]
	}

	return string(buf)
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
	w.Write([]byte("ok"))
}
func ko(w http.ResponseWriter) {
	w.Write([]byte("ko"))
}

func SetToken(w http.ResponseWriter, token string) error {
	if err := cookie.Set(w, "auth-token", token); err != nil {
		return Err(err)
	}
	return nil
}

func UnsetToken(w http.ResponseWriter) {
	cookie.Unset(w, "auth-token")
}

func VerifyToken(r *http.Request) (string, error) {
	token, err := cookie.Get(r, "auth-token")

	if err != nil || !CheckToken(token, Auth.Key) {
		return "", MouldyCookie
	}

	return token, nil
}

func SetInfo(w http.ResponseWriter, msg string) {
	cookie.Set(w, "auth-info", msg)
}

func GetInfo(r *http.Request) string {
	msg, _ := cookie.Get(r, "auth-info")
	return msg

}

func UnsetInfo(w http.ResponseWriter) {
	cookie.Unset(w, "auth-info")
}

func SetError(w http.ResponseWriter, err error) {
	SetInfo(w, "Error: "+err.Error())
}

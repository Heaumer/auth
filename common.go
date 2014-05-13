package main

import (
	"errors"
	"runtime"
	"strconv"
	"time"
)

type User struct {
	Id			int32
	Name		string
	Email		string
	Admin		bool
}

type Service struct {
	Id			int32
	Name		string
	Url			string
	Key			string
	Mode		bool
	Address		string
	Email		string
}

var db			*Database
var	services	map[string]*Service

const (
	Automatic = iota
	Manual
	Disabled
)

var ServiceMode = Manual

type InternalError struct {
	Date		time.Time
	File		string
	Line		int
	Msg			string
}

func (e *InternalError) Error() string {
	return e.Date.String()+" "+e.File+":"+strconv.Itoa(e.Line)+
		" "+e.Msg
}

func MkIErr(err error) *InternalError {
	_, file, line, _ := runtime.Caller(1)

	return &InternalError{ time.Now(), file, line, err.Error() }
}

var (
	WrongLogin		= errors.New("Wrong name/email")
	WrongToken		= errors.New("Wrong Token")
	NonSense		= errors.New("This Sense Makes No Action")
	NeedEmail		= errors.New("Email is mandatory")
	LongEmail		= errors.New("Email is too long (maxsize: "+
						strconv.Itoa(LenToken-1)+")")
	WrongEmail		= errors.New("Wrong Email format (you@provider)")
	NeedName		= errors.New("Name is mandatory")
	LongName		= errors.New("Name is too long (maxsize: "+
						strconv.Itoa(LenToken-1)+")")
	WrongName		= errors.New("Invalid characters in name (no whites, @)")

	WrongUser		= errors.New("User name or password already in use")

	SMTPErr			= errors.New("Email not send. Contact an admin.")

	MouldyCookie	= errors.New("Mouldy Cookie, Sour Tea!")
	NotAdminErr		= errors.New("Can't go there.")
)

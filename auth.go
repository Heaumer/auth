package main

import (
	"crypto/sha512"
	"encoding/base64"
	"strings"
	"net/smtp"
	"log"
)

// sendEmail sends an email to an user.
// XXX use several SMTP according to the destination email
// provider to speed things up.
func sendEmail(to, subject, msg string) error {
	body := "To: " + to + "\r\nSubject: " +
		subject + "\r\n\r\n" + msg

	auth := smtp.PlainAuth("", C.AuthEmail, C.AuthPasswd, C.SMTPServer)

	if err := smtp.SendMail(C.SMTPServer+":"+C.SMTPPort, auth, C.AuthEmail,
			[]string{to},[]byte(body)); err != nil {
		return Err(err)
	}

	return nil
}

// sendEmail sends a token via email to an user.
func sendToken(email string, token *Token) error {
	s := db.GetService2(token.Key)

	err := sendEmail(email, "Token for "+s.Name,
		"Hi there,\r\n"+
		"Here is your token for "+s.Name+" ("+s.Url+")"+": "+token.Token)

	if err != nil { log.Println(err); return SMTPErr }

	return nil
}

func checkName(name string) error {
	switch {
	case name == "":							return NoNameErr
	case len(name) >= C.LenToken:				return LongNameErr
	case strings.Contains(name, "@ \t\n\r"):	return NameFmtErr
	}

	return nil
}

func checkEmail(email string) error {
	switch {
	case email == "":						return NoEmailErr
	case len(email) >= C.LenToken:			return LongEmailErr
	case !strings.Contains(email, "@"):		return EmailFmtErr
	}

	return nil
}

// isToken check whether the login is a token or a name/email.
func isToken(login string) bool { return len(login) == C.LenToken }

// isEmail check whether the login is a name or an email
func isEmail(login string) bool { return strings.Contains(login, "@") }

// Register add a new user to both database and cache.
// If the registration succeeds, a(n activation) token is
// sent to the user.
func Register(name, email, passwd string) error {
	if err := checkName(name); err != nil {
		return err
	}
	if err := checkEmail(email); err != nil {
		return err
	}

	if passwd != "" {
		h := sha512.New()
		h.Write([]byte(passwd))
		passwd = base64.StdEncoding.EncodeToString(h.Sum(nil))
	}

	u := User{ -1, name, email, passwd, false }

	if err := db.AddUser(&u); err != nil {
		log.Println(err)
		return WrongUser
	}

	return sendToken(email, NewToken(u.Id, Auth.Key))
}

func Login(login, passwd string) (string, error) {
	// login with token
	if isToken(login) {
		ntoken := UpdateToken(login)
		if ntoken == "" {
			return "", NoSuchTErr
		}
		return ntoken, nil
	}

	// get user associated with login
	u, err := db.GetUser2(login)
	if err != nil {
		return "", NoSuchErr
	}

	// login with password
	if passwd != "" {
		h := sha512.New()
		h.Write([]byte(passwd))
		passwd = base64.StdEncoding.EncodeToString(h.Sum(nil))

		if passwd == u.Passwd {
			return NewToken(u.Id, Auth.Key).Token, nil
		} else {
			return "", BadPasswd
		}
	}

	// 2-steps login (sending token through token)
	return "", sendToken(u.Email, NewToken(u.Id, Auth.Key))
}

func Logout(token string) {
	RemoveToken(token)
}

/*func Unregister() {
}*/

func IsAdmin(token string) bool {
	return db.IsAdmin(OwnerToken(token))
}

func AddService(name, url, address, email string) (string, error) {
	if name == "" || url == "" {
		return "", EmptyFieldsErr
	}

	if ServiceMode == Disabled { return "ko", nil }

	s := Service{ -1, name, url, randomString(C.LenKey), false, address, email }
	if err := db.AddService(&s); err != nil {
		return "", err
	}

	if ServiceMode == Automatic {
		db.SetMode(s.Id, true)
		return s.Key, nil
	}

	// Manual
	SendAdmin("New Service "+s.Name,
			"Hi there,\r\n"+
			s.Name + " ("+s.Address+", "+s.Url+") asks for landing.")

	return "ok", nil

}

func CheckService(key, address string) bool {
	s := db.GetService2(key)
	if s == nil  { return false }

	return s.Address == address && s.Mode
}

func SendAdmin(subject, msg string) {
	if emails, err := db.GetAdminMail(); err != nil {
		log.Println(err)
		return
	} else {
		for _, to := range emails {
			if err := sendEmail(to, subject, msg); err != nil {
				log.Println(err)
			}
		}
	}
}

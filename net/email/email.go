package email

import (
    "errors"
    "net/smtp"
    "strings"
)

var mailers = map[string]*Mailer{}

type Mailer struct {
    Host string
    Port string
    User string
    Pass string
    auth smtp.Auth
}

func MailerRegister(name, host, port, user, pass string) {

    if ml, ok := mailers[name]; ok {

        if host == ml.Host &&
            port == ml.Port &&
            user == ml.User &&
            pass == ml.Pass {
            return
        }
    }

    mailers[name] = NewMailer(host, port, user, pass)
}

func MailerPull(name string) (*Mailer, error) {

    if v, ok := mailers[name]; ok {
        return v, nil
    }

    return nil, errors.New("No Mailer Found")
}

func NewMailer(host, port, user, pass string) *Mailer {
    return &Mailer{
        User: user,
        Pass: pass,
        Host: host,
        Port: port,
        auth: smtp.PlainAuth("", user, pass, host),
    }
}

func (m Mailer) SendMail(to, subject, body string) error {

    mailtype := "plain"
    if len(body) > 20 && body[:5] == "<html" {
        mailtype = "html"
    }

    msg := "From: " + m.User + "<" + m.User + ">\r\n"
    msg += "To: " + to + "\r\n"
    msg += "Subject: " + subject + "\r\n"
    msg += "Content-Type: text/" + mailtype + "; charset=UTF-8\r\n\r\n"
    msg += body

    send_to := strings.Split(to, ";")

    return smtp.SendMail(m.Host+":"+m.Port, m.auth, m.User, send_to, []byte(msg))
}

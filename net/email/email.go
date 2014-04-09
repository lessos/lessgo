package email

import (
    "errors"
    "net/smtp"
    "strings"
)

var mailers = map[string]*Mailer{}

type Mailer struct {
    User string
    Pass string
    Host string
    auth smtp.Auth
}

func MailerRegister(name, user, password, host string) {

    if _, ok := mailers[name]; ok {
        return
    }

    mailers[name] = NewMailer(user, password, host)
}

func MailerPull(name string) (*Mailer, error) {

    if v, ok := mailers[name]; ok {
        return v, nil
    }

    return nil, errors.New("No Mailer Found")
}

func NewMailer(host, user, password string) *Mailer {

    hs := strings.Split(host, ":")

    return &Mailer{
        User: user,
        Pass: password,
        Host: host,
        auth: smtp.PlainAuth("", user, password, hs[0]),
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

    return smtp.SendMail(m.Host, m.auth, m.User, send_to, []byte(msg))
}

package utils

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/scorredoira/email"
	"net"
	"net/mail"
	"net/smtp"
)

func SendEmail(emailHost, emailUser, emailPass, fromName, sub, body string, toUserList []string) {
	m := email.NewHTMLMessage(sub, body)
	m.From = mail.Address{Name: fromName, Address: "git@gittab.com"}
	m.To = toUserList
	serverName := emailHost
	host, _, _ := net.SplitHostPort(serverName)
	auth := smtp.PlainAuth("", emailUser, emailPass, host)
	log.Info(fmt.Sprintf("开始发送邮件,收件人是%s", m.To))
	if err := email.Send(serverName, auth, m); err != nil {
		log.Fatal(err)
	} else {
		log.Info("发送成功")
	}
}

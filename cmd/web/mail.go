package main

import (
	_"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
	"github.com/yusuf/track-space/pkg/model"
)

func SendMail(m model.Email, password string){
	mailServer := mail.NewSMTPClient()
	mailServer.Host = "smtp.gmail.com"
	mailServer.Port= 465
	mailServer.Password = password
	mailServer.Username = m.Sender
	mailServer.Encryption = mail.EncryptionSSLTLS
	mailServer.ConnectTimeout = 100 * time.Second
	mailServer.SendTimeout= 100 * time.Second

	mailClient, err := mailServer.Connect()
	if err != nil{
		log.Panicln(err)
	}

	msg := mail.NewMSG()
	msg.SetFrom(m.Sender).AddTo(m.Receiver).SetSubject(m.Subject)
	if m.Template == ""{
		msg.SetBody(mail.TextHTML, m.Content)
	} else{
		data, err := ioutil.ReadFile(fmt.Sprintf("./mail-template/%s", m.Template))
		if err != nil{
			log.Panic(err)
		}
		temp := string(data)
		mailToSend := strings.Replace(temp, "[%body%]",m.Content, 1)
		msg.SetBody(mail.TextHTML, mailToSend)
	}
	if err := msg.Send(mailClient); err != nil{
		log.Println(err)
	} else{
		log.Println("Email Sent Successfully")
	}

}


func ListenToMailChannel(password string) {
	func() {
		for {
			mailMsg := <-app.MailChan
			SendMail(mailMsg, password)
		}
	}()
}

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
	"github.com/yusuf/track-space/pkg/model"
)

func SendMailToUser(m model.Email) {
	mailServer := mail.NewSMTPClient()
	mailServer.Host = "localhost"
	mailServer.Port = 1025
	// mailServer.Encryption = mail.EncryptionSTARTTLS
	mailServer.ConnectTimeout = 100 * time.Second
	mailServer.SendTimeout = 100 * time.Second

	//connecting to mailhog server
	mailClient, err := mailServer.Connect()
	if err != nil{
		log.Panicln(err)
	}

	email := mail.NewMSG()
	email.SetFrom(m.Sender).AddTo(m.Receiver).SetSubject(m.Subject)
	if m.Template == ""{
		email.SetBody(mail.TextHTML, m.Content)
	}else{
		byteData, err := ioutil.ReadFile(fmt.Sprintf("./mail-template/%s",m.Template))
		// fs.ReadFile("./mail-template",m.Template)
		if err != nil{
			log.Panic(err)
		}
		template := string(byteData)
		mailSend := strings.Replace(template,"[%body%]", m.Content,1)
		email.SetBody(mail.TextHTML, mailSend)
	}

	if err :=email.Send(mailClient); err != nil{
		log.Println(err)
	} else{
		log.Println("Email Sent Successfully")
	}

}

func ListenToMailChannel()  {
	go func() {
		for {
			mailMsg := <-app.MailChan
			SendMailToUser(mailMsg)
		}
	}()
}
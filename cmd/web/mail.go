package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/yusuf/track-space/pkg/model"
	"gopkg.in/gomail.v2"
)

func SendMail(m model.Email){
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.Sender)
	msg.SetHeader("To", m.Receiver)
	msg.SetHeader("Subject", m.Subject)
	if m.Template == ""{
		msg.SetBody("text/html", m.Content)
	}else{
		d, err := ioutil.ReadFile(fmt.Sprintf("./mail-template/%s", m.Template))
		if err != nil{
			log.Panic(err)
		}
		template := string(d)
		mailToSend := strings.Replace(template, "[%body%]", m.Content, 1)
		msg.SetBody("text/html", mailToSend)
		
	}
	dailMail := gomail.NewDialer("smtp.gmail.com", 587, m.Sender, "Akinleye123")
	dailMail.TLSConfig = &tls.Config{
		InsecureSkipVerify: false,
	}
	if err := dailMail.DialAndSend(msg); err != nil{
		log.Println(err)
		panic(err)
	}else{
		log.Println("Email Sent Successfully")
		return
	}
}


func ListenToMailChannel() {
	func() {
		for {
			mailMsg := <-app.MailChan
			SendMail(mailMsg)
		}
	}()
}
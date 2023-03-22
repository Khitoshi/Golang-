package main

import (
	"fmt"
	"net/smtp"

	"github.com/jhillyerd/enmime"
)

func main() {
	fmt.Printf("Start\n")

	smtpHost := "my-mail-server:25"
	smtpAuth := smtp.PlainAuth(
		"example.com",
		"example.user",
		"example.password",
		"auth.example.com",
	)

	sender := enmime.NewSMTP(smtpHost, smtpAuth)
	//builder patternで作られている
	master := enmime.Builder().
		From("送信太郎", "taroA@example.com").
		Subject("メールタイプ").
		Text([]byte("テキストメール本文")).
		HTML([]byte("<p>HTML メール<b>本文</b></p>")).
		AddFileAttachment("document.pdf")
	msg := master.To("宛先太郎", "taroB@example.com")
	err := msg.Send(sender)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nend")
}

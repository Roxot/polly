package main

import "net/smtp"

const (
	cFrom              = "pssrea@gmail.com"
	cPass              = "golangftw"
	cMailServer        = "smtp.gmail.com"
	cMailServerAndPort = "smtp.gmail.com:587"
)

/* Connects to the gmail server and sends mail. */
func SendMail(to string, body []byte) error {
	auth := smtp.PlainAuth("", cFrom, cPass, cMailServer)
	err := smtp.SendMail(cMailServerAndPort, auth, cFrom, []string{cTo}, body)
	return err
}

func VerTokenBody(string verToken) []byte {

}

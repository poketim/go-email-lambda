package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"net"
	"net/smtp"
)

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("unknown from server")
		}
	}
	return nil, nil
}

func HandleRequest() {
	smtpHost := ""
	port := ""
	from := ""
	password := ""
	to := []string{""}
	subject := ""
	body := ""

	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to[0]
	headers["Subject"] = subject

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	tlsConfig := &tls.Config{
		ServerName: smtpHost,
	}

	conn, err := net.Dial("tcp", smtpHost+":"+port)
	if err != nil {
		fmt.Println("tls.Dial Error: ", err)
		return
	}

	c, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		fmt.Println("smtp.NewClient Error: ", err)
		return
	}

	if err = c.StartTLS(tlsConfig); err != nil {
		fmt.Println("c.StartTLS Error: ", err)
		return
	}

	if err = c.Auth(LoginAuth(from, password)); err != nil {
		fmt.Println("c.Auth Error: ", err)
		return
	}

	if err = c.Mail(from); err != nil {
		fmt.Println("c.Mail Error: ", err)
		return
	}

	for i := 0; i < len(to); i++ {
		if err = c.Rcpt(to[0]); err != nil {
			fmt.Println("c.Rcpt Error: ", err)
			return
		}
	}

	w, err := c.Data()
	if err != nil {
		fmt.Println("c.Data Error: ", err)
		return
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	if err = w.Close(); err != nil {
		fmt.Println("reader Error: ", err)
		return
	}

	if err = c.Quit(); err != nil {
		fmt.Println("quit error: ", err)
		return
	}

	fmt.Println("email sent successfully!")
}

func main() {
	lambda.Start(HandleRequest)
}

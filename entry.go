package main

import (
	"errors"
	"fmt"
	"gopkg.in/gomail.v2"
	"os"
)

const NAME = "easy-mail"
const VERSION = "v0.1"

func parseAndRun(input []string) error {
	raw, e := parseArgs(input)
	if e != nil {
		return e
	}
	args, e := tidyArgs(raw)
	if e != nil {
		return e
	}
	if args.GenerateAuth {
		return saveAuth(args.AuthPath, args.From, args.Password, fmt.Sprintf("%s:%d", args.SMTPHosts[0], args.SMTPPorts[0]))
	}

	if args.ShowVersion {
		fmt.Println(VERSION)
		return nil
	}
	var client gomail.SendCloser
	for _, host := range args.SMTPHosts {
		for _, port := range args.SMTPPorts {
			d := gomail.NewDialer(host, port, args.From, args.Password)
			client, e = d.Dial()
			if e == nil {
				break
			}
		}
		if client != nil {
			break
		}
	}
	if client == nil {
		return errors.New("failed to connect smtp server")
	}
	defer client.Close()

	msg := gomail.NewMessage()
	msg.SetHeader("From", args.From)
	msg.SetHeader("To", args.To...)
	msg.SetHeader("Subject", args.Subject)
	if args.Content != "" {
		msg.SetBody(args.ContentType, args.Content)
	}
	for _, attach := range args.Attaches {
		msg.Attach(attach)
	}
	return gomail.Send(client, msg)
}

func main() {
	e := parseAndRun(os.Args)
	if e != nil {
		fmt.Println(e.Error())
	}
}

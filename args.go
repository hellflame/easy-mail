package main

import (
	"errors"
	"fmt"
	"github.com/akamensky/argparse"
	"net"
	"os"
	"strconv"
)

type RawArgs struct {
	From string
	To []string
	Subject string
	ContentType string
	Content string
	ContentPath string
	Attaches []string
	SMTPServer string
	Password string
	AuthPath string
	GenerateAuth bool
}

type TidyArgs struct {
	From string
	To []string
	Subject string
	Content string
	ContentType string
	Attaches []string
	SMTPHosts []string
	SMTPPorts []int
	Password string
	GenerateAuth bool
	AuthPath string
}


func parseArgs(input []string) (args *RawArgs, e error) {
	parser := argparse.NewParser("mailall", "send mail from command line")
	from := parser.String("", "from", &argparse.Options{Help:"email send from"})
	to := parser.StringList("", "to", &argparse.Options{Help:"recv address list"})
	subject := parser.String("s", "subject", &argparse.Options{Help:"email title"})
	content := parser.String("c", "content", &argparse.Options{Help: "simple email content"})
	contentPath := parser.String("", "content-path", &argparse.Options{Help:"email content path"})
	contentType := parser.String("", "content-type", &argparse.Options{Help:"email content type"})
	attach := parser.StringList("", "attach", &argparse.Options{Help: "attach file path list"})
	smtp := parser.String("", "smtp", &argparse.Options{Help:"manually set smtp address like: smtp.abc.com:587 it can be auto find if not set"})
	password := parser.String("", "password", &argparse.Options{Help:"email password"})
	generateAuth := parser.Flag("g", "generate", &argparse.Options{Help: "generate auth file to simple use"})
	authPath := parser.String("a", "auth", &argparse.Options{Help:"auth file path"})
	if len(input) == 1 {
		input = append(input, "-h")
	}
	e = parser.Parse(input)
	if e != nil {
		return
	}
	args = &RawArgs{
		From: *from,
		To: *to,
		Subject: *subject,
		Attaches: *attach,
		SMTPServer: *smtp,
		Password: *password,
		AuthPath: *authPath,
		GenerateAuth: *generateAuth,
		Content: *content,
		ContentPath: *contentPath,
		ContentType: *contentType,
	}
	return
}


func tidyArgs(args *RawArgs) (*TidyArgs, error) {
	var tidyResult TidyArgs

	if args.From != "" {
		mailBox, valid := validateEmailAddress(args.From)
		if !valid {
			return nil, fmt.Errorf("invalid from address format")
		}
		tidyResult.From = args.From
		if args.SMTPServer == "" {
			hosts := fetchMxHosts(mailBox)
			if len(hosts) == 0 {
				return nil, fmt.Errorf("can't find mx servers for %s", mailBox)
			}
			tidyResult.SMTPHosts = hosts
			tidyResult.SMTPPorts = []int{25, 465, 587}
		}
	}
	if args.Password != "" {
		tidyResult.Password = args.Password
	}
	if args.SMTPServer != "" {
		host, p, e := net.SplitHostPort(args.SMTPServer)
		if e != nil {
			return nil, fmt.Errorf("invalid smtp server: %s", e.Error())
		}
		port, e := strconv.Atoi(p)
		tidyResult.SMTPHosts = []string{host}
		tidyResult.SMTPPorts = []int{port}
	}
	tidyResult.GenerateAuth = args.GenerateAuth
	if args.GenerateAuth {
		return &tidyResult, nil
	}

	if tidyResult.From == "" && tidyResult.Password == "" && len(tidyResult.SMTPHosts) == 0 {
		from, password, host, port, e := loadAuth(args.AuthPath)
		if e != nil {
			return nil, fmt.Errorf("failed to load auth: %s", e.Error())
		}
		fmt.Printf("using auths from storage: %s\n", from)
		tidyResult.From = from
		tidyResult.Password = password
		tidyResult.SMTPHosts = []string{host}
		tidyResult.SMTPPorts = []int{port}
	}

	if tidyResult.From == "" || tidyResult.Password == "" || len(tidyResult.SMTPHosts) == 0 {
		return nil, errors.New("failed to set user credentials")
	}

	if len(args.To) == 0 {
		return nil, errors.New("no one to send to")
	} else {
		for _, to := range args.To {
			if _, ok := validateEmailAddress(to); !ok {
				return nil, fmt.Errorf("address invalid: %s", to)
			}
		}
		tidyResult.To = args.To
	}

	if args.Subject == "" {
		return nil, errors.New("you need a subject")
	} else {
		tidyResult.Subject = args.Subject
	}

	if args.Content == "" {
		if args.ContentPath != "" {
			content, e := os.ReadFile(args.ContentPath)
			if e != nil {
				return nil, e
			}
			tidyResult.Content = string(content)
		}
	} else {
		tidyResult.Content = args.Content
	}

	tidyResult.Attaches = args.Attaches
	tidyResult.GenerateAuth = args.GenerateAuth
	tidyResult.AuthPath = args.AuthPath

	if args.ContentType != "" {
		tidyResult.ContentType = args.ContentType
	} else {
		tidyResult.ContentType = "text/plain"
	}
	return &tidyResult, nil
}
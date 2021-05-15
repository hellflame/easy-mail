package main

import (
	"errors"
	"fmt"
	"github.com/hellflame/argparse"
	"io/ioutil"
	"net"
	"os"
	"path"
	"strconv"
)

type RawArgs struct {
	From         string
	To           []string
	Subject      string
	ContentType  string
	Content      string
	ContentPath  string
	Attaches     []string
	SMTPServer   string
	Password     string
	AuthPath     string
	GenerateAuth bool
	ShowVersion  bool
}

type TidyArgs struct {
	From         string
	To           []string
	Subject      string
	Content      string
	ContentType  string
	Attaches     []string
	SMTPHosts    []string
	SMTPPorts    []int
	Password     string
	GenerateAuth bool
	AuthPath     string
	ShowVersion  bool
}

func parseArgs(input []string) (args *RawArgs, e error) {
	DefaultAuthPath := ".easy-mail.cred"
	home, e := os.UserHomeDir()
	if e == nil && home != "" {
		DefaultAuthPath = path.Join(home, DefaultAuthPath)
	}
	emailValidator := func(email string) error {
		if _, ok := validateEmailAddress(email); !ok {
			return fmt.Errorf("invalid email '%s'", email)
		}
		return nil
	}
	parser := argparse.NewParser(NAME, "easily send mail from command line",
		&argparse.ParserConfig{EpiLog: "more info @ https://github.com/hellflame/easy-mail", AddShellCompletion: true})
	from := parser.String("f", "from", &argparse.Option{Help: "email send from", Validate: emailValidator})
	to := parser.Strings("t", "to", &argparse.Option{Help: "recv address list", Validate: emailValidator})
	subject := parser.String("s", "subject", &argparse.Option{Help: "email subject"})
	content := parser.String("c", "content", &argparse.Option{Help: "email content"})
	contentPath := parser.String("", "content-path", &argparse.Option{Help: "email content path", Meta: "PATH"})
	contentType := parser.String("", "content-type", &argparse.Option{Help: "email content type", Meta: "TYPE", Default: "text/plain"})
	attach := parser.Strings("", "attach", &argparse.Option{Help: "attach file path list", Meta: "PATH"})
	smtp := parser.String("", "smtp", &argparse.Option{Help: "manually set smtp address like: smtp.abc.com:465 it can be auto find if not set",
		Validate: func(arg string) error {
			_, port, e := net.SplitHostPort(arg)
			if e != nil {
				return fmt.Errorf("invalid smtp server: %s", e.Error())
			}
			_, e = strconv.Atoi(port)
			if e != nil {
				return e
			}
			return nil
		}})
	password := parser.String("p", "password", &argparse.Option{Help: "email password"})
	generateAuth := parser.Flag("g", "generate", &argparse.Option{Help: "save auth to file for simple use"})
	authPath := parser.String("a", "auth", &argparse.Option{Help: "auth file path", Default: DefaultAuthPath, Meta: "PATH"})
	showVersion := parser.Flag("v", "version", &argparse.Option{Help: fmt.Sprintf("show version of %s", NAME)})
	e = parser.Parse(input)
	if e != nil {
		return nil, e
	}
	args = &RawArgs{
		From:         *from,
		To:           *to,
		Subject:      *subject,
		Attaches:     *attach,
		SMTPServer:   *smtp,
		Password:     *password,
		AuthPath:     *authPath,
		GenerateAuth: *generateAuth,
		Content:      *content,
		ContentPath:  *contentPath,
		ContentType:  *contentType,
		ShowVersion:  *showVersion,
	}
	return
}

func tidyArgs(args *RawArgs) (*TidyArgs, error) {
	tidyResult := TidyArgs{
		ShowVersion:  args.ShowVersion,
		Password:     args.Password,
		GenerateAuth: args.GenerateAuth,
		To:           args.To,
		Subject:      args.Subject,
		Content:      args.Content,
		Attaches:     args.Attaches,
		AuthPath:     args.AuthPath,
		ContentType:  args.ContentType,
	}
	if tidyResult.ShowVersion {
		// break point: show version
		return &tidyResult, nil
	}
	if args.From != "" {
		tidyResult.From = args.From
		if args.SMTPServer == "" {
			mailBox, _ := validateEmailAddress(args.From)
			hosts := guessSmtpHosts(mailBox)
			if len(hosts) == 0 {
				return nil, fmt.Errorf("can't find mx servers for %s", mailBox)
			}

			tidyResult.SMTPHosts = hosts
			tidyResult.SMTPPorts = []int{465, 25, 587}
		}
	}
	if args.SMTPServer != "" {
		host, p, _ := net.SplitHostPort(args.SMTPServer)
		port, _ := strconv.Atoi(p)
		tidyResult.SMTPHosts = []string{host}
		tidyResult.SMTPPorts = []int{port}
	}
	if args.GenerateAuth {
		// break point: generate auth
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

	if args.Content == "" {
		if args.ContentPath != "" {
			content, e := ioutil.ReadFile(args.ContentPath)
			if e != nil {
				return nil, e
			}
			tidyResult.Content = string(content)
		}
	}
	if tidyResult.Subject == "" {
		return nil, fmt.Errorf("subject is needed")
	}
	if len(tidyResult.To) == 0 {
		return nil, fmt.Errorf("need target to send email")
	}

	return &tidyResult, nil
}

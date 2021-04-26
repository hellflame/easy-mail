package main

import (
	"io/ioutil"
	"path"
	"testing"
)

func Test_Tidy(t *testing.T) {
	DefaultAuthPath = path.Join(t.TempDir(), "not-exist")
	_, e := tidyArgs(&RawArgs{})
	if e == nil {
		t.Error("failed")
		return
	}
	_, e = tidyArgs(&RawArgs{
		To: []string{"a.b"},
	})
	if e == nil {
		t.Error("failed")
		return
	}
	_, e = tidyArgs(&RawArgs{
		To: []string{"w@a.b"},
	})
	if e == nil {
		t.Error("failed")
		return
	}
	_, e = tidyArgs(&RawArgs{
		To:      []string{"w@a.b"},
		Subject: "this is subject",
	})
	if e == nil {
		t.Error("failed")
		return
	}
	_, e = tidyArgs(&RawArgs{
		To:      []string{"w@a.b"},
		Subject: "this is subject",
		From:    "a.com",
	})
	if e == nil || e.Error() != "invalid from address format" {
		t.Error("failed")
		return
	}
	_, e = tidyArgs(&RawArgs{
		To:      []string{"w@a.b"},
		Subject: "this is subject",
		From:    "9@a.com",
	})
	if e == nil {
		t.Error("failed")
		return
	}
	_, e = tidyArgs(&RawArgs{
		To:       []string{"w@a.b"},
		Subject:  "this is subject",
		From:     "9@a.com",
		Password: "123456",
	})
	if e == nil {
		t.Error("fialed")
		return
	}
	args, e := tidyArgs(&RawArgs{
		To:       []string{"w@a.b"},
		Subject:  "this is subject",
		From:     "9@qq.com",
		Password: "123456",
	})
	if e != nil || args.Password != "123456" {
		t.Error("failed")
		return
	}

	p := path.Join(t.TempDir(), "content")
	ioutil.WriteFile(p, []byte("hellflame is fine"), 0600)

	args, e = tidyArgs(&RawArgs{
		To:          []string{"w@a.b"},
		Subject:     "this is subject",
		From:        "9@qq.com",
		Password:    "123456",
		ContentPath: p,
	})
	if e != nil || args.Content != "hellflame is fine" {
		t.Error("failed to load content")
		return
	}
	args, e = tidyArgs(&RawArgs{
		To:          []string{"w@a.b"},
		Subject:     "this is subject",
		From:        "9@qq.com",
		Password:    "123456",
		Content:     "A",
		ContentType: "text/html",
	})
	if e != nil || args.Content != "A" || args.ContentType == "" {
		t.Error("failed to set content")
		return
	}
	args, e = tidyArgs(&RawArgs{
		To:         []string{"w@a.b"},
		Subject:    "this is subject",
		From:       "9@qq.com",
		Password:   "123456",
		SMTPServer: "smtp.a.b:253",
	})
	if e != nil || args.SMTPHosts[0] != "smtp.a.b" {
		t.Error("failed to set smtp server")
		return
	}
	args, e = tidyArgs(&RawArgs{
		To:         []string{"w@a.b"},
		Subject:    "this is subject",
		From:       "9@qq.com",
		Password:   "123456",
		SMTPServer: "smtp.a.b",
	})
	if e == nil {
		t.Error("failed to tell bad smtp server")
		return
	}
	ioutil.WriteFile(DefaultAuthPath, []byte(`{"From": "a@b.c", "Password": "123", "Server": "a.c:26"}`), 0600)
	args, e = tidyArgs(&RawArgs{
		To:      []string{"w@a.b"},
		Subject: "this is subject",
	})
	if e != nil || args.Password != "123" {
		t.Error("failed to load saved auth")
		return
	}
}

func Test_Parse(t *testing.T) {
	args, e := parseArgs([]string{"start", "-g"})
	if e != nil || !args.GenerateAuth {
		t.Error("failed to parse args")
		return
	}

}

package main

import (
    "fmt"
    "path"
    "strings"
    "testing"
)

func Test_fetchMxHosts(t *testing.T) {
	hosts := fetchMxHosts("qq.com")
	if len(hosts) == 0 {
		t.Error("failed to look up qq.com")
		return
	}
	hosts = fetchMxHosts("abc")
	if len(hosts) != 0 {
		t.Error("should be empty")
		return
	}
}

func Test_fileExist(t *testing.T) {
	if fileExist(fmt.Sprintf("%s/not-exist", t.TempDir())) {
		t.Error("should not exist")
		return
	}
}

func Test_SaveLoadAuth(t *testing.T) {
    p := path.Join(t.TempDir(), "cred")
	if e := saveAuth(p, "hellflame@a.b", "password", "smtp.a.b:25"); e != nil {
		t.Error(e.Error())
		return
	}
	from, password, host, port, e := loadAuth(p)
	if e != nil || from == "" || password == "" || host == "" || port == 0 {
		t.Error("failed to load")
		return
	}
}

func Test_validateEmailAddress(t *testing.T) {
	host, valid := validateEmailAddress("hellflame@66.com")
	if !valid || host != "66.com" {
		t.Error("failed to validate address")
	}
}

func Test_guessSmtpHosts(t *testing.T) {
	hasSmtp := false
	for _, host := range guessSmtpHosts("qq.com") {
		if strings.HasPrefix(host, "smtp.") {
			hasSmtp = true
		}
	}
	if !hasSmtp {
		t.Error("failed to add smtp")
	}
}

package main

import (
	"encoding/json"
	"errors"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
)


var DefaultAuthPath = ".mailall.auth"

func init () {
	home, e := os.UserHomeDir()
	if e == nil && home != "" {
		DefaultAuthPath = path.Join(home, DefaultAuthPath)
	}
}

func validateEmailAddress(address string) (mailbox string, valid bool) {
	splits := strings.Split(address, "@")
	if len(splits) != 2 {
		valid = false
		return
	}
	return splits[1], true
}

func fetchMxHosts(mailBox string) []string {
	rec, e := net.LookupMX(mailBox)
	if e != nil {
		return nil
	}
	var hosts []string
	for _, raw := range rec {
		hosts = append(hosts, raw.Host[:len(raw.Host) - 1])
	}
	return hosts
}


func fileExist(path string) bool {
	if _, e := os.Stat(path); e != nil {
		return false
	}
	return true
}

func saveAuth(path, from, password, smtp string) error {
	if path == "" {
		path = DefaultAuthPath
	}
	if from == "" || password == "" || smtp == "" {
		return errors.New("not enough to saveAuth")
	}
	template := map[string]string{
		"From": from,
		"Password": password,
		"Server": smtp,
	}
	dump, e := json.Marshal(&template)
	if e != nil {
		return e
	}
	return os.WriteFile(path, dump, 0600)
}

func loadAuth(path string) (from, password, host string, port int, e error) {
	if path == "" {
		path = DefaultAuthPath
	}
	read, e := os.ReadFile(path)
	if e != nil {
		return
	}
	result := make(map[string]string)
	e = json.Unmarshal(read, &result)
	if e != nil {
		return
	}
	from = result["From"]
	password = result["Password"]
	host, p, e := net.SplitHostPort(result["Server"])
	if e != nil {
		return
	}
	port, e = strconv.Atoi(p)
	if e != nil {
		return
	}
	return
}
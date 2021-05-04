package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
)


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
		hosts = append(hosts, raw.Host[:len(raw.Host)-1])
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
	if from == "" || password == "" || smtp == "" {
		return errors.New("not enough to saveAuth")
	}
	template := map[string]string{
		"From":     from,
		"Password": password,
		"Server":   smtp,
	}
	dump, e := json.Marshal(&template)
	if e != nil {
		return e
	}
	return ioutil.WriteFile(path, dump, 0600)
}

func loadAuth(path string) (from, password, host string, port int, e error) {
	read, e := ioutil.ReadFile(path)
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

// find possible smtp servers
func guessSmtpHosts(mailBox string) []string {
	hosts := fetchMxHosts(mailBox)
	if len(hosts) == 0 {
		return nil
	}
	straight := fmt.Sprintf("smtp.%s", mailBox)
	hosts = append([]string{straight}, hosts...)
	hostPool := map[string]int8{
		mailBox:  1,
		straight: 1,
	}
	for _, host := range hosts {
		left := strings.Split(host, ".")[1:]
		possibleByMX := strings.Join(left, ".")
		if _, exist := hostPool[possibleByMX]; !exist {
			hosts = append([]string{fmt.Sprintf("smtp.%s", possibleByMX)}, hosts...)
			hostPool[possibleByMX] = 1
		}
	}
	return hosts
}

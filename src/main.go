package main

import (
	"os/exec"
	"strings"
	"net/http"
	"log"
	"bufio"
	"os"
	"os/user"
	"bytes"
	"net/url"
	"encoding/json"
)

var serverChan = make(chan os.Signal, 1)

var configFile, _ = os.Open(func() (string) {
	if user, error := user.Current(); error == nil {
		return user.HomeDir
	}
	if home := os.Getenv("HOME"); home != "" {
		return home
	}

	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err == nil {
		return strings.TrimSpace(stdout.String())
	}
	return ""
}() + "/.muttrc")

var port = func() (string) {
	reader := bufio.NewReader(configFile)
	for {
		if line, _, _ := reader.ReadLine(); line != nil {
			paramAndValue := string(line)
			if (strings.Contains(paramAndValue, "server_port")) {
				return paramAndValue[strings.Index(paramAndValue, "=") + 1:]
			}
		}else {
			break
		}
	}

	return "65530"
}()

var logFile, _ = os.OpenFile("email-server.log", os.O_RDWR | os.O_APPEND | os.O_CREATE, os.ModePerm)
var logger = log.New(logFile, "[email-server] ", log.Ldate | log.Ltime | log.Lshortfile)

func main() {

	go http.ListenAndServe(":" + port, http.HandlerFunc(handle))
	logger.Printf("------------ email server started at port %s---------", port)
	<-serverChan
}

func handle(w http.ResponseWriter, req *http.Request) {
	if (req.Method != "POST") {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if (!strings.HasPrefix(req.Header.Get("Authorization"), "Wecash")) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	uri, _ := url.ParseRequestURI(req.RequestURI)
	if (uri.Path != "/") {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	req.ParseForm()
	if (req.Header.Get("Content-Type") == "application/json" || req.Header.Get("content-type") == "application/json") {
		buf := new(bytes.Buffer)
		for {
			bt := make([]byte, 1024)
			if n, error := req.Body.Read(bt); error == nil {
				if (n <= 0) {
					break;
				}
				buf.Write(bt[:n])
			}else {
				buf.Write(bt[:n])
				break
			}
		}
		var body1 map[string]string
		json.Unmarshal(buf.Bytes(), &body1)
		for key := range body1 {
			req.Form.Add(key, body1[key])
		}
		req.Body.Close()
	}

	var receivers []string
	var ok bool
	if receivers, ok = req.Form["receivers"]; !ok {
		w.WriteHeader(http.StatusNoContent)
		return;
	}
	var subject []string
	var content []string
	if subject, ok = req.Form["subject"]; !ok {
		subject = []string{""}
	}
	if content, ok = req.Form["content"]; !ok {
		content = []string{""}
	}
	go send(strings.Join(subject, ""), strings.Join(content, ""), receivers)
	defer func() {
		if r := recover(); r != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Printf("server internal exception:%s", r)
		}
	}()
}

func send(subject, content string, receivers []string) {
	logger.Printf("send eamil use subjec:%s content:%s receivers:%s", subject, content, strings.Join(receivers, " "))

	//echoCmd := exec.Command("echo", content);
	//muttCmd := exec.Command("mutt", strings.Join(receivers, " "), "-s", subject)
	//muttCmd.Stdin, _ = echoCmd.StdoutPipe()
	//muttCmd.Start()
	//echoCmd.Run()
	//muttCmd.Wait()

	exec.Command("/bin/sh", "-c", strings.Join([]string{"echo", content, "|", "mutt", strings.Join(receivers, " "), "-s", subject}, " ")).Start()
}
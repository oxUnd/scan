package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func dialTimeout(network, addr string) (net.Conn, error) {
	var timeout_ = time.Duration(1 * time.Second)
	return net.DialTimeout(network, addr, timeout_)
}

func Head(url string) (*http.Response, error) {

	transport := http.Transport{
		Dial: dialTimeout,
	}

	client := http.Client{
		Transport: &transport,
	}

	return client.Head(url)
}

type Task struct {
	ip     string
	url_   string
	result map[string]string
	status int
}

func (t *Task) Run(r chan (Task)) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Sprintln(os.Stderr, e)
		}
	}()
	url_ := t.GetUrl()
	log.Println("scan ", url_)
	res, err := Head(url_)
	if err != nil {
		r <- *t
		panic(err)
		return
	}
	t.status = res.StatusCode
	t.result["server"] = res.Header.Get("Server")
	r <- *t
}

func (t Task) Ip() string {
	return t.ip
}

func (t *Task) GetUrl() string {
	if len(t.url_) > 0 {
		return t.url_
	}
	return "http://" + t.ip + "/"
}

func (t Task) Status() int {
	return t.status
}

func (t Task) Result() map[string]string {
	return t.result
}

var ips string

func init() {
	flag.StringVar(&ips, "ips", "10.10.10.1-255", "please enter IP range.")
}

func ParseIps(ipRange string) []string {
	ips := []string{}
	p := strings.Index(ipRange, "-")
	if p < 0 {
		ips = append(ips, ipRange)
	} else {
		reg := regexp.MustCompile("(\\d+\\.\\d+\\.\\d+\\.)(\\d+)-(\\d+)")
		m := reg.FindStringSubmatch(ipRange)
		if len(m) > 0 {
			prefix := m[1]
			start, err := strconv.Atoi(m[2])
			if err != nil {
				panic(err)
			}
			end, err := strconv.Atoi(m[3])
			if err != nil {
				panic(err)
			}

			if start >= end {
				panic("error")
			}

			for v := start; v <= end; v++ {
				ips = append(ips, prefix+strconv.Itoa(v))
			}
		}
	}
	return ips
}

func main() {
	flag.Parse()
	range_ := ParseIps(ips)
	count := len(range_)

	if count == 0 {
		panic("Must given a ip.")
	}

	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	var task Task
	r := make(chan (Task))
	for _, v := range range_ {
		task = Task{
			ip:     v,
			result: make(map[string]string),
		}

		go func(task Task) {
			task.Run(r)
		}(task)
	}

	resultFile, err := os.OpenFile("./res.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}

	for i := 0; i < count; i++ {
		ret := <-r
		result := ret.Result()
		if server, ok := result["server"]; ok {
			resultFile.WriteString(fmt.Sprintln(ret.Ip(), " ", server))
		}
	}

	resultFile.Close()
}

package main

import (
	"./core"
	"fmt"
	opt "github.com/docopt/docopt-go"
	"regexp"
	"strconv"
	"strings"
)

func usage() string {
	return `

Usage:
    scan --ip <ip> [--port <port>][--timeout <timeout>]
    scan --ip-range <range> [--port <port>][--timeout <timeout>]

Options:
  -h --help             Show help
  --version             Show version
  --ip <ip>             Will scan ip
  --ip-range <range>    Will scan ip in ip range
  --port <port>         scan port
  --timeout <timeout>   timeout, default 1s
    `
}

func getPort(s string) []int {
	sep := ","
	if strings.Index(s, sep) == -1 {
		result, _ := strconv.Atoi(s)
		return []int{result}
	}
	splits := strings.Split(s, sep)
	ret := []int{}
	for _, ss := range splits {
		result, _ := strconv.Atoi(ss)
		ret = append(ret, result)
	}
	return ret
}

func getAllIp(ipRange string) []string {
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

func run(id int, ip string, ports []int, reports *[]core.Report, timeout string) {
	task := core.NewTask(id, ip, ports, false)
	worker := core.NewWorker(1, &task, reports, timeout)
	worker.ScanHttp()
}

func main() {

	arguments, err := opt.Parse(usage(), nil, true, "2.0", false)

	if err != nil {
		panic(err)
	}

	ports := []int{80}
	if arguments["--port"] != nil {
		ports = getPort(arguments["--port"].(string))
	}

	timeout := "1s"
	if arguments["--timeout"] != nil {
		timeout = arguments["--timeout"].(string)
	}

	reports := []core.Report{}
	if arguments["--ip"] == nil {
		//ip range
		ipRange := arguments["--ip-range"].(string)
		ips := getAllIp(ipRange)
		fmt.Println(ips)
		howManyIp := len(ips)
		wait := make(chan (bool))
		for idx, ip := range ips {
			go func(idx int, ip string) {
				run(idx, ip, ports, &reports, timeout)
				wait <- true
			}(idx, ip)
		}
		//wait all child task done.
		for i := 0; i < howManyIp; i++ {
			<-wait
		}

	} else {
		//ip
		run(0, arguments["--ip"].(string), ports, &reports, timeout)
	}

	for _, r := range reports {
		r.Dump()
	}

	fmt.Println("done")
}

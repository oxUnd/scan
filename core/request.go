package core

import (
	"time"
	"net"
	"net/http"
)

// use: HttpHeadRequest("http://www.baidu.com/", "1s")
func HttpHeadRequest(url string, timeout string) (*http.Response, error) {

	transport := http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			timeout_, err := time.ParseDuration(timeout)
			if err != nil {
				panic(err)
			}
			return net.DialTimeout(network, addr, timeout_)

		},
	}

	client := http.Client{
		Transport: &transport,
	}

	return client.Head(url)
}
